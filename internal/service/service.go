package service

import (
	"context"
	"fmt"
	"time"

	dto "github.com/Piccadilly98/subscription_service/internal/models/dto"
	"github.com/Piccadilly98/subscription_service/internal/models/entities"
	"github.com/Piccadilly98/subscription_service/internal/storage"
)

type Service struct {
	storage *storage.Storage
}

func NewService(st *storage.Storage) (*Service, error) {
	if st == nil {
		return nil, fmt.Errorf("storage connot be nil")
	}
	return &Service{
		storage: st,
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
	s.storage.Cache.AddToCacheByID(id)
	entitie, err := s.storage.Db.GetSubscriptionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	res := dto.FromEntityToSubResp(entitie)
	return res, nil
}

func (s *Service) GetExsistBySubID(ctx context.Context, id string) (bool, error) {
	if s.storage.Cache.CheckBySubID(id) {
		fmt.Println("взяли из кэша")
		return true, nil
	}

	exists, err := s.storage.Db.GetExsistBySubID(ctx, id)
	if err != nil {
		return false, err
	}

	if exists {
		s.storage.Cache.AddToCacheByID(id)
	}

	return exists, nil
}

func (s *Service) GetSubInfoDTOByID(ctx context.Context, id string) (*dto.SubscriptionResponse, error) {
	entitie, err := s.storage.Db.GetSubscriptionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	res := dto.FromEntityToSubResp(entitie)
	return res, nil
}

func (s *Service) UpdateSubscription(ctx context.Context, body *dto.UpdateSubscriptionRequest, id string) (*entities.UpdateSubscription, error) {
	err := body.Validate()
	if err != nil {
		return nil, err
	}
	update, err := body.ToEntitie()
	if err != nil {
		return nil, err
	}
	read, err := s.storage.Db.GetSubscriptionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// защита того что старую подписку мы не меняем - опционально
	if read.IsEnded {
		return nil, fmt.Errorf("cannot update ended subscription")
	}

	if update.EndDate != nil {
		// чекаем что энд позже старт
		if read.StartDate.Compare(*update.EndDate) == 1 {
			return nil, fmt.Errorf("invalid end_data: end_data cannot be before start_date")
		}
		// если энд в прошлом то подписка остановлена
		if update.EndDate.Compare(time.Now()) == -1 {
			update.Ended = getBoolPtr(true)
		}
	}

	if update.Ended != nil {

		// откинули заранее завершенную подписку
		if read.StartDate.Compare(time.Now()) == 1 {
			return nil, fmt.Errorf("can't stop a subscription that hasn't started")
		}
		// откидываем если подписка закончилась
		if read.IsEnded {
			return nil, fmt.Errorf("can't stop ended subscription")
		}

		update.EndDate = GetTimePtr(time.Now())
	}
	// err = s.storage.Db.UpdateSubscription(ctx, update, id)
	// if err != nil {
	// 	return nil, err
	// }
	return update, nil
}
