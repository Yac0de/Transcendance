<template>
  <div class="lobby-container">
    <LeaveLobbyButton @leave-lobby="handleLeaveLobby" />
    <div class="players-container">
      <div class="player-column">
        <PlayerItem :is-left="true" :is-challenged="isAcceptingPlayer"
          :challenged-friend-id="isAcceptingPlayer ? userStore.getId : challengedFriendId" />
        <ReadyCheck v-if="bothPlayerPresent" @ready-changed="handlePlayer1Ready" />
      </div>
      <div class="versus">VS</div>
      <div class="player-column">
        <component :is="challengedFriend ? PlayerItem : PlayerScrolldown" :is-left="false"
          :is-challenged="!isAcceptingPlayer" :challenged-friend-id="challengedFriendId"
          @friend-selected="handleFriendSelected" />
        <ReadyCheck v-if="bothPlayerPresent" @ready-changed="handlePlayer2Ready" />
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
import { useOnlineUsersStore } from '../../stores/onlineUsers';
import { useUserStore } from '../../stores/user'
import { Friend, UserData } from '../../../types/models';
import { eventBus } from '../../events/eventBus'
import { fetchUserById } from '../../utils/fetch'

const online_users_store = useOnlineUsersStore();
const userStore = useUserStore();
const player1Ready = ref<boolean>(false)
const player2Ready = ref<boolean>(false)
let lobbyId: string = '';
const challengedFriend = ref<UserData | null>(null);
let challengedFriendId = ref<number>(0);
const isAcceptingPlayer = ref<boolean>(false);

const bothPlayerPresent = computed(() => {
  return challengedFriend.value !== null;
});

const handleLeaveLobby = () => {
  console.log('Leaving lobby...')
}

const handleFriendSelected = (friend: Friend) => {
  if (userStore.getWebSocketService?.isConnected()) {
    userStore.getWebSocketService?.inviteFriendToLobbyMessage(friend);
    console.log("SENT WS SOCKET MESSAGE TO INVITE A FRIEND TO A GAME");
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
  eventBus.on('LOBBY_CREATED', async (message) => {
    console.log('Lobby created event received: ', message);
    lobbyId = message.lobbyId;
    isAcceptingPlayer.value = message.receiverId === userStore.getId;
    challengedFriendId.value = isAcceptingPlayer.value ? message.senderId : message.receiverId;
    console.log(challengedFriendId.value);

    challengedFriend.value = await fetchUserById(challengedFriendId.value);
  })
})

onUnmounted(() => {
  eventBus.off('LOBBY_CREATED')
})
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
