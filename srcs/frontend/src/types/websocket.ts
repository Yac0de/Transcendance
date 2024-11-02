export interface BaseMessage {
  Type: string;
}

export interface ChatMessage extends BaseMessage {
  Type: 'CHAT';
  Data: string;
  SenderID: string;
  ReceiverID: string;
}

export interface OnlineUsersMessage extends BaseMessage {
  Type: 'ONLINE_USERS';
  UsersOnline: string[];
}

export interface UserStatusMessage extends BaseMessage {
  Type: 'USER_DISCONNECTED | NEW_CONNECTION';
  User: string;
}
