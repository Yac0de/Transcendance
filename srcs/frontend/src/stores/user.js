import { defineStore } from 'pinia'
import api from '../services/api'

export const useUserStore = defineStore('user', {
  state: () => ({
    id: null,
    nickname: null,
    displayname: null,
    avatar: null,
  }),
  getters: {
    getNickname: (state) => state.nickname,
    getDisplayName: (state) => state.displayname,
    getAvatarPath: (state) => state.avatar,
    isSignedIn: (state) => !!state.id
  },
  actions: {
    setUser(userData) {
      this.$patch(userData)
      localStorage.setItem('userData', JSON.stringify(userData))
    },
    async fetchUser() {
      try {
        const userData = await api.user.getUserData()
        if (userData) {
          this.setUser(userData)
        } else {
          this.clearUser()
        }
      } catch (error) {
        console.error('Not connected, no data fetched, user store empty', error)
        this.clearUser()
      }
    },
    clearUser() {
      this.$reset()
      localStorage.removeItem('userData')
    },
    loadUserFromStorage() {
      const storedUser = localStorage.getItem('userData')
      if (storedUser) {
        const userData = JSON.parse(storedUser)
        this.$patch(userData)
        return true
      }
      return false
    },
    async initializeStore() {
      const loadedFromStorage = this.loadUserFromStorage()
      return loadedFromStorage
    }
  },
})

export async function ensureStoreInitialized() {
  const store = useUserStore()
  await store.initializeStore()
  return store
}
