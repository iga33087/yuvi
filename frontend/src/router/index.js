import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Home',
    meta: {
      title:'Home',
    },
    component: () => import('../views/Home.vue')
  },
  {
    path: '/result',
    name: 'Result',
    meta: {
      title:'Result',
    },
    component: () => import('../views/result.vue')
  },
  {
    path: '/edit',
    name: 'Edit',
    meta: {
      title:'Edit',
    },
    component: () => import('../views/Edit.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router