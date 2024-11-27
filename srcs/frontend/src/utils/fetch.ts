import { UserData } from '../types/models'
import api from '../services/api'

export const fetchUserById = async (userId: number): Promise<UserData | null> => {
  try {
    const user = await api.user.getOtherUserData(userId);
    if (!user) {
      console.log('User not found or unauthorized');
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
