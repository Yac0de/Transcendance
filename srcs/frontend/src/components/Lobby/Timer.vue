<!-- Timer.vue -->
<template>
 <div class="timer">
    {{ remainingTime }} 
 </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { eventBus } from '../../events/eventBus'

const remainingTime = ref<number>(60);

onMounted(() => {
  eventBus.on('LOBBY_PREGAME_REMAINING_TIME', (message: LobbyPregameRemainingTime) => {
    console.log("SEC REMAINING = ", message.remainingTime);
    remainingTime.value = message.remainingTime;
  })
})

onUnmounted(() => {
  eventBus.off('LOBBY_PREGAME_REMAINING_TIME')
})
</script>

<style scoped>
.timer {
 position: absolute;
 top: 80px; /* Adjust this value to position above VS */
 left: 50%;
 transform: translateX(-50%);
 font-size: 24px;
 font-weight: bold;
 color: black;
}
</style>
