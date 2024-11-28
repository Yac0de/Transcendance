<!-- TournamentTree.vue -->
<template>
  <div class="tournament-container">
    <div class="timer">{{ remainingSeconds }}</div>
    <div class="bracket">
      <!-- Final -->
      <div class="match-winner">
        <p>Tournament Winner</p>
      </div>
      <div class="match-connections">
        <!-- Semi-Final 1 -->
        <div class="match-branch">
          <div class="bracket">
            <div class="match-winner">
              <p v-if="!UsersInFinal"> Semi 1 Winner</p>
              <p v-else> {{ UsersInFinal[0]?.displayname }} </p>
            </div>
            <div class="match-connections">
              <div class="match-branch">
                <p>Player 1:  {{ UsersInSemis1[0]?.displayname }} </p>
              </div>
              <div class="match-branch">
                <p>Player 2:  {{ UsersInSemis1[1]?.displayname }} </p>
              </div>
            </div>
          </div>
        </div>
        <!-- Semi-Final 2 -->
        <div class="match-branch">
          <div class="bracket">
            <div class="match-winner">
              <p v-if="!UsersInFinal"> Semi 1 Winner</p>
              <p v-else> {{ UsersInFinal[1]?.displayname }} </p>
            </div>
            <div class="match-connections">
              <div class="match-branch">
                <p>Player 3:  {{ UsersInSemis2[0]?.displayname }} </p>
              </div>
              <div class="match-branch">
                <p>Player 4:  {{ UsersInSemis2[1]?.displayname }} </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { UserData } from '../../types/models'
import { fetchMultipleUsers } from '../../utils/fetch'
import { eventBus } from '../../events/eventBus'
import { useRouter } from 'vue-router';

const UsersInSemis1 = ref<(UserData | null)[]>([null, null]); 
const UsersInSemis2 = ref<(UserData | null)[]>([null, null]); 
const UsersInFinal = ref<(UserData | null)[]>([null, null]); 

const semis1Score = ref<number[]>([])
const semis2Score = ref<number[]>([])
const finalScore = ref<number[]>([])

const remainingSeconds = ref<number>(16);
const router = useRouter();

//const props = defineProps<{
//  semi1array: number[]; 
//  semi2array: number[]; 
//}>();

onMounted(async () => {

  eventBus.on('TOURNAMENT_TIMER', (message: TournamentTimer) => {
    console.log("MSG TREE  = ", message);
    remainingSeconds.value = message.remainingTime;
  })

  eventBus.on('TOURNAMENT_TREE_STATE', async (message: TournamentStart) => {
    console.log("<- TOUR TREE STATE RECEIVED", message)
    if (message.semi1) {
      UsersInSemis1.value = await fetchMultipleUsers([message.semi1.player1id, message.semi1.player2id]); 
    }
    if (message.semi2) {
      UsersInSemis2.value = await fetchMultipleUsers([message.semi2.player1id, message.semi2.player2id]); 
    }
    if (message.final) {
      UsersInFinal.value = await fetchMultipleUsers([message.final.player1id, message.final.player2id]); 
    }
  })

  eventBus.on('TOURNAMENT_GAME', (message: TournamentGame) => {
    console.log("TOURNAMENT GAME READY TO START: ", message);
    router.push({
      path: '/game', 
      query: { lobbyId: message.lobbyId }
    });
  })
})

onUnmounted(() => {
  eventBus.off('TOURNAMENT_TIMER')
  eventBus.off('TOURNAMENT_TREE_STATE')
  eventBus.off('TOURNAMENT_GAME')
})
</script>

<style scoped>
.tournament-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 600px;
}

.timer {
  margin-bottom: 20px;
  font-size: 24px;
  font-weight: bold;
}

.bracket {
  display: flex;
  flex-direction: row-reverse;
}

.bracket p {
  padding: 20px;
  margin: 0;
  background-color: #f5f5f5;
  border-radius: 4px;
  min-width: 120px;
  text-align: center;
}

.match-winner {
  position: relative;
  margin-left: 50px;
  display: flex;
  align-items: center;
}

.match-winner::after {
  position: absolute;
  content: '';
  width: 25px;
  height: 2px;
  left: 0;
  top: 50%;
  background-color: #e0e0e0;
  transform: translateX(-100%);
}

.match-connections {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.match-branch {
  display: flex;
  align-items: flex-start;
  justify-content: flex-end;
  margin-top: 10px;
  margin-bottom: 10px;
  position: relative;
}

.match-branch::before {
  content: '';
  position: absolute;
  background-color: #e0e0e0;
  right: 0;
  top: 50%;
  transform: translateX(100%);
  width: 25px;
  height: 2px;
}

.match-branch::after {
  content: '';
  position: absolute;
  background-color: #e0e0e0;
  right: -25px;
  height: calc(50% + 22px);
  width: 2px;
  top: 50%;
}

.match-branch:last-child::after {
  transform: translateY(-100%);
}

.match-branch:only-child::after {
  display: none;
}
</style>
