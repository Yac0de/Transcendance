<template>
  <div class="friend-list-popover">
    <div class="friend-list-header">
      <h3>Friends</h3>
      <button @click="toggleFriendList" class="close-button">&times;</button>
    </div>
    <div class="friend-list-content">
      <div v-if="loadingFriends" class="loading-spinner">Loading friends...</div>
      <div v-else>
        <div v-if="friends.length === 0" class="no-friends-message">
          <p>You have no friends yet 😢</p>
        </div>
        <div v-else>
          <FriendItem v-for="friend in friends" :key="friend.id" :friend="friend"
            :deleteFriendFromList="deleteFriendFromList" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import api from '../../../services/api';
import FriendItem from './FriendItem.vue';
import { Friend } from '../../../types/models';

const { toggleFriendList } = defineProps<{
  toggleFriendList: () => void;
}>();

const friends = ref<Friend[]>([]);
const loadingFriends = ref(false);

const fetchFriendList = async () => {
  loadingFriends.value = true;
  try {
    const fetchedFriends = await api.friendlist.getFriendList();
    if (fetchedFriends) {
      friends.value = fetchedFriends;
    }
  } catch (error) {
    console.error('Failed to fetch friend list', error);
  } finally {
    loadingFriends.value = false;
  }
};

const deleteFriendFromList = async (friendId: number) => {
  try {
    await api.friendlist.deleteFromFriendList(friendId);
    friends.value = friends.value.filter((friend: Friend) => friend.id !== friendId);
  } catch (error) {
    console.error('Failed to delete friend', error);
  }
};

onMounted(fetchFriendList);
</script>

<style scoped>
.friend-list-popover {
  position: fixed;
  bottom: 80px;
  right: 20px;
  width: 300px;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  z-index: 1000;
}

.friend-list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 15px;
  border-bottom: 1px solid #e0e0e0;
  position: relative;
}

.friend-list-header h3 {
  margin: 0;
  font-size: 18px;
}

.close-button {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #666;
}

.friend-list-content {
  padding: 15px;
  max-height: 400px;
  overflow-y: auto;
}
</style>
