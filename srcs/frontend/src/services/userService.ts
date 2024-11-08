import { apiRequest, API_BASE_URL } from './apiUtils';
import { UserData } from '../types/models';

export default {
    async getUserData(): Promise<UserData | null> {
        try {
            return await apiRequest('/users', { credentials: "include" });
        } catch (error: unknown) {
            if ((error as any).message === 'Unauthorized') {
                return null;
            }
            throw new Error('Fetching current user data failed');
        }
    },

    async getProfileData(nickname: string): Promise<UserData | null> {
        const params = new URLSearchParams({ nickname }).toString();
        try {
            return await apiRequest(`/users/?${params}`, { credentials: "include" });
        } catch (error: unknown) {
            if ((error as any).message === 'Unauthorized') {
                return null;
            }
            throw new Error('Fetching other user data failed');
        }
    },

    async updateUserProfile(userData: UserData, avatarFile: File | null): Promise<UserData> {
        const formData = new FormData();
        formData.append('nickname', userData.nickname);
        formData.append('displayname', userData.displayname);
        if (avatarFile) {
            formData.append('avatar', avatarFile);
        }

        return apiRequest('/users/update-profile', {
            method: 'PUT',
            credentials: 'include',
            body: formData,
        });
    },

    getAvatarUrl(avatarPath: string): string {
        const defaultAvatarPath = 'default.png';
        const finalAvatarPath = avatarPath || defaultAvatarPath;
        return `${API_BASE_URL}/users/avatar/${finalAvatarPath}`;
    },

    async deleteUserAccount(password: string): Promise<void> {
        return apiRequest('/users/delete-account', {
            method: 'DELETE',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify({ password }),
        });
    },

    async changePassword(currentPassword: string, newPassword: string): Promise<void> {
        console.log("Sending request with currentPassword:", currentPassword, "and newPassword:", newPassword);
        return apiRequest('/users/change-password', {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify({
                currentPassword,
                newPassword
            })
        });
    },
};
