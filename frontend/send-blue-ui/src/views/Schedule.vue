
<template>
  <div class="page">
    <div class="card form-card">
      <label class="field-label">Phone Number</label>
      <div class="lead-row">
        <select class="lead-select" v-model="phone" @change="onLeadChange">
          <option value="">Select a lead...</option>
          <option class="lead-option" v-for="l in leads" :key="l.leadNumber" :value="l.leadNumber">
            {{ l.leadNumber }}
          </option>
        </select>
      </div>

      <!-- working date + time picker -->
      <label class="field-label">When</label>
      <VueDatePicker
        v-model="dateTime"
        :is24="false"
        enable-time-picker
        placeholder="Select Date and Time"
        class="date-picker"
      />

      <label class="field-label">Message</label>
      <textarea
        v-model="message"
        class="message-input"
        placeholder="Enter your message here..."
        rows="5"
      ></textarea>

      <button class="schedule-btn" :disabled="!canSchedule" @click="scheduleMessage">
        <svg class="icon-plane" viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="1.6">
          <path d="M22 2L11 13"></path>
          <path d="M22 2l-7 20-4-9-9-4 20-7z"></path>
        </svg>
        <span>Schedule Message</span>
      </button>
    </div>

    <ScheduledMessages :scheduledMessages="scheduledMessages" :formatDate="formatDate" />
  </div>
</template>

<script lang="ts">
import type { ClientLead, MessageQueue } from '../data.ts'
import ScheduledMessages from './ScheduledMessages.vue'
import { VueDatePicker } from '@vuepic/vue-datepicker'
import { fetchClientLeads, fetchMessageQueue } from '../data.ts'
import '@vuepic/vue-datepicker/dist/main.css'
import { clientStore } from '@/stores/client.ts'

export default {
  name: 'Schedule',
  components: { ScheduledMessages, VueDatePicker },
  props: {
    clientId: {
      type: String,
      required: true,
    },
  },
  watch: {
    clientId: {
      handler(newVal: string) {
        console.log('clientId changed to:', newVal)
        if (newVal) {
          fetchClientLeads(newVal).then((leads) => {
            this.leads = leads
          })
          fetchMessageQueue(newVal).then((messages) => {
            this.scheduledMessages = messages
          })
        }
      },
      deep: true,
    },
  },
  created() {
    fetchClientLeads(this.clientId).then((leads) => {
      this.leads = leads
    })
    // default date/time to one hour from now
    this.dateTime = new Date(Date.now() + 60 * 60 * 1000)
  },
  mounted() {
    fetchClientLeads(this.clientId).then((leads) => {
        this.leads = leads
      })
  },
  data() {
    return {
      leads: [] as ClientLead[],
      selectedLeadNumber: '',
      message: '',
      scheduledMessages: [] as MessageQueue[],
      phone: '',
      dateTime: new Date(),
      clientsLoading: false,
      clientsError: false,
    }
  },
  computed: {
    canSchedule() {
      if (!this.phone || !this.message || !this.dateTime) return false
      // require future date/time
      const dt = new Date(this.dateTime)
      return dt.getTime() > Date.now()
    },
  },
  methods: {
    async scheduleMessageToBackend(msgToQueue: MessageQueue) {
      const base = import.meta.env.VITE_API_BASE || ''
      const res = await fetch(`${base}/client/schedule_message`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify(msgToQueue),
          })
      if (!res.ok) throw new Error('Failed to schedule message')
          const data = await res.json()
      if (data.error) {
        console.error('Backend error scheduling message:', data.error)
        throw new Error(data.error)
      } else{
          console.log('Message scheduled successfully:', data)
          fetchMessageQueue(this.clientId || '').then((messages) => {
            this.scheduledMessages = messages
            clientStore().setCurrClient(this.clientId)
          })

      }
    },
    onLeadChange() {
      // keep phone in sync with selected lead (leadNumber is a phone here)
      this.phone = this.phone || ''
    },
    scheduleMessage() {
      if (!this.canSchedule) {
        return
      } 
      const whenIso = this.dateTime instanceof Date ? this.dateTime.toISOString() : new Date(this.dateTime).toISOString()
      //create MeessageQueue entry in backend here
      let msgToQueue: MessageQueue = {
        msgUid:'',
        messageBody: this.message,
        fromClientId: this.clientId || '',
        toClientLead: this.phone,
        scheduledSendTime: this.dateTime,
        timeSent: null,
        status: 'QUEUED',
      } 
      this.scheduleMessageToBackend(msgToQueue).catch((err) => {
        console.error('Error scheduling message:', err)
      }).then(() => {
        // Message scheduled successfully, UI updated in scheduleMessageToBackend
        this.phone = ''
        this.message = ''
        this.dateTime = new Date(Date.now() + 60 * 60 * 1000)
      })
    //   this.scheduledMessages.push(msgToQueue);
    },
    formatDate(iso: string) {
      const d = new Date(iso)
      return d.toLocaleString(undefined, {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
        hour: 'numeric',
        minute: '2-digit',
      })
    },
  },
}
</script>

