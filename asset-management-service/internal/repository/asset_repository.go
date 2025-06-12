package repository

import (
	"asset-management-service/internal/models"
	"context"
	"github.com/jmoiron/sqlx"
)

type AssetRepository interface {
	GetAll(ctx context.Context) ([]models.Asset, error)
	GetByIP(ctx context.Context, ip string) (*models.Asset, error)
	Create(ctx context.Context, asset *models.Asset) error
	Update(ctx context.Context, ip string, asset *models.Asset) error
	Delete(ctx context.Context, ip string) error
	SetManaged(ctx context.Context, ip string, managed bool) error

	GetServices(ctx context.Context, ip string) ([]models.Service, error)
	AddService(ctx context.Context, ip string, service *models.Service) error
	UpdateService(ctx context.Context, serviceID int, service *models.Service) error
	DeleteService(ctx context.Context, serviceID int) error
	SetServiceManaged(ctx context.Context, serviceID int, managed bool) error
	CheckAllServicesManaged(ctx context.Context, ip string) (bool, error)
}

type assetRepo struct {
	db *sqlx.DB
}

func NewAssetRepository(db *sqlx.DB) AssetRepository {
	return &assetRepo{db: db}
}

func (r *assetRepo) GetAll(ctx context.Context) ([]models.Asset, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT ip_address, mac_address, hostname, asset_type, location, operating_system, first_seen, last_seen, is_managed FROM assets`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []models.Asset
	for rows.Next() {
		var a models.Asset
		err := rows.Scan(&a.IPAddress, &a.MACAddress, &a.Hostname, &a.AssetType, &a.Location, &a.OperatingSystem, &a.FirstSeen, &a.LastSeen, &a.IsManaged)
		if err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}
	return assets, nil
}

func (r *assetRepo) GetByIP(ctx context.Context, ip string) (*models.Asset, error) {
	var a models.Asset
	err := r.db.QueryRowContext(ctx, `SELECT ip_address, mac_address, hostname, asset_type, location, operating_system, first_seen, last_seen, is_managed FROM assets WHERE ip_address = $1`, ip).
		Scan(&a.IPAddress, &a.MACAddress, &a.Hostname, &a.AssetType, &a.Location, &a.OperatingSystem, &a.FirstSeen, &a.LastSeen, &a.IsManaged)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *assetRepo) Create(ctx context.Context, asset *models.Asset) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO assets(ip_address, mac_address, hostname, asset_type, location, operating_system, first_seen, last_seen, is_managed) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		asset.IPAddress, asset.MACAddress, asset.Hostname, asset.AssetType, asset.Location, asset.OperatingSystem, asset.FirstSeen, asset.LastSeen, asset.IsManaged)
	return err
}

func (r *assetRepo) Update(ctx context.Context, ip string, asset *models.Asset) error {
	_, err := r.db.ExecContext(ctx, `UPDATE assets SET mac_address=$1, hostname=$2, asset_type=$3, location=$4, operating_system=$5, first_seen=$6, last_seen=$7, is_managed=$8 WHERE ip_address=$9`,
		asset.MACAddress, asset.Hostname, asset.AssetType, asset.Location, asset.OperatingSystem, asset.FirstSeen, asset.LastSeen, asset.IsManaged, ip)
	return err
}

func (r *assetRepo) Delete(ctx context.Context, ip string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM assets WHERE ip_address=$1`, ip)
	return err
}

func (r *assetRepo) SetManaged(ctx context.Context, ip string, managed bool) error {
	_, err := r.db.ExecContext(ctx, `UPDATE assets SET is_managed=$1 WHERE ip_address=$2`, managed, ip)
	return err
}

// Service Methods

func (r *assetRepo) GetServices(ctx context.Context, ip string) ([]models.Service, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, port, protocol, description, is_managed FROM services WHERE asset_ip = $1`, ip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []models.Service
	for rows.Next() {
		var s models.Service
		err := rows.Scan(&s.ID, &s.Name, &s.Port, &s.Protocol, &s.Description, &s.IsManaged)
		if err != nil {
			return nil, err
		}
		s.AssetIP = ip
		services = append(services, s)
	}
	return services, nil
}

func (r *assetRepo) AddService(ctx context.Context, ip string, service *models.Service) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO services(asset_ip, name, port, protocol, description, is_managed) VALUES ($1, $2, $3, $4, $5, $6)`,
		ip, service.Name, service.Port, service.Protocol, service.Description, service.IsManaged)
	return err
}

func (r *assetRepo) UpdateService(ctx context.Context, serviceID int, service *models.Service) error {
	_, err := r.db.ExecContext(ctx, `UPDATE services SET name=$1, port=$2, protocol=$3, description=$4, is_managed=$5 WHERE id=$6`,
		service.Name, service.Port, service.Protocol, service.Description, service.IsManaged, serviceID)
	return err
}

func (r *assetRepo) DeleteService(ctx context.Context, serviceID int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM services WHERE id=$1`, serviceID)
	return err
}

func (r *assetRepo) SetServiceManaged(ctx context.Context, serviceID int, managed bool) error {
	_, err := r.db.ExecContext(ctx, `UPDATE services SET is_managed=$1 WHERE id=$2`, managed, serviceID)
	return err
}

func (r *assetRepo) CheckAllServicesManaged(ctx context.Context, ip string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM services WHERE asset_ip = $1 AND is_managed = FALSE`, ip).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}
