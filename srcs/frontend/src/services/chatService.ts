import { io, Socket } from 'socket.io-client';

interface ChatMessage {
    sender: string;
    content: string;
    timestamp: string;
}

type Event = {
    type: string
    senderId: number
    receiverId: number
    data: { text: string }
}

export class WebSocketService {
    private ws: WebSocket | null = null;
    private clientId: string;

    constructor(clientId: string) {
        this.clientId = clientId;
    }

    public connect(): void {
        try {
            const url = `ws://localhost:4001/ws?id=${this.clientId}`
            console.log(url);
            this.ws = new WebSocket(url);

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
                    console.log(message);
                } catch (e) {
                    console.error('Error parsing message, invalid format ?', e);
                }
            };

        } catch (error) {
            console.error('Could not connect to the ws: ', error);
        }
    }


    public sendMessage(content: string): void {
        //const message: Event = {
        //    type: "message",
        //    senderId: parseInt(this.clientId),
        //    receiverId: -1, // Get receiver id from args of this function
        //    data: {
        //        text: content
        //    }
        //}
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

    public disconnect(): void {
        if (this.ws) {
            this.ws.close();
            this.ws = null;
        }
    }

    public isConnected(): bool {
        return this.ws !== null && this.ws.readyState === WebSocket.OPEN;
    }
}
