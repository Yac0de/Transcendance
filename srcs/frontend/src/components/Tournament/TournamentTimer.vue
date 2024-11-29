<!-- TournamentTimer.vue -->
<template>
  <div class="tournament-timer">
    {{ remainingSeconds }}
  </div>
</template>

<script setup lang="ts">
defineProps<{
  remainingSeconds: number;
}>();
</script>

<style scoped>
.tournament-timer {
  position: absolute;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
}
</style>
