<template>
  <div class="player-tile">
    <div class="player-content">
      <div class="avatar-container">
        <div class="avatar-wrapper">
          <img :src="api.user.getAvatarUrl(user.getAvatarPath)" :alt="user.getNickname + '\'s avatar'"
            class="avatar-image" />
        </div>
      </div>
      <div class="player-name">{{ userName }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useUserStore } from '../../stores/user';
import { ref } from 'vue';
import api from '../../services/api';

const user = useUserStore();
const userName: string = user.getNickname;
</script>

<style scoped>
.player-tile {
  background-color: #f8f9fa;
  border-radius: 8px;
  padding: 20px;
  width: 250px;
  height: 60px;
  display: flex;
  justify-content: center;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  margin: 10px;
}

.player-content {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 20px;
  /* Space between avatar and name */
  width: 100%;
}

.avatar-container {
  text-align: center;
  margin-bottom: 0;
  /* Removed margin-bottom since we're horizontal now */
}

.avatar-wrapper {
  position: relative;
  display: inline-block;
  height: 4rem;
  width: 4rem;
}

.avatar-image {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
}

.player-name {
  font-size: 24px;
  font-weight: bold;
  color: #2c3e50;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
