const API_BASE_URL = 'http://localhost:4000';

interface Friend {
    nickname: string;
}

export default {
    async getFriendList(): Promise<Friend[]> {
        const response = await fetch(`${API_BASE_URL}/user/friendlist`, {
            method: 'GET',
            credentials: 'include',
        });
        if (!response.ok) {
            throw new Error('Failed to fetch friend list');
        }
        return response.json();
    },

    async sendFriendRequest(friendNickname: string): Promise<void> {
        const response = await fetch(`${API_BASE_URL}/user/addfriend`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify({ friendNickname }),
        });
        if (!response.ok) {
            throw new Error('Failed to send friend request');
        }
    },
};
