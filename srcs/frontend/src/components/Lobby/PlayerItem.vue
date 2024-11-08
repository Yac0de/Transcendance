<template>
  <div class="player-tile" :class="{ 'player-left': isLeft }">
    <div class="player-name">{{ playerName }}</div>
  </div>
</template>

<script>
export default {
  name: 'PlayerItem',
  props: {
    playerName: {
      type: String,
      required: true
    },
    isLeft: {
      type: Boolean,
      default: true
    }
  }
}
</script>

<style scoped>
.player-tile {
  background-color: #f8f9fa;
  border-radius: 8px;
  padding: 20px;
  width: 200px;
  height: 150px;
  display: flex;
  justify-content: center;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  margin: 10px;
}

.player-name {
  font-size: 24px;
  font-weight: bold;
  color: #2c3e50;
}
</style>
