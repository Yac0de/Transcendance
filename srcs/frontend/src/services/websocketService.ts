import { OnlineUsersMessage, UserStatusMessage } from '../types/connection_status';
import { ChatMessage } from '../types/chat';
import { UserData } from '../types/models';
import { LobbyInvitationToFriend, LobbyInvitationFromFriend, LobbyAcceptFromFriend, LobbyDenyFromFriend, LobbyCreated, LobbyPlayerStatus, LobbyPregameRemainingTime, LobbyTerminate, LobbyDestroyed } from '../types/lobby';
import { GameEvent } from '../types/game';
import { useOnlineUsersStore } from '../stores/onlineUsers';
import { eventBus } from '../events/eventBus';
import { useChatStore } from '../stores/chatStore.ts';

const WS_URL = import.meta.env.PROD
  ? 'wss://localhost:8443/ws'  // Production through Nginx
  : 'ws://localhost:4001/ws'      // Development direct access

//This is needed because we can't get the return type of userStore inside the constructor of a class that is an attribute of
//this very store, because it creates circular dependencies, so we create an interface that helps up set the return type of
//our userStore
interface IUserStore {
    id: number | null;
    nickname: string | null;
    displayname: string | null;
    avatar: string | null;
    webSocketService: WebSocketService | null;
    isRedirectPending: boolean,

    getId: number | null;
    getNickname: string | null;
    getDisplayName: string | null;
    getAvatarPath: string | null;
    isSignedIn: boolean;
    getWebSocketService: WebSocketService | null;

    setUser: (userData: UserData) => void;
    setWebSocketService: (userId: number) => void;
    fetchUser: () => Promise<void>;
    clearUser: () => void;
    loadUserFromStorage: () => boolean;
    initializeStore: () => Promise<boolean>;
}

type MessageHandler<T> = (message: T) => void;
type MessageHandlers = {
    [key: string]: MessageHandler<any>;
};

export class WebSocketService {
    public ws: WebSocket | null = null;
    public clientId: number;
    public onlineUsersStore: ReturnType<typeof useOnlineUsersStore>;
    public userStore: IUserStore;
    public messageHandlers: MessageHandlers = {};

    constructor(clientId: number, onlineUsersStore: ReturnType<typeof useOnlineUsersStore>, userStore: IUserStore) {
        this.clientId = clientId;
        this.onlineUsersStore = onlineUsersStore;
        this.userStore = userStore;
        this.initMessageHandlers();
    }

    public initMessageHandlers(): void {
        this.setMessageHandler<ChatMessage>('CHAT', (message: ChatMessage) => {
            const conversationId = message.senderID === this.clientId
                ? message.receiverID
                : message.senderID;

            if (!conversationId) {
                console.warn('Unable to determine conversation ID for message:', message);
                return;
            }

            const chatStore = useChatStore();
            if (chatStore.selectedFriendId !== conversationId) {
                console.log(`Adding unread message for conversation ID: ${conversationId}`);
                chatStore.addUnreadMessage(conversationId);
            }
        });

        this.setMessageHandler<OnlineUsersMessage>('ONLINE_USERS', (message: OnlineUsersMessage) => {
            this.onlineUsersStore.setOnlineUsers(message.usersOnline);
        });

        this.setMessageHandler<UserStatusMessage>('USER_DISCONNECTED', (message: UserStatusMessage) => {
            this.onlineUsersStore.removeOnlineUser(message.user);
        });

        this.setMessageHandler<UserStatusMessage>('NEW_CONNECTION', (message: UserStatusMessage) => {
            this.onlineUsersStore.addOnlineUser(message.user);
        });
        this.setMessageHandler<LobbyInvitationFromFriend>('LOBBY_INVITATION_FROM_FRIEND', (message: LobbyInvitationFromFriend) => {
            eventBus.emit('LOBBY_INVITATION_FROM_FRIEND', message);
        });
        this.setMessageHandler<LobbyInvitationToFriend>('LOBBY_INVITATION_TO_FRIEND', (message: LobbyInvitationToFriend) => {
            eventBus.emit('LOBBY_INVITATION_TO_FRIEND', message);
        });
        this.setMessageHandler<LobbyCreated>('LOBBY_CREATED', (message: LobbyCreated) => {
            eventBus.emit('LOBBY_CREATED', message);
        });
        this.setMessageHandler<LobbyPlayerStatus>('LOBBY_PLAYER_STATUS', (message: LobbyPlayerStatus) => {
            eventBus.emit('LOBBY_PLAYER_STATUS', message);
        });
        this.setMessageHandler<LobbyPregameRemainingTime>('LOBBY_PREGAME_REMAINING_TIME', (message: LobbyPregameRemainingTime) => {
            eventBus.emit('LOBBY_PREGAME_REMAINING_TIME', message);
        });
        this.setMessageHandler<LobbyDestroyed>('LOBBY_DESTROYED', () => {
            eventBus.emit('LOBBY_DESTROYED');
        });
        this.setMessageHandler<GameEvent>('GAME_EVENT',(message: GameEvent)  => {
            eventBus.emit('GAME_EVENT', message);
        });
        this.setMessageHandler<TournamentJoin>('TOURNAMENT_JOIN_WITH_CODE', (message: TournamentJoin) => {
            eventBus.emit('TOURNAMENT_JOIN_WITH_CODE', message);
        });
        this.setMessageHandler<TournamentCreate>('TOURNAMENT_CREATE', (message: TournamentCreate) => {
            eventBus.emit('TOURNAMENT_CREATE', message);
        })
        this.setMessageHandler<TournamentEvent>('TOURNAMENT_EVENT', (message: TournamentEvent) => {
            eventBus.emit('TOURNAMENT_EVENT', message);
        })
    }

