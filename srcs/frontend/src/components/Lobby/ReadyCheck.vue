<template>
  <div class="ready-check">
    <ToggleButton
      :activeLabel="'Ready!'"
      :inactiveLabel="'Not Ready'"
      :isActive="isPlayerReady"
      :disabled="disabled || bothPlayersReady"
      @toggle="handleToggle"
    />
  </div>
</template>

<script setup lang="ts">
import ToggleButton from './ToggleButton.vue';
import { useUserStore } from '../../stores/user';

const userStore = useUserStore();

interface Props {
  disabled?: boolean;
  lobbyId?: string;
  isPlayerReady: boolean;
  bothPlayersReady?: boolean;
  isGameMode : boolean;
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
  lobbyId: '',
  isPlayerReady: false,
  bothPlayersReady: false,
});

const handleToggle = (isActive: boolean) => {
  if (props.disabled && props.bothPlayersReady) return;
  if (isActive) {
    toggleReady();
  } else {
    toggleUnready();
  }
};

const toggleReady = () => {
  if (userStore.getWebSocketService?.isConnected()) {
    userStore.getWebSocketService?.sendPlayerReadyMessage(props.lobbyId, props.isGameMode);
  } else {
    console.error('WebSocket is not connected');
  }
};

const toggleUnready = () => {
  if (userStore.getWebSocketService?.isConnected()) {
    userStore.getWebSocketService?.sendPlayerUnreadyMessage(props.lobbyId);
  } else {
    console.error('WebSocket is not connected');
  }
};
</script>
<<<<<<< HEAD
=======

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
  background: linear-gradient(to right, var(--secondary-dark-color), color-mix(in srgb, var(--secondary-dark-color) 75%, white));
  color: white;
}

.ready-button.ready {
  background: linear-gradient(to right, var(--secondary-bright-color), color-mix(in srgb, var(--secondary-bright-color) 75%, white));
  color: white;
}

.ready-button.opponent {
  cursor: not-allowed;
  opacity: 0.6;
}

</style>
>>>>>>> main
