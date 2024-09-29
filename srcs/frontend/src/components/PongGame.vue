<template>
  <div class="game-container">
    <canvas ref="gameCanvas" width="1000" height="500" class="game-canvas"></canvas>
    <div class="controls">
      <button @click="startGame" v-if="!gameStarted && !isPaused">Play</button>
      <button @click="togglePause" v-if="gameStarted">{{ isPaused && !isFrozen ? 'Resume' : 'Pause' }}</button>
    </div>
    <div v-if="countdown > 0" class="countdown">{{ countdown }}</div>
    <div class="scoreboard">
      <span>{{ scorePlayer1 }}</span> - <span>{{ scorePlayer2 }}</span>
    </div>
    <div class="commands">
      <p>Player 1: W/S | Player 2: Up/Down</p>
    </div>
    <div v-if="message" class="message">{{ message }}</div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onBeforeUnmount } from 'vue'

interface Paddle {
  x: number
  y: number
  width: number
  height: number
  speed: number
}

interface Ball {
  x: number
  y: number
  radius: number
  speedX: number
  speedY: number
}

interface Keys {
  w: boolean
  s: boolean
  up: boolean
  down: boolean
}

const gameCanvas = ref<HTMLCanvasElement | null>(null)
const ctx = ref<CanvasRenderingContext2D | null>(null)

// Raquettes
const paddle1 = reactive<Paddle>({
  x: 30,
  y: 200,
  width: 10,
  height: 100,
  speed: 5
})

const paddle2 = reactive<Paddle>({
  x: 960,
  y: 200,
  width: 10,
  height: 100,
  speed: 5
})

// Balle
const ball = reactive<Ball>({
  x: 500,
  y: 250,
  radius: 10,
  speedX: 3,
  speedY: 3
})

// Points
const scorePlayer1 = ref(0)
const scorePlayer2 = ref(0)

// Variables pour le jeu
const keys = reactive<Keys>({
  w: false,
  s: false,
  up: false,
  down: false
})

let isPaused = ref(false)
let isFrozen = ref(false)
let gameStarted = ref(false)
let countdown = ref(0)
let exchanges = ref(0)
let message = ref('')

// Fonction d'initialisation du jeu
const initGame = () => {
  if (gameCanvas.value) {
    ctx.value = gameCanvas.value.getContext('2d')
    gameLoop()
  }
}

const gameLoop = () => {
  if (!isPaused.value && gameStarted.value && countdown.value === 0) {
    update()
    draw()
  }
  requestAnimationFrame(gameLoop)
}

const update = () => {
  if (!gameCanvas.value || isFrozen.value) return

  // Déplacement des raquettes
  if (keys.w && paddle1.y > 0) paddle1.y -= paddle1.speed
  if (keys.s && paddle1.y < gameCanvas.value.height - paddle1.height) paddle1.y += paddle1.speed
  if (keys.up && paddle2.y > 0) paddle2.y -= paddle2.speed
  if (keys.down && paddle2.y < gameCanvas.value.height - paddle2.height) paddle2.y += paddle2.speed

  // Déplacement de la balle
  ball.x += ball.speedX
  ball.y += ball.speedY

  // Collisions avec le haut et le bas du canvas
  if (ball.y + ball.radius > gameCanvas.value.height || ball.y - ball.radius < 0) {
    ball.speedY *= -1
  }

  // Collision avec la raquette de gauche (paddle1)
  if (
    ball.x - ball.radius < paddle1.x + paddle1.width &&
    ball.y > paddle1.y &&
    ball.y < paddle1.y + paddle1.height
  ) {
    ball.speedX *= -1
    ball.x = paddle1.x + paddle1.width + ball.radius
    exchanges.value++
    checkSpeedIncrease()
  }

  // Collision avec la raquette de droite (paddle2)
  if (
    ball.x + ball.radius > paddle2.x &&
    ball.y > paddle2.y &&
    ball.y < paddle2.y + paddle2.height
  ) {
    ball.speedX *= -1
    ball.x = paddle2.x - ball.radius
    exchanges.value++
    checkSpeedIncrease()
  }

  // Si la balle dépasse les bords gauche ou droit, ajouter un point et réinitialiser la balle
  if (ball.x - ball.radius < 0) {
    scorePlayer2.value++
    resetBall()
  }

  if (ball.x + ball.radius > gameCanvas.value.width) {
    scorePlayer1.value++
    resetBall()
  }
}

