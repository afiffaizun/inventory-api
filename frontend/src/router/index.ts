import { createRouter, createWebHistory } from 'vue-router'
import ItemView from '../views/ItemView.vue'
import CategoriesView from '../views/CategoriesView.vue'
import ReportsView from '../views/ReportsView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'items',
      component: ItemView,
    },
    {
      path: '/categories',
      name: 'categories',
      component: CategoriesView,
    },
    {
      path: '/reports',
      name: 'reports',
      component: ReportsView,
    },
  ],
})

export default router
