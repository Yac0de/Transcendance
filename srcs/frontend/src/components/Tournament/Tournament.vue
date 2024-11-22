<template>
  <div class="tournament-container">
    <button 
      v-if="currentView !== 'menu'"
      class="back-button"
      @click="handleBack"
    >
      ← Back
    </button>

    <!-- Main Menu View -->
    <div v-if="currentView === 'menu'" class="menu-view">
      <h1 class="tournament-title">Tournament Hub</h1>
      
      <div class="tournament-buttons">
        <button 
          class="tournament-button create"
          @click="handleCreateTournament"
        >
          Create Tournament
        </button>
        
        <button 
          class="tournament-button join"
          @click="handleJoinTournament"
        >
          Join Tournament
        </button>
      </div>
    </div>

    <!-- Join Tournament View -->
    <div v-else-if="currentView === 'join'" class="join-view">
      <h2 class="view-title">Join Tournament</h2>
      <JoinTournamentMenu />
    </div>

    <!-- Create Tournament View (placeholder) -->
    <div v-else-if="currentView === 'create'" class="create-view">
      <h2 class="view-title">Create Tournament</h2>
      <!-- CreateTournament component will go here -->
      <TournamentLobby :tournament-code="tournamentCode"/>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../../stores/user'
import JoinTournamentMenu from './JoinTournamentMenu.vue'
import TournamentLobby from './TournamentLobby.vue'
import { eventBus } from '../../events/eventBus'

type ViewState = 'menu' | 'join' | 'create'

interface Tournament {
  id: string
  name: string
  participants: number
  status: 'pending' | 'active' | 'completed'
}

const userStore = useUserStore();
const router = useRouter()
const currentView = ref<ViewState>('menu')
const tournamentCode = ref<string>('')

const handleCreateTournament = (): void => {
  if (userStore.getWebSocketService?.isConnected()) {
    userStore.getWebSocketService?.createTournamentLobby(tournamentCode.value)
  } else {
    console.error('WebSocket is not connected');
  }
  currentView.value = 'create'
}

const handleJoinTournament = (): void => {
  currentView.value = 'join'
}

const handleBack = (): void => {
  currentView.value = 'menu'
}

onMounted(() => {
  eventBus.on('CREATE_TOURNAMENT_LOBBY', (message: joinTournamentWithCode) => {
    console.log("TOURNAMENT LOBBY CREATED WITH SUCCESS, CODE = ", message.code);
    tournamentCode.value = message.code
    currentView.value = 'create'
  })
})

onUnmounted(() => {
  eventBus.off('CREATE_TOURNAMENT_LOBBY');
})
</script>

<style scoped>
.tournament-container {
  position: relative;
  max-width: 800px;
  margin: 2rem auto;
  padding: 2rem;
  text-align: center;
}

.back-button {
  position: absolute;
  left: 2rem;
  top: 2rem;
  padding: 0.5rem 1rem;
  font-size: 1rem;
  background: none;
  border: 1px solid #ccc;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
}

.back-button:hover {
  background-color: #f5f5f5;
  border-color: #999;
}

.tournament-title {
  margin-bottom: 2rem;
  font-size: 2rem;
  color: #333;
}

.view-title {
  margin-bottom: 2rem;
  font-size: 1.8rem;
  color: #333;
}

.tournament-buttons {
  display: flex;
  gap: 1.5rem;
  justify-content: center;
}

.tournament-button {
  padding: 1rem 2rem;
  font-size: 1.1rem;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: transform 0.2s, background-color 0.2s;
}

.tournament-button:hover {
  transform: translateY(-2px);
}

.tournament-button.create {
  background-color: #4CAF50;
  color: white;
}

.tournament-button.create:hover {
  background-color: #45a049;
}

.tournament-button.join {
  background-color: #2196F3;
  color: white;
}

.tournament-button.join:hover {
  background-color: #1e88e5;
}

/* View transitions */
.menu-view,
.join-view,
.create-view {
  animation: fadeIn 0.3s ease-in-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
