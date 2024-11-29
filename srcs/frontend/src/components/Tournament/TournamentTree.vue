<!-- TournamentTree.vue -->
<template>
  <div class="tournament-container">
    <div v-if="tournamentStatusMessage" class="tournament-status">
      {{ tournamentStatusMessage }}
    </div>
    <div v-if="remainingSeconds != -1 && !hasLost" class="timer">{{ remainingSeconds }}</div>
    <div class="bracket">
      <!-- Final -->
      <div class="match-winner">
        <p v-if="!winner"> Winner</p>
        <p v-else> {{ winner?.displayname }} </p>
      </div>
      <div class="match-connections">
        <!-- Semi-Final 1 -->
        <div class="match-branch">
          <div class="bracket">
            <div class="match-winner">
              <p v-if="!UsersInFinal[0]"> Semi 1 Winner</p>
              <p v-else> {{ UsersInFinal[0]?.displayname }} </p>
            </div>
            <div class="match-connections">
              <div class="match-branch">
                <p> {{ UsersInSemis1[0]?.displayname }} </p>
              </div>
              <div class="match-branch">
                <p> {{ UsersInSemis1[1]?.displayname }} </p>
              </div>
            </div>
          </div>
        </div>
        <!-- Semi-Final 2 -->
        <div class="match-branch">
          <div class="bracket">
            <div class="match-winner">
              <p v-if="!UsersInFinal[1]"> Semi 2 Winner</p>
              <p v-else> {{ UsersInFinal[1]?.displayname }} </p>
            </div>
            <div class="match-connections">
              <div class="match-branch">
                <p> {{ UsersInSemis2[0]?.displayname }} </p>
              </div>
              <div class="match-branch">
                <p> {{ UsersInSemis2[1]?.displayname }} </p>
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
import { TournamentTimer, TournamentTreeState, TournamentGame } from '../../types/tournament'
import { fetchMultipleUsers, fetchUserById } from '../../utils/fetch'
import { eventBus } from '../../events/eventBus'
import { useRouter } from 'vue-router';
import { useUserStore } from '../../stores/user';

const userStore = useUserStore()

const UsersInSemis1 = ref<(UserData | null)[]>([null, null]); 
const UsersInSemis2 = ref<(UserData | null)[]>([null, null]); 
const UsersInFinal = ref<(UserData | null)[]>([null, null]); 

const winner = ref<UserData | null>(null);

const hasLost = ref<boolean>(false);
const remainingSeconds = ref<number>(-1);
const tournamentStatusMessage = ref<string>(''); 
const router = useRouter();

onMounted(async () => {

  eventBus.on('TOURNAMENT_TIMER', (message: TournamentTimer) => {
    remainingSeconds.value = message.remainingTime;
  })

  eventBus.on('TOURNAMENT_TREE_STATE', async (message: TournamentTreeState) => {
    if (message.semi1) {
      UsersInSemis1.value = await fetchMultipleUsers([message.semi1.player1id, message.semi1.player2id]); 
    }
    if (message.semi2) {
      UsersInSemis2.value = await fetchMultipleUsers([message.semi2.player1id, message.semi2.player2id]); 
    }
    if (message.final?.player1id !== 0 || message.final?.player2id !== 0) {
      const finalPlayer1Id = message.final?.player1id ?? null
      const finalPlayer2Id = message.final?.player2id ?? null
      if (finalPlayer1Id === userStore.getId || finalPlayer2Id === userStore.getId) {
        tournamentStatusMessage.value =  'Congratulations, you are qualified in the final'
      } else {
        tournamentStatusMessage.value =  'You lost, better luck next time !'
        hasLost.value = true
      }
      UsersInFinal.value = await fetchMultipleUsers([finalPlayer1Id, finalPlayer2Id]); 
    }

    if (message.final.isFinished) {
      if (message.final.score[0] > message.final.score[1]) {
        winner.value = await fetchUserById(message.final.player1id)
        if (message.final.player1id === userStore.getId) {
          tournamentStatusMessage.value =  'Congratulations, you won the final !'
        } else {
          tournamentStatusMessage.value =  'You lost, better luck next time !'
          hasLost.value = true
        }
      } else {
        winner.value = await fetchUserById(message.final.player2id)
        if (message.final.player2id === userStore.getId) {
          tournamentStatusMessage.value =  'Congratulations, you won the final !'
        } else {
          tournamentStatusMessage.value =  'You lost, better luck next time !'
          hasLost.value = true
        }
      }
    }
  })

  eventBus.on('TOURNAMENT_GAME', (message: TournamentGame) => {
    router.push({
      path: '/game', 
      query: { lobbyId: message.lobbyId }
    });
  })
})

onUnmounted(() => {

  if (userStore.getWebSocketService?.isConnected()) {
    console.log("-> TOURNAMENT_LEAVE");
    userStore.getWebSocketService?.sendLeaveTournament()
  } else {
    console.error('WebSocket is not connected');
  }

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

.status-message {
  margin-bottom: 20px;
  font-size: 24px;
  font-weight: bold;
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
