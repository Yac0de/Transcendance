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
