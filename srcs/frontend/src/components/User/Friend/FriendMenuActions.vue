<template>
  <div v-if="showFriendMenu">
    <button @click="toggleFriendList" class="icon-button list-button">
      <i class="fas fa-list"></i>
    </button>
    <button @click="toggleAddFriend" class="icon-button add-friend-button">
      <i class="fas fa-user-plus"></i>
    </button>
    <button @click="toggleFriendRequests" class="icon-button friend-requests-button">
      <i class="fas fa-bell"></i>
      <span v-if="friendRequests && friendRequests.length > 0" class="notification-badge">
        {{ friendRequests.length }}
      </span>
    </button>
    <button @click="toggleOffFriendMenu" class="icon-button xmark-button">
      <i class="fas fa-times"></i>
    </button>
  </div>
</template>

<script setup lang="ts">
import { Friend } from '../../../types/models';

const {
  toggleFriendList,
  toggleAddFriend,
  toggleFriendRequests,
  toggleOffFriendMenu,
  showFriendMenu,
  friendRequests
} = defineProps<{
  toggleFriendList: () => void;
  toggleAddFriend: () => void;
  toggleFriendRequests: () => void;
  toggleOffFriendMenu: () => void;
  showFriendMenu: boolean;
  friendRequests: Friend[] | null;
}>();
</script>

<style scoped>
.icon-button {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  border: none;
  background-color: var(--main-extra-color);
  color: white;
  font-size: 20px;
  cursor: pointer;
  margin-left: 10px;
  transition: background-color 0.3s;
  position: relative;
}

.icon-button:hover {
  background-color: color-mix(in srgb, var(--main-extra-color) 85%, white);
}

.icon-button.xmark-button {
  background-color: var(--secondary-extra-color);
  color: white;
}

.icon-button.xmark-button:hover {
  background-color: color-mix(in srgb, var(--secondary-extra-color) 95%, white);
}

.notification-badge {
  position: absolute;
  top: -5px;
  right: -5px;
  background-color: #dc3545;
  color: white;
  border-radius: 50%;
  padding: 2px 6px;
  font-size: 12px;
}
</style>
