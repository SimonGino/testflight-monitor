package main

import (
	"log"

	"tf-monitor/internal/api"
	"tf-monitor/internal/config"
	"tf-monitor/internal/model"
	"tf-monitor/internal/repository"
	"tf-monitor/internal/service/scheduler"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	if err := repository.InitDB(cfg.Database.Path); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	proxyURL := ""
	if cfg.Proxy.Enabled {
		proxyURL = cfg.Proxy.URL
	}

	sched := scheduler.GetScheduler()
	sched.Init(proxyURL)

	var telegramCfg model.TelegramConfig
	if repository.GetDB().First(&telegramCfg).Error == nil && telegramCfg.Enabled {
		sched.UpdateNotifier(telegramCfg.BotToken, telegramCfg.ChatID)
	}

	sched.Start()
	defer sched.Stop()

	r := gin.Default()

	r.Static("/assets", "./web/dist/assets")
	r.StaticFile("/", "./web/dist/index.html")
	r.StaticFile("/favicon.ico", "./web/dist/favicon.ico")
	r.NoRoute(func(c *gin.Context) {
		c.File("./web/dist/index.html")
	})

	handler := api.NewHandler(proxyURL)
	handler.RegisterRoutes(r)

	log.Printf("Server starting on :%s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
