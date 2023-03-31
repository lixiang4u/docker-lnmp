import { createRouter, createWebHistory } from 'vue-router'
import RootComponent from "../views/project/ListView.vue";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: RootComponent,
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
      path: '/container/list',
      name: 'containerList',
      component: () => import('../views/container/ListView.vue')
    },
    {
      path: '/container/logs/:id',
      name: 'containerLogs',
      component: () => import('../views/container/LogsView.vue')
    },
    {
      path: '/image/list',
      name: 'imageList',
      component: () => import('../views/image/ListView.vue')
    },
    {
      path: '/project/list',
      name: 'projectList',
      component: () => import('../views/project/ListView.vue')
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
