import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../components/HomeView.vue'
import SignInForm from '../components/SignInForm.vue'
import SignUpForm from '../components/SignUpForm.vue'
import PongGame from '../components/PongGame.vue'
import Account from '../components/Account.vue'
import NotFound from '../components/NotFound.vue'
import api from '../services/api'

const routes = [
  { path: '/', component: HomeView },
  { path: '/signin',
    component: SignInForm,
    meta: { requiresGuest: true } // Empêche l'accès si l'utilisateur est connecté
  },
  { path: '/signup',
    component: SignUpForm,
    meta: { requiresGuest: true }
  },
  {
    path: '/pong', 
    component: PongGame, 
    meta: { requiresAuth: true } // Indique que l'authentification est nécessaire
  },
  { 
    path: '/account', 
    component: Account, 
    meta: { requiresAuth: true }
  },
  {
    path: '/:pathMatch(.*)*', 
    component: NotFound
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation Guard (executé avant chaque changement de route)
router.beforeEach(async (to, _from, next) => {
  const isAuthenticated = await api.isAuthenticated() // Vérifie si l'utilisateur est connecté

  // Si la route nécessite une authentification et que l'utilisateur n'est pas connecté
  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (!isAuthenticated) {
      next('/signin') // Redirige vers la page de connexion
    } else {
      next() // Continue vers la route demandée
    }
  }
  // Si la route nécessite d'être un invité (non connecté) et que l'utilisateur est connecté
  else if (to.matched.some(record => record.meta.requiresGuest)) {
    if (isAuthenticated) {
      next('/account') // Redirige vers la page du compte si l'utilisateur est connecté
    } else {
      next() // Continue vers la route demandée
    }
  } else {
    next() // Si aucune condition spéciale, continue
  }
})

export default router
