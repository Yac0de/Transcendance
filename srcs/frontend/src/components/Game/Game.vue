<template>
  <div class="container">
    <div class="game-container">
      <GameHeader :player1id="player1Id" :player2id="player2Id":state="currentGameState"/>
      <div class="canvas-wrapper">
        <canvas
        id="gameCanvas"
        width="800"
        height="600"
        ref="canvasRef"
        class="game-canvas"
        >
        </canvas>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { GameEvent, GameState, GameFinished } from '../../types/game'
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { eventBus } from '../../events/eventBus';
import { drawPaddle, drawBall, drawEndGame, drawBoostStatus } from '../../services/gamerender';
import { useUserStore } from '../../stores/user';
import { useGameSettingsStore } from '../../stores/gameSettings.js';
import { useRoute, useRouter } from 'vue-router'
import GameHeader from './GameHeader.vue';

const route = useRoute()
const router = useRouter()

const userStore = useUserStore();
const gameSettingsStore = useGameSettingsStore();
const canvasRef = ref<HTMLCanvasElement | null>(null)

const player1Id = ref<number | null>(null)
const player2Id = ref<number | null>(null)
let isTournamentGame: boolean =  false;

const currentGameState: GameState = reactive({
    ball: { x: 0, y: 0 },  // Assuming Ball has x, y properties
    score: { 
        player1: 0, 
        player2: 0 
    },
    isGameMode: false,
    winner: 0,  // or 0, depending on how you represent no winner
    isPaused: false,
    pauseTime: '',  // or null, depending on how you handle empty time
    remainingTime: 300,
    player1boost: {
      ballhit: 0,
      boostready:false,
      isboostactive:false,
    },
    player2boost: {
      ballhit: 0,
      boostready:false,
      isboostactive:false,
    }
})

const handlePressUp = (event: KeyboardEvent): void => {
  if (event.code === 'ArrowUp' || event.code === 'KeyW') {
    event.preventDefault();
    if (userStore.getWebSocketService?.isConnected()) {
      const gameEvent: GameEvent = {
        type: 'GAME_EVENT',
        lobbyId: route.query.lobbyId as string,
        userId: userStore.getId!,
        keyPressed: 'UP',
      };
      userStore.getWebSocketService?.sendGameEvent(gameEvent);
    } else {
      console.error('WebSocket is not connected');
    }
  }
}

const handleReleaseUp = (event: KeyboardEvent): void => {
  if (event.code === 'ArrowUp' || event.code === 'KeyW') {
    event.preventDefault();
    if (userStore.getWebSocketService?.isConnected()) {
      const gameEvent: GameEvent = {
        type: 'GAME_EVENT',
        lobbyId: route.query.lobbyId as string,
        userId: userStore.getId!,
        keyPressed: 'STOP'
      };
      userStore.getWebSocketService?.sendGameEvent(gameEvent);
    } else {
      console.error('WebSocket is not connected');
    }
  }
}

const handlePressDown = (event: KeyboardEvent): void => {
  if (event.code === 'ArrowDown' || event.code === 'KeyS') {
    event.preventDefault();
    if (userStore.getWebSocketService?.isConnected()) {
      const gameEvent: GameEvent = {
        type: 'GAME_EVENT',
        lobbyId: route.query.lobbyId as string,
        userId: userStore.getId!,
        keyPressed: 'DOWN'
      };
      userStore.getWebSocketService?.sendGameEvent(gameEvent);
    } else {
      console.error('WebSocket is not connected');
    }
  }
}

const handleReleaseDown = (event: KeyboardEvent): void => {
  if (event.code === 'ArrowDown' || event.code === 'KeyS') {
    event.preventDefault();
    if (userStore.getWebSocketService?.isConnected()) {
      const gameEvent: GameEvent = {
        type: 'GAME_EVENT',
        lobbyId: route.query.lobbyId as string,
        userId: userStore.getId!,
        keyPressed: 'STOP'
      };
      userStore.getWebSocketService?.sendGameEvent(gameEvent);
    } else {
      console.error('WebSocket is not connected');
    }
  }
}

const handleSpace = (event: KeyboardEvent): void => {
 if (event.code === 'Space') {
   event.preventDefault();
   if (userStore.getWebSocketService?.isConnected()) {
     const gameEvent: GameEvent = {
       type: 'GAME_EVENT', 
       lobbyId: route.query.lobbyId as string,
       userId: userStore.getId!,
       keyPressed: 'SPACE'
     };
     userStore.getWebSocketService?.sendGameEvent(gameEvent);
   } else {
     console.error('WebSocket is not connected');
   }
 }
}

