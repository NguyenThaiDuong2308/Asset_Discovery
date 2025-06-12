package service

import (
	"asset-management-service/internal/models"
	"asset-management-service/internal/repository"
	"context"
	"errors"
	"time"
)

type AssetService interface {
	GetAllAssets(ctx context.Context) ([]models.Asset, error)
	GetAssetByIP(ctx context.Context, ip string) (*models.Asset, error)
	CreateAsset(ctx context.Context, asset *models.Asset) error
	UpdateAsset(ctx context.Context, ip string, asset *models.Asset) error
	DeleteAsset(ctx context.Context, ip string) error
	ManageAsset(ctx context.Context, ip string, managed bool) error

	GetServices(ctx context.Context, ip string) ([]models.Service, error)
	AddService(ctx context.Context, ip string, service *models.Service) error
	UpdateService(ctx context.Context, id int, service *models.Service) error
	DeleteService(ctx context.Context, id int) error
	ManageService(ctx context.Context, ip string, id int, managed bool) error
}

type assetService struct {
	repo repository.AssetRepository
}

func NewAssetService(repo repository.AssetRepository) AssetService {
	return &assetService{repo: repo}
}

func (s *assetService) GetAllAssets(ctx context.Context) ([]models.Asset, error) {
	return s.repo.GetAll(ctx)
}

func (s *assetService) GetAssetByIP(ctx context.Context, ip string) (*models.Asset, error) {
	return s.repo.GetByIP(ctx, ip)
}

func (s *assetService) CreateAsset(ctx context.Context, asset *models.Asset) error {
	asset.FirstSeen = time.Now()
	asset.LastSeen = time.Now()
	return s.repo.Create(ctx, asset)
}

func (s *assetService) UpdateAsset(ctx context.Context, ip string, asset *models.Asset) error {
	asset.LastSeen = time.Now()
	return s.repo.Update(ctx, ip, asset)
}

func (s *assetService) DeleteAsset(ctx context.Context, ip string) error {
	return s.repo.Delete(ctx, ip)
}

func (s *assetService) ManageAsset(ctx context.Context, ip string, managed bool) error {
	if managed {
		allManaged, err := s.repo.CheckAllServicesManaged(ctx, ip)
		if err != nil {
			return err
		}
		if !allManaged {
			return errors.New("Cannot mark asset as managed while it contains unmanaged services")
		}
	}
	return s.repo.SetManaged(ctx, ip, managed)
}

func (s *assetService) GetServices(ctx context.Context, ip string) ([]models.Service, error) {
	return s.repo.GetServices(ctx, ip)
}

func (s *assetService) AddService(ctx context.Context, ip string, service *models.Service) error {
	return s.repo.AddService(ctx, ip, service)
}

func (s *assetService) UpdateService(ctx context.Context, id int, service *models.Service) error {
	return s.repo.UpdateService(ctx, id, service)
}

func (s *assetService) DeleteService(ctx context.Context, id int) error {
	return s.repo.DeleteService(ctx, id)
}

func (s *assetService) ManageService(ctx context.Context, ip string, id int, managed bool) error {
	err := s.repo.SetServiceManaged(ctx, id, managed)
	if err != nil {
		return err
	}
	if managed {
		allManaged, err := s.repo.CheckAllServicesManaged(ctx, ip)
		if err != nil {
			return err
		}
		if allManaged {
			return s.repo.SetManaged(ctx, ip, true)
		}
	} else {
		return s.repo.SetManaged(ctx, ip, false)
	}
	return nil
}
