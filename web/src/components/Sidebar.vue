<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import type { TelegramConfig } from '../types'
import type { Messages } from '../i18n'

const props = defineProps<{
  telegramConfig: TelegramConfig | null
  loading: boolean
  t: Messages
}>()

const emit = defineEmits<{
  (e: 'add-monitor', data: { urls: string; interval: number; duration: number; notifyMode: string; autoStart: boolean }): void
  (e: 'update-telegram', config: TelegramConfig): void
  (e: 'test-telegram', config: { botToken: string; chatId: string }): void
}>()

const urls = ref('')
const interval = ref(60)
const duration = ref(2)
const customDuration = ref('')
const notifyMode = ref('loop')
const autoStart = ref(true)

const durationOptions = [
  { label: '2h', value: 2 },
  { label: '8h', value: 8 },
  { label: '12h', value: 12 },
  { label: '24h', value: 24 },
  { label: '∞', value: 0 },
]

const localTelegram = reactive<TelegramConfig>({
  botToken: '',
  chatId: '',
  enabled: false,
})

watch(
  () => props.telegramConfig,
  (newVal) => {
    if (newVal) {
      localTelegram.botToken = newVal.botToken
      localTelegram.chatId = newVal.chatId
      localTelegram.enabled = newVal.enabled
    }
  },
  { immediate: true }
)

const setDuration = (hours: number) => {
  duration.value = hours
  customDuration.value = ''
}

const handleCustomDuration = () => {
  const hours = parseFloat(customDuration.value)
  if (!isNaN(hours) && hours > 0) {
    duration.value = Math.round(hours)
  }
}

const handleAdd = () => {
  if (!urls.value.trim()) return

  emit('add-monitor', {
    urls: urls.value.trim(),
    interval: interval.value,
    duration: duration.value,
    notifyMode: notifyMode.value,
    autoStart: autoStart.value,
  })

  urls.value = ''
}

const saveTelegram = () => {
  emit('update-telegram', { ...localTelegram })
}
</script>

<template>
  <aside class="sidebar">
    <div class="sidebar-header">
      <h2>{{ t.sidebar.newMonitor }}</h2>
    </div>

    <div class="sidebar-content">
      <div class="form-group">
        <label>{{ t.sidebar.urlsLabel }}</label>
        <textarea v-model="urls" :placeholder="t.sidebar.urlsPlaceholder" rows="5"></textarea>
        <p class="hint">{{ t.sidebar.urlsHint }}</p>
      </div>

      <div class="form-group">
        <label>{{ t.sidebar.intervalLabel }}</label>
        <input type="number" v-model.number="interval" min="10" />
      </div>

      <div class="form-group">
        <label>{{ t.sidebar.durationLabel }}</label>
        <div class="duration-presets">
          <button
            v-for="opt in durationOptions"
            :key="opt.value"
            type="button"
            class="preset-btn"
            :class="{ active: duration === opt.value }"
            @click="setDuration(opt.value)"
          >
            {{ opt.label }}
          </button>
        </div>
        <div class="custom-duration">
          <input type="number" v-model="customDuration" placeholder="Custom (hours)" @input="handleCustomDuration" />
        </div>
      </div>

      <div class="form-group checkbox-row">
        <label class="switch-label">
          <input type="checkbox" v-model="autoStart" />
          <span>{{ t.sidebar.autoStart }}</span>
        </label>
      </div>

      <button class="primary-btn" @click="handleAdd" :disabled="!urls.trim()">
        {{ t.sidebar.addButton }}
      </button>

      <div class="divider"></div>

      <div class="telegram-section">
        <h3>{{ t.sidebar.telegram }}</h3>
        <div class="form-group">
          <label>{{ t.sidebar.botToken }}</label>
          <input type="password" v-model="localTelegram.botToken" placeholder="123456:ABC-DEF..." />
        </div>
        <div class="form-group">
          <label>{{ t.sidebar.chatId }}</label>
          <input type="text" v-model="localTelegram.chatId" placeholder="-100123456789" />
        </div>
        <div class="form-group checkbox-row">
          <label class="switch-label">
            <input type="checkbox" v-model="localTelegram.enabled" />
            <span>{{ t.sidebar.enableNotify }}</span>
          </label>
        </div>
        <div class="telegram-actions">
          <button class="secondary-btn" @click="saveTelegram">{{ t.sidebar.save }}</button>
          <button class="text-btn" @click="$emit('test-telegram', { botToken: localTelegram.botToken, chatId: localTelegram.chatId })">{{ t.sidebar.testSend }}</button>
        </div>
      </div>
    </div>
  </aside>
</template>

<style scoped>
.sidebar {
  width: 320px;
  background: white;
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  height: 100%;
  flex-shrink: 0;
  overflow-y: auto;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid var(--gray-light);
}

.sidebar-header h2 {
  font-size: 18px;
  font-weight: 600;
}

.sidebar-content {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
}

input[type='text'],
input[type='password'],
input[type='number'],
textarea {
  padding: 10px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: var(--bg-color);
  color: var(--text-primary);
  width: 100%;
  transition: border-color 0.2s, box-shadow 0.2s;
}

input:focus,
textarea:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(0, 122, 255, 0.1);
  background: white;
}

textarea {
  resize: vertical;
  min-height: 80px;
}

.hint {
  font-size: 12px;
  color: var(--text-secondary);
}

.duration-presets {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
}

.preset-btn {
  flex: 1;
  padding: 6px;
  font-size: 12px;
  border-radius: 6px;
  background: var(--bg-color);
  color: var(--text-primary);
  border: 1px solid transparent;
  cursor: pointer;
}

.preset-btn.active {
  background: var(--primary);
  color: white;
}

.checkbox-row {
  flex-direction: row;
  align-items: center;
}

.switch-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 14px;
  user-select: none;
}

.primary-btn {
  background: var(--primary);
  color: white;
  padding: 12px;
  border-radius: 10px;
  font-weight: 600;
  margin-top: 8px;
  border: none;
  cursor: pointer;
}

.primary-btn:hover {
  background: var(--primary-hover);
}

.primary-btn:disabled {
  background: var(--gray);
  opacity: 0.5;
  cursor: not-allowed;
}

.divider {
  height: 1px;
  background: var(--border-color);
  margin: 0 -20px;
}

.telegram-section h3 {
  font-size: 15px;
  font-weight: 600;
  margin-bottom: 16px;
}

.telegram-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.secondary-btn {
  background: var(--bg-color);
  border: 1px solid var(--border-color);
  padding: 8px 16px;
  border-radius: 8px;
  font-weight: 500;
  cursor: pointer;
}

.secondary-btn:hover {
  background: var(--gray-light);
}

.text-btn {
  background: none;
  border: none;
  color: var(--primary);
  font-size: 13px;
  cursor: pointer;
}

.text-btn:hover {
  text-decoration: underline;
}

@media (max-width: 768px) {
  .sidebar {
    width: 100%;
    border-right: none;
    border-bottom: 1px solid var(--border-color);
    height: auto;
    max-height: 50vh;
  }
}
</style>
