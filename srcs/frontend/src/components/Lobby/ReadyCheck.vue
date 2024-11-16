<template>
  <div class="ready-check">
    <button 
      class="ready-button" 
      :class="{ 'ready': isPlayerReady }" 
      @click="isPlayerReady ? toggleUnready() : toggleReady()"
    >
      {{ isPlayerReady ? 'Ready!' : 'Not Ready' }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { useUserStore } from '../../stores/user'
import { UserData } from '../../types/models';

const userStore = useUserStore();

interface Props {
  disabled?: boolean;
  lobbyId?: string;
  isPlayerReady: boolean;
  challengedFriend?: UserData | null;
  isAccepting: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
  lobbyId: '',
  isPlayerReady: false,
  challengedFriend: null,
  isAccepting: false
})

const toggleReady = () => {
  if (props.disabled) return
  if (userStore.getWebSocketService?.isConnected()) {
    userStore.getWebSocketService?.sendPlayerReadyMessage(props.isAccepting, props.challengedFriend.id, props.lobbyId);
  } else {
    console.error('WebSocket is not connected');
  }
}

const toggleUnready = () => {
  if (props.disabled) return
  if (userStore.getWebSocketService?.isConnected()) {
    userStore.getWebSocketService?.sendPlayerUnreadyMessage(props.isAccepting, props.challengedFriend.id, props.lobbyId);
  } else {
    console.error('WebSocket is not connected');
  }
}

</script>

<style scoped>
.ready-check {
  width: 120px;
  margin: 10px;
}

.ready-button {
  width: 100%;
  padding: 10px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: bold;
  transition: all 0.3s ease;
}

.ready-button:not(.ready) {
  background-color: #dc3545;
  color: white;
}

.ready-button.ready {
  background-color: #28a745;
  color: white;
}
</style>
