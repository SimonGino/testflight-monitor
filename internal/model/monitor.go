package model

import (
	"time"

	"gorm.io/gorm"
)

// MonitorStatus represents the status of a TestFlight app
type MonitorStatus string

const (
	StatusAvailable MonitorStatus = "available" // has open slots
	StatusFull      MonitorStatus = "full"      // no slots available
	StatusChecking  MonitorStatus = "checking"  // currently checking
	StatusError     MonitorStatus = "error"     // check failed
	StatusExpired   MonitorStatus = "expired"   // monitor duration expired
)

// NotifyMode represents how notifications should be sent
type NotifyMode string

const (
	NotifyLoop          NotifyMode = "loop"           // notify every time
	NotifyOnce          NotifyMode = "once"           // notify once then stop
	NotifyOnlyAvailable NotifyMode = "only_available" // notify only when status changes to available
)

// Monitor represents a TestFlight app being monitored
type Monitor struct {
	gorm.Model
	AppID         string        `json:"appId" gorm:"index"`          // TestFlight app ID
	AppName       string        `json:"appName"`                     // App name (fetched from TestFlight)
	IconURL       string        `json:"iconUrl"`                     // App icon URL
	TestFlightURL string        `json:"testFlightUrl" gorm:"unique"` // Original TestFlight URL
	Status        MonitorStatus `json:"status" gorm:"default:checking"`
	Interval      int           `json:"interval" gorm:"default:30"` // Check interval in seconds (min: 10)
	Duration      int           `json:"duration" gorm:"default:24"` // Monitor duration in hours
	NotifyMode    NotifyMode    `json:"notifyMode" gorm:"default:once"`
	Enabled       bool          `json:"enabled" gorm:"default:true"`   // Whether monitoring is active
	Notified      bool          `json:"notified" gorm:"default:false"` // Whether notification has been sent (for once mode)
	LastCheck     *time.Time    `json:"lastCheck"`                     // Last check timestamp
	LastError     string        `json:"lastError"`                     // Last error message
	ExpireAt      *time.Time    `json:"expireAt"`                      // When monitoring expires
}

// TelegramConfig stores Telegram notification settings
type TelegramConfig struct {
	gorm.Model
	BotToken string `json:"botToken"`
	ChatID   string `json:"chatId"`
	Enabled  bool   `json:"enabled" gorm:"default:true"`
}

// ProxyConfig stores proxy settings
type SystemConfig struct {
	gorm.Model
	Key   string `json:"key" gorm:"uniqueIndex"`
	Value string `json:"value"`
}
