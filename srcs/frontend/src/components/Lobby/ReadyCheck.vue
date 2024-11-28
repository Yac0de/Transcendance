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
    userStore.getWebSocketService?.sendPlayerReadyMessage(props.lobbyId);
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
