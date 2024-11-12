export interface LobbyInvitationToFriend extends BaseMessage {
  type: 'LOBBY_INVITATION_TO_FRIEND';
  userID: number;
  senderID: number;
  receiverID: number;
  lobbyID: string;
}

export interface LobbyInvitationFromFriend extends BaseMessage {
  type: 'LOBBY_INVITATION_FROM_FRIEND';
  userID: number;
  senderID: number;
  receiverID: number;
  lobbyID: string;
}

export interface LobbyAcceptFromFriend extends BaseMessage {
  type: 'LOBBY_ACCEPT_FROM_FRIEND';
  userID: number;
  senderID: number;
  receiverID: number;
  lobbyID: string;
}

export interface LobbyDenyFromFriend extends BaseMessage {
  type: 'LOBBY_DENY_FROM_FRIEND';
  userID: number;
  senderID: number;
  receiverID: number;
  lobbyID: string;
}

export interface LobbyTerminate extends BaseMessage {
  type: 'LOBBY_TERMINATE';
  data: string;
  senderID: number;
  receiverID: number;
}
