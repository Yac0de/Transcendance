import { GameState } from '../types/game'
import { fetchUserById } from '../utils/fetch'
import { UserData } from '../types/models'
import i18n from '../../i18n.ts';

let animationTime = 0;

function getCSSVariable(variableName: string): string {
    return getComputedStyle(document.documentElement).getPropertyValue(variableName).trim();
}

function hexToRgba(hex: string, alpha: number): string {
    // Supprimer le '#' si présent
    hex = hex.replace('#', '');

    // Extraire les composantes RGB
    const r = parseInt(hex.substring(0, 2), 16);
    const g = parseInt(hex.substring(2, 4), 16);
    const b = parseInt(hex.substring(4, 6), 16);

    // Retourner au format RGBA
    return `rgba(${r}, ${g}, ${b}, ${alpha})`;
}


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


function drawFireBall(ctx: CanvasRenderingContext2D, state: GameState) {
    // Récupérer la couleur hexadécimale de la variable CSS
    const colorExtract = getCSSVariable('--secondary-bright-color');

    ctx.globalCompositeOperation = 'lighter';

    for (let i = 0; i < 10; i++) {
        const radius = Math.random() * 15 + 5;
        const angle = Math.random() * Math.PI * 2;
        const offsetX = Math.cos(angle) * (Math.random() * 10);
        const offsetY = Math.sin(angle) * (Math.random() * 10);

        const gradient = ctx.createRadialGradient(
            state.ball.x + offsetX,
            state.ball.y + offsetY,
            0,
            state.ball.x + offsetX,
            state.ball.y + offsetY,
            radius
        );

        // Ajouter les étapes du dégradé avec transparence
        gradient.addColorStop(0, hexToRgba(colorExtract, 0.8)); // Opacité à 0.8
        gradient.addColorStop(0.4, hexToRgba(colorExtract, 0.4)); // Opacité à 0.4
        gradient.addColorStop(1, hexToRgba(colorExtract, 0)); // Opacité à 0

        ctx.beginPath();
        ctx.fillStyle = gradient;
        ctx.arc(state.ball.x + offsetX, state.ball.y + offsetY, radius, 0, Math.PI * 2);
        ctx.fill();
    }
}


export function drawBall(ctx: CanvasRenderingContext2D, state: GameState) {
    // Dessiner l'effet de feu
    drawFireBall(ctx, state);
    
    // Restaurer le mode de composition normal
    ctx.globalCompositeOperation = 'source-over';
    
    // Dessiner la balle principale (toujours en rouge maintenant)
    ctx.beginPath();
    ctx.fillStyle = 'rgba(255, 255, 255 , 0.7)';
    ctx.arc(state.ball.x, state.ball.y, 10, 0, Math.PI * 2);
    ctx.fill();
    
    // Dessiner les indicateurs de boost
    ctx.globalCompositeOperation = 'source-over';
}


export function drawBoostStatus(ctx: CanvasRenderingContext2D, state: GameState) {
    const statusHeight = 30;
    const margin = 20;
    const width = 100;
    const borderRadius = 15;
    
    const y = ctx.canvas.height - statusHeight - margin;
    
    function roundRect(x: number, y: number, w: number, h: number) {
        ctx.beginPath();
        ctx.moveTo(x + borderRadius, y);
        ctx.lineTo(x + w - borderRadius, y);
        ctx.quadraticCurveTo(x + w, y, x + w, y + borderRadius);
        ctx.lineTo(x + w, y + h - borderRadius);
        ctx.quadraticCurveTo(x + w, y + h, x + w - borderRadius, y + h);
        ctx.lineTo(x + borderRadius, y + h);
        ctx.quadraticCurveTo(x, y + h, x, y + h - borderRadius);
        ctx.lineTo(x, y + borderRadius);
        ctx.quadraticCurveTo(x, y, x + borderRadius, y);
        ctx.closePath();
    }
    
    // Fonction pour déterminer la couleur et le texte selon l'état
    function getBoostInfo(boost: any) {
        if (boost.isboostactive) {
            return { 
                color: '#ff0000', 
                text: 'ACTIVE'
            };  
        } else if (boost.boostReady) {
            return { 
                color: getCSSVariable('--secondary-bright-color'), 
                text: 'READY'
            };
        } else {
            return { 
                color: '#333333', 
                text: 'LOCKED'
            };
        }
    }

    // Player 1 boost
    const player1Info = getBoostInfo(state.player1boost);
    ctx.fillStyle = player1Info.color;
    roundRect(margin, y, width, statusHeight);
    ctx.fill();
    ctx.strokeStyle = '#ffffff';
    ctx.stroke();
    
    // Player 2 boost
    const player2Info = getBoostInfo(state.player2boost);
    ctx.fillStyle = player2Info.color;
    roundRect(ctx.canvas.width - width - margin, y, width, statusHeight);
    ctx.fill();
    ctx.stroke();
    
    // Texte
    ctx.fillStyle = '#ffffff';
    ctx.font = '14px Arial';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    
    ctx.fillText(player1Info.text, margin + width/2, y + statusHeight/2);
    ctx.fillText(player2Info.text, ctx.canvas.width - width/2 - margin, y + statusHeight/2);
}
 

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
    const winnerText = i18n.global.t('winnerText', { displayname: winner?.displayname });

    ctx.fillText(winnerText, 0, -40);

    // Score final
    const finalScoreText = i18n.global.t('finalScoreText', { score1: state.score.player1, score2: state.score.player2 });

    ctx.font = '32px Arial';
    ctx.fillStyle = '#FFFFFF';
    ctx.fillText(finalScoreText, 0, 20);

    // Message de redirection
    const alpha = (Math.sin(animationTime * 4) + 1) / 2;
    const backToMenuText = i18n.global.t('backToMenuText');

    ctx.fillStyle = `rgba(255, 255, 255, ${alpha})`;
    ctx.font = '24px Arial';
    ctx.fillText(backToMenuText, 0, 80);

    ctx.restore();
}

export async function drawSemiEndGame(
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
