const API_BASE_URL = 'http://localhost:4000';

interface Credentials {
    nickname: string;
    password: string;
}

interface UserData {
    nickname: string;
    email: string;
    password: string;
    avatar: string;
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

        const response = await fetch(`${API_BASE_URL}/users/upload-avatar`, {
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

    async getUserData(): Promise<UserData> {
        const response = await fetch(`${API_BASE_URL}/users`, {
            credentials: "include",
        });

        if (!response.ok) {
            throw new Error('Fetching user data failed');
        }

        const data: UserData = await response.json();
        return data;
    },

    getAvatarUrl(avatarPath: string): string {
        console.log("Avatar: ", avatarPath)
        return `${API_BASE_URL}/users/avatar/${avatarPath}`;
    }
};
