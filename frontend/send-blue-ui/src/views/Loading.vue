th
<template>
    <div class="loading-page">
        <div class="loading-card">
            <div class="spinner-container">
                <svg class="spinner" viewBox="0 0 50 50">
                    <circle class="path" cx="25" cy="25" r="20" fill="none" stroke-width="4"></circle>
                </svg>
            </div>
            <p class="loading-text">Loading clients...</p>
        </div>
    </div>
</template>

<script lang="ts">
import { fetchClients } from '../data'

export default {
    name: 'Loading',
    emits: ['clients-loaded'],
    async created() {
        console.log('Loading.vue created hook')
        try {
            const clients = await fetchClients()
            const clientUid = clients.length > 0 && clients[0] ? clients[0].uid : ''
            this.$emit('clients-loaded', clients)
        } catch (error) {
            console.error('Error loading clients:', error)
            this.$emit('clients-loaded', [])
        }
    }
}

</script>

<style scoped>
.loading-page {
    width: 100%;
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 32px 20px;
    background: linear-gradient(180deg, #071329 0%, #0b1220 100%);
    font-family: Inter, system-ui, -apple-system, "Segoe UI", Roboto, "Helvetica Neue", Arial;
}

.loading-card {
    background: linear-gradient(180deg, rgba(99, 80, 80, 0.02) 0%, rgba(255,255,255,0.01) 100%);
    border-radius: 18px;
    padding: 48px 64px;
    border: 1px solid rgba(255,255,255,0.08);
    box-shadow: 0 20px 60px rgba(33, 86, 255, 0.3), 0 0 40px rgba(107, 61, 241, 0.2), inset 0 1px 0 rgba(255,255,255,0.1);
    backdrop-filter: blur(10px);
    text-align: center;
}

.spinner-container {
    display: flex;
    justify-content: center;
    margin-bottom: 24px;
}

.spinner {
    animation: rotate 2s linear infinite;
    width: 64px;
    height: 64px;
}

.spinner .path {
    stroke: #2156ff;
    stroke-linecap: round;
    animation: dash 1.5s ease-in-out infinite;
}

@keyframes rotate {
    100% {
        transform: rotate(360deg);
    }
}

@keyframes dash {
    0% {
        stroke-dasharray: 1, 150;
        stroke-dashoffset: 0;
    }
    50% {
        stroke-dasharray: 90, 150;
        stroke-dashoffset: -35;
    }
    100% {
        stroke-dasharray: 90, 150;
        stroke-dashoffset: -124;
    }
}

.loading-text {
    color: #cbd5e1;
    font-size: 16px;
    font-weight: 500;
    margin: 0;
}
</style>