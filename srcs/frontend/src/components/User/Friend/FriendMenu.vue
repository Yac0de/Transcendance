<template>
  <div class="friend-list-container">
    <div class="friend-icon-container">
      <div v-if="showFriendIcon" class="friend-icon" @click="toggleOnFriendMenu">
        <i class="fas fa-user-friends"></i>
      </div>
      <FriendMenuActions :toggleFriendList="toggleFriendList" :toggleAddFriend="toggleAddFriend"
        :toggleFriendRequests="toggleFriendRequests" :toggleOffFriendMenu="toggleOffFriendMenu"
        :showFriendMenu="showFriendMenu" :friendRequests="friendRequests || []" />
    </div>

    <FriendListPopover v-if="showFriendList" :toggleFriendList="toggleFriendList" />

    <AddFriendPopover v-if="showAddFriend" :toggleAddFriend="toggleAddFriend"
      :fetch-friend-requests="fetchFriendRequests" />

    <FriendRequestPopover v-if="showFriendRequests" :toggleFriendRequests="toggleFriendRequests"
      :fetch-friend-requests="fetchFriendRequests" :friendRequests="friendRequests" />
  </div>
</template>


<script setup lang="ts">
import { ref, watch } from 'vue';
import FriendListPopover from './FriendListPopover.vue';
import AddFriendPopover from './AddFriendPopover.vue';
import FriendRequestPopover from './FriendRequestPopover.vue';
import FriendMenuActions from './FriendMenuActions.vue';
import api from '../../../services/api';

interface Friend {
  id: string;
  nickname: string;
  avatar: string;
}

interface FriendRequest {
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
const friendRequests = ref<FriendRequest[]>([]);
const friendsLoaded = ref(false);
const friendRequestsLoaded = ref(false);
const errorMessage = ref('');
const successMessage = ref('');

const loadingFriends = ref(false);
const loadingFriendRequests = ref(false);

watch(showAddFriend, (newValue) => {
  if (!newValue) {
    errorMessage.value = '';
    successMessage.value = '';
    newFriendNickname.value = '';
  }
});

const resetMessages = () => {
  errorMessage.value = '';
  successMessage.value = '';
};

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

const toggleFriendList = async () => {
  showFriendList.value = !showFriendList.value;
  if (showFriendList.value) {
    showAddFriend.value = false;
    showFriendRequests.value = false;
    await fetchFriendList(); // Load only when necessary
  }
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
  if (friendRequestsLoaded.value) return;

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

.friend-icon-container:hover {
  display: flex;
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

@media (max-width: 600px) {

  .friend-icon {
    width: 40px;
    height: 40px;
  }

  .friend-icon i {
    font-size: 20px;
  }
}
</style>
