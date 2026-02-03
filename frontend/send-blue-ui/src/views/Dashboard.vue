<template>
    <div class="dashboard">
        <div class="dashboard-container">
           <div class="table-section">
                <h2>Messages</h2>
                <table class="data-table">
                    <thead>
                        <tr>
                            <th>To Lead</th>
                            <th>Message Body</th>
                            <th>Scheduled Time</th>
                            <th>Time Sent</th>
                            <th>Status</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="item in combinedMessages" :key="item.msgUid" @mouseenter="onRowHover(item)" @mouseleave="onRowLeave">
                            <td>{{ item.toClientLead }}</td>
                            <td>{{ item.messageBody }}</td>
                            <td>{{ formatDate(item.scheduledSendTime) }}</td>
                            <td>{{ formatDate(item.timeSent) }}</td>
                            <td>{{ item.status }}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
                <div class="details-section">
        <h2>Message Details</h2>
        <div v-if="selectedMessage" class="details-content">
            <p><strong>UID:</strong> {{ selectedMessage.uid }}</p>
            <p><strong>From Client:</strong> {{ selectedMessage.fromClientId }}</p>
            <p><strong>To Lead:</strong> {{ selectedMessage.toClientLead }}</p>
            <p><strong>Message:</strong> {{ selectedMessage.messageBody }}</p>
            <p><strong>Status:</strong> {{ selectedMessage.status }}</p>

            <h3>Event History</h3>
            <div v-if="historyLoading" class="loading">Loading history...</div>
            <div v-else-if="messageHistory.length > 0" class="history-list">
                <div v-for="event in messageHistory" :key="event.time_stamp" class="history-item">
                    <span class="timestamp">{{ formatDate(event.time_stamp) }}</span>
                    <span class="status-change">→ {{ event.curr_status }}</span>
                </div>
            </div>
            <div v-else class="empty-state">No history available</div>
        </div>
        <div v-else class="empty-state">
            Hover over a message to view details
        </div>
    </div>
        </div>
    </div>
</template>

<script setup>
import { computed, onMounted, watch, ref } from 'vue';
import { clientStore } from '@/stores/client';
import { fetchMessageEventHistory } from '@/data.ts';
onMounted(() => {
});
const store = clientStore();
watch(() => [store.currClient.name, store.currClient.messageQueue.length, store.currClient.allMessagesSent.length], () => {
});
const combinedMessages = computed(() => {
    return [
        ...store.currClient.allMessagesSent,
        ...store.currClient.messageQueue
    ];
});

const formatDate = (date) => {
    if (!date) return '-';
    return new Date(date).toLocaleString();
};
const selectedMessage = ref(null);
const messageHistory = ref([]);
const historyLoading = ref(false);

const onRowHover = async (item) => {
    selectedMessage.value = item;
    historyLoading.value = true;
    messageHistory.value = (await fetchMessageEventHistory(item.uid)).values;
    historyLoading.value = false;
};
const onRowLeave = () => {
    selectedMessage.value = null;
    messageHistory.value = [];
    historyLoading.value = false;
};


</script>

<style scoped>
.dashboard {
    width: 100%;
    padding: 0;
}

.dashboard-container {
    display: grid;
    grid-template-columns: 2fr 1fr;
    gap: 30px;
    align-items: start;
}

.data-table tbody tr:hover {
    background-color: #666363;
    cursor: pointer;
}

.table-section,
.chart-section {
    background: rgb(31, 28, 28);
    padding: 20px;
    border-radius: 8px;
    width: 100%;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.data-table {
    width: 100%;
    border-collapse: collapse;
}

.data-table th,
.data-table td {
    padding: 12px;
    text-align: left;
    border-bottom: 1px solid #312e2e;
}

.data-table th {
    background-color: #2e2c2c;
    font-weight: 600;
}

h2 {
    margin-top: 0;
    margin-bottom: 20px;
}

.details-content {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.details-content p {
    margin: 0;
    padding: 8px;
    background: rgba(255, 255, 255, 0.05);
    border-radius: 4px;
    word-break: break-word;
}

.details-content strong {
    color: #4686e7;
}

.empty-state {
    color: #999;
    font-style: italic;
    text-align: center;
    padding: 40px 20px;
}
.history-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
    margin-top: 16px;
}

.history-item {
    padding: 16px;
    border-radius: 8px;
    border-left: 4px solid #4686e7;
    background: linear-gradient(135deg, rgba(70, 134, 231, 0.15) 0%, rgba(70, 134, 231, 0.05) 100%);
    transition: all 0.3s ease;
    display: flex;
    justify-content: space-between;
    align-items: center;
    box-shadow: 0 4px 12px rgba(70, 134, 231, 0.1);
}
</style>