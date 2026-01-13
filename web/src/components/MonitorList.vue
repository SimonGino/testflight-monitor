<script setup lang="ts">
import MonitorCard from './MonitorCard.vue'
import type { Monitor } from '../types'
import type { Messages } from '../i18n'

defineProps<{
  monitors: Monitor[]
  t: Messages
}>()

defineEmits<{
  (e: 'toggle', id: number): void
  (e: 'delete', id: number): void
  (e: 'update', id: number, data: { interval?: number; duration?: number }): void
}>()
</script>

<template>
  <div class="monitor-grid">
    <MonitorCard
      v-for="monitor in monitors"
      :key="monitor.id"
      :monitor="monitor"
      :t="t"
      @toggle="$emit('toggle', $event)"
      @delete="$emit('delete', $event)"
      @update="(id, data) => $emit('update', id, data)"
    />

    <div v-if="monitors.length === 0" class="empty-state">
      <div class="empty-icon">📱</div>
      <h3>{{ t.empty.title }}</h3>
      <p>{{ t.empty.desc }}</p>
    </div>
  </div>
</template>

<style scoped>
.monitor-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  padding: 20px;
}

.empty-state {
  grid-column: 1 / -1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
  color: var(--text-secondary);
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.empty-state h3 {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 8px;
}
</style>
