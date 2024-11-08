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
        <PlayerItem :player-name="player2Name" :is-left="false" />
        <ReadyCheck @ready-changed="handlePlayer2Ready" />
      </div>
    </div>
  </div>
</template>

<script>
import LeaveLobbyButton from './LeaveLobbyButton.vue'
import PlayerItem from './PlayerItem.vue'
import ReadyCheck from './ReadyCheck.vue'

export default {
  name: 'Lobby',
  components: {
    LeaveLobbyButton,
    PlayerItem,
    ReadyCheck
  },
  data() {
    return {
      player1Name: 'Player 1',
      player2Name: 'Player 2',
      player1Ready: false,
      player2Ready: false
    }
  },
  methods: {
    handleLeaveLobby() {
      // Implement lobby leaving logic
      console.log('Leaving lobby...')
    },
    handlePlayer1Ready(isReady) {
      this.player1Ready = isReady
    },
    handlePlayer2Ready(isReady) {
      this.player2Ready = isReady
    }
  }
}
</script>

<style scoped>
.lobby-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
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
