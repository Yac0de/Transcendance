export interface GameEvent {
    type: 'GAME_EVENT';
    lobbyId: string;
    userId: number;
    state?: GameState;
    keyPressed: string;
    player1id?: number;
    player2id?: number;
}

export interface GameStart {
    type: 'GAME_START';
}

export interface GameFinished {
    type: 'GAME_FINISHED';
    state?: GameState;
    isTournamentGame: boolean;
}

export interface Ball {
    x: number;
    y: number;
}

export interface Paddle {
    width: number;
    height: number;
    player1X: number;
    player1Y: number;
    player2X: number;
    player2Y: number;
    player1Direction: number;
    player2Direction: number;
}

export interface Score {
    player1: number;
    player2: number;
}

export interface GameState {
    ball: Ball;
    paddle?: Paddle;
    score: Score;
    isGameMode: boolean;
    winner: number;
    isPaused: boolean;
    pauseTime: string; 
    player1boost: BoostState;  
    player2boost: BoostState;
    remainingTime: number;
}

export interface BoostState {
    ballhit: number;
    boostready:boolean;
    isboostactive:boolean;
}

export interface Player {
    id: number;
    position: number;
}

export interface Game {
    player1: Player;
    player2: Player;
    state: GameState;
    status: string;
}
