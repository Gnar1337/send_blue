<template>
	<div class="scheduled-section">
	      <div class="scheduled-header">
	        <span class="dot"></span>
	        <span class="title">Scheduled Messages</span>
	        <span class="count">({{ scheduledMessages.length }})</span>
	      </div>
	
	      <div class="list">
	        <div v-for="(m, i) in scheduledMessages" :key="i" class="message-card">
	          <div class="left-icon">
	            <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="#3b6cf8" stroke-width="1.6">
	              <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.86 19.86 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6A19.86 19.86 0 0 1 2.09 4.18 2 2 0 0 1 4 2h3a2 2 0 0 1 2 1.72c.12 1.02.5 2.02 1.12 2.9a2 2 0 0 1-.45 2.54L9.7 11.7a15.05 15.05 0 0 0 6 6l1.54-1.01a2 2 0 0 1 2.54-.45c.88.62 1.88 1 2.9 1.12A2 2 0 0 1 22 16.92z"></path>
	            </svg>
	          </div>
	
	          <div class="content">
	            <div class="phone">{{ m.toClientLead }}</div>
	            <div class="text">{{ m.messageBody }}</div>
	            <div class="meta">
	              <svg class="icon-clock" viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="#6b7280" stroke-width="1.6">
	                <circle cx="12" cy="12" r="9"></circle>
	                <path d="M12 7v6l4 2"></path>
	              </svg>
	              <span class="time">{{ formatDate(m.scheduledSendTime) }}</span>
	            </div>
	          </div>
	        </div>
	      </div>
	    </div>
</template>

<script setup lang="ts">
// import { m } from 'vue-router/dist/router-CWoNjPRp.mjs';
import { MessageQueue } from '../d.ts'
defineProps<{
	// scheduledMessages: { phone: string; message: string; when: string; }[];
	scheduledMessages: MessageQueue[];
	formatDate: (iso: any) => string;
}>()
</script>
<style scoped>
/* Scheduled list */
.scheduled-section {
  padding-top: 6px;
}

.scheduled-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 14px;
  color: #cbd5e1;
}

.dot {
  width: 6px;
  height: 6px;
  background: #60a5fa;
  border-radius: 50%;
  box-shadow: 0 0 0 6px rgba(96,165,250,0.06);
}

.title {
  font-size: 15px;
  font-weight: 500;
}

.count {
  margin-left: auto;
  color: #94a3b8;
  font-size: 14px;
}

/* Message card */
.message-card {
  display: flex;
  gap: 16px;
  align-items: flex-start;
  background: linear-gradient(180deg, rgba(255,255,255,0.02), rgba(255,255,255,0.01));
  border-radius: 12px;
  padding: 18px;
  margin-bottom: 16px;
  border: 1px solid rgba(255,255,255,0.03);
  box-shadow: 0 8px 18px rgba(2, 6, 23, 0.55);
}

.left-icon {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: linear-gradient(180deg, #072a54 0%, #061232 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: 0 4px 8px rgba(13,36,103,0.28);
}

.left-icon svg {
  stroke: #60a5fa;
}

.content {
  flex: 1;
}

.phone {
  font-weight: 600;
  color: #ffffff;
  margin-bottom: 8px;
}

.text {
  color: #cbd5e1;
  margin-bottom: 12px;
  line-height: 1.45;
}

.meta {
  display: flex;
  gap: 8px;
  align-items: center;
  color: #94a3b8;
  font-size: 13px;
}

.icon-clock {
  stroke: #94a3b8;
}
.time {
  font-size: 13px;
  color: #94a3b8;
}
</style>