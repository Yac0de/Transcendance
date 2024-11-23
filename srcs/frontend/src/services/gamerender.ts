import { GameState } from '../types/game'

export function drawPaddle(ctx: CanvasRenderingContext2D, state: GameState): void {
    ctx.fillStyle = 'white';
    ctx.fillRect(
        state.paddle?.player1X ?? 0,
        state.paddle?.player1Y ?? 0,
        state.paddle?.width ?? 0,
        state.paddle?.height ?? 0,
    );
    ctx.fillRect(
        state.paddle?.player2X ?? 0,
        state.paddle?.player2Y ?? 0,
        state.paddle?.width ?? 0,
        state.paddle?.height ?? 0,
    );
}

export function drawBall(ctx: CanvasRenderingContext2D, state: GameState) : void {
    ctx.beginPath();
    ctx.arc(state.ball.x, state.ball.y, 10, 0, Math.PI * 2);
    ctx.fillStyle = 'green';
    ctx.fill();
    ctx.closePath();
}
