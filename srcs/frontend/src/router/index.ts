import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../components/HomeView.vue'
import SignInForm from '../components/SignInForm.vue'
import SignUpForm from '../components/SignUpForm.vue'
import PongGame from '../components/PongGame.vue'
import Account from '../components/Account.vue'

const routes = [
  { path: '/', component: HomeView },
  { path: '/signin', component: SignInForm },
  { path: '/signup', component: SignUpForm },
  { path: '/pong', component: PongGame },
  { path: '/account', component: Account }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
