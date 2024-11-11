export interface BaseMessage {
  type: string;
}

export interface ChatMessage extends BaseMessage {
  type: 'CHAT';
  data: string;
  senderID: number;
  receiverID: number;
}

export interface GameInvitationMessage extends BaseMessage {
  type: 'GAME_INVITATION_TO_FRIEND';
  friendId: number;
}

export interface OnlineUsersMessage extends BaseMessage {
  type: 'ONLINE_USERS';
  usersOnline: number[];
}

export interface UserStatusMessage extends BaseMessage {
  type: 'USER_DISCONNECTED | NEW_CONNECTION';
  user: number;
}
