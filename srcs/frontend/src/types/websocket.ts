export interface BaseMessage {
  type: string;
}

export interface ChatMessage extends BaseMessage {
  type: 'CHAT';
  data: string;
  senderID: number;
  receiverID: number;
}

export interface LobbyInvitationToFriend extends BaseMessage {
  type: 'LOBBY_INVITATION_TO_FRIEND';
  userID: number;
  senderID: number;
  receiverID: number;
}

export interface LobbyInvitationFromFriend extends BaseMessage {
  type: 'LOBBY_INVITATION_FROM_FRIEND';
  userID: number;
  senderID: number;
  receiverID: number;
}

export interface LobbyAcceptFromFriend extends BaseMessage {
  type: 'LOBBY_ACCEPT_FROM_FRIEND';
  userID: number;
  senderID: number;
  receiverID: number;
  lobbyID: number;
}

export interface LobbyDenyFromFriend extends BaseMessage {
  type: 'LOBBY_DENY_FROM_FRIEND';
  data: string;
  senderID: number;
  receiverID: number;
  lobbyID: number;
}

export interface LobbyTerminate extends BaseMessage {
  type: 'LOBBY_TERMINATE';
  data: string;
  senderID: number;
  receiverID: number;
}

export interface OnlineUsersMessage extends BaseMessage {
  type: 'ONLINE_USERS';
  usersOnline: number[];
}

export interface UserStatusMessage extends BaseMessage {
  type: 'USER_DISCONNECTED | NEW_CONNECTION';
  user: number;
}
