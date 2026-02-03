<script setup lang="ts">
import { RouterLink, RouterView, useRoute } from 'vue-router'
import type { Client } from './data.ts'
import { fetchClients, fetchClientData } from './data.ts'
import { clientStore } from './stores/client.ts';
import router from './router/index.ts';
const store = clientStore()
</script>
<template>
  <div class="app-container">

    <!-- TOP BANNER -->
    <header class="top-banner">
      <img alt="Vue logo" class="logo" src="@/assets/send_blue_logo.jpg" width="80" />

      <nav class="nav-links">
        <RouterLink to="/">Refresh</RouterLink>
        <RouterLink :to="`/dash/${selectedClientId}`">Dashboard</RouterLink>
        <RouterLink :to="`/schedule/${selectedClientId}`">Schedule</RouterLink>
      </nav>

      <div class="client-select-wrapper">
        login as:
        <select class="client-select" v-model="selectedClientId" @change="onClientChange">
          <option value="">Select a client...</option>
          <option v-for="c in clients" :key="c.uid" :value="c.uid">
            {{ c.name }}
          </option>
        </select>
      </div>
    </header>

    <!-- ROUTER CONTENT BELOW -->
    <main class="route-container">
      <RouterView :key="update" @clients-loaded="onClientsLoaded" />
    </main>

  </div>
</template>
<script lang="ts">
export default {
  name: 'App',
  data() {
    return {
      selectedClientId: '',
      update: 0,
      clients: [] as Client[],
      clientsLoading: false,
      clientsError: false,
      schedule: true,
    };
  },
  created() {
    console.log('App created hook');
  },
  methods: {
    onClientsLoaded(clients: Client[]) {
      this.clients = clients;
      if (this.clients.length > 0) {
        this.selectedClientId = this.clients[0]?.uid ?? '';
        var basePath = '/schedule/' + this.selectedClientId
        router.push({ path: basePath });
        clientStore().setCurrClient(this.selectedClientId)
      } else{
        router.push({ path: '/' });
      }
    },
    onClientChange() {
      var basePath = '/' + router.currentRoute.value.path.split('/')[1]
      router.push({ path: `${basePath}/${this.selectedClientId}` });
      clientStore().setCurrClient(this.selectedClientId).finally( () => {
        this.update++
      });
    },
  },
};
</script>
<style scoped>
.app-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
}

/* FIXED TOP BANNER */
.top-banner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 20px;
  background: rgba(255,255,255,0.05);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(255,255,255,0.1);
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  z-index: 1000;
}

.logo {
  height: 60px;
  width: 60px;
}

/* NAV LINKS */
.nav-links a {
  margin: 0 10px;
  color: #4686e7;
  text-decoration: none;
  font-weight: 600;
}

.nav-links a.router-link-exact-active {
  color: white;
}

/* CLIENT DROPDOWN */
.client-select-wrapper {
  min-width: 200px;
}

.client-select {
  width: 100%;
  padding: 8px 12px;
  border-radius: 8px;
  background: rgba(255,255,255,0.1);
  color: #4686e7;
  border: 1px solid rgba(255,255,255,0.2);
}

/* ROUTER CONTENT BELOW HEADER */
.route-container {
  margin-top: 100px; /* height of header */
  overflow-y: auto;
  padding: 0;
  flex: 1;
  display: block;
  width: 100%;
  max-width: 100%;
}
</style>
