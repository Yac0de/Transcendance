<template>
  <div class="player-tile" :class="{ empty: !player }">
    <div class="player-content">
      <div v-if="player" class="player-info">
        <div class="player-avatar">
          <img 
            v-if="player.avatar"
            :src="api.user.getAvatarUrl(player.avatar)"
            :alt="`${player.displayname}'s avatar`"
            class="avatar-image"
          />
          <span v-else>
            {{ player.displayname.charAt(0).toUpperCase() }}
          </span>
        </div>
        <div class="player-details">
          <span class="player-name">{{ player.displayname }}</span>
        </div>
      </div>
      <div v-else class="waiting-text">
        {{ $t('waitingForPlayer') }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { UserData } from '../../types/models';
import api from '../../services/api';

defineProps<{
  player?: UserData | null;
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
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
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
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.2rem;
  font-weight: bold;
  overflow: hidden;
  background: linear-gradient(to right, var(--secondary-dark-color), 
    color-mix(in srgb, var(--secondary-dark-color) 75%, white));
  color: white;
}

.avatar-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
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

.waiting-text {
  color: #666;
  font-style: italic;
}
</style>
