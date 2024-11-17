<!-- Timer.vue -->
<template>
 <div class="timer">
    {{ remainingSeconds }} 
 </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const remainingTime = ref<number>(60);

const props = defineProps<{
  remainingSeconds: number;
}>();

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
