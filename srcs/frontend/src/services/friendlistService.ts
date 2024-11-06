import { apiRequest } from './apiUtils';
import { Friend } from '../types/models';

export default {
    async getFriendList(): Promise<Friend[]> {
        const response = await apiRequest('/users/friendships/list', {
            method: 'GET',
            credentials: 'include',
        });
        if (!response || !Array.isArray(response)) {
            return [];
        }
        return response;
    },

    async deleteFromFriendList(friendId: string): Promise<void> {
        return apiRequest(`/users/friendships/delete/${friendId}`, {
            method: 'DELETE',
            credentials: 'include',
        });
    },

    async sendFriendRequest(friendNickname: string): Promise<void> {
        return apiRequest(`/users/friendships/add?nickname=${encodeURIComponent(friendNickname)}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
        });
    },

    async getFriendRequests(): Promise<Friend[]> {
        return apiRequest('/users/friendships/requests', {
            method: 'GET',
            credentials: 'include',
        });
    },

    async acceptFriendRequest(friendId: string): Promise<void> {
        return apiRequest(`/users/friendships/accept/${friendId}`, {
            method: 'PUT',
            credentials: 'include',
        });
    },

    async denyFriendRequest(friendId: string): Promise<void> {
        return apiRequest(`/users/friendships/deny/${friendId}`, {
            method: 'PUT',
            credentials: 'include',
        });
    },
};
