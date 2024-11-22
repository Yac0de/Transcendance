<template>
  <div class="game-container">
    <GameHeader :state="currentGameState"/>
    <div class="canvas-wrapper">
      <canvas
        id="gameCanvas"
        width="800"
        height="600"
        ref="canvasRef"
        class="game-canvas"
      ></canvas>
    </div>
  </div>
</template>

<script setup lang="ts">
import { GameEvent, GameData } from '../../types/game'
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { eventBus } from '../../events/eventBus';
import { drawPaddle, drawBall } from '../../services/gamerender';
import { useUserStore } from '../../stores/user';
import { useRoute } from 'vue-router'
import GameHeader from './GameHeader.vue';

const route = useRoute()

const userStore = useUserStore();
const canvasRef = ref<HTMLCanvasElement | null>(null)

const currentGameState = reactive({
  score: { player1: 0, player2: 0 },
  remainingTime: 300,
  paddle: { player1X: 30, player1Y: 240, player2X: 740, player2Y: 240 },
  ball: { x: 400, y: 300 },
});

const handlePressUp = (event: KeyboardEvent): void => {
  if (event.code === 'ArrowUp' || event.code === 'KeyW') {
    if (userStore.getWebSocketService?.isConnected()) {
      const gameEvent: GameEvent = {
        type: 'GAME_EVENT',
        lobbyId: route.query.lobbyId as string,
        userId: userStore.getId!,
        keyPressed: 'UP'
      };
      userStore.getWebSocketService?.sendGameEvent(gameEvent);
    } else {
      console.error('WebSocket is not connected');
    }
  }
}

const handleReleaseUp = (event: KeyboardEvent): void => {
  if (event.code === 'ArrowUp' || event.code === 'KeyW') {
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

onMounted(() => {
  // Add key listener
  window.addEventListener('keydown', handlePressUp)
  window.addEventListener('keydown', handlePressDown)
  window.addEventListener('keyup', handleReleaseUp)
  window.addEventListener('keyup', handleReleaseDown)

  eventBus.on('GAME_EVENT', async (message: GameEvent) => {
    if (canvasRef.value) {
      const ctx:CanvasRenderingContext2D = canvasRef.value.getContext('2d') as CanvasRenderingContext2D
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
})

onUnmounted(() => {
  // Remove key listener
  window.removeEventListener('keydown', handlePressUp)
  window.removeEventListener('keydown', handlePressDown)
  window.removeEventListener('keyup', handleReleaseUp)
  window.removeEventListener('keyup', handleReleaseDown)

  eventBus.off('GAME_EVENT')
})
</script>

<style scoped>
body {
  margin: 0;
  overflow: hidden;
}
.game-container {
  width: 100%;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background-color: #f0f0f0;
}

.canvas-wrapper {
  border: 2px solid #333;
  border-radius: 4px;
  background-color: white;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.game-canvas {
  display: block; /* Removes bottom margin/spacing */
}
</style>
