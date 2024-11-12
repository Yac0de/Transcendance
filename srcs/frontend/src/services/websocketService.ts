import { OnlineUsersMessage, UserStatusMessage } from '../types/connection_status';
import { ChatMessage} from '../types/chat';
import { useOnlineUsersStore } from '../stores/onlineUsers';
import { useUserStore } from '../stores/user';
import { eventBus } from '../events/eventBus';

type MessageHandler<T> = (message: T) => void;
type MessageHandlers = {
    [key: string]: MessageHandler<ChatMessage | OnlineUsersMessage | UserStatusMessage>;
};

export class WebSocketService {
    public ws: WebSocket | null = null;
    public clientId: number;
    public onlineUsersStore: ReturnType<typeof useOnlineUsersStore>;
    public userStore: ReturnType<typeof useUserStore>;
    public messageHandlers: MessageHandlers = {};

    constructor(clientId: number, onlineUsersStore: ReturnType<typeof useOnlineUsersStore>, userStore: ReturnType<typeof useUserStore>) {
        this.clientId = clientId;
        this.onlineUsersStore = onlineUsersStore;
        this.userStore = userStore;
        this.initMessageHandlers();
    }

    public initMessageHandlers(): void {
        this.setMessageHandler<OnlineUsersMessage>('ONLINE_USERS', (message: OnlineUsersMessage) => {
            this.onlineUsersStore.setOnlineUsers(message.usersOnline);
        });

        this.setMessageHandler<UserStatusMessage>('USER_DISCONNECTED', (message: UserStatusMessage) => {
            this.onlineUsersStore.removeOnlineUser(message.user);
        });

        this.setMessageHandler<UserStatusMessage>('NEW_CONNECTION', (message: UserStatusMessage) => {
            this.onlineUsersStore.addOnlineUser(message.user);
        });
        this.setMessageHandler<UserStatusMessage>('LOBBY_INVITATION_FROM_FRIEND', (message) => {
            eventBus.emit('lobby-invitation', () => {

            })
        });
    }

    public setMessageHandler<T>(type: string, handler: MessageHandler<T>): void {
        this.messageHandlers[type] = handler as MessageHandler<ChatMessage | OnlineUsersMessage | UserStatusMessage>;
    }

    public connect(): void {
        try {
            const url = `ws://localhost:4001/ws?id=${this.clientId}`
            console.log(url);
            this.ws = new WebSocket(url);
            this.ws.onopen = () => {
                console.log('Websocket connected!');
                console.log('WS ready state: ', this.ws?.readyState);
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
                    const handler = this.messageHandlers[message.type];
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

    public sendChatMessage(content: string, senderID: number, receiverID: number): void {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            const message: ChatMessage = {
                type: 'CHAT',
                data: content,
                senderID: senderID,
                receiverID: receiverID
            };
            this.ws.send(JSON.stringify(message));
        } else {
            console.warn("Can't send a message, ws is not connected");
        }
    }

    public inviteFriendToLobbyMessage(friendId: number): void {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            const message: LobbyInvitationToFriend  = {
                Type: 'LOBBY_INVITATION_TO_FRIEND',
                userID: this.userStore.getId,
                senderID: this.userStore.getId,
                receiverID: friendId
            };
            console.log("MSG SENT");
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
