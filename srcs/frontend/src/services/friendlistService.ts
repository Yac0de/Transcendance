import { apiRequest } from './apiUtils';

interface Friend {
    id: string;
    displayname: string;
    nickname: string;
    avatar: string;
}

export default {
    async getFriendList(): Promise<Friend[]> {
        return apiRequest('/users/friendships/list', {
            method: 'GET',
            credentials: 'include',
        });
    },

    async deleteFromFriendList(friendId: string): Promise<void> {
        return apiRequest(`/users/friendships/delete/${friendId}`, {
            method: 'POST',
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
            method: 'POST',
            credentials: 'include',
        });
    },

    async denyFriendRequest(friendId: string): Promise<void> {
        return apiRequest(`/users/friendships/deny/${friendId}`, {
            method: 'POST',
            credentials: 'include',
        });
    },
};
