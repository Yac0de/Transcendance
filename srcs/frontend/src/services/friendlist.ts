const API_BASE_URL = 'http://localhost:4000';

interface Friend {
    id: string;
    displayname: string;
    nickname: string;
    avatar: string;
}

export default {
    async getFriendList(): Promise<Friend[]> {
        const response = await fetch(`${API_BASE_URL}/users/friendships/list`, {
            method: 'GET',
            credentials: 'include',
        });
        if (!response.ok) {
            throw new Error('Failed to fetch friend list');
        }
        return response.json();
    },

    async getFriendRequests(): Promise<Friend[]> {
        const response = await fetch(`${API_BASE_URL}/users/friendships/requests`, {
            method: 'GET',
            credentials: 'include',
        });
        if (!response.ok) {
            throw new Error('Failed to fetch friend list');
        }
        return response.json();
    },

    async sendFriendRequest(friendUsername: string): Promise<void> {
        const response = await fetch(`${API_BASE_URL}/users/friendships/add?username=${encodeURIComponent(friendUsername)}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
        });

        console.log('Response status:', response.status);
        console.log('Response data:', response);

        if (!response.ok) {
            throw new Error('Failed to send friend request');
        }
    },
};
