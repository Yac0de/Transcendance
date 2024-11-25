export interface TournamentJoinWithCode {
  type: 'TOURNAMENT_JOIN_WITH_CODE';
  userId: number;
  code: string;
}

export interface TournamentCreate {
  type: 'TOURNAMENT_CREATE';
  userId: number;
  code: number;
}

export interface TournamentGame {
  type: 'TOURNAMENT_GAME';
  lobbyId: number;
}

export interface TournamentLeave {
  type: 'TOURNAMENT_LEAVE_WAITING_ROOM';
  userId: number;
  code: number;
}

export interface TournamentTimer {
  type: 'TOURNAMENT_TIMER';
  remainingTime: number;
}

export interface TournamentStart {
  type: 'TOURNAMENT_START';
  userId: number;
  code: number;
  player1id: number;
  player2id: number;
  player3id: number;
  player4id: number;
}

export interface TournamentEvent {
  type: string;
  creatorId: number;
  code: number;
  player1id: number;
  player2id: number;
  player3id: number;
  player4id: number;
}
