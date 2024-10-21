import { io, Socket } from 'socket.io-client';
import { API_BASE_URL } from './apiUtils';

interface   ChatMessage {
    sender: string;
    content: string;
    timestamp: string;
}

class   ChatService {
    private socket: Socket | null = null;

    public  connectToWebsocket(): void {
        this.socket = io(API_BASE_URL, {
        transports: ['websocket'],
        autoConnect: false
        });

        this.socket.on('connect', this.onConnect); 
        this.socket.on('disconnect', this.onDisconnect); 
        this.socket.on('chat message', this.onChatMessage); 

        this.socket.connect();
    }

    public  sendMessage(content: string): void {
        if (this.socket?.connected) {
            const   message: ChatMessage = {
                sender: 'User',
                content: content,
                timestamp: new Date().toISOString()
            }
            this.socket.emit('chat message');
        } else {
            console.log("Socket is not connected"); 
        }
    }

    public  disconnect(): void {
        if (this.socket) {
            this.socket.close();
            this.socket = null;
        }
    }

    private onConnect = (): void => {
        console.log('Connected to the server !');
    };

    private onDisconnect = (reason: string): void => {
        console.log('Disconnected from the server:  !', reason);
    };

    private onChatMessage = (message: string): void => {
        console.log('Received message: ', message);
    };
}
