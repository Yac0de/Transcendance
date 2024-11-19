import { GameState } from '../types/game'

export function drawPaddle(ctx: HTMLCanvasElement, state: GameState): void {
    ctx.fillStyle = 'white';
    console.log("paddle", state.paddle);
    ctx.fillRect(state.paddle.Player1X, state.paddle.Player1Y, state.paddle.width, state.paddle.height);
    ctx.fillRect(state.paddle.Player2X, state.paddle.Player2Y, state.paddle.width, state.paddle.height);
}

export function drawBall(ctx: HTMLCanvasElement, state: GameState) : void {
    ctx.beginPath();
    ctx.arc(state.ball.x, state.ball.y, 10, 0, Math.PI * 2);
    ctx.fillStyle = 'green';
    ctx.fill();
    ctx.closePath();
}
