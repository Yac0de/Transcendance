<template>
  <div v-if="show" class="invite-popup">
    <div class="popup-content">
      <p class="text-message">{{ inviter?.nickname }} invited you to play!</p>
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
import { UserData } from '../../types/models'
import { LobbyInvitationFromFriend } from '../../types/lobby'
import { fetchUserById } from '../../utils/fetch'

const show = ref(false)

const userStore = useUserStore()
const router = useRouter()

let lobbyId: string = '';
let inviter: UserData | null = null;

const accept = () => {
  if (userStore.isRedirectPending) {
    show.value = false;
    return
  }

  const wsService = userStore.getWebSocketService
  if (wsService) {
    if (inviter) {
      wsService.acceptInviteFromFriend(lobbyId, inviter.id);
    }
    router.push('/lobby')
  }
  show.value = false
}

const decline = () => {
  const wsService = userStore.getWebSocketService
  if (wsService) {
    if (inviter) {
      wsService.denyInviteFromFriend(lobbyId, inviter.id);
      console.log('WebSocket service found, would send DECLINE for inviterId:', inviter.id)
    }
  }
  show.value = false
}

onMounted(() => {
  eventBus.on('LOBBY_INVITATION_FROM_FRIEND', async (message: LobbyInvitationFromFriend) => {
    lobbyId = message.lobbyId;
    inviter = await fetchUserById(message.sender.id);
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
  right: 20px;
  z-index: 2000;
  width: 300px;
  /* Fixed width for the tile */
}

.popup-content {
  background-color: #f3f4f6;
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
}

button:hover {
  transform: scale(1.02);
}

.accept-button {
  background: linear-gradient(to right, var(--secondary-bright-color), color-mix(in srgb, var(--secondary-bright-color) 75%, white));
  color: white;
  border: none;
}

.accept-button:hover {
  background: linear-gradient(to right, var(--secondary-bright-color), color-mix(in srgb, var(--secondary-bright-color) 85%, white));
}

.decline-button {
  background: linear-gradient(to right, var(--secondary-dark-color), color-mix(in srgb, var(--secondary-dark-color) 75%, white));
  color: white;
  border: none;
}

.decline-button:hover {
  background: linear-gradient(to right, var(--secondary-dark-color), color-mix(in srgb, var(--secondary-dark-color) 85%, white));
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
