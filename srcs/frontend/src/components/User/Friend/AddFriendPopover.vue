<template>
  <div class="add-friend-popover">
    <div class="add-friend-header"> 
      <h3>Add Friend</h3>
      <button @click="toggleAddFriend" class="close-button">&times;</button>
    </div>
    <div class="add-friend-content">
      <input v-model="newFriendNickname" type="text" placeholder="Enter friend's name" class="friend-input" />
      <div v-if="errorMessage" class="error-message">{{ errorMessage }}</div>
      <div v-if="successMessage" class="success-message">{{ successMessage }}</div>
      <button @click="addFriend" class="add-friend-btn" :disabled="loadingAddFriend">
        <span v-if="loadingAddFriend">Adding...</span>
        <span v-else>Add Friend</span>
      </button>
    </div>
  </div>
</template>
    
<script setup lang="ts">
import { ref } from 'vue';
import api from '../../../services/api';

const newFriendNickname = ref('');
const errorMessage = ref('');
const successMessage = ref('');
const loadingAddFriend = ref(false);

const props = defineProps<{
toggleAddFriend: (event: MouseEvent) => void;
fetchFriendRequests: () => void;
}>();

const resetMessages = () => {
  errorMessage.value = '';
  successMessage.value = '';
};

const addFriend = async () => {
  resetMessages();

  if (newFriendNickname.value.trim() === "") {
    errorMessage.value = "Please enter a friend's name.";
    return;
  }

  if (newFriendNickname.value.length < 3) {
    errorMessage.value = "Friend's nickname must be at least 3 characters long.";
    return;
  }

  loadingAddFriend.value = true;
  try {
    await api.friendlist.sendFriendRequest(newFriendNickname.value.trim());
    newFriendNickname.value = '';
    await props.fetchFriendRequests();
    successMessage.value = 'Friend request sent successfully!';
    errorMessage.value = '';
  } catch (error: any) {
    handleFriendRequestError(error);
  } finally {
    loadingAddFriend.value = false;
  }
};

const handleFriendRequestError = (error: any) => {
  if (error.status) {
    switch (error.status) {
      case 409:
        errorMessage.value = "This friendship already exists. You already sent a friend request to this user, or you have a pending request from them.";
        break;
      case 404:
        errorMessage.value = "The user with this nickname does not exist.";
        break;
      case 400:
        errorMessage.value = "Invalid request. Please check the nickname.";
        break;
      default:
        errorMessage.value = "An unexpected error occurred. Please try again.";
    }
  } else {
    errorMessage.value = "Network error or server is unreachable.";
  }
};
</script>
  
<style scoped>
.add-friend-popover {
  position: fixed;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  z-index: 1000;
  bottom: 70px;
  right: 10px;
  width: calc(100% - 20px);
  max-width: 300px;
}

.add-friend-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 15px;
  border-bottom: 1px solid #e0e0e0;
  position: relative;
}

.add-friend-header h3 {
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

.add-friend-content {
  padding: 15px;
  max-height: 400px;
  overflow-y: auto;
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
  