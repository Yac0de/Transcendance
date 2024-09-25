const API_BASE_URL = 'http://localhost:4000/auth';

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
        const response = await fetch(`${API_BASE_URL}/signin`, {
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
        const response = await fetch(`${API_BASE_URL}/signup`, {
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
};
