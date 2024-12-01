export interface TournamentJoinWithCode {
  type: 'TOURNAMENT_JOIN_WITH_CODE';
  userId: number;
  code: string;
}

export interface TournamentTreeState {
  type: 'TOURNAMENT_TREE_STATE';
  userId: number;
  code: string;
  final?: { 
    player1id: number;
    player2id: number;
    score: number[];
    isFinished: boolean
  };
  semi1?: { 
    player1id: number;
    player2id: number;
    score: number[];
    isFinished: boolean;
  };
  semi2?: { 
    player1id: number;
    player2id: number;
    score: number[];
    isFinished: boolean;
  };
}

export interface TournamentCreate {
  type: 'TOURNAMENT_CREATE';
  userId: number;
  code: string;
  player1id?: number;
  player2id?: number;
  player3id?: number;
  player4id?: number;
}

export interface TournamentGame {
  type: 'TOURNAMENT_GAME';
  lobbyId: number;
}

export interface TournamentLeave {
  type: 'TOURNAMENT_LEAVE';
  userId: number;
}

export interface TournamentLeaveWaitingRoom {
  type: 'TOURNAMENT_LEAVE_WAITING_ROOM';
  userId: number;
}

export interface TournamentTerminate {
  type: 'TOURNAMENT_TERMINATE';
}

export interface TournamentError {
  type: 'TOURNAMENT_ERROR';
  error: string;
}

export interface TournamentTimer {
  type: 'TOURNAMENT_TIMER';
  remainingTime: number;
}

export interface TournamentStart {
  type: 'TOURNAMENT_START';
  userId: number;
  code: string;
  player1id?: number;
  player2id?: number;
  player3id?: number;
  player4id?: number;
}

export interface TournamentEvent {
  type: string;
  creatorId: number;
  code: string;
  player1id: number;
  player2id: number;
  player3id: number;
  player4id: number;
}
