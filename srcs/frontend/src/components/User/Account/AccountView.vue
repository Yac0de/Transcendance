<template>
  <div class="account-view">
    <div class="avatar-container">
      <div class="avatar-wrapper">
        <img :src="avatarUrl" alt="User Avatar" class="avatar-image" />
      </div>
    </div>
    <div class="account-info">
      <p><strong>Nickname:</strong> {{ user.nickname }}</p>
      <p><strong>Display Name:</strong> {{ user.displayname }}</p>
    </div>
    <div v-if="isOwnProfile" class="account-actions">
      <button class="action-button" @click="$emit('startEditing')">Edit Profile</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import api from '../../../services/api';
import { UserData } from '../../../types/models';

interface Props {
  user: UserData;
  isOwnProfile: boolean;
}

const props = defineProps<Props>();
const emit = defineEmits(['startEditing']);

const avatarUrl = computed(() => {
  return props.user.avatar ? api.user.getAvatarUrl(props.user.avatar) : api.user.getAvatarUrl('default.png');
});

</script>

<style scoped>
.account-view {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.avatar-container {
  text-align: center;
  margin-bottom: 20px;
}

.avatar-wrapper {
  position: relative;
  display: inline-block;
  height: 8rem;
  width: 8rem;
}

.avatar-image {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
}

.account-info {
  width: 100%;
  margin-bottom: 20px;
}

.account-info p {
  color: white;
  font-size: 16px;
  margin-bottom: 10px;
}

.account-info p strong {
  text-shadow: 1px 1px 2px black;
}

.account-actions {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.action-button {
  width: 100%;
  padding: 10px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  color: white;
  font-size: 14px;
  transition: all 0.3s;
  background: linear-gradient(to right, var(--secondary-dark-color), color-mix(in srgb, var(--secondary-dark-color) 75%, white));
}

.action-button:hover {
  background: linear-gradient(to right, var(--secondary-dark-color), color-mix(in srgb, var(--secondary-dark-color) 85%, white));
  transform: scale(1.02);
}
</style>
