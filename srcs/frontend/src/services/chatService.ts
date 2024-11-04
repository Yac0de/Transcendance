import { apiRequest } from './apiUtils';
import { WebSocketService } from './websocketService';
import { ChatMessage } from '../types/websocket';

export default {
    async getChatHistory(friendId: string): Promise<ChatMessage | null> {
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
