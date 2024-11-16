<template>
  <div class="leave-button-container">
    <button class="nav-button leave-lobby-button" @click="leaveLobby">Leave Lobby</button>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { eventBus } from '../../events/eventBus'
import { useUserStore } from '../../stores/user'

const router = useRouter();

let lobbyId: string = '';
let inviterId: number = 0;
const userStore = useUserStore()

const leaveLobby = () => {
  const wsService = userStore.getWebSocketService
  if (wsService && lobbyId !== '') {
    wsService.leaveAndTerminateLobby(lobbyId);
    console.log('terminating the lobby: ', lobbyId);
  }
  router.push('/');
}

onMounted(() => {
  console.log('Leave lobby  component mounted')
  eventBus.on('LOBBY_INVITATION_TO_FRIEND', (message) => {
    console.log('Game invite event sent with success: ', message.lobbyId);
    lobbyId = message.lobbyID;
    inviterId = message.sender.id;
  })
})

onUnmounted(() => {
  eventBus.off('LOBBY_INVITATION_TO_FRIEND')
})
</script>

<style scoped>
.leave-button-container {
  position: absolute;
  top: 10px;
  left: 20px;
  z-index: 100;
}

.leave-lobby-button {
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  text-decoration: none;
  transition: background-color 0.3s;
  padding: 5px 10px;
  font-size: 16px;
}

.leave-lobby-button:hover {
  background-color: #0056b3;
}

@media (max-width: 600px) {
  .leave-button-container {
    top: 120px;
    /* Additional space for collapsed navbar on mobile */
  }
}
</style>
