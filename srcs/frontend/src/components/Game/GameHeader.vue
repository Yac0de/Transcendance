<template>
    <div class="game-header">
      <div class="scoreP1">
        {{ player1?.nickname || 'Player 1' }} {{ state.score.player1 }}
      </div>
      <div class="timer-container">
        <div class="timer-title">TIME
          <div class="timer">{{ state.remainingTime }}s</div>
        </div>
      </div>
      <div class="scoreP2">
        {{ player2?.nickname || 'Player 2' }} {{ state.score.player2 }}
      </div>
    </div>
</template>
  
<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { GameState } from '../../types/game';
import { UserData } from '../../types/models'
import { fetchUserById } from '../../utils/fetch'

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
    console.log(player2.value.nickname)
  }
}
onMounted(fetchPlayerData)

// Watch for changes in player IDs and refetch data when they change
watch([() => props.player1Id, () => props.player2Id], fetchPlayerData)
</script>

<style scoped>
  .game-header {
    position: relative;
    display: flex;
    justify-content:space-around;
    align-items: center;
    background-color: #99307a;
    color: black;
    width: 800px;
    height: 50px;
    font-size: x-large;
    font-weight: bolder;
  }

  .timer-container {
    position: absolute;
    top: 50%; left: 50%;
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

  .timer-title{
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    color: #5b3c54;
    letter-spacing: .2rem;
  }

  .timer{
    letter-spacing: normal;
    color: #f7ddef;
  }
</style>
