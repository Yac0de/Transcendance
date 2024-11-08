import { apiRequest } from './apiUtils';
import { WebSocketService } from './websocketService';
import { ChatHistory } from '../types/models';

export default {
    async getChatHistory(friendId: number): Promise<ChatHistory | null> {
        try {
            return await apiRequest(`/conversation/${friendId}`, { credentials: "include" });
        } catch (error: unknown) {
            if ((error as any).message === 'Unauthorized') {
                return null;
            }
            throw new Error('Fetching chat history failed');
        }
    },
}

export { WebSocketService };