<style scoped>
.page {
  max-width: 920px;
  margin: 32px auto;
  padding: 0 20px;
  color: #e6eef8;
  font-family: Inter, system-ui, -apple-system, "Segoe UI", Roboto, "Helvetica Neue", Arial;
  background: linear-gradient(180deg, #071329 0%, #0b1220 100%);
}

.card {
  background: linear-gradient(180deg, rgba(255,255,255,0.02) 0%, rgba(255,255,255,0.01) 100%);
  border-radius: 12px;
  padding: 28px;
  border: 1px solid rgba(255,255,255,0.03);
  box-shadow: 0 18px 30px rgba(2, 6, 23, 0.6);
}

.form-card {
  margin-bottom: 28px;
}

.field-label {
  display: block;
  font-weight: 500;
  margin-bottom: 8px;
  color: #cbd5e1;
}

.phone-input,
.message-input {
  width: 100%;
  border: 1px solid rgba(255,255,255,0.04);
  background: rgba(255,255,255,0.02);
  border-radius: 10px;
  padding: 14px 16px;
  font-size: 15px;
  color: #e6eef8;
  outline: none;
  box-sizing: border-box;
  margin-bottom: 20px;
}

.phone-input::placeholder,
.message-input::placeholder {
  color: rgba(230,238,248,0.45);
}

.message-input {
  min-height: 110px;
  resize: vertical;
}

.schedule-btn {
  width: 100%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  background: linear-gradient(90deg, #2156ff 0%, #6b3df1 100%);
  color: #fff;
  padding: 14px 18px;
  border-radius: 12px;
  border: none;
  font-weight: 600;
  font-size: 16px;
  cursor: pointer;
  box-shadow: 0 12px 22px rgba(33, 86, 255, 0.28);
  transition: transform 120ms ease, box-shadow 120ms ease;
}

.schedule-btn:disabled {
  filter: grayscale(0.6) opacity(0.6);
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
  background: linear-gradient(90deg, rgba(33,86,255,0.18), rgba(107,61,241,0.18));
}

.schedule-btn:active:not(:disabled) {
  transform: translateY(1px);
}

.icon-plane {
  stroke: #ffffff;
  opacity: 0.95;
}

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
.lead-row {
  display: flex;
  gap: 10px;
  color: rgb(48, 113, 187);
   background: rgba(240, 236, 236, 0.02);
  align-items: center;
}

.lead-select {
  flex: 1;
  padding: 12px 14px;
  border-radius: 10px;
  border: 1px solid rgba(255,255,255,0.04);
  background: rgba(10, 2, 2, 0.02);
  color: #3c7bc4;
  outline: none;
}
.lead-option {
  flex: 1;
  padding: 12px 14px;
  border-radius: 10px;
  border: 1px solid rgba(255,255,255,0.04);
  background: rgba(0, 0, 0, 0.02);
  color: #4589f0;
  outline: none;
}

/* datepicker wrapper */
.date-picker {
  margin-bottom: 20px;
  display: block;
}

/* style datepicker input & popup (target internal classes with v-deep) */
:deep(.vdp__input) {
  width: 100%;
  border: 1px solid rgba(15, 12, 12, 0.04);
  background: rgba(12, 11, 11, 0.02);
  border-radius: 10px;
  padding: 12px 14px;
  font-size: 15px;
  color: #0b76f8;
  box-sizing: border-box;
}

/* calendar popup */
:deep(.vdp__calendar) {
  background: linear-gradient(180deg, rgba(17, 14, 43, 0.02), rgba(4, 23, 77, 0.01));
  color: #e6eef8;
  border-radius: 8px;
  border: 1px solid rgba(255,255,255,0.03);
}

/* time picker rows */
:deep(.vdp__time) {
  color: #e6eef8;
}
</style>
