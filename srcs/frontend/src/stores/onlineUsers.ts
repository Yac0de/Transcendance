import { defineStore } from 'pinia'

export const useOnlineUsersStore = defineStore('onlineUsers', {
  state: () => {
    return {
      onlineUsers: [] as number[]
    }
  },
  
  actions: {
    setOnlineUsers(users: number[]) {
      this.onlineUsers = users || [];
    },
    
    addOnlineUser(userId: number) {
      if (!this.onlineUsers.includes(userId)) {
        this.onlineUsers.push(userId);
      }
    },
    
    removeOnlineUser(userId: number) {
      this.onlineUsers = this.onlineUsers.filter((id: number) => id !== userId);
    }
  },
  
  getters: {
    isUserOnline: (state) => {
      return (userId: number) => state.onlineUsers?.includes(userId)
    },
    getOnlineUsers: (state) => {
      return state.onlineUsers
    }
  }
})
