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

export interface TournamentEvent {
  type: string;
  creatorId: number;
  code: number;
  player1id: number;
  player2id: number;
  player3id: number;
  player4id: number;
}
