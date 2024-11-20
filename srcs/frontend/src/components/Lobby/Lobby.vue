<template>
  <div class="lobby-container" :class="{ 'no-clicks': showNotification }">
    <div v-if="showNotification" class="modal-overlay">
      <div class="notification">
        Lobby has been destroyed, you will be redirected to homepage. 
        Maybe your opponent was scared ?
      </div>
    </div>
    <LeaveLobbyButton :lobby-id="lobbyId"/>
    <Timer :remaining-seconds="remainingSeconds" v-if="showTimer"/>
    <div class="players-container">
      <div class="player-column">
        <PlayerItem :is-left="true"
          :challenged-friend="challengedFriend" />
        <ReadyCheck :both-players-ready="player1Ready && player2Ready" :isPlayerReady="player1Ready" :challenged-friend="challengedFriend" :is-accepting="isAcceptingPlayer" :lobbyId="lobbyId" :disabled="false" v-if="bothPlayerPresent && showReadyChecks" @ready-changed="handlePlayer1Ready" />
      </div>
      <div class="versus">VS</div>
      <div class="player-column">
        <component :is="challengedFriend ? PlayerItem : PlayerScrolldown" :is-left="false"
          :challenged-friend="challengedFriend"
          @friend-selected="handleFriendSelected" />
        <ReadyCheck :both-players-ready="player1Ready && player2Ready" :isPlayerReady="player2Ready" :challenged-friend="challengedFriend" :is-accepting="isAcceptingPlayer" :lobbyId="lobbyId" :disabled="true" v-if="bothPlayerPresent && showReadyChecks" @ready-changed="handlePlayer2Ready" />
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
import { UserData } from '../../types/models';
import { LobbyPlayerStatus, LobbyCreated, LobbyPregameRemainingTime } from '../../types/lobby';
import { eventBus } from '../../events/eventBus'
import { fetchUserById } from '../../utils/fetch'
import { useRouter } from 'vue-router'

const router = useRouter()
const userStore = useUserStore();

const player1Ready = ref<boolean>(false)
const player2Ready = ref<boolean>(false)

let lobbyId: string = '';
const challengedFriend = ref<UserData | null>(null);
let challengedFriendId = ref<number>(0);
let isAcceptingPlayer: boolean = false;

const showNotification = ref<boolean>(false);

const bothPlayerPresent = computed(() => {
  return challengedFriendId.value !== 0;
});

const remainingSeconds = ref<number>(0);
const showReadyChecks = ref<boolean>(true);
const showTimer = ref<boolean>(false);

const handleFriendSelected = (friendId: number) => {
  if (userStore.getWebSocketService?.isConnected()) {
    userStore.getWebSocketService?.inviteFriendToLobbyMessage(friendId);
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
    challengedFriendId.value = isAcceptingPlayer ? message.sender.id : message.receiver.id;
    challengedFriend.value = await fetchUserById(challengedFriendId.value);
  })

  eventBus.on('LOBBY_PLAYER_STATUS', (message: LobbyPlayerStatus) => {
    if (message.userId === userStore.getId) {
      player1Ready.value = player1Ready.value ? false : true 
    } else if (message.userId === challengedFriendId.value) {
      player2Ready.value = player2Ready.value ? false : true 
    }
  })

  eventBus.on('LOBBY_PREGAME_REMAINING_TIME', async (message: LobbyPregameRemainingTime) => {
    if (showReadyChecks.value === true) {
      await new Promise(resolve => setTimeout(resolve, 1000));
    }
    remainingSeconds.value = message.remainingSecondsToStart;
    showReadyChecks.value = false;
    showTimer.value = true;
  })

  eventBus.on('GAME_EVENT', async () => {
    console.log("GAME EVENT RECEIVED")
    router.push({
      path: '/game', 
      query: {lobbyId: lobbyId }
    });
  });

  eventBus.on('LOBBY_DESTROYED', async () => {
    showNotification.value = true;
    userStore.isRedirectPending = true;
    setTimeout(() => {
      showNotification.value = false;
      userStore.isRedirectPending = false;
      router.push('/'); 
    }, 3000);
  })
})

onUnmounted(() => {
  eventBus.off('LOBBY_CREATED')
  eventBus.off('LOBBY_PLAYER_STATUS')
  eventBus.off('LOBBY_PREGAME_REMAINING_TIME')
  eventBus.off('LOBBY_DESTROYED')
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

.no-clicks {
  pointer-events: none;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  z-index: 999;
  display: flex;
  justify-content: center;
  align-items: center;
  pointer-events: all;
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

.notification {
  background: #2c3e50;
  color: white;
  padding: 12px 24px;
  border-radius: 4px;
  z-index: 1000;
  box-shadow: 0 2px 4px rgba(0,0,0,0.2);
  text-align: center;
}

@media (max-width: 768px) {
  .players-container {
    flex-direction: column;
    gap: 20px;
  }
}
</style>
