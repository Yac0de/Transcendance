import { UserData } from '../types/models'
import api from '../services/api'

export const fetchUserById = async (userId: number): Promise<UserData | null> => {
  try {
    const user = await api.user.getOtherUserData(userId);
    if (!user) {
      return null;
    }
    return user;
  } catch (error) {
    console.error('Error, not connected or invalid inviter id ', error)
  }
  return null;
};

export function getBaseHost(): string {
  if (import.meta.env.PROD) {
    const currentHost = window.location.host;
    return `https://${currentHost}`;
  }
  return 'http://localhost';
}

export const fetchMultipleUsers = async (userIds: number[]) => {
    try {
        const userPromises = userIds.map(id =>
            id !== 0 ? fetchUserById(id).catch(() => null) : Promise.resolve(null)
        );
        return await Promise.all(userPromises);
    } catch (error) {
        console.error("One or more fetches failed:", error);
        return Array(userIds.length).fill(null);
    }
};
