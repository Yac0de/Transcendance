<template>
  <div class="friend-list-container">
    <div class="friend-icon-container">
      <div class="friend-icon">
        <i class="fas fa-user-friends"></i>
      </div>
      <div class="friend-icon-hover">
        <button @click="toggleFriendList" class="icon-button list-button">
          <i class="fas fa-list"></i>
        </button>
        <button @click="toggleAddFriend" class="icon-button add-friend-button">
          <i class="fas fa-user-plus"></i>
        </button>
        <button @click="toggleFriendRequests" class="icon-button friend-requests-button">
          <i class="fas fa-bell"></i>
          <span v-if="friendRequests.length" class="notification-badge">{{ friendRequests.length }}</span>
        </button>
      </div>
    </div>
    <div v-if="showFriendList" class="friend-list-popover">
      <div class="friend-list-header">
        <h3>Friends</h3>
        <button @click="toggleFriendList" class="close-button">&times;</button>
      </div>
      <div class="friend-list-content">
        <div v-for="friend in friends" :key="friend.id" class="friend-item">
          <div class="friend-avatar" :style="{ backgroundColor: friend.color }">
          </div>
          <div class="friend-info">
            <div class="friend-name">{{ friend.name }}</div>
            <div class="friend-status" :class="friend.status.toLowerCase()">
              {{ friend.status }}
            </div>
          </div>
        </div>
      </div>
    </div>
    <div v-if="showAddFriend" class="add-friend-popover">
      <div class="add-friend-header">
        <h3>Add Friend</h3>
        <button @click="toggleAddFriend" class="close-button">&times;</button>
      </div>
      <div class="add-friend-content">
        <input v-model="newFriendNickname" type="text" placeholder="Enter friend's name" class="friend-input" />
        <button @click="addFriend" class="add-friend-btn">Add Friend</button>
      </div>
    </div>
    <div v-if="showFriendRequests" class="friend-requests-popover">
      <div class="friend-requests-header">
        <h3>Friend Requests</h3>
        <button @click="toggleFriendRequests" class="close-button">&times;</button>
      </div>
      <div class="friend-requests-content">
        <div v-for="request in friendRequests" :key="request.id" class="friend-request-item">
          <div class="friend-request-name">{{ request.name }}</div>
          <div class="friend-request-actions">
            <button @click="acceptFriendRequest(request.id)" class="accept-button">
              <i class="fas fa-check"></i>
            </button>
            <button @click="denyFriendRequest(request.id)" class="deny-button">
              <i class="fas fa-times"></i>
            </button>
          </div>
        </div>
        <div v-if="friendRequests.length === 0" class="no-requests">
          No pending friend requests
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import friendlistAPI from '../services/friendlist';

const showFriendList = ref(false);
const showAddFriend = ref(false);
const newFriendNickname = ref('');
const friends = ref<{ nickname: string }[]>([]);
const loading = ref(false);
const showFriendRequests = ref(false);


const friendRequests = ref([
  { id: 1, name: 'Faris' },
  { id: 2, name: 'Farisa' },
]);

const toggleFriendList = async () => {
  showFriendList.value = !showFriendList.value;
  if (showFriendList.value) {
    showAddFriend.value = false;
    showFriendRequests.value = false;
    await fetchFriends();
  }
};

const toggleAddFriend = () => {
  showAddFriend.value = !showAddFriend.value;
  if (showAddFriend.value) {
    showFriendList.value = false;
    showFriendRequests.value = false;
  }
};

const toggleFriendRequests = () => {
  showFriendRequests.value = !showFriendRequests.value;
  if (showFriendRequests.value) {
    showFriendList.value = false;
    showAddFriend.value = false;
  }
};

const fetchFriends = async () => {
  loading.value = true;
  try {
    friends.value = await friendlistAPI.getFriendList();
  } catch (error) {
    console.error('Failed to fetch friends:', error);
  } finally {
    loading.value = false;
  }
};

const addFriend = async () => {
  if (newFriendNickname.value.trim()) {
    try {
      await friendlistAPI.sendFriendRequest(newFriendNickname.value.trim());
      newFriendNickname.value = '';
      showAddFriend.value = false;
      await fetchFriends();
    } catch (error) {
      console.error('Failed to add friend:', error);
    }
  }
};

const getRandomColor = (nickname: string) => {
  let hash = 0;
  for (let i = 0; i < nickname.length; i++) {
    hash = nickname.charCodeAt(i) + ((hash << 5) - hash);
  }
  return `hsl(${hash % 360}, 70%, 60%)`;
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
  display: none;
  position: absolute;
  bottom: 0;
  right: 0;
}

.friend-icon-container:hover .friend-icon {
  display: none;
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
  display: flex;
  justify-content: center;
  align-items: center;
  font-weight: bold;
  color: white;
  margin-right: 10px;
}

.friend-info,
.friend-request-name {
  flex-grow: 1;
}

.friend-name {
  font-weight: bold;
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

.friend-request-actions {
  display: flex;
  gap: 10px;
}

.accept-button,
.deny-button {
  padding: 5px 10px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.3s;
}

.accept-button {
  background-color: #28a745;
  color: white;
}

.accept-button:hover {
  background-color: #218838;
}

.deny-button {
  background-color: #dc3545;
  color: white;
}

.deny-button:hover {
  background-color: #c82333;
}

.no-requests {
  text-align: center;
  color: #666;
  padding: 20px 0;
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
