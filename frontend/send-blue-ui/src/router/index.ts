import { createRouter, createWebHistory } from 'vue-router'
// import HomeView from '../views/HomeView.vue'
import Schedule from '../views/Schedule.vue'
import Dash from '../views/Dash.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
        {
      path: '/dash/:clientId?',
      name: 'dash',
      component: Dash,
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
