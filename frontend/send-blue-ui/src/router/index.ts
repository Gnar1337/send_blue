import { createRouter, createWebHistory } from 'vue-router'
// import HomeView from '../views/HomeView.vue'
import Schedule from '../views/Schedule.vue'
import Dashboard from '../views/Dashboard.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
        {
      path: '/dash/:clientId?',
      name: 'dash',
      component: Dashboard,
    },
    {
      path: '/schedule/:clientId?',
      name: 'schedule',
      component: Schedule,
      props: (route) => ({
        clientId: route.params.clientId || '0' // Use '0' if param is empty or missing
      })
    },
  ],
})

export default router
