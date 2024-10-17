<template>
  <div class="friend-requests-popover">
    <div class="friend-requests-header">
      <h3>Friend Requests</h3>
      <button @click="toggleFriendRequests" class="close-button">&times;</button>
    </div>
    <div class="friend-requests-content">
      <template v-if="friendRequests && friendRequests.length > 0">
        <div v-if="loadingFriendRequests" class="loading-spinner">Loading friend requests...</div>
        <div v-else>
          <FriendRequestItem
            v-for="request in friendRequests"
            :key="request.id"
            :request="request"
            :acceptFriend="acceptFriend"
            :denyFriend="denyFriend"
          />
        </div>
      </template>
      <div v-else>
        <p>No pending friend requests</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import api from '../../../services/api';
import FriendRequestItem from './FriendRequestItem.vue';

interface FriendRequest {
  id: string;
  nickname: string;
}

const props = defineProps<{
  toggleFriendRequests: () => void;
}>();

const friendRequests = ref<FriendRequest[]>([]);
const loadingFriendRequests = ref(false);

const fetchFriendRequests = async () => {
  loadingFriendRequests.value = true;
  try {
    friendRequests.value = await api.friendlist.getFriendRequests();
  } catch (error) {
    console.error('Failed to fetch friend requests', error);
  } finally {
    loadingFriendRequests.value = false;
  }
};

const acceptFriend = async (friendId: string) => {
  try {
    await api.friendlist.acceptFriendRequest(friendId);
    friendRequests.value = friendRequests.value.filter(request => request.id !== friendId);
    await fetchFriendRequests();
  } catch (error) {
    console.error('Failed to accept friend request', error);
  }
};

const denyFriend = async (requestId: string) => {
  try {
    await api.friendlist.denyFriendRequest(requestId);
    friendRequests.value = friendRequests.value.filter(request => request.id !== requestId);
    await fetchFriendRequests();
    props.toggleFriendRequests();
  } catch (error) {
    console.error('Failed to deny friend request', error);
  }
};

onMounted(fetchFriendRequests);
</script>
  
<style scoped>
.friend-requests-popover {
  position: fixed;
  bottom: 70px;
  right: 10px;
  width: 300px;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  z-index: 1000;
  width: calc(100% - 20px);
  max-width: 300px;
}

.friend-requests-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 15px;
  border-bottom: 1px solid #e0e0e0;
  position: relative;
}

.friend-requests-header h3 {
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

.friend-requests-content {
  padding: 15px;
  max-height: 400px;
  overflow-y: auto;
}

.error-message,
.success-message {
  position: relative;
  bottom: 100%;
  left: 0;
  right: 0;
  padding: 8px;
  border-radius: 4px;
  font-size: 14px;
  text-align: center;
  z-index: 1002;
  margin-bottom: 8px;
  max-height: none;
  overflow-wrap: break-word;
  word-wrap: break-word;
}

.error-message {
  background-color: #ffebee;
  color: #d32f2f;
  border: 1px solid #ef9a9a;
}

.success-message {
  background-color: #e8f5e9;
  color: #388e3c;
  border: 1px solid #a5d6a7;
}
</style>