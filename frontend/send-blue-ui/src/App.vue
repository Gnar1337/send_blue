<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import {Client} from './.d.ts'
import router from './router/index.ts';
import Schedule from './views/Schedule.vue';
import Dash from './views/Dash.vue';
</script>

<template>
  <header>
    <img alt="Vue logo" class="logo" src="@/assets/send_blue_logo.jpg" width="125" height="125" />

    <div class="wrapper">
      <!-- <HelloWorld msg="Send Blue" /> -->
  <div v-if="clientsLoading" class="card form-card">
        Loading clients...
      </div>
         <div else class="card form-card">
      <label class="field-label">Client</label>
      <div class="client-row">
        <select class="client-select" v-model="selectedClientId" @change="onClientChange">
          <option value="">Select a client...</option>
          <option v-for="c in clients" :key="c.name" :value="c.uid">
            {{ c.name }}
          </option>
        </select>
      </div>
         </div>
      <nav>
        <!-- <RouterLink to="/">Home</RouterLink>
        <RouterLink to="/Dashboard">Dashboard</RouterLink>
         <RouterLink to="/schedule">Schedule</RouterLink> -->
      </nav>
    </div>
  </header>
  <Schedule :clientId="selectedClientId" />
  <!-- <Dash v-else /> -->
  <!-- <RouterView /> -->
</template>
<script lang="ts">
export default {
  name: 'App',
  data() {
    return {
      // clients: [
      //   { id: '1', name: 'Client A' },
      //   { id: '2', name: 'Client B' },
      // ],
      selectedClientId: '',
      clients: [] as Client[],
      clientsLoading: false,
      clientsError: false,
      schedule: true,
    };
    
  },
  created() {
    console.log('App created hook');
    this.fetchClients();
  },
  methods: {
      async fetchClients() {
      this.clientsLoading = true
      this.clientsError = false
      try {
        const base = import.meta.env.VITE_API_BASE || ''
        const res = await fetch(`${base}/clients`)
        if (!res.ok) throw new Error('Network')
        const data = await res.json()
        console.log('Fetched clients:', data)
        // Expecting [ { id, name, phone } ]
        console.log('Clients set to:', this.clients)
        this.clients = data.clients || []
      } catch (e) {
        this.clientsError = true
        console.error('fetchClients error', e)
      } finally {
        this.clientsLoading = false
      }
    },
    onClientChange() {
      console.log('Selected client ID:', this.selectedClientId);
      router.push({ path: `/schedule/${this.selectedClientId}` });
    },
  },
};
</script>
<style scoped>
.client-row {
  display: flex;
  gap: 10px;
  color: rgb(66, 148, 224);
  background: rgba(255, 255, 255, 0.02);
  align-items: center;
}

.client-select {
  flex: 1;
  padding: 12px 14px;
  border-radius: 10px;
  border: 1px solid rgba(255,255,255,0.04);
  background: rgba(255,255,255,0.02);
  color: #4686e7;
  outline: none;
}
header {
  line-height: 1.5;
  max-height: 100vh;
}

.logo {
  display: block;
  margin: 0 auto 2rem;
}

nav {
  width: 100%;
  font-size: 12px;
  text-align: center;
  margin-top: 2rem;
}

nav a.router-link-exact-active {
  color: var(--color-text);
}

nav a.router-link-exact-active:hover {
  background-color: transparent;
}

nav a {
  display: inline-block;
  padding: 0 1rem;
  border-left: 1px solid var(--color-border);
}

nav a:first-of-type {
  border: 0;
}

@media (min-width: 1024px) {
  header {
    display: flex;
    place-items: center;
    padding-right: calc(var(--section-gap) / 2);
  }

  .logo {
    margin: 0 2rem 0 0;
  }

  header .wrapper {
    display: flex;
    place-items: flex-start;
    flex-wrap: wrap;
  }

  nav {
    text-align: left;
    margin-left: -1rem;
    font-size: 1rem;

    padding: 1rem 0;
    margin-top: 1rem;
  }
}
</style>
