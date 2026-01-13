export interface Monitor {
  id: number
  appId: string
  appName: string
  iconUrl: string
  testFlightUrl: string
  status: 'available' | 'full' | 'checking' | 'error' | 'expired'
  interval: number
  duration: number
  notifyMode: 'loop' | 'once' | 'only_available'
  enabled: boolean
  lastCheck: string | null
  lastError: string
  expireAt: string | null
  createdAt: string
}

export interface TelegramConfig {
  botToken: string
  chatId: string
  enabled: boolean
}

export interface CreateMonitorParams {
  urls: string
  interval: number
  duration: number
  notifyMode: 'loop' | 'once' | 'only_available'
  autoStart: boolean
}

export interface StatusResponse {
  activeJobs: number
  nextCheckAt: string | null
}
