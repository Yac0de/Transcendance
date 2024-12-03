<template>
  <div class="player-tile">
    <div class="player-content">
      <div class="avatar-container">
        <div class="avatar-wrapper">
          <img
            :src="isLeft ? api.user.getAvatarUrl(user_store.getAvatarPath) : api.user.getAvatarUrl(challengedFriend?.avatar ?? null)"
            :alt="(isLeft ? user_store.nickname : challengedFriend?.nickname) + '\'s avatar'" class="avatar-image" />
        </div>
      </div>
      <div class="player-name" :style="nameStyle">
        {{ isLeft ? user_store.nickname : challengedFriend?.nickname }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useUserStore } from '../../stores/user';
import api from '../../services/api';
import type { UserData } from '../../types/models';

const props = defineProps<{
  isLeft: boolean;
  challengedFriend?: UserData | null;
}>();

const user_store = useUserStore();

const nameStyle = computed(() => {
  const nickname = (props.isLeft ? user_store.nickname : props.challengedFriend?.nickname) || '';
  const minFontSize = 11;
  const maxFontSize = 24;
  const minLength = 3;

  const length = nickname.length;
  const fontSize = Math.max(minFontSize, Math.min(maxFontSize, maxFontSize - (length - minLength) * 1.5));

  return {
    fontSize: `${fontSize}px`,
  };
});
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
  padding: 40px 10px;
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
  font-weight: bold;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
