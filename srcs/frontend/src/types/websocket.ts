export interface BaseMessage {
  Type: string;
}

export interface ChatMessage extends BaseMessage {
  Type: 'CHAT';
  Data: string;
  SenderID: number;
  ReceiverID: number;
}

export interface OnlineUsersMessage extends BaseMessage {
  Type: 'ONLINE_USERS';
  UsersOnline: number[];
}

export interface UserStatusMessage extends BaseMessage {
  Type: 'USER_DISCONNECTED | NEW_CONNECTION';
  User: number;
}
