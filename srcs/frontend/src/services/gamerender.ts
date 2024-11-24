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

const trailLength = 15; // Nombre de positions précédentes à conserver
let previousPositions: {x: number, y: number}[] = [];

export function drawBall(ctx: CanvasRenderingContext2D, state: GameState) {
    // Ajouter la position actuelle
    previousPositions.push({x: state.ball.x, y: state.ball.y});
    
    // Garder seulement les N dernières positions
    if (previousPositions.length > trailLength) {
        previousPositions.shift();
    }

    // Dessiner la traînée
    previousPositions.forEach((pos, index) => {
        const alpha = (index + 1) / trailLength;
        const radius = 10 * ((index + 1) / trailLength); // Le rayon diminue pour les anciennes positions
        ctx.beginPath();
        ctx.fillStyle = state.isBoostActive 
            ? `rgba(255, 0, 0, ${alpha * 0.3})`  // Rouge pour le boost
            : `rgba(255, 255, 255, ${alpha * 0.3})`; // Blanc normal
        ctx.arc(pos.x, pos.y, radius, 0, Math.PI * 2);
        ctx.fill();
    });

    // Dessiner la balle principale
    ctx.beginPath();
    ctx.fillStyle = state.isBoostActive ? 'red' : 'rgba(94, 84, 142)';
    ctx.arc(state.ball.x, state.ball.y, 10, 0, Math.PI * 2);
    ctx.fill();
}