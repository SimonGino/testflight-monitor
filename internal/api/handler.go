package api

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"tf-monitor/internal/model"
	"tf-monitor/internal/repository"
	"tf-monitor/internal/service/monitor"
	"tf-monitor/internal/service/notify"
	"tf-monitor/internal/service/scheduler"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	proxyURL string
}

func NewHandler(proxyURL string) *Handler {
	return &Handler{proxyURL: proxyURL}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/monitors", h.ListMonitors)
		api.POST("/monitors", h.CreateMonitor)
		api.GET("/monitors/:id", h.GetMonitor)
		api.PUT("/monitors/:id", h.UpdateMonitor)
		api.DELETE("/monitors/:id", h.DeleteMonitor)
		api.POST("/monitors/:id/toggle", h.ToggleMonitor)

		api.GET("/telegram", h.GetTelegramConfig)
		api.PUT("/telegram", h.UpdateTelegramConfig)
		api.POST("/telegram/test", h.TestTelegram)

		api.GET("/proxy", h.GetProxyConfig)
		api.PUT("/proxy", h.UpdateProxyConfig)

		api.GET("/status", h.GetStatus)
	}
}

type CreateMonitorRequest struct {
	URLs       string `json:"urls"`
	Interval   int    `json:"interval"`
	Duration   int    `json:"duration"`
	NotifyMode string `json:"notifyMode"`
	AutoStart  bool   `json:"autoStart"`
}

type MonitorResponse struct {
	ID            uint       `json:"id"`
	AppID         string     `json:"appId"`
	AppName       string     `json:"appName"`
	IconURL       string     `json:"iconUrl"`
	TestFlightURL string     `json:"testFlightUrl"`
	Status        string     `json:"status"`
	Interval      int        `json:"interval"`
	Duration      int        `json:"duration"`
	NotifyMode    string     `json:"notifyMode"`
	Enabled       bool       `json:"enabled"`
	LastCheck     *time.Time `json:"lastCheck"`
	LastError     string     `json:"lastError"`
	ExpireAt      *time.Time `json:"expireAt"`
	CreatedAt     time.Time  `json:"createdAt"`
}

func toMonitorResponse(m *model.Monitor) MonitorResponse {
	return MonitorResponse{
		ID:            m.ID,
		AppID:         m.AppID,
		AppName:       m.AppName,
		IconURL:       m.IconURL,
		TestFlightURL: m.TestFlightURL,
		Status:        string(m.Status),
		Interval:      m.Interval,
		Duration:      m.Duration,
		NotifyMode:    string(m.NotifyMode),
		Enabled:       m.Enabled,
		LastCheck:     m.LastCheck,
		LastError:     m.LastError,
		ExpireAt:      m.ExpireAt,
		CreatedAt:     m.CreatedAt,
	}
}

func (h *Handler) ListMonitors(c *gin.Context) {
	var monitors []model.Monitor
	repository.GetDB().Order("created_at desc").Find(&monitors)

	result := make([]MonitorResponse, len(monitors))
	for i, m := range monitors {
		result[i] = toMonitorResponse(&m)
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *Handler) CreateMonitor(c *gin.Context) {
	var req CreateMonitorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	urls := strings.Split(strings.TrimSpace(req.URLs), "\n")
	created := []MonitorResponse{}
	errors := []string{}

	checker := monitor.NewChecker(h.proxyURL)

	for _, url := range urls {
		url = strings.TrimSpace(url)
		if url == "" {
			continue
		}

		appID, err := monitor.ParseURL(url)
		if err != nil {
			errors = append(errors, url+": "+err.Error())
			continue
		}

		var existing model.Monitor
		if repository.GetDB().Where("app_id = ?", appID).First(&existing).Error == nil {
			errors = append(errors, url+": already exists")
			continue
		}

		interval := req.Interval
		if interval < 10 {
			interval = 30
		}

		duration := req.Duration

		notifyMode := model.NotifyMode(req.NotifyMode)
		if notifyMode == "" {
			notifyMode = model.NotifyOnce
		}

		var expireAt *time.Time
		if duration > 0 {
			t := time.Now().Add(time.Duration(duration) * time.Hour)
			expireAt = &t
		}

		m := model.Monitor{
			AppID:         appID,
			TestFlightURL: url,
			Interval:      interval,
			Duration:      duration,
			NotifyMode:    notifyMode,
			Enabled:       req.AutoStart,
			ExpireAt:      expireAt,
		}

		info, err := checker.Check(appID)
		if err == nil {
			m.AppName = info.AppName
			m.IconURL = info.IconURL
			if info.Available {
				m.Status = model.StatusAvailable
			} else {
				m.Status = model.StatusFull
			}
		}

		if err := repository.GetDB().Create(&m).Error; err != nil {
			errors = append(errors, url+": "+err.Error())
			continue
		}

		created = append(created, toMonitorResponse(&m))

		if req.AutoStart {
			scheduler.GetScheduler().StartJob(m.ID)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"created": created,
		"errors":  errors,
	})
}

func (h *Handler) GetMonitor(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var m model.Monitor
	if err := repository.GetDB().First(&m, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "monitor not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": toMonitorResponse(&m)})
}

