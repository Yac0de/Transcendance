<!-- PlayerTile.vue -->
<template>
  <div class="player-tile" :class="{ empty: !player }">
    <div class="player-content">
      <div v-if="player" class="player-info">
        <div class="player-avatar">
          {{ player.displayname.charAt(0).toUpperCase() }}
        </div>
        <div class="player-details">
          <span class="player-name">{{ player.displayname }}</span>
        </div>
      </div>
      <div v-else class="waiting-text">
        Waiting for player...
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Player {
  username: string;
  status: 'ready' | 'not-ready';
}

defineProps<{
  player?: Player;
}>();
</script>

<style scoped>
.player-tile {
  width: 200px;
  height: 100px;
  background-color: #ffffff;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  padding: 1rem;
  transition: all 0.2s;
}

.player-tile.empty {
  background-color: #f5f5f5;
  border-style: dashed;
}

.player-content {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.player-info {
  display: flex;
  align-items: center;
  gap: 1rem;
  width: 100%;
}

.player-avatar {
  width: 40px;
  height: 40px;
  background-color: #2196F3;
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.2rem;
  font-weight: bold;
}

.player-details {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.player-name {
  font-weight: bold;
  color: #333;
}

.player-status {
  font-size: 0.9rem;
}

.player-status.ready {
  color: #4CAF50;
}

.player-status.not-ready {
  color: #f44336;
}

.waiting-text {
  color: #666;
  font-style: italic;
}
</style>

