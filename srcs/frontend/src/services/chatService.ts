import { io, Socket } from 'socket.io-client';
import { API_BASE_URL } from './apiUtils';

interface   ChatMessage {
    sender: string;
    content: string;
    timestamp: string;
}

export  class   WebSocketService {
    private ws: WebSocket | null = null;
    private clientId: string;

    constructor(clientId: string) {
        this.clientId = clientId;
    }

    public  connect(): void {
        try {
            console.log(`ws://${API_BASE_URL}/ws?=${this.clientId}`);
            this.ws = new WebSocket(`ws://${API_BASE_URL}/ws?${this.clientId}`);

            this.ws.onopen = () => {
                console.log('Websocket connected!');
            };

            this.ws.onclose = (event) => {
                console.log('Disconnected to Websocket!, ', event.reason);
            };

            this.ws.onerror = (error) => {
                console.error('Websocket error, ', error);
            };

            this.ws.onmessage = (event) => {
                try {
                    const message = JSON.parse(event.data);
                    console.log('');
                } catch (e) {
                    console.error('Error parsing message, invalid format ?', e); 
                }
            };

        } catch (error) {
                console.error('Could not connect to the ws: ', error);
        }
    }

    public  sendMessage(content: string): void {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            const message = ChatMessage = {
                sender: 'User',
                content: content,
                timestamp: new Date().toISOString()
            };
            this.ws.send(JSON.stringify(message));
        } else {
            console.warn("Can't send a message, ws is not connected");
        }
    }

    public  disconnect(): void {
        if (this.ws) {
            this.ws.close();
            this.ws = null;
        }
    }

    public  isConnected(): bool {
        return this.ws !== null && this.ws.readyState === WebSocket.OPEN;
    }
}
