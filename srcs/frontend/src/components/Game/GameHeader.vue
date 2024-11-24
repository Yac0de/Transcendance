<template>
    <div class="game-header">
      <div class="scoreP1">
        <div class="player-info">
          <div class="avatar-container">
          <img :src="api.user.getAvatarUrl(player1?.avatar ?? null)" class="avatar" alt="Player 1 avatar">
          </div>
          <span>{{ player1?.nickname || 'Player 1' }} {{ state.score.player1 }}</span>
        </div>
      </div>
      <div class="timer-container">
        <div class="timer-title">TIME
          <div class="timer">{{ state.remainingTime }}s</div>
        </div>
      </div>
      <div class="scoreP2">
        <div class="player-info">
          <span>{{ player2?.nickname || 'Player 2' }} {{ state.score.player2 }}</span>
          <div class="avatar-container">
            <img :src="api.user.getAvatarUrl(player2?.avatar ?? null)" class="avatar" alt="Player 2 avatar">
          </div>
        </div>
      </div>
    </div>
</template>
  
<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { GameState } from '../../types/game';
import { UserData } from '../../types/models'
import { fetchUserById } from '../../utils/fetch'
import api from '../../services/api';

const props = defineProps<{
  state: GameState
  player1id: number | null
  player2id: number | null
}>()

const player1 = ref<UserData | null>(null)
const player2 = ref<UserData | null>(null)

const fetchPlayerData = async () => {
  if (props.player1id) {
    player1.value = await fetchUserById(props.player1id ?? null)
  }
  if (props.player2id) {
    player2.value = await fetchUserById(props.player2id ?? null)
  }
}

onMounted(fetchPlayerData)

// Watch for changes in player IDs and refetch data when they change
watch([() => props.player1id, () => props.player2id], fetchPlayerData)
</script>

<style scoped>
.game-header {
  position: relative;
  display: flex;
  justify-content: space-around;
  align-items: center;
  background-color: #99307a;
  color: black;
  width: 800px;
  height: 50px;
  font-size: x-large;
  font-weight: bolder;
}

.player-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.avatar-container {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  border: 2px solid #f7ddef;
  background-color: #5b3c54;
}

.avatar {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.timer-container {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: black;
  clip-path: polygon(0% 0, 100% 0, 80% 100%, 20% 100%);
  width: 25%;
  height: 98%;
  font-size: 1rem;
}

.timer-title {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #5b3c54;
  letter-spacing: .2rem;
}

.timer {
  letter-spacing: normal;
  color: #f7ddef;
}

.scoreP1 .player-info {
  flex-direction: row;
}

.scoreP2 .player-info {
  flex-direction: row-reverse;
}

/* Optional: Add hover effect on avatar */
.avatar-container:hover {
  border-color: white;
  transform: scale(1.05);
  transition: all 0.2s ease-in-out;
}
</style>
