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
    // In production, we can use window.location to get the current host
    // This will automatically handle IP addresses or domain names
    const currentHost = window.location.host; // This includes host:port
    return `https://${currentHost}`;
  }
  // In development, use localhost with direct ports
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
        return Array(userIds.length).fill(null); // Return null users instead of throwing
    }
};
