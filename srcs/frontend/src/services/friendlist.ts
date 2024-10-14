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

    async deleteFromFriendList(friendId: string): Promise<void> {
        const response = await fetch(`${API_BASE_URL}/users/friendships/delete/${friendId}`, {
            method: 'POST',
            credentials: 'include',
        });

        console.log('Response status:', response.status);
        console.log('Response data:', response);

        if (!response.ok) {
            throw new Error('Failed to delete from friend list');
        }
    },

    async sendFriendRequest(friendNickname: string): Promise<void> {
        const response = await fetch(`${API_BASE_URL}/users/friendships/add?nickname=${encodeURIComponent(friendNickname)}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
        });

        console.log('Response status:', response.status);
        console.log('Response data:', response);

        if (!response.ok) {
            const error: any = new Error('Failed to send friend request'); 
            error.status = response.status
            throw error;
        }
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


    async acceptFriendRequest(friendId: string): Promise<void> {
        const response = await fetch(`${API_BASE_URL}/users/friendships/accept/${friendId}`, {
            method: 'POST',
            credentials: 'include',
        });

        console.log('Response status:', response.status);
        console.log('Response data:', response);

        if (!response.ok) {
            throw new Error('Failed to accept friend request');
        }
    },

    async denyFriendRequest(friendId: string): Promise<void> {
        const response = await fetch(`${API_BASE_URL}/users/friendships/deny/${friendId}`, {
            method: 'POST',
            credentials: 'include',
        });

        console.log('Response status:', response.status);
        console.log('Response data:', response);

        if (!response.ok) {
            throw new Error('Failed to deny friend request');
        }
    },
};
