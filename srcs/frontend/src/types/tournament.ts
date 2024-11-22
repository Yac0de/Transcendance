export interface joinTournamentWithCode {
  type: 'JOIN_TOURNAMENT_WITH_CODE';
  userId: number;
  code: string;
}

export interface createTournamentLobby {
  type: 'CREATE_TOURNAMENT_LOBBY';
  userId: number;
  code: number;
}

export interface TournamentLobbyState {
  type: 'TOURNAMENT_LOBBY_STATE';
  creatorId: number;
  code: number;
  player1Id: number;
  player2Id: number;
  player3Id: number;
  player4Id: number;
}
