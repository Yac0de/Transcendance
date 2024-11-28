<template>
  <div class="match-history-container">
    <!-- Menu de navigation style LoL -->
    <div class="history-title">
      <h1>HISTORY</h1>
    </div>

    <!-- Zone de filtres -->
    <div class="filters">
      <button class="filter-button">
        <span>Time Period</span>
        <span>Last 5 matches</span>
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
    <div class="match-details">
      <span class="match-result">{{ game.is_winner ? 'Victory' : 'Defeat' }}</span>
    </div>
  </div>

  <!-- Scores -->
  <div class="match-stats">
    <div class="score">
      <!-- Modifions l'ordre des scores en fonction de si vous êtes player1 ou player2 -->
      <template v-if="game.player1_id === userStore.id">
        {{ game.score1 }} - {{ game.score2 }}
      </template>
      <template v-else>
        {{ game.score2 }} - {{ game.score1 }}
      </template>
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
import { fetchUserById } from '../../../utils/fetch';


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
const games = ref<Game[]>([])

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
.history-title {
  text-align: center;
  margin-bottom: 20px;
  font-weight: 600;
  font-size: 1.4rem;
}

.match-history-container {
  width: 100%;  /* Utilise toute la largeur disponible */
  min-height: 100vh;  /* Prend au moins toute la hauteur de la vue */
  margin: 0;  /* Enlève les marges */
  padding: 40px;  /* Augmente le padding */
  background: #1a1a1a;  /* Fond plus sombre pour un meilleur contraste */
  color: #fff;  /* Texte blanc */
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

.match-details {
  display: flex;
  align-items: center;
  gap: 8px;
}

.player-name {
  font-size: 0.9em;
  opacity: 0.8;
  font-weight: normal;
}

.match-result {
  margin-right: 4px;
}

</style>