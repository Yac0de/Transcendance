import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../components/General/HomeView.vue'
import SignInForm from '../components/Auth/SignInForm.vue'
import SignUpForm from '../components/Auth/SignUpForm.vue'
import PongGame from '../components/Game/PongGame.vue'
import Account from '../components/User/Account/Account.vue'
import NotFound from '../components/General/NotFound.vue'
import api from '../services/api'
import { useUserStore } from '../stores/user'

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
    path: '/:nickname', 
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

  const userStore = useUserStore();
  const nickname = userStore.getNickname;

  if (requiresAuth) {
    const isAuthenticated = await api.auth.isAuthenticated();
    if (!isAuthenticated) {
      return next('/signin');
    }
  }

  if (requiresGuest) {
    const isAuthenticated = await api.auth.isAuthenticated();
    if (isAuthenticated) {
      return next(`/${nickname}`);
    }
  }

  next();
});

export default router;
