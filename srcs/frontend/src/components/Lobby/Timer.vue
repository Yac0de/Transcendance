<template>
 <div class="timer">
    {{ remainingSeconds }} 
 </div>
</template>

<script setup lang="ts">

defineProps<{
  remainingSeconds: number;
}>();

</script>

<style scoped>
.timer {
 position: absolute;
 top: 80px;
 left: 50%;
 transform: translateX(-50%);
 font-size: 24px;
 font-weight: bold;
 color: black;
}
</style>
