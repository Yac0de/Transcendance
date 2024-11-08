import { defineStore } from 'pinia'
import api from '../services/api'
import { WebSocketService } from '../services/websocketService';
import { useOnlineUsersStore } from '../stores/onlineUsers';
import { UserState, UserData } from '../types/models';

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    id: null,
    nickname: null,
    displayname: null,
    avatar: null,
    webSocketService: null, 
  }),
  getters: {
    getId: (state: UserState): number | null => state.id,
    getNickname: (state: UserState): string | null => state.nickname,
    getDisplayName: (state: UserState): string | null => state.displayname,
    getAvatarPath: (state: UserState): string | null => state.avatar,
    isSignedIn: (state: UserState): boolean => !!state.id,
    getWebSocketService: (state: UserState): InstanceType<typeof WebSocketService> | null => state.webSocketService
  },
  actions: {
    setUser(userData: UserData) {
      this.$patch({
        id: userData.id,
        nickname: userData.nickname,
        displayname: userData.displayname,
        avatar: userData.avatar
      });
      const storageData: UserData = {
        id: userData.id,
        nickname: userData.nickname,
        displayname: userData.displayname,
        avatar: userData.avatar
      };
      localStorage.setItem('userData', JSON.stringify(storageData))
    },

    setWebSocketService(userId: number) {
      const store = useOnlineUsersStore();
      const webSocketService: WebSocketService = new WebSocketService(userId, store);
      webSocketService.connect();
      this.webSocketService = webSocketService;
    },

    async fetchUser() {
      try {
        const userData: UserData | null = await api.user.getUserData()
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
        this.webSocketService?.disconnect();
        this.webSocketService = null;
      }
      this.$reset()
      localStorage.removeItem('userData')
    },

    loadUserFromStorage() {
      const storedUser = localStorage.getItem('userData')
      if (storedUser) {
        const userData: UserData = JSON.parse(storedUser)
        this.$patch(userData)

        const store = useOnlineUsersStore();
        const webSocketService = new WebSocketService(userData.id, store);
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
