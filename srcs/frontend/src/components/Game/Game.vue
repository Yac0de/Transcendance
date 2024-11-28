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
        <div class="vertical-line"></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { GameEvent, GameState } from '../../types/game'
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { eventBus } from '../../events/eventBus';
import { drawPaddle, drawBall, drawEndGame } from '../../services/gamerender';
import { useUserStore } from '../../stores/user';
import { useRoute, useRouter } from 'vue-router'
import GameHeader from './GameHeader.vue';

const route = useRoute()
const router = useRouter()

const userStore = useUserStore();
const canvasRef = ref<HTMLCanvasElement | null>(null)

const player1Id = ref<number | null>(null)
const player2Id = ref<number | null>(null)

const currentGameState: GameState = reactive({
    ball: { x: 0, y: 0 },  // Assuming Ball has x, y properties
    score: { 
        player1: 0, 
        player2: 0 
    },
    isActive: false,
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

onMounted(() => {
  // Add key listener
  window.addEventListener('keydown', handlePressUp)
  window.addEventListener('keydown', handlePressDown)
  window.addEventListener('keyup', handleReleaseUp)
  window.addEventListener('keyup', handleReleaseDown)
  window.addEventListener('keydown', handleSpace)
  const ctx:CanvasRenderingContext2D = canvasRef.value.getContext('2d') as CanvasRenderingContext2D

  eventBus.on('GAME_EVENT', async (message: GameEvent) => {

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

        drawPaddle(ctx, message.state!);
        drawBall(ctx, message.state!);
      }
    }
  })

  eventBus.on('GAME_FINISHED', async (message: GameEvent) => {
    console.log("<- GAME FINISHED: ", message);
    console.log("PLAYERS = ", player1Id.value, player2Id.value);
    drawEndGame(ctx, message.state!, player1Id.value, player2Id.value);

    if (message.isTournamentGame === false) {
      window.setTimeout(() => {
        router.push('/');
      }, 3000)
    } else {
      window.setTimeout(() => {
        router.push({
        path:'/tournament',
        query: { view: 'tournament-tree' }
        });
      })
    }
  })

})

onUnmounted(() => {
  // Remove key listener
  window.removeEventListener('keydown', handlePressUp)
  window.removeEventListener('keydown', handlePressDown)
  window.removeEventListener('keyup', handleReleaseUp)
  window.removeEventListener('keyup', handleReleaseDown)
  window.removeEventListener('keydown', handleSpace)
  eventBus.off('GAME_EVENT')
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
  background: #e5c4dc;
  border: solid 8px #e5c4dc;
  border-radius: 1% 1% 1% 1%;
}

.canvas-wrapper {
  position: relative;
}

.vertical-line {
  position: absolute;
  top: 0;
  left: 50%;
  width: 1px;
  height: 100%;
  background-color: #5b3c54;
  transform: translateX(-50%);
}

.game-canvas {
  position: relative;
  display: block;
  border-radius: 0% 0% 1% 1%;
}
</style>
