import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { ensureStoreInitialized } from './stores/user'

const app = createApp(App)
const pinia = createPinia()
app.use(pinia)

ensureStoreInitialized().then(() => {
  app.use(router)
  app.mount('#app')
})
