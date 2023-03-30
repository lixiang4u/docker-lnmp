import { createRouter, createWebHistory } from 'vue-router'
import ComposeStatus from "../views/compose/StatusView.vue";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: ComposeStatus
    },
    {
      path: '/host/list',
      name: 'hostList',
      component: () => import('../views/host/ListView.vue')
    },
    {
      path: '/host/show/:domain',
      name: 'hostShow',
      component: () => import('../views/host/ShowView.vue')
    },
    {
      path: '/host/create',
      name: 'hostCreate',
      component: () => import('../views/host/ShowView.vue')
    },
    {
      path: '/host/update/:domain',
      name: 'hostUpdate',
      component: () => import('../views/host/ShowView.vue')
    },
    {
      path: '/compose/status',
      name: 'composeStatus',
      component: () => import('../views/compose/StatusView.vue')
    },
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/AboutView.vue')
    }
  ]
})

export default router
