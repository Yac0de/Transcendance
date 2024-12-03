<template>
  <div v-if="userExist" class="match-history-container">
    <!-- Menu de navigation -->
    <div class="history-title">
      <h1>{{ $t('history') }}</h1>
      <WinrateCircle :percentage="winrate" />
      <RankIcon :winrate="winrate" />
    </div>

    <!-- Zone de filtres -->
    <div class="filters">
      <button class="filter-button">
        <span>{{ $t('timePeriod') }}</span>
        <span>{{ $t('last6Matches') }}</span>
      </button>
    </div>

    <!-- Liste des matchs -->
    <div class="matches-list" v-if="games.length > 0">
      <div 
        v-for="game in limitedGames" 
        :key="game.id" 
        :class="['match-card', { 'victory': game.is_winner, 'defeat': !game.is_winner }]"
      >
        <!-- Type de match + Résultat -->
        <div class="match-info">
          <span class="match-type">{{ $t('classic') }}</span>
          <div class="match-details">
            <span class="match-result">{{ game.is_winner ? $t('victory') : $t('defeat') }}</span>
          </div>
        </div>

        <!-- Scores -->
        <div class="match-stats">
          <div class="score">
            <template v-if="game.player1.nickname === targetNickname">
              <span class="player-name">{{ game.player1.nickname }}</span>
              {{ game.score1 }} - {{ game.score2 }}
              <span class="player-name">{{ game.player2.nickname }}</span>
            </template>
            <template v-else>
              <span class="player-name">{{ game.player2.nickname }}</span>
              {{ game.score2 }} - {{ game.score1 }}
              <span class="player-name">{{ game.player1.nickname }}</span>
            </template>
          </div>
        </div>

        <!-- Date -->
        <div class="match-date">
          {{ formatDate(game.CreatedAt) }}
        </div>
      </div>
    </div>

    <!-- Message si pas de matchs -->
    <div v-else class="no-matches">
      {{ $t('noMatches') }}
    </div>
  </div>
  <NotFound v-else />
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import  gameHistoryService  from '../../../services/gameHistoryService'
import { useUserStore } from '../../../stores/user'
import WinrateCircle from './WinrateCircle.vue';
import NotFound from '../../General/NotFound.vue';
import { GameHistory} from '../../../types/models';
import RankIcon from './RankIcon.vue'

const winrate =ref(0)
const route = useRoute()
const userStore = useUserStore()
const games =  ref<GameHistory[]>([])
const userExist = ref<boolean>(true)


const limitedGames = computed(() => {
  return games.value.slice(0, 6);
})

const targetNickname = ref(route.params.nickname as string)

const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
  })
}

onMounted(async () => {
  console.log("test " + targetNickname.value)
  try {
    console.log("NEW:", route.params);
    const userId = userStore.id;
    if (userId) {
      
      const history = await gameHistoryService.getUserHistory(route.params.nickname as string);
      console.log("History received:", history);
      if (history) {
        // Assigner directement l'historique reçu
        games.value = history;
        const victories = history.filter(game =>game.is_winner).length;
        winrate.value = (victories / history.length) * 100;
      }
    }
  } catch (error) {
    console.error('Failed to fetch game history:', error);
    userExist.value = false;
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
.history-title {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
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

.score {
  display: flex;
  align-items: center;
  gap: 8px;  /* Espace entre les éléments */
  font-size: 1.2em;
  font-weight: bold;
}

.player-name {
  font-size: 0.8em;  /* Taille plus petite pour les noms */
  font-weight: normal;  /* Police normale pour les noms */
  color: #fffdfd;  /* Couleur grise pour les noms */
}
</style>
