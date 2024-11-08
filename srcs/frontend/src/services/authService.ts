import { apiRequest } from './apiUtils';
import { useUserStore } from '../stores/user';
import { Credentials } from '../types/models';

export default {
    async signin(credentials: Credentials): Promise<any> {
        return apiRequest('/auth/signin', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: "include",
            body: JSON.stringify(credentials),
        });
    },

    async signup(credentials: Credentials): Promise<any> {
        return apiRequest('/auth/signup', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(credentials),
        });
    },

    async signout(): Promise<void> {
        return apiRequest('/auth/signout', {
            method: 'POST',
            credentials: 'include',
        });
    },

    async isAuthenticated(): Promise<boolean> {
        const   userStore = useUserStore();
        return userStore.isSignedIn;
    }
};
