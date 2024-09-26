<template>
  <div class="game-container">
    <canvas ref="gameCanvas" width="400" height="400" class="game-canvas"></canvas>
    <p>Use WASD keys to move the square</p>
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, onBeforeUnmount, ref, reactive } from 'vue'

interface Square {
  x: number
  y: number
  size: number
  speed: number
}

interface Keys {
  w: boolean
  a: boolean
  s: boolean
  d: boolean
}

export default defineComponent({
  name: 'PongGame',
  setup() {
    const gameCanvas = ref<HTMLCanvasElement | null>(null)
    const ctx = ref<CanvasRenderingContext2D | null>(null)

    const square = reactive<Square>({
      x: 175,
      y: 175,
      size: 50,
      speed: 5
    })

    const keys = reactive<Keys>({
      w: false,
      a: false,
      s: false,
      d: false
    })

    const initGame = () => {
      if (gameCanvas.value) {
        ctx.value = gameCanvas.value.getContext('2d')
        gameLoop()
      }
    }

    const gameLoop = () => {
      update()
      draw()
      requestAnimationFrame(gameLoop)
    }

    const update = () => {
      if (!gameCanvas.value) return

      if (keys.w && square.y > 0) square.y -= square.speed
      if (keys.s && square.y < gameCanvas.value.height - square.size) square.y += square.speed
      if (keys.a && square.x > 0) square.x -= square.speed
      if (keys.d && square.x < gameCanvas.value.width - square.size) square.x += square.speed
    }

    const draw = () => {
      if (!ctx.value || !gameCanvas.value) return

      // Clear the canvas
      ctx.value.clearRect(0, 0, gameCanvas.value.width, gameCanvas.value.height)
      
      // Draw the square
      ctx.value.fillStyle = 'blue'
      ctx.value.fillRect(square.x, square.y, square.size, square.size)
    }

    const handleKeyDown = (e: KeyboardEvent) => {
      const key = e.key.toLowerCase()
      if (key in keys) {
        keys[key as keyof Keys] = true
      }
    }

    const handleKeyUp = (e: KeyboardEvent) => {
      const key = e.key.toLowerCase()
      if (key in keys) {
        keys[key as keyof Keys] = false
      }
    }

    onMounted(() => {
      initGame()
      window.addEventListener('keydown', handleKeyDown)
      window.addEventListener('keyup', handleKeyUp)
    })

    onBeforeUnmount(() => {
      window.removeEventListener('keydown', handleKeyDown)
      window.removeEventListener('keyup', handleKeyUp)
    })

    return {
      gameCanvas
    }
  }
})
</script>

<style scoped>
.game-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: 20px;
}
.game-canvas {
  border: 1px solid black;
}
</style>
