<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import type { TelegramConfig } from '../types'
import type { Locale, Messages } from '../i18n'

const props = defineProps<{
  visible: boolean
  telegramConfig: TelegramConfig | null
  locale: Locale
  t: Messages
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'update-telegram', config: TelegramConfig): void
  (e: 'test-telegram'): void
  (e: 'update-locale', locale: Locale): void
  (e: 'update-proxy', config: { enabled: boolean; url: string }): void
}>()

const localTelegram = reactive<TelegramConfig>({
  botToken: '',
  chatId: '',
  enabled: false,
})

const proxyEnabled = ref(false)
const proxyUrl = ref('')

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

const saveTelegram = () => {
  emit('update-telegram', { ...localTelegram })
}

const saveProxy = () => {
  emit('update-proxy', { enabled: proxyEnabled.value, url: proxyUrl.value })
}

const handleOverlayClick = (e: MouseEvent) => {
  if ((e.target as HTMLElement).classList.contains('modal-overlay')) {
    emit('close')
  }
}
</script>

<template>
  <div v-if="visible" class="modal-overlay" @click="handleOverlayClick">
    <div class="modal-content">
      <div class="modal-header">
        <h2>{{ t.settings.title }}</h2>
        <button class="close-btn" @click="$emit('close')">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>

      <div class="modal-body">
        <section class="settings-section">
          <h3>{{ t.settings.language }}</h3>
          <div class="language-selector">
            <button
              :class="{ active: locale === 'zh-CN' }"
              @click="$emit('update-locale', 'zh-CN')"
            >
              简体中文
            </button>
            <button
              :class="{ active: locale === 'en-US' }"
              @click="$emit('update-locale', 'en-US')"
            >
              English
            </button>
          </div>
        </section>

        <section class="settings-section">
          <h3>{{ t.settings.telegram }}</h3>
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
          <div class="button-row">
            <button class="secondary-btn" @click="saveTelegram">{{ t.sidebar.save }}</button>
            <button class="text-btn" @click="$emit('test-telegram')">{{ t.sidebar.testSend }}</button>
          </div>
        </section>

        <section class="settings-section">
          <h3>{{ t.settings.proxy }}</h3>
          <div class="form-group checkbox-row">
            <label class="switch-label">
              <input type="checkbox" v-model="proxyEnabled" />
              <span>{{ t.settings.proxyEnabled }}</span>
            </label>
          </div>
          <div class="form-group" v-if="proxyEnabled">
            <label>{{ t.settings.proxyUrl }}</label>
            <input type="text" v-model="proxyUrl" :placeholder="t.settings.proxyPlaceholder" />
          </div>
          <button class="secondary-btn" @click="saveProxy" v-if="proxyEnabled">{{ t.sidebar.save }}</button>
        </section>
      </div>
    </div>
  </div>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 16px;
  width: 90%;
  max-width: 480px;
  max-height: 80vh;
  overflow-y: auto;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid var(--border-color);
}

.modal-header h2 {
  font-size: 18px;
  font-weight: 600;
}

.close-btn {
  background: none;
  border: none;
  padding: 4px;
  cursor: pointer;
  color: var(--text-secondary);
  border-radius: 6px;
}

.close-btn:hover {
  background: var(--bg-color);
}

.modal-body {
  padding: 24px;
}

.settings-section {
  margin-bottom: 28px;
}

.settings-section:last-child {
  margin-bottom: 0;
}

.settings-section h3 {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-secondary);
  margin-bottom: 16px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.language-selector {
  display: flex;
  gap: 8px;
}

.language-selector button {
  flex: 1;
  padding: 10px 16px;
  border-radius: 8px;
  background: var(--bg-color);
  border: 2px solid transparent;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.language-selector button.active {
  background: rgba(0, 122, 255, 0.1);
  border-color: var(--primary);
  color: var(--primary);
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
  margin-bottom: 6px;
}

.form-group input[type="text"],
.form-group input[type="password"] {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 14px;
  background: var(--bg-color);
}

.form-group input:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(0, 122, 255, 0.1);
  background: white;
}

.checkbox-row {
  flex-direction: row;
}

.switch-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 14px;
}

.button-row {
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
</style>
