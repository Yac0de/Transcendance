<template>
  <div class="leave-button-container">
    <button class="nav-button leave-lobby-button" @click="leaveLobby">Leave Lobby</button>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useUserStore } from '../../stores/user'

interface Props {
  lobbyId?: string;
}

const router = useRouter();
const userStore = useUserStore()

const props = withDefaults(defineProps<Props>(), {
  lobbyId: '',
})


const leaveLobby = () => {
  const wsService = userStore.getWebSocketService
  if (wsService && props.lobbyId !== '') {
    wsService.leaveAndTerminateLobby(props.lobbyId);
  }
  router.push('/');
}
</script>

<style scoped>
.leave-button-container {
  display: flex;
  align-items: start;
}

.leave-lobby-button {
  background-color: var(--secondary-dark-color);
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
  background-color: color-mix(in srgb, var(--secondary-dark-color) 95%, white);
}

@media (max-width: 600px) {
  .leave-button-container {
    top: 120px;
    /* Additional space for collapsed navbar on mobile */
  }
}
</style>
