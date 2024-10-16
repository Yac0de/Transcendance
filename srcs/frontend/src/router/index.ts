import { createRouter, createWebHistory, routeRecordRaw } from 'vue-router'
import HomeView from '../components/HomeView.vue'
import SignInForm from '../components/SignInForm.vue'
import SignUpForm from '../components/SignUpForm.vue'
import PongGame from '../components/PongGame.vue'
import Account from '../components/User/Account.vue'
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
    path: '/account/:nickname', 
    component: Account, 
    meta: { requiresAuth: true },
    props: true 
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
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
  const requiresGuest = to.matched.some(record => record.meta.requiresGuest);

  if (requiresAuth) {
    const isAuthenticated = await api.auth.isAuthenticated();
    if (!isAuthenticated) {
      return next('/signin'); // Redirige vers la page de connexion
    }
  }

  if (requiresGuest) {
    const isAuthenticated = await api.auth.isAuthenticated();
    if (isAuthenticated) {
      return next('/account'); // Redirige vers la page du compte si l'utilisateur est connecté
    }
  }

  next(); // Si aucune condition spéciale, continue vers la route demandée
});


export default router
