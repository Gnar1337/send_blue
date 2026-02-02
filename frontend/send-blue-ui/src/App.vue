<script setup lang="ts">
import { RouterLink, RouterView, useRoute } from 'vue-router'
import type { Client } from './data.ts'
import { fetchClients, fetchClientData } from './data.ts'
import { clientStore } from './stores/client.ts';
import router from './router/index.ts';
const store = clientStore()
</script>

<template>
  <header>
    <img alt="Vue logo" class="logo" src="@/assets/send_blue_logo.jpg" width="125" height="125" />

    <div class="wrapper">
    <div v-if="clientsLoading" class="card form-card">
        Loading clients...
    </div>
    <div v-else class="card form-card">
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
        <RouterLink :to="`/dash/${selectedClientId}`">Dashboard</RouterLink>
        <RouterLink :to="`/schedule/${selectedClientId}`">Schedule</RouterLink>
      </nav>
    </div>
  </header>
  <RouterView v-if="selectedClientId != '' && !clientsLoading" :key="update" />
</template>
<script lang="ts">
export default {
  name: 'App',
  data() {
    return {
      selectedClientId: '',
      update: 0,
      clients: [] as Client[],
      clientsLoading: true,
      clientsError: false,
      schedule: true,
    };
  },
  created() {
    console.log('App created hook');
    fetchClients().then(data => {
      this.clients = data;
      this.clientsLoading = false;
    }).catch(err => {
      console.error('Error fetching clients:', err);
      this.clientsError = true;
      this.clientsLoading = false;
    })
  },
  watch: {
  clients(newVal) {
    if (newVal.length > 0) {
      this.selectedClientId = newVal[0].uid;
      this.onClientChange
    }
  }
}, 
  methods: {
    onClientChange() {
      var basePath = '/' + router.currentRoute.value.path.split('/')[1]
      console.log('Selected path = ', basePath + '/' + this.selectedClientId);
      router.push({ path: `${basePath}/${this.selectedClientId}` });
      clientStore().setCurrClient(this.selectedClientId).finally( () => {
        console.log(this.selectedClientId)
        this.update++
      });
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
