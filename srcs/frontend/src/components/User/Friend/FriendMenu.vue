<template>
  <div class="friend-list-container">
    <div class="friend-icon-container">
      <div v-if="showFriendIcon" class="friend-icon" @click="toggleOnFriendMenu">
        <i class="fas fa-user-friends"></i>
      </div>
      <FriendMenuActions :toggleFriendList="toggleFriendList" :toggleAddFriend="toggleAddFriend"
        :toggleFriendRequests="toggleFriendRequests" :toggleOffFriendMenu="toggleOffFriendMenu"
        :showFriendMenu="showFriendMenu" :friendRequests="friendRequests" />
    </div>

    <FriendListPopover v-if="showFriendList" :toggleFriendList="toggleFriendList" />

    <AddFriendPopover v-if="showAddFriend" :toggleAddFriend="toggleAddFriend"
      :fetch-friend-requests="fetchFriendRequests" />

    <FriendRequestsPopover v-if="showFriendRequests" :toggleFriendRequests="toggleFriendRequests"
      :fetchFriendList="fetchFriendList" />

    <div v-if="showFriendRequests" class="friend-requests-popover">
      <div class="friend-requests-header">
        <h3>Friend Requests</h3>
        <button @click="toggleFriendRequests" class="close-button">&times;</button>
      </div>
      <div class="friend-requests-content">
        <template v-if="friendRequests && friendRequests.length > 0">
          <div v-if="loadingFriendRequests" class="loading-spinner">Loading friend requests...</div>
          <div v-else>
            <div v-for="request in friendRequests" :key="request.id" class="friend-request-item">
              <div class="friend-request-name">{{ request.nickname }}</div>
              <div class="friend-request-actions">
                <button @click="acceptFriend(request.id)" class="accept-button">
                  <i class="fas fa-check"></i>
                </button>
                <button @click="denyFriend(request.id)" class="deny-button">
                  <i class="fas fa-times"></i>
                </button>
              </div>
            </div>
          </div>
        </template>
        <div v-else>
          <p>No pending friend requests</p>
        </div>
        <div v-if="successMessage" class="success-message">{{ successMessage }}</div>
        <div v-if="errorMessage" class="error-message">{{ errorMessage }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import api from '../../../services/api'
import FriendListPopover from './FriendListPopover.vue';
import AddFriendPopover from './AddFriendPopover.vue';
import FriendRequestsPopover from './FriendRequestsPopover.vue'
import FriendMenuActions from './FriendMenuActions.vue'

interface Friend {
  id: string;
  nickname: string;
  avatar: string;
}

const showFriendIcon = ref(true);
const showFriendMenu = ref(false);
const showFriendList = ref(false);
const showAddFriend = ref(false);
const showFriendRequests = ref(false);
const newFriendNickname = ref('');
const friends = ref<Friend[]>([]);
const friendRequests = ref<Friend[]>([]);
const friendsLoaded = ref(false);
const friendRequestsLoaded = ref(false);
const errorMessage = ref('');
const successMessage = ref('');

const loadingFriends = ref(false);
const loadingFriendRequests = ref(false);
const loadingAddFriend = ref(false);
const loadingDeleteFriend = ref<string | null>(null);

const toggleOnFriendMenu = async () => {
  await fetchFriendRequests();
  showFriendMenu.value = !showFriendMenu.value;
  showFriendIcon.value = false;
}

const toggleOffFriendMenu = () => {
  showFriendMenu.value = false;
  showFriendIcon.value = true;
  showFriendList.value = false;
  showAddFriend.value = false;
  showFriendRequests.value = false;

  // Reset states to force reloading next time
  friendsLoaded.value = false;
  friendRequestsLoaded.value = false;
};

watch(showAddFriend, (newValue) => {
  if (!newValue) {
    errorMessage.value = '';
    successMessage.value = '';
    newFriendNickname.value = '';
  }
});

const toggleFriendList = async () => {
  showFriendList.value = !showFriendList.value;
  if (showFriendList.value) {
    showAddFriend.value = false;
    showFriendRequests.value = false;
    await fetchFriendList(); // Load only when necessary
  }
};

const resetMessages = () => {
  errorMessage.value = '';
  successMessage.value = '';
};

const toggleFriendRequests = async () => {
  resetMessages();
  showFriendRequests.value = !showFriendRequests.value;
  if (showFriendRequests.value) {
    showFriendList.value = false;
    showAddFriend.value = false;
    await fetchFriendRequests();
  }
};

const toggleAddFriend = () => {
  resetMessages();
  showAddFriend.value = !showAddFriend.value;
  if (showAddFriend.value) {
    showFriendList.value = false;
    showFriendRequests.value = false;
  }
};


const fetchFriendList = async () => {
  loadingFriends.value = true;
  try {
    const fetchedFriends = await api.friendlist.getFriendList();
    if (fetchedFriends) {
      friends.value = fetchedFriends;
      friendsLoaded.value = true;
    }
  } catch (error) {
    console.error('Failed to fetch friend list:', error);
    errorMessage.value = 'Failed to fetch friend list.';
  } finally {
    loadingFriends.value = false;
  }
};

const fetchFriendRequests = async () => {
  if (friendRequestsLoaded.value) return; // Do not recharge if already charged

  loadingFriendRequests.value = true;
  try {
    friendRequests.value = await api.friendlist.getFriendRequests();
    friendRequestsLoaded.value = true; // Mark requests as loaded
  } catch (error) {
    console.error('Failed to fetch friend requests:', error);
  } finally {
    loadingFriendRequests.value = false;
  }
};

const acceptFriend = async (friendId: string) => {
  try {
    await api.friendlist.acceptFriendRequest(friendId);
    friendRequests.value = friendRequests.value.filter(request => request.id !== friendId);
    await fetchFriendList();
    await fetchFriendRequests();
    successMessage.value = 'Friend request accepted!';
  } catch (error) {
    errorMessage.value = 'Failed to accept friend request';
    console.error('Failed to accept friend request:', error);
  }
};

const denyFriend = async (requestId: string) => {
  try {
    await api.friendlist.denyFriendRequest(requestId);

    friendRequests.value = friendRequests.value.filter(request => request.id !== requestId);

    successMessage.value = 'Friend request denied!';
    await fetchFriendRequests();
    toggleFriendRequests();
  } catch (error) {
    errorMessage.value = 'Failed to deny friend request';
    console.error('Failed to deny friend request:', error);
  }
};

</script>

<style scoped>
.friend-list-container {
  position: fixed;
  bottom: 20px;
  right: 20px;
  z-index: 1001;
}

.friend-icon-container {
  position: relative;
}

.friend-icon {
  width: 50px;
  height: 50px;
  background-color: #007bff;
  border-radius: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: background-color 0.3s;
}

.friend-icon i {
  color: white;
  font-size: 24px;
}

.friend-icon-hover {
  position: absolute;
  bottom: 0;
  right: 0;
}

.friend-icon-container:hover .friend-icon-hover {
  display: flex;
}

.icon-button {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  border: none;
  background-color: #007bff;
  color: white;
  font-size: 20px;
  cursor: pointer;
  margin-left: 10px;
  transition: background-color 0.3s;
  position: relative;
}

.icon-button:hover {
  background-color: #0056b3;
}

.icon-button.xmark-button {
  background-color: red;
  color: white;
}

.icon-button.xmark-button:hover {
  background-color: darkred;
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

.friend-list-popover,
.add-friend-popover,
.friend-requests-popover {
  position: fixed;
  bottom: 80px;
  right: 20px;
  width: 300px;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  z-index: 1000;
}

.friend-list-header,
.add-friend-header,
.friend-requests-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 15px;
  border-bottom: 1px solid #e0e0e0;
  position: relative;
}

.friend-list-header h3,
.add-friend-header h3,
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

.friend-list-content,
.add-friend-content,
.friend-requests-content {
  padding: 15px;
  max-height: 400px;
  overflow-y: auto;
}

.friend-item,
.friend-request-item {
  display: flex;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #e0e0e0;
}

.friend-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  margin-right: 10px;
  flex-shrink: 0;
}

