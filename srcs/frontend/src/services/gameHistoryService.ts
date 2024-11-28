import { apiRequest } from './apiUtils';
import type { GameHistory } from '../types/models';

interface GameHistoryResponse {
    data: GameHistory[];
}

export default {
    async getUserHistory(userId: number): Promise<GameHistory[] | null> {
        try {
            const response = await apiRequest<GameHistoryResponse>(`/api/game-history/${userId}`, { 
                credentials: "include" 
            });
            return response.data;
        } catch (error: unknown) {
            if ((error as any).message === 'Unauthorized') {
                return null;
            }
            throw new Error('Fetching game history failed');
        }
    }
};