const API_BASE_URL = 'http://localhost:4000';

interface Credentials {
    nickname: string;
    password: string;
}

interface UserData {
    nickname: string;
    email: string;
    password: string;
}

export default {
    async login(credentials: Credentials): Promise<any> {
        const response = await fetch(`${API_BASE_URL}/auth/signin`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: "include",
            body: JSON.stringify(credentials),
        });
        if (!response.ok) {
            throw new Error('Login failed');
        }
        return response.json();
    },

    async signup(userData: UserData): Promise<any> {
        const response = await fetch(`${API_BASE_URL}/auth/signup`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: "include",
            body: JSON.stringify(userData),
        });
        if (!response.ok) {
            throw new Error('Signup failed');
        }
        return response.json();
    },

    async uploadAvatar(file: File): Promise<string> {
        const formData = new FormData();
        formData.append('avatar', file);

        const response = await fetch(`${API_BASE_URL}/upload-avatar`, {
            method: 'POST',
            credentials: "include",
            body: formData,
        });

        if (!response.ok) {
            throw new Error('Avatar upload failed');
        }

        const data = await response.json();
        return data.avatarPath;
    },

    getAvatarUrl(avatarPath: string): string {
        return `${API_BASE_URL}/uploads/${avatarPath}`;
    }
};
