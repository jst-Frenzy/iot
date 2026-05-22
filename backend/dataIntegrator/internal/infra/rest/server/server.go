package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jst-Frenzy/iot/backend/dataIntegrator/internal/usecase"
)

type processor interface {
	GetDevices() ([]usecase.Device, error)
	GetTelemetryByPeriod(deviceName string, from time.Time, to time.Time) ([]usecase.Telemetry, error)
	ChangeFanMode() error
	ChangePumpMode() error
}

type Handler struct {
	processor processor
}

func New(processor processor) *Handler {
	return &Handler{
		processor: processor,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	api := router.Group("/api")
	{
		devices := api.Group("/devices")
		{
			devices.GET("/", h.getDevices)
			devices.GET("", h.getDevices)
			devices.GET("/changeFanMode", h.ChangeFanMode)
			devices.GET("/changePumpMode", h.ChangePumpMode)
		}

		api.GET("/telemetry", h.getTelemetryByPeriod)

	}

	return router
}

func (h *Handler) getDevices(c *gin.Context) {
	devices, err := h.processor.GetDevices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"devices": devices,
	})
}

func (h *Handler) getTelemetryByPeriod(c *gin.Context) {
	deviceName := c.Query("device_name")
	fromStr := c.Query("from")
	toStr := c.Query("to")

	if deviceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "device_name is required",
		})
		return
	}

	if fromStr == "" || toStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "from and to are required",
		})
		return
	}

	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid from format",
		})
		return
	}

	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid to format",
		})
		return
	}

	telemetry, err := h.processor.GetTelemetryByPeriod(
		deviceName,
		from,
		to,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"telemetry": telemetry,
	})
}

func (h *Handler) ChangeFanMode(c *gin.Context) {
	err := h.processor.ChangeFanMode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) ChangePumpMode(c *gin.Context) {
	err := h.processor.ChangePumpMode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
