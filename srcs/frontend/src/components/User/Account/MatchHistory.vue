<template>
  <div class="game-history">
    <h2 class="text-2xl font-bold mb-6">Game History</h2>
    <div class="space-y-6">
      <!-- Individual game row -->
      <div 
        v-for="game in games" 
        :key="game.id" 
        class="game-row"
      >
        <div class="flex justify-between items-center">
          <!-- Left side: Game info -->
          <div class="flex items-center space-x-6">
            <div 
              class="game-result" 
              :class="game.result === 'Victory' ? 'bg-green-100 text-green-800 border-green-200' : 'bg-red-100 text-red-800 border-red-200'"
            >
              {{ game.result }}
            </div>
            <div class="game-details">
              <h3 class="font-semibold text-gray-800">{{ game.gameType }}</h3>
              <p class="text-gray-600">{{ formatDate(game.date) }}</p>
            </div>
          </div>
          
          <!-- Right side: Stats -->
          <div class="flex space-x-8">
            <div class="stat-item">
              <span class="stat-label">Score</span>
              <span class="stat-value">{{ game.score }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">Duration</span>
              <span class="stat-value">{{ game.duration }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">Rank</span>
              <span class="stat-value">{{ game.rank }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface GameStats {
  id: number
  result: 'Victory' | 'Defeat'
  gameType: string
  date: string
  score: string
  duration: string
  rank: string
}

const games = ref<GameStats[]>([
  {
    id: 1,
    result: 'Victory',
    gameType: 'Ranked Match',
    date: '2024-03-25T14:30:00',
    score: '15/3/7',
    duration: '32:45',
    rank: 'Gold II'
  },
  {
    id: 2,
    result: 'Defeat',
    gameType: 'Casual Match',
    date: '2024-03-25T13:15:00',
    score: '8/5/4',
    duration: '28:30',
    rank: 'Gold II'
  },
  {
    id: 3,
    result: 'Victory',
    gameType: 'Ranked Match',
    date: '2024-03-25T11:45:00',
    score: '12/2/9',
    duration: '35:20',
    rank: 'Gold II'
  }
])

const formatDate = (dateString: string): string => {
  return new Date(dateString).toLocaleString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>

<style scoped>
.game-history {
  max-width: 900px;
  margin: 0 auto;
  padding: 20px;
}

.game-row {
  background-color: #f8f9fa;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  border: 1px solid #e5e7eb;
  transition: all 0.2s ease-in-out;
}

.game-row:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
  background-color: #ffffff;
}

.game-result {
  padding: 8px 16px;
  border-radius: 12px;
  font-weight: 600;
  font-size: 0.875rem;
  border: 2px solid;
  transition: all 0.2s ease;
}

.game-result:hover {
  transform: scale(1.05);
}

.game-details {
  padding: 4px 0;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 8px 16px;
  background-color: #ffffff;
  border-radius: 8px;
  min-width: 100px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.stat-label {
  font-size: 0.75rem;
  color: #6b7280;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 4px;
}

.stat-value {
  font-weight: 600;
  color: #374151;
  font-size: 1rem;
}
</style>
