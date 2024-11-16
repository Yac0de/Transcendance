<template>
  <div class="lobby-container">
    <LeaveLobbyButton @leave-lobby="handleLeaveLobby" />
    <Timer v-if="showTimer"/>
    <div class="players-container">
      <div class="player-column">
        <PlayerItem :is-left="true"
          :challenged-friend="challengedFriend" />
        <ReadyCheck :isPlayerReady="player1Ready" :challenged-friend="challengedFriend" :is-accepting="isAcceptingPlayer" :lobbyId="lobbyId" :disabled="false" v-if="bothPlayerPresent && showReadyChecks" @ready-changed="handlePlayer1Ready" />
      </div>
      <div class="versus">VS</div>
      <div class="player-column">
        <component :is="challengedFriend ? PlayerItem : PlayerScrolldown" :is-left="false"
          :challenged-friend="challengedFriend"
          @friend-selected="handleFriendSelected" />
        <ReadyCheck :isPlayerReady="player2Ready" :challenged-friend="challengedFriend" :is-accepting="isAcceptingPlayer" :lobbyId="lobbyId" :disabled="true" v-if="bothPlayerPresent && showReadyChecks" @ready-changed="handlePlayer2Ready" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import LeaveLobbyButton from './LeaveLobbyButton.vue'
import PlayerItem from './PlayerItem.vue'
import PlayerScrolldown from './PlayerScrolldown.vue'
import ReadyCheck from './ReadyCheck.vue'
import Timer from './Timer.vue'
import { useUserStore } from '../../stores/user'
import { Friend, UserData } from '../../types/models';
import { LobbyPlayerStatus, LobbyCreated } from '../../types/lobby';
import { eventBus } from '../../events/eventBus'
import { fetchUserById } from '../../utils/fetch'

const userStore = useUserStore();

const player1Ready = ref<boolean>(false)
const player2Ready = ref<boolean>(false)

let lobbyId: string = '';
const challengedFriend = ref<UserData | null>(null);
let challengedFriendId = ref<number>(0);
let isAcceptingPlayer: boolean = false;

const bothPlayerPresent = computed(() => {
  return challengedFriendId.value !== 0;
});

const showReadyChecks = ref<boolean>(true);
const showTimer = ref<boolean>(false);

const handleLeaveLobby = () => {
  console.log('Leaving lobby...')
}

const handleFriendSelected = (friend: Friend) => {
  if (userStore.getWebSocketService?.isConnected()) {
    userStore.getWebSocketService?.inviteFriendToLobbyMessage(friend);
  } else {
    console.error('WebSocket is not connected');
  }
}

const handlePlayer1Ready = (isReady: boolean) => {
  player1Ready.value = isReady
}

const handlePlayer2Ready = (isReady: boolean) => {
  player2Ready.value = isReady
}

onMounted(() => {
  eventBus.on('LOBBY_CREATED', async (message: LobbyCreated) => {
    lobbyId = message.lobbyId;
    isAcceptingPlayer = message.receiver.id === userStore.getId;
    console.log("IS CLIENT THE ACCEPTING PLAYER", isAcceptingPlayer);
    challengedFriendId.value = isAcceptingPlayer ? message.sender.id : message.receiver.id;
    console.log("ID OF CHALL", challengedFriendId.value);
    challengedFriend.value = await fetchUserById(challengedFriendId.value);
    console.log("CHALL = ", challengedFriend.value.nickname);
  })

  eventBus.on('LOBBY_PLAYER_STATUS', (message: LobbyPlayerStatus) => {
    console.log('LOBBY PLAYER STATUS UPDATE, ID: ', message.userId);
    if (message.userId === userStore.getId) {
      player1Ready.value = player1Ready.value ? false : true 
    } else if (message.userId === challengedFriendId.value) {
      player2Ready.value = player2Ready.value ? false : true 
    }
  })

  eventBus.on('LOBBY_PREGAME_REMAINING_TIME', (message: LobbyPregameRemainingTime) => {
    showReadyChecks.value = false;
    showTimer.value = true;
  })
})

onUnmounted(() => {
  eventBus.off('LOBBY_CREATED')
  eventBus.off('LOBBY_PLAYER_STATUS')
})
</script>

<style scoped>
.lobby-container {
  min-height: 100vh;
  display: flex;
  position: relative;
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

}
</style>
