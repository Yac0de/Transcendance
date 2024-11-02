import { defineStore } from 'pinia'

export const useOnlineUsersStore = defineStore('onlineUsers', {
  state: () => ({
    onlineUsers: [] as string[]
  }),

  actions: {
    setOnlineUsers(users: string[]) {
      this.onlineUsers = users
    },

    addOnlineUser(userId: string) {
      if (!this.onlineUsers.includes(userId)) {
        this.onlineUsers.push(userId)
      }
    },

    removeOnlineUser(userId: string) {
      this.onlineUsers = this.onlineUsers.filter(id => id !== userId)
    }
  },

  getters: {
    isUserOnline: (state) => {
      return (userId: string) => state.onlineUsers.includes(userId)
    },
    getOnlineUsers: (state) => state.onlineUsers
  }
})
