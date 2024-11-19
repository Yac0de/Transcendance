export interface LobbyUserState {
  id: number;
  isReady: boolean;
}

export interface LobbyInvitationToFriend {
  type: 'LOBBY_INVITATION_TO_FRIEND';
  userId: number;
  sender: LobbyUserState;
  receiver: LobbyUserState;
}

export interface LobbyInvitationFromFriend {
  type: 'LOBBY_INVITATION_FROM_FRIEND';
  userId: number;
  sender: LobbyUserState;
  receiver: LobbyUserState;
  lobbyId: string;
}

export interface LobbyAcceptFromFriend {
  type: 'LOBBY_ACCEPT_FROM_FRIEND';
  user: number;
  sender: LobbyUserState;
  receiver: LobbyUserState;
  lobbyId: string;
}

export interface LobbyDenyFromFriend {
  type: 'LOBBY_DENY_FROM_FRIEND';
  user: number;
  sender: LobbyUserState;
  receiver: LobbyUserState;
  lobbyId: string;
}

export interface LobbyPlayerStatus {
  type: 'LOBBY_PLAYER_READY_STATUS' | 'LOBBY_PLAYER_UNREADY_STATUS';
  userId: number;
  lobbyId: string;
}

export interface LobbyCreated {
  type: 'LOBBY_CREATED';
  sender: LobbyUserState;
  receiver: LobbyUserState;
  lobbyId: string;
}

export interface LobbyPregameRemainingTime {
  type: 'LOBBY_PREGAME_REMAINING_TIME';
  remainingSecondsToStart: number;
  lobbyId: string;
}

export interface LobbyTerminate {
  type: 'LOBBY_TERMINATE';
  sender: LobbyUserState;
  lobbyId: string;
}

export interface LobbyDestroyed {
  type: 'LOBBY_DESTROYED';
  lobbyId: string;
}
