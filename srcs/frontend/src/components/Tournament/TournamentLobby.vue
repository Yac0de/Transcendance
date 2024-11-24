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
        v-for="(user, index) in users" 
        :key="index"
        :player="user"
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
import { eventBus } from '../../events/eventBus'
import { fetchUserById } from '../../utils/fetch'

const userStore = useUserStore();

defineProps<{
  tournamentCode: string;
}>();

const users = ref<(UserData | null)[]>([null, null, null, null]); 

const fetchMultipleUsers = async (userIds: number[]) => {
    try {
        const userPromises = userIds.map(id =>
        id !== 0 ? fetchUserById(id) : Promise.resolve(null)
      );
        const users = await Promise.all(userPromises);
        return users;
    } catch (error) {
        console.error("One or more fetches failed:", error);
        throw error;
    }
};

const handleStartTournament = () => {
  // Implement tournament start logic
  console.log('Starting tournament...', tournamentCode.value);
};

onMounted(() => {
  eventBus.on('TOURNAMENT_EVENT', async (message: tournamentLobbyState) => {
    console.log("TOURNAMENT EVENT: ", message);

    const playerIds = [
      message.player1id,
      message.player2id,
      message.player3id,
      message.player4id
    ];

    users.value = await fetchMultipleUsers(playerIds);
  })

  eventBus.on('TOURNAMENT_CREATE', async (message: tournamentCreate) => {
    console.log("TOURNAMENT CREATE EVENT: ", message);

    const playerIds = [
      message.player1id,
      message.player2id,
      message.player3id,
      message.player4id
    ];

    users.value = await fetchMultipleUsers(playerIds);
  })
})

onUnmounted(() => {
  eventBus.off('TOURNAMENT_EVENT');
  eventBus.off('TOURNAMENT_CREATE');
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
