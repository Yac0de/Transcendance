import { BaseMessage, ChatMessage, OnlineUsersMessage, UserStatusMessage } from '../types/websocket';
import { useOnlineUsersStore } from '../stores/onlineUsers';
import { useOnlineUsersStore } from '../stores/onlineUsers';

export class WebSocketService {
    private ws: WebSocket | null = null;
    private clientId: string;
    private onlineUsersStore = useOnlineUsersStore();

    private messageHandlers = {
    'CHAT': (message: ChatMessage) => {
        }
    }

    constructor(clientId: string) {
        this.clientId = clientId;
        this.initMessageHandlers();
    }

    public initMessageHandlers(): void {
        this.setMessageHandler('ONLINE_USERS', (message: OnlineUsersMessage) => {
            this.onlineUsersStore.setOnlineUsers(message.UsersOnline);
            console.log("RECEIVED THE LIST");
        });

        this.setMessageHandler('USER_DISCONNECTED', (message: UserStatusMessage) => {
            this.onlineUsersStore.removeOnlineUser(message.User);
            console.log("USER DISCO");
        });

        this.setMessageHandler('NEW_CONNECTION', (message: UserStatusMessage) => {
            this.onlineUsersStore.addOnlineUser(message.User);
            console.log("USER CO");
        });
    }
 
    public setMessageHandler(type: string, handler: (message: any) => void): void {
        console.log("Setting handler for type: ", type); 
        this.messageHandlers[type] = handler;
    }

    public connect(): void {
        try {
            const url = `ws://localhost:4001/ws?id=${this.clientId}`
            console.log(url);
            this.ws = new WebSocket(url);
            this.ws.onopen = () => {
                console.log('Websocket connected!');
                console.log('WS ready state: ', this.ws.readyState);
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
                    console.log("Message returned from the server: ", message);
                    const handler = this.messageHandlers[message.Type];
                    if (handler) {
                        handler(message);
                    }
                } catch (e) {
                    console.error('Error parsing message, invalid format ?', e);
                }
            };
        } catch (error) {
            console.error('Could not connect to the ws: ', error);
        }
    }
    public sendMessage(content: string, senderID: string, receiverID: string): void {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            const message: ChatMessage = {
                Type: 'CHAT',
                Data: content,
                SenderID: senderID,
                ReceiverID: receiverID
            };
            console.log("MESSAGE SENT = ", message);
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
    public isConnected(): boolean {
        return this.ws !== null && this.ws.readyState === WebSocket.OPEN;
    }
}

