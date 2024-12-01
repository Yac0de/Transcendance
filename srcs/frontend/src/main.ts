import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { ensureStoreInitialized } from './stores/user'
import i18n from '../i18n.ts'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(i18n);

ensureStoreInitialized().then(() => {
  app.use(router)
  app.mount('#app')
})