// Fonction pour augmenter la vitesse de la balle après un certain nombre d'échanges
const checkSpeedIncrease = () => {
  if (exchanges.value % 5 === 0 && exchanges.value > 0) {
    ball.speedX *= 1.2
    ball.speedY *= 1.2
    message.value = 'Ball speed increased!'
    setTimeout(() => {
      message.value = ''
    }, 2000) // Le message disparaît après 2 secondes
  }
}

const resetBall = () => {
  exchanges.value = 0
  isFrozen.value = true
  setTimeout(() => {
    ball.x = 500
    ball.y = 250
    ball.speedX = (Math.random() > 0.5 ? 1 : -1) * 3
    ball.speedY = (Math.random() > 0.5 ? 1 : -1) * 3
    isFrozen.value = false
  }, 1000) // Pause de 1 seconde avant de relancer la balle
}

// Bouton pour mettre le jeu en pause ou le reprendre
const togglePause = () => {
  if (!isFrozen.value) {
    isPaused.value = !isPaused.value
  }
}

// Démarrer le jeu avec un décompte de 3 secondes
const startGame = () => {
  gameStarted.value = true
  countdown.value = 3
  const countdownInterval = setInterval(() => {
    countdown.value--
    if (countdown.value === 0) {
      clearInterval(countdownInterval)
    }
  }, 1000)
}

const draw = () => {
  if (!ctx.value || !gameCanvas.value) return

  // Effacer le canvas
  ctx.value.clearRect(0, 0, gameCanvas.value.width, gameCanvas.value.height)

  // Dessiner les raquettes
  ctx.value.fillStyle = 'blue'
  ctx.value.fillRect(paddle1.x, paddle1.y, paddle1.width, paddle1.height)
  ctx.value.fillRect(paddle2.x, paddle2.y, paddle2.width, paddle2.height)

  // Dessiner la balle
  ctx.value.fillStyle = 'red'
  ctx.value.beginPath()
  ctx.value.arc(ball.x, ball.y, ball.radius, 0, Math.PI * 2)
  ctx.value.fill()

  // Dessiner le score
  ctx.value.font = '30px Arial'
  ctx.value.fillText(`${scorePlayer1.value}`, 50, 50)
  ctx.value.fillText(`${scorePlayer2.value}`, 920, 50)
}

const handleKeyDown = (e: KeyboardEvent) => {
  const key = e.key.toLowerCase()
  if (key === 'w') keys.w = true
  if (key === 's') keys.s = true
  if (key === 'arrowup') keys.up = true
  if (key === 'arrowdown') keys.down = true
}

const handleKeyUp = (e: KeyboardEvent) => {
  const key = e.key.toLowerCase()
  if (key === 'w') keys.w = false
  if (key === 's') keys.s = false
  if (key === 'arrowup') keys.up = false
  if (key === 'arrowdown') keys.down = false
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
</script>

<style scoped>
.game-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: 20px;
  position: relative;
}

.game-canvas {
  border: 1px solid black;
}

.controls {
  margin-top: 20px;
}

.scoreboard {
  margin-top: 10px;
  font-size: 30px;
  position: absolute;
  top: 10px;
  left: 50%;
  transform: translateX(-50%);
}

.commands {
  margin-top: 10px;
  font-size: 20px;
}

.countdown {
  font-size: 50px;
  font-weight: bold;
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

.message {
  position: absolute;
  top: 20%;
  left: 50%;
  transform: translateX(-50%);
  font-size: 30px;
  color: red;
}
</style>