    public setMessageHandler<T>(type: string, handler: MessageHandler<T>): void {
        this.messageHandlers[type] = handler;
    }

    public connect(): void {
        try {
            const url = WS_URL + `?id=${this.clientId}`
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
                    console.log(event.data)
                    const message = JSON.parse(event.data);
                    console.log(message)
                    const handler = this.messageHandlers[message.type];
                    if (handler) {
                        handler(message);
                    } else
                        console.warn(`No handler found for message type: ${message.type}`);
                } catch (e) {
                    console.error('Error parsing WebSocket message:', e);
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
            const message: LobbyInvitationToFriend = {
                type: 'LOBBY_INVITATION_TO_FRIEND',
                userId: this.userStore.getId!,
                sender: {
                    id: this.userStore.getId!,
                    isReady: false
                },
                receiver: {
                    id: friendId,
                    isReady: false
                },
            };
            this.ws.send(JSON.stringify(message));
        } else {
            console.warn("Can't send a message, ws is not connected");
        }
    }

    public acceptInviteFromFriend(lobbyId: string, inviterId: number): void {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            const message: LobbyAcceptFromFriend = {
                type: 'LOBBY_ACCEPT_FROM_FRIEND',
                user: this.userStore.getId!,
                sender: {
                    id: inviterId,
                    isReady: false
                },
                receiver: {
                    id: this.userStore.getId!,
                    isReady: false
                },
                lobbyId: lobbyId
            };
            this.ws.send(JSON.stringify(message));
        } else {
            console.warn("Can't send a message, ws is not connected");
        }
    }

    public denyInviteFromFriend(lobbyId: string, inviterId: number): void {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            const message: LobbyDenyFromFriend = {
                type: 'LOBBY_DENY_FROM_FRIEND',
                user: this.userStore.getId!,
                sender: {
                    id: this.userStore.getId!,
                    isReady: false
                },
                receiver: {
                    id: inviterId,
                    isReady: false
                },
                lobbyId: lobbyId
            };
            this.ws.send(JSON.stringify(message));
        } else {
            console.warn("Can't send a message, ws is not connected");
        }
    }

    public sendPlayerReadyMessage(lobbyId: string): void {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            const message: LobbyPlayerStatus = {
                type: 'LOBBY_PLAYER_READY_STATUS',
                userId: this.userStore.getId!,
                lobbyId: lobbyId,
            };
            this.ws.send(JSON.stringify(message));
        } else {
            console.warn("Can't send a message, ws is not connected");
        }
    }

    public sendPlayerUnreadyMessage(lobbyId: string): void {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            const message: LobbyPlayerStatus = {
                type: 'LOBBY_PLAYER_UNREADY_STATUS',
                userId: this.userStore.getId!,
                lobbyId: lobbyId
            };
            this.ws.send(JSON.stringify(message));
        } else {
            console.warn("Can't send a message, ws is not connected");
        }
    }

    public leaveAndTerminateLobby(lobbyId: string): void {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            const message: LobbyTerminate = {
                type: 'LOBBY_TERMINATE',
                sender: {
                    id: this.userStore.getId!,
                    isReady: false
                },
                lobbyId: lobbyId
            };
            this.ws.send(JSON.stringify(message));
        } else {
            console.warn("Can't send a message, ws is not connected");
        }
    }

    public joinTournamentWithCode(code: string): void {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            const message: joinTournamentWithCode = {
                type: 'TOURNAMENT_JOIN_WITH_CODE',
                userId: this.userStore.getId!,
                code: code
            };
            console.log("JOIN TOURNAMENT WITH CODE -> ", message);
            this.ws.send(JSON.stringify(message));
        } else {
            console.warn("Can't send a message, ws is not connected");
        }
    }

    public createTournamentWaitingRoom(): void {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            const message: TournamentCreate = {
                type: 'TOURNAMENT_CREATE',
                userId: this.userStore.getId!,
                code:'' 
            };
            console.log("CREATE TOURNAMENT LOBBY -> ", message);
            this.ws.send(JSON.stringify(message));
        }
    }

    public leaveTournamentWaitingRoom(code: string): void {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            const message: LeaveTournamentWaitingRoom = {
                type: 'LEAVE_TOURNAMENT_WAITING_ROOM',
                userId: this.userStore.getId!,
                code: code
            };
            console.log("LEAVE TOURNEY WAITING ROOM -> ", message);
            this.ws.send(JSON.stringify(message));
        }
    }

    public sendTournamentStart(code: string): void {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            const message: TournamentStart = {
                type: 'TOURNAMENT_START',
                userId: this.userStore.getId!,
                code: code
            };
            console.log("TOURNAMENT START -> ", message);
            this.ws.send(JSON.stringify(message));
        }
    }

    public sendGameEvent(game_event: GameEvent): void {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify(game_event));
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
