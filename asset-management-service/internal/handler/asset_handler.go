package handler

import (
	"asset-management-service/internal/models"
	"asset-management-service/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type AssetHandler struct {
	Service service.AssetService
}

func NewAssetHandler(s service.AssetService) *AssetHandler {
	return &AssetHandler{Service: s}
}

func (h *AssetHandler) GetAssets(c *gin.Context) {
	assets, err := h.Service.GetAllAssets(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assets)
}

func (h *AssetHandler) GetAssetByIP(c *gin.Context) {
	ip := c.Param("ip")
	asset, err := h.Service.GetAssetByIP(c.Request.Context(), ip)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
		return
	}
	c.JSON(http.StatusOK, asset)
}

func (h *AssetHandler) CreateAsset(c *gin.Context) {
	var a models.Asset
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Service.CreateAsset(c.Request.Context(), &a); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, a)
}

func (h *AssetHandler) UpdateAsset(c *gin.Context) {
	ip := c.Param("ip")
	var a models.Asset
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Service.UpdateAsset(c.Request.Context(), ip, &a); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, a)
}

func (h *AssetHandler) DeleteAsset(c *gin.Context) {
	ip := c.Param("ip")
	if err := h.Service.DeleteAsset(c.Request.Context(), ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *AssetHandler) ManageAsset(c *gin.Context) {
	ip := c.Param("ip")
	if err := h.Service.ManageAsset(c.Request.Context(), ip, true); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *AssetHandler) GetServices(c *gin.Context) {
	ip := c.Param("ip")
	services, err := h.Service.GetServices(c.Request.Context(), ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, services)
}

func (h *AssetHandler) AddService(c *gin.Context) {
	ip := c.Param("ip")
	var s models.Service
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Service.AddService(c.Request.Context(), ip, &s); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, s)
}

func (h *AssetHandler) UpdateService(c *gin.Context) {
	serviceID, err := strconv.Atoi(c.Param("service_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}
	var s models.Service
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Service.UpdateService(c.Request.Context(), serviceID, &s); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *AssetHandler) DeleteService(c *gin.Context) {
	serviceID, err := strconv.Atoi(c.Param("service_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}
	if err := h.Service.DeleteService(c.Request.Context(), serviceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *AssetHandler) ManageService(c *gin.Context) {
	serviceID, err := strconv.Atoi(c.Param("service_id"))
	ip := c.Param("ip")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}
	if err := h.Service.ManageService(c.Request.Context(), ip, serviceID, true); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
