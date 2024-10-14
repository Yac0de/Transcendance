import { apiRequest } from './apiUtils';
import userService from './userService';

interface Credentials {
    nickname: string;
    password: string;
}

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
        try {
            const userData = await userService.getUserData();
            return userData !== null;
        } catch (error: any) {
            if (error.response && error.response.status === 401) {
                return false;
            }
            return false;
        }
    }
};
