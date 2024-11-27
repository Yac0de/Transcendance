<template>
  <div class="friend-item">
    <div class="friend-avatar" v-if="friend" :class="{ 'friend-online': isOnline }">
      <img :src="api.user.getAvatarUrl(friend.avatar)" :alt="friend.nickname + '\'s avatar'" />
    </div>
    <div class="friend-info" v-if="friend">
      <router-link :to="`/${friend.nickname}`" class="friend-nickname">{{ friend.nickname }}</router-link>
      <div class="friend-actions">
        <button class="friend-action-btn delete-btn" @click="deleteFriend">
          <i class="fas fa-trash-alt"></i>
          <span v-if="loadingDeleteFriend" class="loading-spinner"></span>
        </button>
        <button class="friend-action-btn" @click="openChat">
          <i class="fas fa-comment"></i>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import api from '../../../services/api';
import { useOnlineUsersStore } from '../../../stores/onlineUsers';
import { useChatStore } from '../../../stores/chatStore';
import { Friend } from '../../../types/models';

const OnlineUsersStore = useOnlineUsersStore();
const chatStore = useChatStore();

const props = defineProps<{
  friend: Friend | null;
  deleteFriendFromList: (id: number) => Promise<void>;
}>();

const isOnline = computed(() => {
  if (!props.friend) {
    return false;
  }
  return OnlineUsersStore.isUserOnline(props.friend.id);
});

const loadingDeleteFriend = ref(false);

const deleteFriend = async () => {
  if (!props.friend) return;
  loadingDeleteFriend.value = true;
  try {
    await props.deleteFriendFromList(props.friend.id);
  } catch (error) {
    console.error('Failed to delete friend', error);
  } finally {
    loadingDeleteFriend.value = false;
  }
};

const openChat = () => {
  if (props.friend) {
    chatStore.toggleFriend(props.friend.id);
  }
};
</script>

<style scoped>
.friend-item {
  display: flex;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid var(--main-extra-color);
}

.friend-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  margin-right: 10px;
  border: 2px solid #ff0000;
}

.friend-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.friend-online {
  border: 2px solid #2ecc71;
}

;

.friend-info {
  flex-grow: 1;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.friend-nickname {
  font-weight: bold;
  text-decoration: none;
  color: black;
}

.friend-actions {
  display: flex;
  gap: 10px;
}

.friend-action-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 5px;
  font-size: 16px;
  color: #666;
  transition: color 0.3s;
}

.friend-action-btn:hover {
  color: #007bff;
}

.friend-action-btn.delete-btn:hover {
  color: #dc3545;
}
</style>
