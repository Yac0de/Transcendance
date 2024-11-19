<template>
  <div class="game-container">
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
import { GameEvent, Ball, Paddle, Player, Score, GameState, Game, GameCommand } from '../../types/game'
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { eventBus } from '../../events/eventBus';
import { drawPaddle, drawBall } from '../../services/gamerender';

const canvasRef = ref<HTMLCanvasElement | null>(null)

onMounted(() => {
  eventBus.on('GAME_EVENT', async (message: GameState) => {
    console.log(message.type)
    if (canvasRef.value) {
      const ctx = canvasRef.value.getContext('2d')
      if (ctx) {
        // You can add initial canvas setup here
        ctx.fillStyle = 'black'
        ctx.fillRect(0, 0, canvasRef.value.width, canvasRef.value.height)
        drawPaddle(ctx, message.state);
        drawBall(ctx, message.state);
      }
    }
  })
  // Initialize canvas context if needed
})

onUnmounted(() => {
  eventBus.off('GAME_EVENT')
})
</script>

<style scoped>
.game-container {
  width: 100%;
  min-height: 100vh;
  display: flex;
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
