<script setup lang="ts">
import { computed, ref } from 'vue'
import { useTimeAgo } from '@vueuse/core'
import type { Monitor } from '../types'
import type { Messages } from '../i18n'

const props = defineProps<{
  monitor: Monitor
  t: Messages
}>()

const emit = defineEmits<{
  (e: 'toggle', id: number): void
  (e: 'delete', id: number): void
  (e: 'update', id: number, data: { interval?: number; duration?: number }): void
}>()

const timeAgo = useTimeAgo(new Date(props.monitor.lastCheck || Date.now()))

const isEditing = ref(false)
const editInterval = ref(props.monitor.interval)
const editDuration = ref(props.monitor.duration)

const durationOptions = [
  { label: '2h', value: 2 },
  { label: '8h', value: 8 },
  { label: '12h', value: 12 },
  { label: '24h', value: 24 },
  { label: '∞', value: 0 },
]

const startEdit = () => {
  editInterval.value = props.monitor.interval
  editDuration.value = props.monitor.duration
  isEditing.value = true
}

const cancelEdit = () => {
  isEditing.value = false
}

const saveEdit = () => {
  emit('update', props.monitor.id, {
    interval: editInterval.value,
    duration: editDuration.value,
  })
  isEditing.value = false
}

const statusClass = computed(() => {
  switch (props.monitor.status) {
    case 'available':
      return 'success'
    case 'full':
      return 'danger'
    case 'checking':
      return 'gray'
    case 'error':
      return 'warning'
    case 'expired':
      return 'gray'
    default:
      return 'gray'
  }
})

const statusLabel = computed(() => {
  switch (props.monitor.status) {
    case 'available':
      return props.t.monitor.available
    case 'full':
      return props.t.monitor.full
    case 'checking':
      return props.t.monitor.checking
    case 'error':
      return props.t.monitor.error
    case 'expired':
      return props.t.monitor.expired
    default:
      return 'Unknown'
  }
})

const durationLabel = computed(() => {
  if (props.monitor.duration === 0) {
    return props.t.monitor.forever
  }
  return `${props.monitor.duration}h`
})
</script>

<template>
  <div class="monitor-card">
    <div class="card-header">
      <div class="app-info">
        <img :src="monitor.iconUrl || 'https://placehold.co/64'" alt="App Icon" class="app-icon" />
        <div class="app-details">
          <h3 class="app-name">{{ monitor.appName || t.monitor.loading }}</h3>
          <a :href="monitor.testFlightUrl" target="_blank" class="app-link" @click.stop>
            {{ monitor.testFlightUrl }}
          </a>
        </div>
      </div>
      <div :class="['badge', statusClass]">{{ statusLabel }}</div>
    </div>

    <div class="card-meta">
      <div class="meta-item">
        <span class="label">{{ t.monitor.interval }}:</span>
        <span class="value">{{ monitor.interval }}s</span>
      </div>
      <div class="meta-item">
        <span class="label">{{ t.monitor.duration }}:</span>
        <span class="value">{{ durationLabel }}</span>
      </div>
      <div class="meta-item" v-if="monitor.lastCheck">
        <span class="value">{{ timeAgo }}</span>
      </div>
    </div>

    <div v-if="isEditing" class="edit-form">
      <div class="edit-row">
        <label>{{ t.monitor.interval }} ({{ t.monitor.seconds }})</label>
        <input type="number" v-model.number="editInterval" min="10" />
      </div>
      <div class="edit-row">
        <label>{{ t.monitor.duration }}</label>
        <div class="duration-btns">
          <button
            v-for="opt in durationOptions"
            :key="opt.value"
            type="button"
            class="duration-btn"
            :class="{ active: editDuration === opt.value }"
            @click="editDuration = opt.value"
          >
            {{ opt.label }}
          </button>
        </div>
      </div>
      <div class="edit-actions">
        <button class="save-btn" @click="saveEdit">{{ t.sidebar.save }}</button>
        <button class="cancel-btn" @click="cancelEdit">{{ t.monitor.cancel }}</button>
      </div>
    </div>

    <div class="card-actions" v-else>
      <div class="left-actions">
        <button class="action-btn" :class="monitor.enabled ? 'pause' : 'resume'" @click="emit('toggle', monitor.id)">
          {{ monitor.enabled ? t.monitor.pause : t.monitor.resume }}
        </button>
        <button class="action-btn edit" @click="startEdit">{{ t.monitor.edit }}</button>
      </div>
      <button class="icon-btn delete" @click="emit('delete', monitor.id)" :title="t.monitor.delete">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="16"
          height="16"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <polyline points="3 6 5 6 21 6"></polyline>
          <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
        </svg>
      </button>
    </div>
  </div>
