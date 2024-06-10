import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import Error500View from '../views/Error500View.vue'
import db from '../database/index'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView
    },
    {
      path: '/error-500',
      name: 'error500',
      component: Error500View
    },
    // {
    //   path: '/about',
    //   name: 'about',
    //   // route level code-splitting
    //   // this generates a separate chunk (About.[hash].js) for this route
    //   // which is lazy-loaded when the route is visited.
    //   component: () => import('../views/AboutView.vue')
    // }
  ]
})

router.beforeEach((to, from, next) => {
  if (to.name === 'error500') {
    next()
    return
  }

  if (to.name === 'login' && db.authenticated()) {
    next({ name: 'home' })
    return
  }

  if (to.name !== 'login' && !db.authenticated()) {
    next({ name: 'login' })
  }
  else {
    next()
  }
})

export default router
