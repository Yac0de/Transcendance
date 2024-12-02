<template>
  <div class="tournament-lobby">
    <div class="code-display">
      <h3>{{ $t('tournamentCode') }}</h3>
      <div class="code-box" @click="copyToClipboard">
        {{ tournamentCode.slice(0, 8) }}
        <i class="fas fa-clipboard copy-icon"></i>
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
      :disabled="users.includes(null)" 
      v-if="creatorId === clientId"
      class="start-button"
      @click="handleStartTournament"
    >
      {{ $t('startTournament') }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import PlayerTile from './PlayerTile.vue';
import { UserData } from '../../types/models';
import { useUserStore } from '../../stores/user'
import { eventBus } from '../../events/eventBus'
import { fetchMultipleUsers } from '../../utils/fetch'
import { useRouter } from 'vue-router';
import { TournamentCreate, TournamentEvent } from '../../types/tournament';
import { useI18n } from 'vue-i18n';

const userStore = useUserStore();
const router = useRouter();
const tournamentCode = ref<string>('');
const creatorId = ref<number>(0);
const clientId = ref<number | null>(userStore.getId);
const users = ref<(UserData | null)[]>([null, null, null, null]); 

const { t } = useI18n();

const copyToClipboard = (event: MouseEvent) => {
  const element = event.currentTarget as HTMLElement;
  element.classList.add('click-animation');
  
  setTimeout(() => {
    element.classList.remove('click-animation');
  }, 150); 

  navigator.clipboard.writeText(tournamentCode.value.slice(0, 8))
    .then(() => {
      console.log(t('copySuccess')); // Message de confirmation dans la console
    })
    .catch(err => console.error('Failed to copy code:', err));
};

const handleStartTournament = () => {
  if (userStore.getWebSocketService?.isConnected()) {
    userStore.getWebSocketService?.sendTournamentStart(tournamentCode.value)
  } else {
    console.error('WebSocket is not connected');
  }
};

onMounted(() => {
  eventBus.on('TOURNAMENT_EVENT', async (message: TournamentEvent) => {
    try {
      creatorId.value = message.player1id ?? 0;
      const playerIds = [
        message.player1id ?? 0,
        message.player2id ?? 0,
        message.player3id ?? 0,
        message.player4id ?? 0
      ];
      tournamentCode.value = String(message.code);
      users.value = await fetchMultipleUsers(playerIds);
    } catch (error) {
      console.error("Failed to handle tournament event:", error);
      users.value = [null, null, null, null];
    }
  });

  eventBus.on('TOURNAMENT_CREATE', async (message: TournamentCreate) => {
    try {
      creatorId.value = message.player1id ?? 0;
      const playerIds = [
        message.player1id ?? 0,
        message.player2id ?? 0,
        message.player3id ?? 0,
        message.player4id ?? 0
      ];
      tournamentCode.value = String(message.code);
      users.value = await fetchMultipleUsers(playerIds);
    } catch (error) {
      console.error("Failed to handle tournament event:", error);
      users.value = [null, null, null, null];
    }
  })

  eventBus.on('TOURNAMENT_TERMINATE', async () => {
    router.push('/');
  })
})

onUnmounted(() => {
  eventBus.off('TOURNAMENT_EVENT');
  eventBus.off('TOURNAMENT_CREATE');
  eventBus.off('TOURNAMENT_TERMINATE');
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
  color: white;
}

.code-box {
  display: flex;
  align-items: center;
  gap: 10px;
  background-color: #f5f5f5;
  padding: 0.75rem 1.5rem;
  border-radius: 4px;
  font-family: monospace;
  font-size: 1.2rem;
  letter-spacing: 2px;
  color: #333;
  border: 1px solid #e0e0e0;
  cursor: pointer;
  transition: transform 0.2s ease;
}

/* Animation class for click effect */
.click-animation {
  transform: scale(0.95);
}

.copy-icon {
  font-size: 1rem;
  color: #666;
  margin-left: 5px;
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
  font-weight: 600;
  background: linear-gradient(to right, var(--secondary-bright-color), 
  color-mix(in srgb, var(--secondary-bright-color) 75%, white));
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.start-button:hover:not(:disabled) {
  background: linear-gradient(to right, var(--secondary-bright-color), 
  color-mix(in srgb, var(--secondary-bright-color) 85%, white));
  transform: translateY(-2px);
}

.start-button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}
</style>
