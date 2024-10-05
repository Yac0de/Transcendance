const API_BASE_URL = 'http://localhost:4000';

interface Credentials {
    nickname: string;
    password: string;
}

interface UserData {
    nickname: string;
    email: string;
    avatar: string;
}

interface SignUpData {
    nickname: string;
    email: string;
    password: string;
}

export default {
    async signin(credentials: Credentials): Promise<any> {
        const response = await fetch(`${API_BASE_URL}/auth/signin`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: "include",
            body: JSON.stringify(credentials),
        });
        if (!response.ok) {
            throw new Error('Sign in failed');
        }
        return response.json();
    },

    async signup(userData: SignUpData): Promise<any> {
        const response = await fetch(`${API_BASE_URL}/auth/signup`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(userData),
        });
        if (!response.ok) {
            throw new Error('Signup failed');
        }
        return response.json();
    },

    async getUserData(): Promise<UserData | null> {
        const response = await fetch(`${API_BASE_URL}/users`, {
            credentials: "include",
        });
        if (!response.ok) {
            if (response.status === 401) {
                return null;
            }
            throw new Error('Fetching user data failed');
        }
        return response.json();
    },

    async updateUserProfile(userData: UserData, avatarFile: File | null): Promise<UserData> {
        const formData = new FormData();
        formData.append('nickname', userData.nickname);
        formData.append('email', userData.email);
        
        if (avatarFile) {
            formData.append('avatar', avatarFile);
        }

        const response = await fetch(`${API_BASE_URL}/users/update-profile`, {
            method: 'PUT',
            credentials: "include",
            body: formData,
        });

        if (!response.ok) {
            throw new Error('Profile update failed');
        }

        return response.json();
    },

    getAvatarUrl(avatarPath: string): string {
        return `${API_BASE_URL}/users/avatar/${avatarPath}`;
    },

    async   signout(): Promise<void> {
        const   response = await fetch(`${API_BASE_URL}/auth/signout`, {
            method: 'POST',
            credentials: "include",
        })
        if (!response.ok) {
            throw new Error('Sign out failed');
        }
    },

    async   isAuthenticated(): Promise<boolean> {
        const   userData = await this.getUserData();
        return userData !== null;
    }
};
