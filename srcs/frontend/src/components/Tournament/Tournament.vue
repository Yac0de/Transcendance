<template>
  <div class="tournament-container">
    <!-- Main Menu View -->
    <div v-if="currentView === 'menu'" class="menu-view">
      <h1 class="tournament-title">Tournament</h1>
      
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
      <JoinTournamentMenu :error="error"/>
    </div>

    <!-- Waiting room tournament view (placeholder) -->
    <div v-else-if="currentView === 'waiting-room'" class="create-view">
      <!-- CreateTournament component will go here -->
      <TournamentWaitingRoom/>
    </div>
    <div v-else-if="currentView === 'tournament-tree'" class="create-view">
      <!-- CreateTournament component will go here -->
      <TournamentTree/>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '../../stores/user'
import JoinTournamentMenu from './JoinTournamentMenu.vue'
import TournamentWaitingRoom from './TournamentWaitingRoom.vue'
import TournamentTree from './TournamentTree.vue'
import { eventBus } from '../../events/eventBus'
import { TournamentCreate, TournamentEvent, TournamentError } from '../../types/tournament'

type ViewState = 'menu' | 'join' | 'waiting-room' | 'tournament-tree'

const userStore = useUserStore();
const currentView = ref<ViewState>('menu')
const tournamentCode = ref<string>('')
const route = useRoute();

const error = ref<string>('')

const handleCreateTournament = (): void => {
  if (userStore.getWebSocketService?.isConnected()) {
    userStore.getWebSocketService?.createTournamentWaitingRoom()
  } else {
    console.error('WebSocket is not connected');
  }
  currentView.value = 'waiting-room'
}

const handleJoinTournament = (): void => {
  currentView.value = 'join'
}

onMounted(() => {
  if (route.query.view === 'tournament-tree') {
    currentView.value = 'tournament-tree'
  }

  eventBus.on('TOURNAMENT_CREATE', (message: TournamentCreate) => {
    tournamentCode.value = String(message.code)
    currentView.value = 'waiting-room'
  })

  eventBus.on('TOURNAMENT_JOIN_WITH_CODE', () => {
    currentView.value = 'waiting-room'
  })

  eventBus.on('TOURNAMENT_EVENT', (message: TournamentEvent) => {
    tournamentCode.value = String(message.code)
  })

  eventBus.on('TOURNAMENT_START', () => {
    currentView.value = 'tournament-tree'
    eventBus.emit('CHAT_FROM_TOURNAMENT_MASTER', "You just started a tournament, good luck ..");
  })

  eventBus.on('TOURNAMENT_ERROR', (message: TournamentError) => {
    error.value = message.error
  })
})

onUnmounted(() => {
  if (tournamentCode.value) {
    if (userStore.getWebSocketService?.isConnected()) {
      userStore.getWebSocketService?.leaveTournamentWaitingRoom(tournamentCode.value)
    } else {
      console.error('WebSocket is not connected');
    }
  }

  eventBus.off('TOURNAMENT_CREATE');
  eventBus.off('TOURNAMENT_JOIN_WITH_CODE');
  eventBus.off('TOURNAMENT_START');
  eventBus.off('TOURNAMENT_ERROR');
})
</script>

<style scoped>
.tournament-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: space-around;
  text-align: center;
  height: 65vh;
  min-height: 200px;
  min-width: 300px;
  padding: 2vh 5vw;
  border-radius: 20px;
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.5);
  overflow: hidden;
  background: var(--main-color);
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
  color: white;
  text-shadow: 1px 1px 10px black;
  font-size: 3rem;
  margin-bottom: 2vh;
}

.view-title {
  color: white;
  text-shadow: 1px 1px 10px black;
  font-size: 2.5rem;
  margin-bottom: 2vh;
}

.tournament-buttons {
  display: flex;
  gap: 1.5rem;
  justify-content: center;
}

.tournament-button {
  width: 100%;
  padding: 0.5rem 1rem;
  font-size: 1.1rem;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.5);
  transition: transform 0.2s, background-color 0.2s;
}

.tournament-button:hover {
  transform: translateY(-2px);
}

.tournament-button.create {
  background: linear-gradient(  to right,
                                var(--secondary-dark-color),
                                color-mix(in srgb, var(--secondary-dark-color) 75%,
                                white));
}

.tournament-button.create:hover {
  background: linear-gradient(  to right,
                                var(--secondary-dark-color),
                                color-mix(in srgb, var(--secondary-dark-color) 85%,
                                white));
}

.tournament-button.join {
  background: linear-gradient(  to right,
                                var(--secondary-bright-color), 
                                color-mix(in srgb, var(--secondary-bright-color) 75%,
                                white));
}

.tournament-button.join:hover {
  background: linear-gradient(  to right,
                                var(--secondary-bright-color), 
                                color-mix(in srgb, var(--secondary-bright-color) 85%,
                                white));
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
