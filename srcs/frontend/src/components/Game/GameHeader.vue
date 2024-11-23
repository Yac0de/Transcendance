<template>
  <div class="game-header">
    <div class="score">
      {{ player1?.nickname || 'Player 1' }} {{ state.score.player1 }} - 
      {{ player2?.nickname || 'Player 2' }} : {{ state.score.player2 }}
    </div>
    <div class="timer">Time Left: {{ state.remainingTime }}s</div>
  </div>
</template>
  
<script setup lang="ts">
import { PropType } from 'vue';
import { GameState } from '../../types/game';
import { UserData } from '../../types/models'
import { fetchUserById } from '../../utils/fetch'
import { ref, onMounted, watch } from 'vue'

const props = defineProps<{
  state: GameState
  player1Id: number | null
  player2Id: number | null
}>()

const player1 = ref<UserData | null>(null)
const player2 = ref<UserData | null>(null)

const fetchPlayerData = async () => {
  if (props.player1Id) {
    player1.value = await fetchUserById(props.player1Id)
  }
  if (props.player2Id) {
    player2.value = await fetchUserById(props.player2Id)
  }
}

onMounted(fetchPlayerData)

// Watch for changes in player IDs and refetch data when they change
watch([() => props.player1Id, () => props.player2Id], fetchPlayerData)
</script>

<style scoped>
.game-header {
  display: flex;
  justify-content: space-between;
  padding: 10px;
  background-color: #333;
  color: #fff;
}
</style>
