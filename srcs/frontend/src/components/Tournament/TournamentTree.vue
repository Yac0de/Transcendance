<!-- TournamentTree.vue -->
<template>
  <div class="status-bracket-container">
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

const props = defineProps<{
  tournamentCode: string;
}>();

const userStore = useUserStore()

let goingIntoGame: boolean = false;

const UsersInSemis1 = ref<(UserData | null)[]>([null, null]); 
const UsersInSemis2 = ref<(UserData | null)[]>([null, null]); 
const UsersInFinal = ref<(UserData | null)[]>([null, null]); 

const winner = ref<UserData | null>(null);
let hasEmittedFinalMessage: boolean = false;

const hasLost = ref<boolean>(false);
const remainingSeconds = ref<number>(-1);
const tournamentStatusMessage = ref<string>(''); 
const router = useRouter();

const handleGameRouting = async (message: TournamentGame) => {
    console.log("LOBBY EVENT: ", message)
    const lobbyId = message.lobbyId;
    console.log("WE WILL PUSH THE ROUTER WITH THIS IN THE QUERY", lobbyId)
    
    // Set goingIntoGame before starting navigation
    goingIntoGame = true;
    
    try {
        // Wait for the navigation to complete
        await router.push({
            path: '/game',
            query: { lobbyId: lobbyId.toString() } // Ensure lobbyId is a string
        });
    } catch (error) {
        console.error('Navigation failed:', error);
        goingIntoGame = false; // Reset if navigation fails
    }
}

onMounted(() => {
    if (userStore.getWebSocketService?.isConnected()) {
        console.log("TC: ", props.tournamentCode)
        userStore.getWebSocketService?.sendTreeStateMessage(props.tournamentCode)
    } else {
        console.error('WebSocket is not connected');
    }

    // Other event listeners...
    eventBus.on('TOURNAMENT_GAME', handleGameRouting)
});

onUnmounted(() => {
    // Only attempt to leave if we're not transitioning to a game
    if (!goingIntoGame) {
        if (userStore.getWebSocketService?.isConnected()) {
            console.log("-> TOURNAMENT_LEAVE");
            userStore.getWebSocketService?.sendLeaveTournament()
        } else {
            console.error('WebSocket is not connected');
        }
    }

    // Clean up event listeners
    eventBus.off('TOURNAMENT_GAME', handleGameRouting)
    eventBus.off('TOURNAMENT_TIMER')
    eventBus.off('TOURNAMENT_TREE_STATE')
});


onMounted(async () => {
  if (userStore.getWebSocketService?.isConnected()) {
    console.log("TC: ", props.tournamentCode)
    userStore.getWebSocketService?.sendTreeStateMessage(props.tournamentCode)
  } else {
    console.error('WebSocket is not connected');
  }

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
        if (finalPlayer1Id === userStore.getId && !hasEmittedFinalMessage) {
          eventBus.emit('CHAT_FROM_TOURNAMENT_MASTER_SEMIS', "You are expected to play in the final next!");
          hasEmittedFinalMessage = true;
        } else if (finalPlayer2Id === userStore.getId && !hasEmittedFinalMessage) {
          eventBus.emit('CHAT_FROM_TOURNAMENT_MASTER_SEMIS', "You are expected to play in the final next!");
          hasEmittedFinalMessage = true;
        }
      } else {
        tournamentStatusMessage.value =  'You lost, better luck next time !'
        hasLost.value = true
      }
      UsersInFinal.value = await fetchMultipleUsers([finalPlayer1Id ?? 0, finalPlayer2Id ?? 0]); 
    }

    if (message.final?.isFinished) {
      if (message.final?.score[0] > message.final?.score[1]) {
        winner.value = await fetchUserById(message.final.player1id)
        if (message.final?.player1id === userStore.getId) {
          tournamentStatusMessage.value =  'Congratulations, you won the final !'
          eventBus.emit('CHAT_FROM_TOURNAMENT_MASTER_FINAL', "You won the tournament !");
        } else {
          tournamentStatusMessage.value =  'You lost, better luck next time !'
          eventBus.emit('CHAT_FROM_TOURNAMENT_MASTER_FINAL', "You lost in the final ... Too bad ..");
          hasLost.value = true
        }
      } else {
        winner.value = await fetchUserById(message.final?.player2id)
        if (message.final?.player2id === userStore.getId) {
          tournamentStatusMessage.value =  'Congratulations, you won the final !'
          eventBus.emit('CHAT_FROM_TOURNAMENT_MASTER_FINAL', "You won the tournament !");
        } else {
          tournamentStatusMessage.value =  'You lost, better luck next time !'
          eventBus.emit('CHAT_FROM_TOURNAMENT_MASTER_FINAL', "You lost in the final ... Too bad ..");
          hasLost.value = true
        }
      }
    }

  })

  
  eventBus.on('TOURNAMENT_GAME', handleGameRouting)

})

onUnmounted(() => {
  if (!goingIntoGame) {
    if (userStore.getWebSocketService?.isConnected()) {
      console.log("-> TOURNAMENT_LEAVE");
      userStore.getWebSocketService?.sendLeaveTournament()
    } else {
      console.error('WebSocket is not connected');
    }
  }

  eventBus.off('TOURNAMENT_GAME', handleGameRouting)
  eventBus.off('TOURNAMENT_TIMER')
  eventBus.off('TOURNAMENT_TREE_STATE')
})
</script>

<style scoped>

.status-bracket-container {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: center;
  gap: 2rem;
}


.tournament-status {  
  color: white;
  font-size: 1.3rem;
  text-shadow: 0.5px 0.5px 1px black;
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
  color: white;
  text-shadow: 0.5px 0.5px 1px black;
}

.bracket {
  display: flex;
  flex-direction: row-reverse;
  color: white;
  text-shadow: 0.5px 0.5px 1px black;
}

.bracket p {
  padding: 20px;
  margin: 0;
  background-color: var(--secondary-bright-color);
  border-radius: 4px;
  min-width: 120px;
  text-align: center;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.3);
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
  background-color: var(--secondary-bright-color);
  transform: translateX(-100%);
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.3);
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
  background-color: var(--secondary-bright-color);
  right: 0;
  top: 50%;
  transform: translateX(100%);
  width: 25px;
  height: 2px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.3);
}

.match-branch::after {
  content: '';
  position: absolute;
  background-color: var(--secondary-bright-color);
  right: -25px;
  height: calc(50% + 22px);
  width: 2px;
  top: 50%;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.3);
}

.match-branch:last-child::after {
  transform: translateY(-100%);
}

.match-branch:only-child::after {
  display: none;
}
</style>
