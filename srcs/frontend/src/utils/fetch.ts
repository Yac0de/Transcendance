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
