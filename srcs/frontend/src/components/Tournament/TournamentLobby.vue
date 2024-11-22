<!-- TournamentLobby.vue -->
<template>
  <div class="tournament-lobby">
    <div class="code-display">
      <h3>Tournament Code</h3>
      <div class="code-box">
        {{ tournamentCode }}
      </div>
    </div>
    
    <div class="players-container">
      <PlayerTile 
        v-for="(player, index) in players" 
        :key="index"
        :player="player"
      />
    </div>
    
    <button 
      class="start-button"
      :disabled="!canStartTournament"
      @click="handleStartTournament"
    >
      Start Tournament
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import PlayerTile from './PlayerTile.vue';
import UserData from '../../types/models';
import { useUserStore } from '../../stores/user'

const userStore = useUserStore();

defineProps<{
  tournamentCode: string;
}>();

const players = ref<(UserData | null)[]>([
  { username: userStore.getNickname, status: 'ready', avatar: userStore.getAvatarPath },
  null,
  null,
  null,
]);

const handleStartTournament = () => {
  // Implement tournament start logic
  console.log('Starting tournament...', tournamentCode.value);
};

onMounted(() => {
  eventBus.on('TOURNAMENT_LOBBY_STATE', (message: tournamentLobbyState) => {
    console.log("TOURNAMENT LOBBY STATE: ", message);
  })
})

onUnmounted(() => {
  eventBus.off('TOURNAMENT_LOBBY_STATE');
})
</script>

<style scoped>
.tournament-lobby {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2rem;
  padding: 2rem;
  max-width: 800px;
  margin: 0 auto;
}

.code-display {
  text-align: center;
}

.code-display h3 {
  margin-bottom: 0.5rem;
  color: #333;
}

.code-box {
  background-color: #f5f5f5;
  padding: 0.75rem 1.5rem;
  border-radius: 4px;
  font-family: monospace;
  font-size: 1.2rem;
  letter-spacing: 2px;
  color: #333;
  border: 1px solid #e0e0e0;
}

.players-container {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 2rem;
  padding: 2rem;
}

.start-button {
  padding: 1rem 2rem;
  font-size: 1.1rem;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.start-button:hover:not(:disabled) {
  background-color: #45a049;
  transform: translateY(-2px);
}

.start-button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}
</style>
