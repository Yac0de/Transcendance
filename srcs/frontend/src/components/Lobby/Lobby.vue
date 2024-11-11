<template>
  <div class="lobby-container">
    <LeaveLobbyButton @leave-lobby="handleLeaveLobby" />
    <div class="players-container">
      <div class="player-column">
        <PlayerItem :player-name="player1Name" :is-left="true" />
        <ReadyCheck @ready-changed="handlePlayer1Ready" />
      </div>
      <div class="versus">VS</div>
      <div class="player-column">
        <PlayerScrolldown :friends="friendsList" @friend-selected="handleFriendSelected"
          :selectedFriend="player2Name" />
        <ReadyCheck @ready-changed="handlePlayer2Ready" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import LeaveLobbyButton from './LeaveLobbyButton.vue'
import PlayerItem from './PlayerItem.vue'
import PlayerScrolldown from './PlayerScrolldown.vue'
import ReadyCheck from './ReadyCheck.vue'
import { useOnlineUsersStore } from '../../stores/onlineUsers';
import { Friend } from '../../../types/models';

const online_users_store = useOnlineUsersStore();

// State
const player1Name = ref<string>('Player 1') // Current player's name
const player2Name = ref<string | null>(null)
const player1Ready = ref<boolean>(false)
const player2Ready = ref<boolean>(false)

// Methods
const handleLeaveLobby = () => {
  // Implement lobby leaving logic
  console.log('Leaving lobby...')
}

const handleFriendSelected = (friend: Friend) => {
  player2Name.value = friend.name
  // Add any additional logic for when a friend is selected
  console.log('Friend selected:', friend)
}

const handlePlayer1Ready = (isReady: boolean) => {
  player1Ready.value = isReady
}

const handlePlayer2Ready = (isReady: boolean) => {
  player2Ready.value = isReady
}
</script>

<style scoped>
.lobby-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  gap: 20px;
}

.players-container {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 40px;
}

.player-column {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

.versus {
  font-size: 48px;
  font-weight: bold;
  color: #2c3e50;
  margin: 0 20px;
}

@media (max-width: 768px) {
  .players-container {
    flex-direction: column;
    gap: 20px;
  }

  .versus {
    margin: 10px 0;
  }
}
</style>