</template>

<style scoped>
.monitor-card {
  background: var(--card-bg);
  border-radius: var(--radius-md);
  padding: 16px;
  box-shadow: var(--shadow-sm);
  display: flex;
  flex-direction: column;
  gap: 16px;
  transition: transform 0.2s, box-shadow 0.2s;
}

.monitor-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
}

.app-info {
  display: flex;
  gap: 12px;
  overflow: hidden;
}

.app-icon {
  width: 48px;
  height: 48px;
  border-radius: 10px;
  object-fit: cover;
  background: var(--gray-light);
  flex-shrink: 0;
}

.app-details {
  display: flex;
  flex-direction: column;
  justify-content: center;
  overflow: hidden;
}

.app-name {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.app-link {
  font-size: 12px;
  color: var(--primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  text-decoration: none;
}

.app-link:hover {
  text-decoration: underline;
}

.card-meta {
  font-size: 12px;
  color: var(--text-secondary);
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.meta-item {
  display: flex;
  gap: 4px;
}

.meta-item .label {
  color: var(--text-secondary);
}

.meta-item .value {
  color: var(--text-primary);
  font-weight: 500;
}

.card-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: auto;
  padding-top: 12px;
  border-top: 1px solid var(--gray-light);
}

.left-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  font-size: 13px;
  font-weight: 500;
  padding: 6px 12px;
  border-radius: 6px;
  background: var(--gray-light);
  color: var(--text-primary);
  border: none;
  cursor: pointer;
}

.action-btn.pause:hover,
.action-btn.edit:hover {
  background: #e5e5ea;
}

.action-btn.resume {
  background: var(--primary);
  color: white;
}

.action-btn.resume:hover {
  background: var(--primary-hover);
}

.icon-btn {
  background: none;
  border: none;
  padding: 8px;
  border-radius: 6px;
  cursor: pointer;
}

.icon-btn.delete {
  color: var(--text-secondary);
}

.icon-btn.delete:hover {
  color: var(--danger);
  background: rgba(255, 59, 48, 0.1);
}

.edit-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 12px;
  background: var(--bg-color);
  border-radius: 8px;
}

.edit-row {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.edit-row label {
  font-size: 12px;
  color: var(--text-secondary);
  font-weight: 500;
}

.edit-row input {
  padding: 8px 10px;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  font-size: 14px;
  width: 100%;
}

.edit-row input:focus {
  outline: none;
  border-color: var(--primary);
}

.duration-btns {
  display: flex;
  gap: 6px;
}

.duration-btn {
  flex: 1;
  padding: 6px 8px;
  font-size: 12px;
  border-radius: 6px;
  background: white;
  border: 1px solid var(--border-color);
  cursor: pointer;
}

.duration-btn.active {
  background: var(--primary);
  color: white;
  border-color: var(--primary);
}

.edit-actions {
  display: flex;
  gap: 8px;
  margin-top: 4px;
}

.save-btn {
  flex: 1;
  padding: 8px;
  background: var(--primary);
  color: white;
  border: none;
  border-radius: 6px;
  font-weight: 500;
  cursor: pointer;
}

.save-btn:hover {
  background: var(--primary-hover);
}

.cancel-btn {
  padding: 8px 12px;
  background: var(--gray-light);
  border: none;
  border-radius: 6px;
  cursor: pointer;
}

.cancel-btn:hover {
  background: #e5e5ea;
}
</style>
