export interface LobbyUserState {
  id: number;
  isReady: boolean;
}

export interface LobbyInvitationToFriend {
  type: 'LOBBY_INVITATION_TO_FRIEND';
  userID: number;
  sender: LobbyUserState;
  receiver: LobbyUserState;
  lobbyID: string;
}

export interface LobbyInvitationFromFriend {
  type: 'LOBBY_INVITATION_FROM_FRIEND';
  userID: number;
  sender: LobbyUserState;
  receiver: LobbyUserState;
  lobbyID: string;
}

export interface LobbyAcceptFromFriend {
  type: 'LOBBY_ACCEPT_FROM_FRIEND';
  userID: number;
  sender: LobbyUserState;
  receiver: LobbyUserState;
  lobbyID: string;
}

export interface LobbyDenyFromFriend {
  type: 'LOBBY_DENY_FROM_FRIEND';
  userID: number;
  sender: LobbyUserState;
  receiver: LobbyUserState;
  lobbyID: string;
}

export interface LobbyPlayerStatus {
  type: 'LOBBY_PLAYER_READY_STATUS' | 'LOBBY_PLAYER_UNREADY_STATUS';
  userID: number;
  lobbyID: string;
}

export interface LobbyCreated {
  type: 'LOBBY_CREATED';
  sender: LobbyUserState;
  receiver: LobbyUserState;
  lobbyID: number;
}

export interface LobbyTerminate {
  type: 'LOBBY_TERMINATE';
  sender: LobbyUserState;
  lobbyID: number;
}
