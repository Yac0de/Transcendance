import { createRouter, createWebHistory, RouteRecordRaw, NavigationGuardNext, RouteLocationNormalized } from 'vue-router';
import HomeView from '../components/General/HomeView.vue'
import SignInForm from '../components/Auth/SignInForm.vue'
import SignUpForm from '../components/Auth/SignUpForm.vue'
import PongGame from '../components/Game/PongGame.vue'
import Game from '../components/Game/Game.vue'
import Account from '../components/User/Account/Account.vue'
import Lobby from '../components/Lobby/Lobby.vue'
import MatchHistory from '../components/User/Account/MatchHistory.vue'
import NotFound from '../components/General/NotFound.vue'
import Tournament from '../components/Tournament/Tournament.vue'
import TwoFa from '../components/Auth/TwoFa.vue'
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
    path: '/tournament', 
    component: Tournament, 
    meta: { requiresAuth: true } // Indique que l'authentification est nécessaire
  },
  {
    path: '/lobby', 
    component: Lobby, 
    meta: { requiresAuth: true } // Indique que l'authentification est nécessaire
  },
  {
    path: '/game',
    component: Game,
    meta: { requiresAuth: true },
    beforeEnter: (to: RouteLocationNormalized, _from: RouteLocationNormalized, next: NavigationGuardNext) => {
      if ('lobbyId' in to.query) {
        next(); 
      } else {
        next('/'); 
      }
    }
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
  },
  { 
    path: '/match_history/:nickname', 
    component: MatchHistory, 
    meta: { requiresAuth: true },
    props: true 
  },
  { 
    path: '/2fa', 
    component: TwoFa, 
    meta: { requiresAuth: true },
    props: true 
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation Guard (executé avant chaque changement de route)
router.beforeEach(async (to: RouteLocationNormalized, _from: RouteLocationNormalized, next: NavigationGuardNext) => {
  const requiresAuth = to.matched.some((record: RouteRecordRaw) => record.meta && record.meta.requiresAuth);
  const requiresGuest = to.matched.some((record: RouteRecordRaw) => record.meta && record.meta.requiresGuest);

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
