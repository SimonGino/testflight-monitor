<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch } from 'vue'
import type { Messages } from '../i18n'

const props = defineProps<{
  activeCount: number
  nextCheckAt: string | null
  t: Messages
}>()

defineEmits<{
  (e: 'open-settings'): void
}>()

const timeRemaining = ref('')
let timer: number | null = null

const updateTimer = () => {
  if (!props.nextCheckAt) {
    timeRemaining.value = props.t.header.waiting
    return
  }

  const target = new Date(props.nextCheckAt).getTime()
  const now = Date.now()
  const diff = Math.max(0, Math.floor((target - now) / 1000))

  if (diff > 0) {
    timeRemaining.value = `${props.t.header.nextCheck} ${diff}s`
  } else {
    timeRemaining.value = props.t.header.checking
  }
}

watch(() => props.t, updateTimer)
watch(() => props.nextCheckAt, updateTimer)

onMounted(() => {
  updateTimer()
  timer = window.setInterval(updateTimer, 1000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<template>
  <header class="app-header">
    <div class="header-left">
      <div class="logo">
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="logo-icon">
          <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
          <polyline points="22 4 12 14.01 9 11.01"></polyline>
        </svg>
        <h1>{{ t.app.title }}</h1>
      </div>
      <div class="status-pill">
        <span class="dot pulse"></span>
        {{ activeCount }} {{ t.app.active }}
      </div>
    </div>

    <div class="header-right">
      <div class="countdown">{{ timeRemaining }}</div>
      <button class="icon-btn settings-btn" @click="$emit('open-settings')" :title="t.app.settings">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="3"></circle>
          <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
        </svg>
      </button>
    </div>
  </header>
</template>

<style scoped>
.app-header {
  height: 64px;
  background: white;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  flex-shrink: 0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 24px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  color: var(--primary);
}

.logo h1 {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

.status-pill {
  background: rgba(52, 199, 89, 0.1);
  color: var(--success);
  padding: 4px 12px;
  border-radius: 100px;
  font-size: 13px;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 8px;
}

.dot {
  width: 8px;
  height: 8px;
  background: currentColor;
  border-radius: 50%;
  display: block;
}

.pulse {
  box-shadow: 0 0 0 0 rgba(52, 199, 89, 0.4);
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% {
    box-shadow: 0 0 0 0 rgba(52, 199, 89, 0.4);
  }
  70% {
    box-shadow: 0 0 0 6px rgba(52, 199, 89, 0);
  }
  100% {
    box-shadow: 0 0 0 0 rgba(52, 199, 89, 0);
  }
}

.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.countdown {
  font-family: monospace;
  font-size: 14px;
  color: var(--text-secondary);
  background: var(--bg-color);
  padding: 4px 8px;
  border-radius: 4px;
}

.settings-btn {
  background: none;
  border: none;
  padding: 8px;
  border-radius: 8px;
  cursor: pointer;
  color: var(--text-secondary);
  transition: all 0.2s;
}

.settings-btn:hover {
  background: var(--bg-color);
  color: var(--text-primary);
}
</style>
