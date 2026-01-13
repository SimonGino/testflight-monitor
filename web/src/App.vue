<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, computed } from 'vue'
import Sidebar from './components/Sidebar.vue'
import MonitorList from './components/MonitorList.vue'
import Header from './components/Header.vue'
import SettingsModal from './components/SettingsModal.vue'
import * as api from './api'
import type { Monitor, TelegramConfig } from './types'
import { getMessages, getStoredLocale, setStoredLocale, type Locale } from './i18n'

const monitors = ref<Monitor[]>([])
const telegramConfig = ref<TelegramConfig | null>(null)
const activeCount = ref(0)
const nextCheckAt = ref<string | null>(null)
const loading = ref(false)
const showSettings = ref(false)
const locale = ref<Locale>(getStoredLocale())

const t = computed(() => getMessages(locale.value))

let pollTimer: number | null = null

const fetchData = async () => {
  try {
    const [monitorsData, statusData] = await Promise.all([
      api.getMonitors(),
      api.getStatus()
    ])
    monitors.value = monitorsData
    activeCount.value = statusData.activeJobs
    nextCheckAt.value = statusData.nextCheckAt
  } catch (err) {
    console.error('Failed to fetch data', err)
  }
}

const fetchTelegram = async () => {
  try {
    telegramConfig.value = await api.getTelegramConfig()
  } catch (err) {
    console.error('Failed to fetch telegram config', err)
  }
}

const handleAddMonitor = async (data: any) => {
  try {
    loading.value = true
    await api.createMonitors(data)
    await fetchData()
  } catch (err) {
    alert(locale.value === 'zh-CN' ? '创建监控失败' : 'Failed to create monitor')
  } finally {
    loading.value = false
  }
}

const handleToggle = async (id: number) => {
  try {
    await api.toggleMonitor(id)
    await fetchData()
  } catch (err) {
    console.error(err)
  }
}

const handleDelete = async (id: number) => {
  if (!confirm(t.value.monitor.confirmDelete)) return
  try {
    await api.deleteMonitor(id)
    await fetchData()
  } catch (err) {
    console.error(err)
  }
}

const handleUpdateMonitor = async (id: number, data: { interval?: number; duration?: number }) => {
  try {
    await api.updateMonitor(id, data)
    await fetchData()
  } catch (err) {
    console.error(err)
  }
}

const handleUpdateTelegram = async (config: TelegramConfig) => {
  try {
    const updated = await api.updateTelegramConfig(config)
    telegramConfig.value = updated
    alert(t.value.settings.saved)
  } catch (err) {
    alert(locale.value === 'zh-CN' ? '保存失败' : 'Failed to save')
  }
}

const handleTestTelegram = async (config: { botToken: string; chatId: string }) => {
  try {
    await api.testTelegram(config)
    alert(t.value.settings.testSuccess)
  } catch (err) {
    alert(t.value.settings.testFailed)
  }
}

const handleUpdateProxy = async (config: { enabled: boolean; url: string }) => {
  try {
    await api.updateProxyConfig(config)
    alert(t.value.settings.saved)
  } catch (err) {
    alert(locale.value === 'zh-CN' ? '保存失败' : 'Failed to save')
  }
}

const handleUpdateLocale = (newLocale: Locale) => {
  locale.value = newLocale
  setStoredLocale(newLocale)
}

const startPolling = () => {
  if (pollTimer) return
  pollTimer = window.setInterval(fetchData, 5000)
}

const stopPolling = () => {
  if (pollTimer) {
    clearInterval(pollTimer)
    pollTimer = null
  }
}

watch(activeCount, (count) => {
  if (count > 0) {
    startPolling()
  } else {
    stopPolling()
  }
})

onMounted(() => {
  fetchData()
  fetchTelegram()
})

onUnmounted(() => {
  stopPolling()
})
</script>

<template>
  <div class="app-container">
    <Sidebar
      :telegram-config="telegramConfig"
      :loading="loading"
      :t="t"
      @add-monitor="handleAddMonitor"
      @update-telegram="handleUpdateTelegram"
      @test-telegram="handleTestTelegram"
    />

    <main class="main-content">
      <Header
        :active-count="activeCount"
        :next-check-at="nextCheckAt"
        :t="t"
        @open-settings="showSettings = true"
      />
      <div class="scroll-area">
        <MonitorList
          :monitors="monitors"
          :t="t"
          @toggle="handleToggle"
          @delete="handleDelete"
          @update="handleUpdateMonitor"
        />
      </div>
      <footer class="app-footer">
        <a href="https://github.com/SimonGino/testflight-monitor" target="_blank" rel="noopener noreferrer">
          <svg height="20" width="20" viewBox="0 0 16 16" fill="currentColor">
            <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/>
          </svg>
          <span>GitHub</span>
        </a>
      </footer>
    </main>

    <SettingsModal
      :visible="showSettings"
      :telegram-config="telegramConfig"
      :locale="locale"
      :t="t"
      @close="showSettings = false"
      @update-telegram="handleUpdateTelegram"
      @test-telegram="handleTestTelegram"
      @update-locale="handleUpdateLocale"
      @update-proxy="handleUpdateProxy"
    />
  </div>
</template>

<style scoped>
.app-container {
  display: flex;
  height: 100vh;
  width: 100vw;
  overflow: hidden;
  background-color: var(--bg-color);
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.scroll-area {
  flex: 1;
  overflow-y: auto;
  padding-bottom: 60px;
}

.app-footer {
  padding: 12px 20px;
  text-align: center;
  border-top: 1px solid var(--border-color);
  background: white;
}

.app-footer a {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 13px;
  transition: color 0.2s;
}

.app-footer a:hover {
  color: var(--text-primary);
}

@media (max-width: 768px) {
  .app-container {
    flex-direction: column;
  }

  .sidebar {
    height: auto;
    max-height: 40vh;
  }

  .main-content {
    flex: 1;
  }
}
</style>
