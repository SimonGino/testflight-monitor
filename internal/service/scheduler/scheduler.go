package scheduler

import (
	"fmt"
	"log"
	"sync"
	"time"

	"tf-monitor/internal/model"
	"tf-monitor/internal/repository"
	"tf-monitor/internal/service/monitor"
	"tf-monitor/internal/service/notify"
)

type Scheduler struct {
	checker     *monitor.Checker
	notifier    notify.Notifier
	proxyURL    string
	mu          sync.RWMutex
	jobs        map[uint]*Job
	stopChan    chan struct{}
	nextCheckAt time.Time
}

type Job struct {
	MonitorID uint
	StopChan  chan struct{}
	Running   bool
}

var instance *Scheduler
var once sync.Once

func GetScheduler() *Scheduler {
	once.Do(func() {
		instance = &Scheduler{
			jobs:     make(map[uint]*Job),
			stopChan: make(chan struct{}),
		}
	})
	return instance
}

func (s *Scheduler) Init(proxyURL string) {
	s.proxyURL = proxyURL
	s.checker = monitor.NewChecker(proxyURL)
}

func (s *Scheduler) UpdateNotifier(botToken, chatID string) {
	s.notifier = notify.NewTelegramNotifier(botToken, chatID, s.proxyURL)
}

func (s *Scheduler) Start() {
	log.Println("Scheduler started")

	var monitors []model.Monitor
	repository.GetDB().Where("enabled = ?", true).Find(&monitors)

	for _, m := range monitors {
		s.StartJob(m.ID)
	}
}

func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, job := range s.jobs {
		if job.Running {
			close(job.StopChan)
			job.Running = false
		}
		delete(s.jobs, id)
	}
	log.Println("Scheduler stopped")
}

func (s *Scheduler) StartJob(monitorID uint) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if job, exists := s.jobs[monitorID]; exists && job.Running {
		close(job.StopChan)
	}

	job := &Job{
		MonitorID: monitorID,
		StopChan:  make(chan struct{}),
		Running:   true,
	}
	s.jobs[monitorID] = job

	go s.runJob(job)
}

func (s *Scheduler) StopJob(monitorID uint) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if job, exists := s.jobs[monitorID]; exists && job.Running {
		close(job.StopChan)
		job.Running = false
		delete(s.jobs, monitorID)
	}
}

func (s *Scheduler) runJob(job *Job) {
	for {
		var m model.Monitor
		if err := repository.GetDB().First(&m, job.MonitorID).Error; err != nil {
			log.Printf("Monitor %d not found, stopping job", job.MonitorID)
			return
		}

		if m.ExpireAt != nil && time.Now().After(*m.ExpireAt) {
			repository.GetDB().Model(&m).Updates(map[string]interface{}{
				"enabled": false,
				"status":  model.StatusExpired,
			})
			log.Printf("Monitor %d expired", job.MonitorID)
			return
		}

		s.performCheck(&m)

		interval := time.Duration(m.Interval) * time.Second
		if interval < 10*time.Second {
			interval = 10 * time.Second
		}

		s.mu.Lock()
		s.nextCheckAt = time.Now().Add(interval)
		s.mu.Unlock()

		select {
		case <-job.StopChan:
			log.Printf("Job for monitor %d stopped", job.MonitorID)
			return
		case <-time.After(interval):
		}
	}
}

func (s *Scheduler) performCheck(m *model.Monitor) {
	now := time.Now()

	repository.GetDB().Model(m).Updates(map[string]interface{}{
		"status":     model.StatusChecking,
		"last_check": now,
	})

	info, err := s.checker.Check(m.AppID)
	if err != nil {
		repository.GetDB().Model(m).Updates(map[string]interface{}{
			"status":     model.StatusError,
			"last_error": err.Error(),
		})
		log.Printf("Check failed for %s: %v", m.AppID, err)
		return
	}

	if m.AppName == "" && info.AppName != "" {
		repository.GetDB().Model(m).Updates(map[string]interface{}{
			"app_name": info.AppName,
			"icon_url": info.IconURL,
		})
	}

	status := model.StatusFull
	if info.Available {
		status = model.StatusAvailable
	}

	prevStatus := m.Status

	repository.GetDB().Model(m).Updates(map[string]interface{}{
		"status":     status,
		"last_error": "",
	})

	if info.Available && s.notifier != nil {
		shouldNotify := false

		switch m.NotifyMode {
		case model.NotifyLoop:
			shouldNotify = true
		case model.NotifyOnce:
			if !m.Notified {
				shouldNotify = true
			}
		case model.NotifyOnlyAvailable:
			if prevStatus != model.StatusAvailable {
				shouldNotify = true
			}
		}

		if shouldNotify {
			title := "🎉 TestFlight 有位了!"
			message := fmt.Sprintf("**%s**\n\n%s\n\n[点击加入](%s)",
				info.AppName, info.Message, m.TestFlightURL)

			if err := s.notifier.Send(title, message); err != nil {
				log.Printf("Failed to send notification: %v", err)
			} else {
				repository.GetDB().Model(m).Update("notified", true)
				log.Printf("Notification sent for %s", info.AppName)
			}
		}
	}

	log.Printf("Checked %s: %s (available: %v)", m.AppID, info.AppName, info.Available)
}

func (s *Scheduler) GetNextCheckTime() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.nextCheckAt
}

func (s *Scheduler) GetActiveJobCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	count := 0
	for _, job := range s.jobs {
		if job.Running {
			count++
		}
	}
	return count
}
