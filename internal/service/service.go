package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	dto "github.com/Piccadilly98/subscription_service/internal/models/dto"
	"github.com/Piccadilly98/subscription_service/internal/storage/cache"
	"github.com/Piccadilly98/subscription_service/internal/storage/data_base"
)

const (
	LevelForChangedSubs = "important"
	LevelForDeleteSubs  = "warning"
)

type Service struct {
	db               *data_base.DataBase
	cache            *cache.Cache
	changeLogger     *log.Logger
	dbCriticalLogger *log.Logger
}

func NewService(db *data_base.DataBase, cache *cache.Cache) (*Service, error) {
	if db == nil {
		return nil, fmt.Errorf("db cannot be nil")
	}
	return &Service{
		db:               db,
		cache:            cache,
		changeLogger:     log.New(os.Stdout, "[UPDATE SUBS INFO] ", log.Ldate|log.Ltime),
		dbCriticalLogger: log.New(os.Stderr, "[DB PING ERROR] ", log.Ldate|log.Ltime),
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

	id, err := s.db.CreateNewSubscribe(ctx, mod)
	if err != nil {
		return nil, err
	}
	if s.cache != nil {
		s.cache.AddToCacheByID(id)
	}
	entitie, err := s.db.GetSubscriptionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	s.changeLogger.Printf("INFO: Create new subscription with id: %s", id)
	res := dto.FromEntityToSubResp(entitie)
	return res, nil
}

func (s *Service) GetExsistBySubID(ctx context.Context, id string) (bool, error) {
	if s.cache != nil {
		if s.cache.CheckBySubID(id) {
			return true, nil
		}
	}
	exists, err := s.db.GetExsistBySubID(ctx, id)
	if err != nil {
		return false, err
	}
	if exists && s.cache != nil {
		s.cache.AddToCacheByID(id)
	}

	return exists, nil
}

func (s *Service) GetSubInfoDTOByID(ctx context.Context, id string) (*dto.SubscriptionResponse, error) {
	entitie, err := s.db.GetSubscriptionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if s.cache != nil {
		s.cache.AddToCacheByID(id)
	}
	res := dto.FromEntityToSubResp(entitie)
	return res, nil
}

func (s *Service) UpdateSubscription(ctx context.Context, body *dto.UpdateSubscriptionRequest, id string) (*dto.SubscriptionResponse, error) {
	err := body.Validate()
	if err != nil {
		return nil, err
	}
	update, err := body.ToEntitie()
	if err != nil {
		return nil, err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		tx.Rollback()

	}()
	read, err := s.db.GetSubscriptionForUpdateTX(ctx, id, tx)
	if err != nil {
		return nil, err
	}

	if read.IsEnded {
		return nil, fmt.Errorf("cannot update ended subscription")
	}

	if update.EndDate != nil {
		if read.StartDate.Compare(*update.EndDate) == 1 {
			return nil, fmt.Errorf("invalid end_date: end_data cannot be before start_date")
		}
		if update.EndDate.Compare(time.Now()) == -1 {
			update.Ended = getBoolPtr(true)
		}
	}
	if update.Ended != nil && update.EndDate == nil {
		if read.StartDate.Compare(time.Now()) == 1 {
			return nil, fmt.Errorf("can't stop a subscription that hasn't started")
		}
		if read.IsEnded {
			return nil, fmt.Errorf("can't stop ended subscription")
		}

		update.EndDate = GetTimePtr(time.Now())
	}
	err = s.db.UpdateSubscriptionTX(ctx, update, id, tx)
	if err != nil {
		return nil, err
	}
	res, err := s.db.GetSubscriptionTX(ctx, id, tx)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}
	s.changeLogger.Printf("INFO: Update subscription with id: %s", id)
	if s.cache != nil {
		s.cache.AddToCacheByID(id)
	}
	return dto.FromEntityToSubResp(res), nil
}

func (s *Service) DeleteSubByID(ctx context.Context, id string) error {
	exist, err := s.GetExsistBySubID(ctx, id)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("invalid subscribe_id")
	}
	err = s.db.DeleteRowBySubID(ctx, id)
	if err != nil {
		return err
	}
	s.changeLogger.Printf("WARNING: delete subscription with id: %s", id)
	if s.cache != nil {
		s.cache.DeleteByID(id)
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

	sum, err := s.db.GetSumaryByParam(ctx, entitie)
	if err != nil {
		return nil, err
	}
	if sum == nil {
		sum = getIntPtr(0)
	}
	res := dto.FromEntityToSummaryResponse(entitie, *sum)
	return res, nil
}

func (s *Service) CheckHealh(ctx context.Context) (*dto.CheckHealth, int) {
	statusServer := "ok"
	statusDB := "ok"
	code := http.StatusOK
	err := s.db.PingWithCtx(ctx)
	if err != nil {
		statusServer = "Service Unavailable"
		statusDB = "does not respond"
		code = http.StatusServiceUnavailable
		s.dbCriticalLogger.Printf("CRITICAL: db ping error: %s\n", err.Error())
	}
	return dto.ToCheckHealthDTO(statusServer, statusDB, err), code
}
