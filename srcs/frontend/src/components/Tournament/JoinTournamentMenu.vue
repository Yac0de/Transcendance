<template>
  <div class="join-tournament">
    <div class="tournament-prompt">
      <input
        v-model="tournamentCode"
        type="text"
        class="tournament-input"
        placeholder="Enter tournament code"
      />
      <button 
        class="join-button"
        @click="handleJoin"
      >
        Join Tournament
      </button>
      <div v-if="error" class="error-message">
        {{ error }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import type { joinTournamentWithCode } from '../types/tournament'
import { useUserStore } from '../../stores/user'
import { eventBus } from '../../events/eventBus'

const tournamentCode = ref<string>('')
const userStore = useUserStore()

defineProps<{
  error: string;
}>();

const handleJoin = () => {
  console.log('Joining tournament with code:', tournamentCode.value)
  if (userStore.getWebSocketService?.isConnected()) {
    userStore.getWebSocketService?.joinTournamentWithCode(tournamentCode.value)
  } else {
    console.error('WebSocket is not connected');
  }
}

</script>

<style scoped>
.join-tournament {
  max-width: 400px;
  margin: 0 auto;
}

.tournament-prompt {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  padding: 2rem;
  border-radius: 8px;
  background-color: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.tournament-input {
  padding: 0.75rem;
  font-size: 1rem;
  border: 1px solid #ccc;
  border-radius: 4px;
}

.tournament-input:focus {
  outline: none;
  border-color: #2196F3;
}

.join-button {
  padding: 0.75rem;
  font-size: 1rem;
  background-color: #2196F3;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.join-button:hover {
  background-color: #1e88e5;
}

.error-message {
  margin-top: 1px;
  padding: 0.75rem;
  background-color: #ffebee;
  color: #c62828;
  border-radius: 4px;
  font-size: 0.875rem;
  text-align: center;
  border: 1px solid #ffcdd2;
}
</style>
