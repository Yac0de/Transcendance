const API_BASE_URL = 'http://localhost:4000';

interface Credentials {
    nickname: string;
    password: string;
}

interface UserData {
    nickname: string;
    displayname: string;
    avatar: string;
}

interface SignUpData {
    nickname: string;
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
        const result = await response.json();

        if (!response.ok) {
            throw new Error(result.error || 'Sign in failed');
        }
        return result;
    },

    async signup(userData: SignUpData): Promise<any> {
        const response = await fetch(`${API_BASE_URL}/auth/signup`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(userData),
        });
        const result = await response.json();

        if (!response.ok) {
            throw new Error(result.error || 'Signup failed');
        }
        return result;
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
        const data = await response.json();
        return data;
    },

    async updateUserProfile(userData: UserData, avatarFile: File | null): Promise<UserData> {
        const formData = new FormData();
        formData.append('nickname', userData.nickname);
        formData.append('displayname', userData.displayname);
        
        if (avatarFile) {
            formData.append('avatar', avatarFile);
        }

        const response = await fetch(`${API_BASE_URL}/users/update-profile`, {
            method: 'PUT',
            credentials: "include",
            body: formData,
        });
        const result = await response.json();

        if (!response.ok) {
            throw new Error(result.error || 'Profile update failed');
        }
        return result;
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
    },

    async deleteUserAccount(password: string): Promise<void> {
        const response = await fetch(`${API_BASE_URL}/users/delete-account`, {
          method: 'DELETE',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
          body: JSON.stringify({ password }),
        });
        if (!response.ok) {
          const errorData = await response.json();
          throw new Error(errorData.error || 'Account deletion failed');
        }
      }
};