.friend-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.friend-info,
.friend-request-item {
  flex-grow: 1;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.friend-nickname,
.friend-request-name {
  font-weight: bold;
}

.friend-actions,
.friend-request-actions {
  display: flex;
  gap: 10px;
}

.friend-action-btn,
.accept-button,
.deny-button {
  background: none;
  border: none;
  cursor: pointer;
  padding: 5px;
  font-size: 16px;
  color: #666;
  transition: color 0.3s, background-color 0.3s;
}

.friend-action-btn:hover,
.accept-button:hover,
.deny-button:hover {
  color: #007bff;
}

.friend-action-btn.delete-btn:hover,
.deny-button:hover {
  color: #dc3545;
}

.accept-button:hover {
  color: #28a745;
}

.friend-status {
  font-size: 12px;
  padding: 2px 6px;
  border-radius: 10px;
  display: inline-block;
  margin-top: 4px;
}

.friend-status.online {
  background-color: #4CAF50;
  color: white;
}

.friend-status.offline {
  background-color: #F44336;
  color: white;
}

.friend-status.away {
  background-color: #FFC107;
  color: black;
}

.friend-input {
  width: 100%;
  padding: 8px;
  margin-bottom: 10px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 14px;
}

.add-friend-btn {
  width: 100%;
  padding: 8px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.3s;
}

.add-friend-btn:hover {
  background-color: #0056b3;
}

.no-requests {
  text-align: center;
  color: #666;
  padding: 20px 0;
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

@media (max-width: 600px) {

  .friend-icon,
  .icon-button {
    width: 40px;
    height: 40px;
  }

  .friend-icon i,
  .icon-button i {
    font-size: 20px;
  }

  .friend-list-popover,
  .add-friend-popover,
  .friend-requests-popover {
    bottom: 70px;
    right: 10px;
    width: calc(100% - 20px);
    max-width: 300px;
  }
}
</style>