func (h *Handler) UpdateMonitor(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var m model.Monitor
	if err := repository.GetDB().First(&m, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "monitor not found"})
		return
	}

	var req struct {
		Interval   *int    `json:"interval"`
		Duration   *int    `json:"duration"`
		NotifyMode *string `json:"notifyMode"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := make(map[string]interface{})
	if req.Interval != nil && *req.Interval >= 10 {
		updates["interval"] = *req.Interval
	}
	if req.Duration != nil {
		updates["duration"] = *req.Duration
		if *req.Duration > 0 {
			expireAt := time.Now().Add(time.Duration(*req.Duration) * time.Hour)
			updates["expire_at"] = expireAt
		} else {
			updates["expire_at"] = nil
		}
	}
	if req.NotifyMode != nil {
		updates["notify_mode"] = *req.NotifyMode
	}

	repository.GetDB().Model(&m).Updates(updates)
	repository.GetDB().First(&m, id)

	c.JSON(http.StatusOK, gin.H{"data": toMonitorResponse(&m)})
}

func (h *Handler) DeleteMonitor(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	scheduler.GetScheduler().StopJob(uint(id))
	repository.GetDB().Delete(&model.Monitor{}, id)

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *Handler) ToggleMonitor(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var m model.Monitor
	if err := repository.GetDB().First(&m, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "monitor not found"})
		return
	}

	m.Enabled = !m.Enabled
	if m.Enabled {
		expireAt := time.Now().Add(time.Duration(m.Duration) * time.Hour)
		m.ExpireAt = &expireAt
		m.Notified = false
	}
	repository.GetDB().Save(&m)

	if m.Enabled {
		scheduler.GetScheduler().StartJob(m.ID)
	} else {
		scheduler.GetScheduler().StopJob(m.ID)
	}

	c.JSON(http.StatusOK, gin.H{"data": toMonitorResponse(&m)})
}

func (h *Handler) GetTelegramConfig(c *gin.Context) {
	var cfg model.TelegramConfig
	repository.GetDB().FirstOrCreate(&cfg, model.TelegramConfig{})
	c.JSON(http.StatusOK, gin.H{
		"botToken": cfg.BotToken,
		"chatId":   cfg.ChatID,
		"enabled":  cfg.Enabled,
	})
}

func (h *Handler) UpdateTelegramConfig(c *gin.Context) {
	var req struct {
		BotToken string `json:"botToken"`
		ChatID   string `json:"chatId"`
		Enabled  bool   `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cfg model.TelegramConfig
	repository.GetDB().FirstOrCreate(&cfg, model.TelegramConfig{})
	cfg.BotToken = req.BotToken
	cfg.ChatID = req.ChatID
	cfg.Enabled = req.Enabled
	repository.GetDB().Save(&cfg)

	scheduler.GetScheduler().UpdateNotifier(cfg.BotToken, cfg.ChatID)

	c.JSON(http.StatusOK, gin.H{"message": "saved"})
}

func (h *Handler) TestTelegram(c *gin.Context) {
	var req struct {
		BotToken string `json:"botToken"`
		ChatID   string `json:"chatId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.BotToken == "" || req.ChatID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "botToken and chatId required"})
		return
	}

	notifier := notify.NewTelegramNotifier(req.BotToken, req.ChatID, h.proxyURL)
	if err := notifier.Send("TestFlight Monitor", "🎉 测试消息发送成功！"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "sent"})
}

func (h *Handler) GetProxyConfig(c *gin.Context) {
	var enabled, url model.SystemConfig
	repository.GetDB().Where("key = ?", "proxy_enabled").First(&enabled)
	repository.GetDB().Where("key = ?", "proxy_url").First(&url)

	c.JSON(http.StatusOK, gin.H{
		"enabled": enabled.Value == "true",
		"url":     url.Value,
	})
}

func (h *Handler) UpdateProxyConfig(c *gin.Context) {
	var req struct {
		Enabled bool   `json:"enabled"`
		URL     string `json:"url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repository.GetDB().Where("key = ?", "proxy_enabled").Assign(model.SystemConfig{
		Key:   "proxy_enabled",
		Value: strconv.FormatBool(req.Enabled),
	}).FirstOrCreate(&model.SystemConfig{})

	repository.GetDB().Where("key = ?", "proxy_url").Assign(model.SystemConfig{
		Key:   "proxy_url",
		Value: req.URL,
	}).FirstOrCreate(&model.SystemConfig{})

	proxyURL := ""
	if req.Enabled {
		proxyURL = req.URL
	}
	scheduler.GetScheduler().Init(proxyURL)

	c.JSON(http.StatusOK, gin.H{"message": "saved"})
}

func (h *Handler) GetStatus(c *gin.Context) {
	sched := scheduler.GetScheduler()
	c.JSON(http.StatusOK, gin.H{
		"activeJobs":  sched.GetActiveJobCount(),
		"nextCheckAt": sched.GetNextCheckTime(),
	})
}
