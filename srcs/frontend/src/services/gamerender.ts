import { GameState } from '../types/game'


export function drawPaddle(ctx: CanvasRenderingContext2D, state: GameState): void {
    ctx.fillStyle = 'white';
    ctx.fillRect(state.paddle.player1X,state.paddle.player1Y , state.paddle.width, state.paddle.height);
    ctx.fillRect(state.paddle.player2X,state.paddle.player2Y , state.paddle.width, state.paddle.height);
}

export function drawBall(ctx: CanvasRenderingContext2D, state: GameState) : void {
    ctx.beginPath();
    ctx.arc(state.ball.x, state.ball.y, 10, 0, Math.PI * 2);
    ctx.fillStyle = 'green';
    ctx.fill();
    ctx.closePath();
}
