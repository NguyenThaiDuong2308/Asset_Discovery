package service

import (
	"asset-management-service/internal/models"
	"asset-management-service/internal/repository"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type LogAsset struct {
	IPAddress  string    `json:"ip_address"`
	MACAddress string    `json:"mac_address"`
	Hostname   string    `json:"hostname"`
	AssetType  string    `json:"asset_type"`
	FirstSeen  time.Time `json:"first_seen"`
	LastSeen   time.Time `json:"last_seen"`
}

type LogService struct {
	ID       int    `json:"id"`
	AssetIP  string `json:"asset_ip"`
	Name     string `json:"name"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

func FetchAssetsFromLogAnalysis() ([]LogAsset, error) {
	resp, err := http.Get("http://log-analysis-service:8080/api/assets")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Assets []LogAsset `json:"assets"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	return result.Assets, err
}
func FetchServicesFromLogAnalysis() ([]LogService, error) {
	resp, err := http.Get("http://log-analysis-service:8080/api/services")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Services []LogService `json:"services"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	return result.Services, err
}

func SyncFromLogService(ctx context.Context, repo repository.AssetRepository) error {
	logAssets, err := FetchAssetsFromLogAnalysis()
	if err != nil {
		return err
	}

	logServices, err := FetchServicesFromLogAnalysis()
	if err != nil {
		return err
	}

	for _, la := range logAssets {
		existing, err := repo.GetByIP(ctx, la.IPAddress)
		if err != nil {
			newAsset := &models.Asset{
				IPAddress:       la.IPAddress,
				MACAddress:      la.MACAddress,
				Hostname:        la.Hostname,
				AssetType:       la.AssetType,
				Location:        "", // Có thể cập nhật sau
				OperatingSystem: "",
				FirstSeen:       la.FirstSeen,
				LastSeen:        la.LastSeen,
				IsManaged:       false,
			}
			if err := repo.Create(ctx, newAsset); err != nil {
				log.Printf("Create asset failed: %v", err)
			}
		} else {
			// Asset đã tồn tại -> cập nhật thông tin (trừ is_managed)
			existing.MACAddress = la.MACAddress
			existing.Hostname = la.Hostname
			existing.AssetType = la.AssetType
			existing.FirstSeen = la.FirstSeen
			existing.LastSeen = la.LastSeen
			if err := repo.Update(ctx, la.IPAddress, existing); err != nil {
				log.Printf("Update asset failed: %v", err)
			}
		}
	}

	// 3. Duyệt qua từng service
	for _, s := range logServices {
		services, err := repo.GetServices(ctx, s.AssetIP)
		if err != nil {
			log.Printf("GetServices failed: %v", err)
			continue
		}

		found := false
		for _, existing := range services {
			if existing.Name == s.Name && existing.Port == s.Port && existing.Protocol == s.Protocol {
				// Đã có -> cập nhật thông tin (trừ is_managed)
				existing.Description = "" // có thể cập nhật nếu có source mô tả
				if err := repo.UpdateService(ctx, existing.ID, &existing); err != nil {
					log.Printf("Update service failed: %v", err)
				}
				found = true
				break
			}
		}

		if !found {
			// Chưa có -> thêm mới
			newService := &models.Service{
				AssetIP:   s.AssetIP,
				Name:      s.Name,
				Port:      s.Port,
				Protocol:  s.Protocol,
				IsManaged: false,
			}
			if err := repo.AddService(ctx, s.AssetIP, newService); err != nil {
				log.Printf("Add service failed: %v", err)
			}
		}
	}

	return nil
}
