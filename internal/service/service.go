package service

import (
	"context"
	"fmt"
	"time"

	logger "github.com/Piccadilly98/subscription_service/internal/loger"
	dto "github.com/Piccadilly98/subscription_service/internal/models/dto"
	"github.com/Piccadilly98/subscription_service/internal/storage"
)

const (
	LevelForChangedSubs = "important"
	LevelForDeleteSubs  = "warning"
)

type Service struct {
	storage      *storage.Storage
	changeLogger *logger.Logger
}

func NewService(st *storage.Storage, changeLogger *logger.Logger) (*Service, error) {
	if st == nil {
		return nil, fmt.Errorf("storage connot be nil")
	}
	return &Service{
		storage:      st,
		changeLogger: changeLogger,
	}, nil
}

func (s *Service) CreateSubsription(ctx context.Context, req *dto.CreateSubscriptionsRequest) (*dto.SubscriptionResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	mod, err := req.ToEntitie()
	if err != nil {
		return nil, err
	}

	id, err := s.storage.Db.CreateNewSubscribe(ctx, mod)
	if err != nil {
		return nil, err
	}
	if s.storage.Cache != nil {
		s.storage.Cache.AddToCacheByID(id)
	}
	entitie, err := s.storage.Db.GetSubscriptionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	s.changeLogger.Printf(LevelForChangedSubs, "create new subscription with id: %s", id)
	res := dto.FromEntityToSubResp(entitie)
	return res, nil
}

func (s *Service) GetExsistBySubID(ctx context.Context, id string) (bool, error) {
	if s.storage.Cache != nil {
		if s.storage.Cache.CheckBySubID(id) {
			return true, nil
		}
	}
	exists, err := s.storage.Db.GetExsistBySubID(ctx, id)
	if err != nil {
		return false, err
	}
	if exists && s.storage.Cache != nil {
		s.storage.Cache.AddToCacheByID(id)
	}

	return exists, nil
}

func (s *Service) GetSubInfoDTOByID(ctx context.Context, id string) (*dto.SubscriptionResponse, error) {
	entitie, err := s.storage.Db.GetSubscriptionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if s.storage.Cache != nil {
		s.storage.Cache.AddToCacheByID(id)
	}
	res := dto.FromEntityToSubResp(entitie)
	return res, nil
}

func (s *Service) UpdateSubscription(ctx context.Context, body *dto.UpdateSubscriptionRequest, id string) error {
	err := body.Validate()
	if err != nil {
		return err
	}
	update, err := body.ToEntitie()
	if err != nil {
		return err
	}
	read, err := s.storage.Db.GetSubscriptionByID(ctx, id)
	if err != nil {
		return err
	}

	if read.IsEnded {
		return fmt.Errorf("cannot update ended subscription")
	}

	if update.EndDate != nil {
		if read.StartDate.Compare(*update.EndDate) == 1 {
			return fmt.Errorf("invalid end_date: end_data cannot be before start_date")
		}
		if update.EndDate.Compare(time.Now()) == -1 {
			update.Ended = getBoolPtr(true)
		}
	}

	if update.Ended != nil {
		if read.StartDate.Compare(time.Now()) == 1 {
			return fmt.Errorf("can't stop a subscription that hasn't started")
		}
		if read.IsEnded {
			return fmt.Errorf("can't stop ended subscription")
		}

		update.EndDate = GetTimePtr(time.Now())
	}
	err = s.storage.Db.UpdateSubscription(ctx, update, id)
	if err != nil {
		return err
	}
	s.changeLogger.Printf(LevelForChangedSubs, "update subscription with id: %s", id)
	if s.storage.Cache != nil {
		s.storage.Cache.AddToCacheByID(id)
	}
	return nil
}

func (s *Service) DeleteSubByID(ctx context.Context, id string) error {
	err := s.storage.Db.DeleteRowBySubID(ctx, id)
	if err != nil {
		return err
	}
	s.changeLogger.Printf(LevelForDeleteSubs, "delete subscription with id: %s", id)
	if s.storage.Cache != nil {
		s.storage.Cache.DeleteByID(id)
	}
	return nil
}

func (s *Service) GetSummarySubs(ctx context.Context, req *dto.QueryParamsSummary) (*dto.SummaryResponse, error) {
	entitie, err := req.ToEntitie()
	if err != nil {
		return nil, err
	}

	if entitie.UserID != nil {
		exist, err := s.GetExsistBySubID(ctx, *entitie.UserID)
		if err != nil {
			return nil, err
		}
		if !exist {
			return nil, fmt.Errorf("invalid user_id")
		}
	}
	if entitie.EndDate.Before(entitie.StartDate) {
		return nil, fmt.Errorf("end_date cannot be before start_date")
	}

	sum, err := s.storage.Db.GetSumaryByParam(ctx, entitie)
	if err != nil {
		return nil, err
	}
	if sum == nil {
		sum = getIntPtr(0)
	}
	res := dto.FromEntityToSummaryResponse(entitie, *sum)
	return res, nil
}

func (s *Service) CheckHealh(ctx context.Context) *dto.CheckHealth {
	statusServer := "ok"
	statusDB := "ok"
	err := s.storage.Db.PingWithCtx(ctx)
	if err != nil {
		statusServer = "Service Unavailable"
		statusDB = "does not respond"
	}
	return dto.ToCheckHealthDTO(statusServer, statusDB, err)
}
