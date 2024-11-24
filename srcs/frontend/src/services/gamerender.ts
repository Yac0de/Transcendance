import { GameState } from '../types/game'
import { fetchUserById } from '../utils/fetch'
import { UserData } from '../types/models'

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
    drawBoostIndicator(ctx, state.boostReady);
}

let time = 0;

function drawBoostIndicator(ctx: CanvasRenderingContext2D, isReady: boolean) {
    const indicatorWidth = 100;
    const indicatorHeight = 20;
    const padding = 10;
    
    // Position en bas au centre
    const x = (ctx.canvas.width - indicatorWidth) / 2;
    const y = ctx.canvas.height - indicatorHeight - padding;

    if (isReady) {
        time += 0.05;
        drawFlame(ctx, x, y, indicatorWidth, indicatorHeight);
    }

    // Fond de l'indicateur
    ctx.fillStyle = isReady ? '#800000' : '#333333'; // Rouge foncé quand prêt
    ctx.fillRect(x, y, indicatorWidth, indicatorHeight);

    // Texte "BOOST"
    ctx.fillStyle = isReady ? '#ff0000' : '#666666'; // Rouge vif quand prêt
    ctx.font = '16px Arial';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    ctx.fillText('BOOST', x + indicatorWidth/2, y + indicatorHeight/2);

    // Bordure
    ctx.strokeStyle = isReady ? '#ff0000' : '#666666';
    ctx.lineWidth = 2;
    ctx.strokeRect(x, y, indicatorWidth, indicatorHeight);
}

function drawFlame(ctx: CanvasRenderingContext2D, x: number, y: number, width: number, height: number) {
    const flameHeight = 40;
    
    ctx.save();
    ctx.beginPath();
    
    // Point de départ à gauche
    ctx.moveTo(x, y + height);

    // Courbe gauche de la flamme
    ctx.bezierCurveTo(
        x - 10 + Math.sin(time * 2) * 5,
        y + height/2,
        x - 20 + Math.sin(time * 3) * 10,
        y - flameHeight/2,
        x + width/2,
        y - flameHeight + Math.sin(time) * 10
    );

    // Courbe droite de la flamme
    ctx.bezierCurveTo(
        x + width + 20 + Math.sin(time * 3) * 10,
        y - flameHeight/2,
        x + width + 10 + Math.sin(time * 2) * 5,
        y + height/2,
        x + width,
        y + height
    );

    // Gradient pour l'effet de flamme réaliste
    const gradient = ctx.createLinearGradient(x + width/2, y + height, x + width/2, y - flameHeight);
    gradient.addColorStop(0, 'rgba(255, 60, 0, 0.8)');  // Orange-rouge à la base
    gradient.addColorStop(0.3, 'rgba(255, 150, 0, 0.7)'); // Orange
    gradient.addColorStop(0.6, 'rgba(255, 220, 0, 0.5)'); // Jaune
    gradient.addColorStop(1, 'rgba(255, 255, 100, 0)');   // Jaune clair qui s'estompe
    
    ctx.fillStyle = gradient;
    ctx.fill();
    ctx.restore();
}

let animationTime = 0;

export async function drawEndGame(
    ctx: CanvasRenderingContext2D, 
    state: GameState, 
    player1id: number | null, 
    player2id: number | null,
) {
    animationTime += 0.02;

    // Fond semi-transparent
    ctx.fillStyle = 'rgba(0, 0, 0, 0.85)';
    ctx.fillRect(0, 0, ctx.canvas.width, ctx.canvas.height);

    // Animation de texte
    const scale = 1 + Math.sin(animationTime * 2) * 0.1;

    ctx.save();
    ctx.translate(ctx.canvas.width / 2, ctx.canvas.height / 2);
    ctx.scale(scale, scale);

    // Texte principal
    ctx.font = 'bold 48px Arial';
    ctx.fillStyle = '#FFD700';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';

    const winnerId = state.winner === player1id ? player1id : player2id ?? 0;
    const winner: UserData | null = await fetchUserById(winnerId);
    ctx.fillText(`${winner?.displayname} GAGNE!`, 0, -40);

    // Score final
    ctx.font = '32px Arial';
    ctx.fillStyle = '#FFFFFF';
    ctx.fillText(`Score Final: ${state.score.player1} - ${state.score.player2}`, 0, 20);

    // Message de redirection
    const alpha = (Math.sin(animationTime * 4) + 1) / 2;
    ctx.fillStyle = `rgba(255, 255, 255, ${alpha})`;
    ctx.font = '24px Arial';
    ctx.fillText("Retour au menu dans quelques secondes...", 0, 80);

    ctx.restore();
}