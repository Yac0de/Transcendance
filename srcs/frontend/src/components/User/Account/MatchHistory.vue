<template>
  <div class="match-history-container">
    <!-- Menu de navigation style LoL -->
    <div class="history-nav">
      <button 
        v-for="tab in tabs" 
        :key="tab.id"
        @click="currentTab = tab.id"
        :class="['nav-button', { active: currentTab === tab.id }]"
      >
        {{ tab.name }}
        <span class="match-count">{{ tab.count }}</span>
      </button>
    </div>

    <!-- Zone de filtres -->
    <div class="filters">
      <button class="filter-button">
        <span>Time Period</span>
        <span>Last 5 matches</span>
      </button>
      <button class="filter-button">
        <span>Game Type</span>
        <span>All</span>
      </button>
    </div>

    <!-- Liste des matchs -->
    <div class="matches-list" v-if="games.length > 0">
      <div 
        v-for="game in games" 
        :key="game.id" 
        :class="['match-card', { 'victory': game.is_winner, 'defeat': !game.is_winner }]"
      >
        <!-- Type de match + Résultat -->
        <div class="match-info">
          <span class="match-type">Classic</span>
          <span class="match-result">{{ game.is_winner ? 'Victory' : 'Defeat' }}</span>
          <span class="match-duration">5:00</span>
        </div>

        <!-- Scores -->
        <div class="match-stats">
          <div class="score">
            {{ game.score1 }} - {{ game.score2 }}
          </div>
        </div>

        <!-- Date -->
        <div class="match-date">
          {{ formatDate(game.created_at) }}
        </div>
      </div>
    </div>

    <!-- Message si pas de matchs -->
    <div v-else class="no-matches">
      No matches found
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import  gameHistoryService  from '../../../services/gameHistoryService'
import { useUserStore } from '../../../stores/user'


interface Game {
  id: number
  player1_id: number
  player2_id: number
  winner_id: number
  score1: number
  score2: number
  created_at: string
  is_winner: boolean
}

const route = useRoute()
const userStore = useUserStore()
const currentTab = ref('all')
const games = ref<Game[]>([])

const tabs = computed(() => [
  { id: 'all', name: 'Overview', count: games.value.length },
  { id: 'ranked', name: 'Ranked', count: 0 },
  { id: 'normal', name: 'Normal', count: games.value.length },
])

const formatDate = (dateString: string): string => {
  const date = new Date(dateString)
  return date.toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
  })
}
console.log("Entire user store state:", userStore.$state); // Ajoutez cette ligne

onMounted(async () => {
  console.log("MatchHistory component mounted");
  try {
    const userId = userStore.id;
    console.log("User ID from store:", userId);
    if (userId) {
      console.log("Attempting to fetch history for user", userId);
      const history = await gameHistoryService.getUserHistory(userId);
      console.log("History received:", history);
      if (history) {
        // Assigner directement l'historique reçu
        games.value = history;
      }
    }
  } catch (error) {
    console.error('Failed to fetch game history:', error);
  }
});

</script>

<style scoped>
.match-history-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  background: #f3f3f3;
}

.history-nav {
  display: flex;
  gap: 4px;
  margin-bottom: 20px;
  background: #1a1a1a;
  padding: 4px;
  border-radius: 4px;
}

.nav-button {
  padding: 8px 16px;
  color: #9f9f9f;
  background: transparent;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: all 0.2s;
}

.nav-button.active {
  color: #fff;
  background: #2f2f2f;
  border-radius: 4px;
}

.match-count {
  background: #3f3f3f;
  padding: 2px 6px;
  border-radius: 10px;
  font-size: 0.8em;
}

.filters {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
}

.filter-button {
  display: flex;
  flex-direction: column;
  padding: 8px 16px;
  background: #fff;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
}

.filter-button span:first-child {
  font-size: 0.8em;
  color: #666;
}

.matches-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.match-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  border-radius: 4px;
  border-left: 4px solid;
}

.victory {
  background: rgba(63, 185, 80, 0.1);
  border-left-color: #3fb950;
}

.defeat {
  background: rgba(218, 54, 51, 0.1);
  border-left-color: #da3633;
}

.match-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.match-type {
  font-size: 0.9em;
  color: #666;
}

.match-result {
  font-weight: bold;
}

.victory .match-result {
  color: #3fb950;
}

.defeat .match-result {
  color: #da3633;
}

.match-duration {
  font-size: 0.8em;
  color: #666;
}

.match-stats {
  display: flex;
  align-items: center;
  gap: 20px;
}

.score {
  font-size: 1.2em;
  font-weight: bold;
}

.match-date {
  color: #666;
  font-size: 0.9em;
}

.no-matches {
  text-align: center;
  padding: 40px;
  color: #666;
  background: #fff;
  border-radius: 4px;
}
</style>