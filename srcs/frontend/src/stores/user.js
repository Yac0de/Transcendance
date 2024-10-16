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
		isLoggedIn: (state) => !!state.id
	},
	actions: {
		setUser(userData) {
			this.$patch(userData)
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
				console.log('Not connected, no data fetched, user store empty', error)
				this.clearUser()
			}
		},
		clearUser() {
			this.$reset()
		}
	}
})