function drawVerticalLine(ctx: CanvasRenderingContext2D): void {
    ctx.strokeStyle = '#5b3c54';
    ctx.lineWidth = 1;
    ctx.beginPath();
    ctx.moveTo(ctx.canvas.width / 2, 0);
    ctx.lineTo(ctx.canvas.width / 2, ctx.canvas.height);
    ctx.stroke();
}


onMounted(() => {
  // Add key listener
  window.addEventListener('keydown', handlePressUp)
  window.addEventListener('keydown', handlePressDown)
  window.addEventListener('keyup', handleReleaseUp)
  window.addEventListener('keyup', handleReleaseDown)
  if(gameSettingsStore.gameMode)
    window.addEventListener('keydown', handleSpace)
  const ctx:CanvasRenderingContext2D = canvasRef.value?.getContext('2d') as CanvasRenderingContext2D

  eventBus.on('GAME_EVENT', async (message: GameEvent) => {

    if (message.isTournamentGame === true) {
      isTournamentGame = true;
    }

    if (player1Id.value === null) {
      player1Id.value = message.player1id ?? null;
    }
    if (player2Id.value === null) {
      player2Id.value = message.player2id ?? null;
    }

    if (canvasRef.value) {
      if (ctx) {
        Object.assign(currentGameState, message.state);

        ctx.clearRect(0, 0, canvasRef.value.width, canvasRef.value.height);

        ctx.fillStyle = 'black'
        ctx.fillRect(0, 0, canvasRef.value.width, canvasRef.value.height);

        drawVerticalLine(ctx);

        drawPaddle(ctx, message.state!);
        drawBall(ctx, message.state!);
        if(gameSettingsStore.gameMode)
          drawBoostStatus(ctx, message.state!);
        if (!message.state!.isGameMode && message.state!.winner !== 0)
          drawEndGame(ctx, message.state!, player1Id.value, player2Id.value);
      }
    }
  })

  eventBus.on('GAME_FINISHED', async (message: GameFinished) => {
    drawEndGame(ctx, message.state!, player1Id.value, player2Id.value);

    gameSettingsStore.gameMode = false;
    if (isTournamentGame === false) {
      window.setTimeout(() => {
        router.push('/');
      }, 3000)
    } else {
      window.setTimeout(() => {
        router.push({ 
          path: '/tournament', 
          query: { view: 'tournament-tree'}
        });
      }, 3000);
    }
  })
})

onUnmounted(() => {
  if (isTournamentGame) {
    if (userStore.getWebSocketService?.isConnected()) {
      console.log("-> TOURNAMENT_LEAVE (IN A GAME)");
      userStore.getWebSocketService?.sendLeaveTournament()
    } else {
      console.error('WebSocket is not connected');
    }
  } else {
    if (userStore.getWebSocketService?.isConnected()) {
      console.log("-> GAME LEAVE (IN A GAME)");
      userStore.getWebSocketService?.sendGameLeave()
    } else {
      console.error('WebSocket is not connected');
    }
  }

  window.removeEventListener('keydown', handlePressUp)
  window.removeEventListener('keydown', handlePressDown)
  window.removeEventListener('keyup', handleReleaseUp)
  window.removeEventListener('keyup', handleReleaseDown)
  if(gameSettingsStore.gameMode)
    window.removeEventListener('keydown', handleSpace)
  eventBus.off('GAME_EVENT')
  eventBus.off('GAME_FINISHED')
})
</script>

<style scoped>
body {
  margin: 0;
  overflow: hidden;
}
.container {
  width: 100%;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.game-container {
  background: linear-gradient(
    to right,
    var(--secondary-bright-color),
    color-mix(in srgb, var(--secondary-bright-color) 75%, white)
  );
  border: none;
  border-radius: 15px;
  overflow: hidden;

  /* Empilement des box-shadow */
  box-shadow: 
    0 0 0 6px var(--secondary-bright-color), /* Bordure simulée */
    0 0 30px rgba(0, 0, 0, 0.85);          /* Ombre réelle */
}



.canvas-wrapper {
  position: relative;
}

.game-canvas {
  position: relative;
  display: block;
}
</style>
