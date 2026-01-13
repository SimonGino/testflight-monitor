import axios from 'axios'
import type { Monitor, CreateMonitorParams, TelegramConfig, StatusResponse } from '../types'

const api = axios.create({
  baseURL: '/api'
})

export const getMonitors = async (): Promise<Monitor[]> => {
  const response = await api.get('/monitors')
  return response.data.data || []
}

export const createMonitors = async (params: CreateMonitorParams): Promise<Monitor[]> => {
  const response = await api.post('/monitors', params)
  return response.data
}

export const toggleMonitor = async (id: number): Promise<Monitor> => {
  const response = await api.post(`/monitors/${id}/toggle`)
  return response.data
}

export const deleteMonitor = async (id: number): Promise<void> => {
  await api.delete(`/monitors/${id}`)
}

export const updateMonitor = async (id: number, data: { interval?: number; duration?: number }): Promise<Monitor> => {
  const response = await api.put(`/monitors/${id}`, data)
  return response.data.data
}

export const getTelegramConfig = async (): Promise<TelegramConfig> => {
  const response = await api.get('/telegram')
  return response.data
}

export const updateTelegramConfig = async (config: TelegramConfig): Promise<TelegramConfig> => {
  const response = await api.put('/telegram', config)
  return response.data
}

export const testTelegram = async (config: { botToken: string; chatId: string }): Promise<void> => {
  await api.post('/telegram/test', config)
}

export const getStatus = async (): Promise<StatusResponse> => {
  const response = await api.get('/status')
  return response.data
}

export const getProxyConfig = async (): Promise<{ enabled: boolean; url: string }> => {
  const response = await api.get('/proxy')
  return response.data
}

export const updateProxyConfig = async (config: { enabled: boolean; url: string }): Promise<void> => {
  await api.put('/proxy', config)
}
