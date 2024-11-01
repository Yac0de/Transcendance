import { defineStore } from 'pinia'
import api from '../services/api'
import { WebSocketService } from '../services/websocketService';

export const useUserStore = defineStore('user', {
  state: () => ({
    id: null,
    nickname: null,
    displayname: null,
    avatar: null,
    webSocketService: null, 
  }),
  getters: {
    getId: (state) => state.id,
    getNickname: (state) => state.nickname,
    getDisplayName: (state) => state.displayname,
    getAvatarPath: (state) => state.avatar,
    isSignedIn: (state) => !!state.id,
    getWebSocketService: (state) => state.webSocketService
  },
  actions: {
    setUser(userData) {
      this.$patch({
        id: userData.id,
        nickname: userData.nickname,
        displayname: userData.displayname,
        avatar: userData.avatar
      });
      const storageData = {
        id: userData.id,
        nickname: userData.nickname,
        displayname: userData.displayname,
        avatar: userData.avatar
      };
      localStorage.setItem('userData', JSON.stringify(storageData))
    },

    setWebSocketService(userId) {
      const webSocketService = new WebSocketService(userId);
      webSocketService.connect();
      this.webSocketService = webSocketService;
    },

    async fetchUser() {
      try {
        const userData = await api.user.getUserData()
        if (userData) {
          const currentWS = this.webSocketService;
          
          this.setUser(userData);
          
          if (currentWS) {
            this.webSocketService = currentWS;
          }
        } else {
          this.clearUser()
        }
      } catch (error) {
        console.error('Not connected, no data fetched, user store empty', error)
        this.clearUser()
      }
    },

    clearUser() {
      if (this.webSocketService) {
        this.webSocketService.disconnect();
        this.webSocketService = null;
      }
      this.$reset()
      localStorage.removeItem('userData')
    },

    loadUserFromStorage() {
      const storedUser = localStorage.getItem('userData')
      if (storedUser) {
        const userData = JSON.parse(storedUser)
        this.$patch(userData)

        const webSocketService = new WebSocketService(userData.id);
        webSocketService.connect();
        this.webSocketService = webSocketService;

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
