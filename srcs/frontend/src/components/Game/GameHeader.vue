<template>
    <div class="game-header">
      <div class="score">
        Player 1: {{ state.score.player1 }} - Player 2: {{ state.score.player2 }}
      </div>
      <div class="timer">Time Left: {{ state.remainingTime }}s</div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { PropType } from 'vue';
  import { GameState } from '../../types/game';
  
  const props = defineProps({
    state: {
      type: Object as PropType<GameState>,
      required: true,
    },
  });
  </script>
  
  <style scoped>
  .game-header {
    display: flex;
    justify-content: space-between;
    padding: 10px;
    background-color: #333;
    color: #fff;
  }
  </style>
  