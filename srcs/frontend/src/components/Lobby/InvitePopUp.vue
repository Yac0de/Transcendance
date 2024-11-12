<template>
  <div v-if="show" class="invite-popup">
    <div class="popup-content">
      <p class="text-message">{{ inviterName || 'Debug User' }} invited you to play!</p>
      <div class="button-container">
        <button @click="accept" class="accept-button">
          Accept
        </button>
        <button @click="decline" class="decline-button">
          Decline
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../../stores/user'
import { eventBus } from '../../events/eventBus'

const show = ref(false)
const inviterName = ref('Debug User')

const userStore = useUserStore()
const router = useRouter()

let lobbyId: string = '';
let inviterId: number = 0;

const accept = () => {
  console.log('Accept clicked')
  const wsService = userStore.getWebSocketService
  if (wsService) {
    wsService.acceptInviteFromFriend(lobbyId, inviterId);
    console.log('WebSocket service found, would send ACCEPT for inviterId:', inviterId)
    router.push('/lobby')
  }
  show.value = false
}

const decline = () => {
  console.log('Decline clicked')
  const wsService = userStore.getWebSocketService
  if (wsService) {
    wsService.denyInviteFromFriend(lobbyId, inviterId);
    console.log('WebSocket service found, would send DECLINE for inviterId:', inviterId)
  }
  show.value = false
}

onMounted(() => {
  console.log('InvitePopUp component mounted')
  eventBus.on('LOBBY_INVITATION_FROM_FRIEND', (message) => {
    console.log('Game invite event received: ', message.lobbyId);
    lobbyId = message.lobbyId;
    inviterId = message.senderId;
    show.value = true;
  })
})

onUnmounted(() => {
  eventBus.off('LOBBY_INVITATION_FROM_FRIEND')
})
</script>

<style scoped>
.invite-popup {
  position: fixed;
  top: 70px;
  /* Positioned just below nav bar - adjust if needed */
  right: 20px;
  z-index: 2000;
  width: 300px;
  /* Fixed width for the tile */
}

.popup-content {
  background-color: #f3f4f6;
  /* Light grey background */
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
}

.text-message {
  color: #374151;
  font-size: 0.95rem;
  margin-bottom: 12px;
  font-weight: 500;
}

.button-container {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
}

.accept-button,
.decline-button {
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 0.875rem;
  font-weight: 500;
  transition: all 0.2s ease;
}

.accept-button {
  background-color: #10B981;
  color: white;
  border: none;
}

.accept-button:hover {
  background-color: #059669;
}

.decline-button {
  background-color: #EF4444;
  color: white;
  border: none;
}

.decline-button:hover {
  background-color: #DC2626;
}

/* Optional: Add animation */
@keyframes slideIn {
  from {
    transform: translateX(100%);
    opacity: 0;
  }

  to {
    transform: translateX(0);
    opacity: 1;
  }
}

.invite-popup {
  animation: slideIn 0.3s ease-out;
}
</style>
