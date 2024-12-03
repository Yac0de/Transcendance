<template>
  <div class="friend-requests-popover">
    <div class="friend-requests-header">
      <h3>{{ $t('friendRequestsTitle') }}</h3>
      <button @click="toggleFriendRequests" class="close-button">&times;</button>
    </div>
    <div class="friend-requests-content">
      <div v-if="friendRequests && friendRequests.length > 0">
        <div v-if="loadingFriendRequests" class="loading-spinner">{{ $t('loadingFriendRequests') }}</div>
        <div v-else>
          <div v-for="request in friendRequests" :key="request.id" class="friend-request-item">
            <div class="friend-avatar">
              <img :src="api.user.getAvatarUrl(request.avatar)" :alt="request.nickname + '\'s avatar'" />
            </div>
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
      </div>
      <div v-else>
        <p>{{ $t('noPendingRequests') }}</p>
      </div>
      <div v-if="successMessage" class="success-message">{{ successMessage }}</div>
      <div v-if="errorMessage" class="error-message">{{ errorMessage }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import api from '../../../services/api';
import { Friend } from '../../../types/models';

const { toggleFriendRequests, fetchFriendRequests, friendRequests } = defineProps<{
  toggleFriendRequests: () => void;
  fetchFriendRequests: () => void;
  friendRequests: Friend[] | null;
}>();

const loadingFriendRequests = ref(false);
const successMessage = ref('');
const errorMessage = ref('');

const { t } = useI18n();

const acceptFriend = async (friendId: number) => {
  try {
    await api.friendlist.acceptFriendRequest(friendId);
    friendRequests?.splice(friendRequests.findIndex(req => req.id === friendId), 1);
    successMessage.value = t('acceptRequest');
    fetchFriendRequests();
  } catch (error) {
    errorMessage.value = t('errorAcceptRequest');
  }
};

const denyFriend = async (requestId: number) => {
  try {
    await api.friendlist.denyFriendRequest(requestId);
    friendRequests?.splice(friendRequests.findIndex(req => req.id === requestId), 1);
    successMessage.value = t('denyRequest');
    fetchFriendRequests();
  } catch (error) {
    errorMessage.value = t('errorDenyRequest');
  }
};

onMounted(fetchFriendRequests);
</script>



<style scoped>
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

.friend-requests-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 15px;
  border-bottom: 1px solid var(--main-extra-color);
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

.friend-request-item {
  display: flex;
  flex-grow: 1;
  justify-content: space-between;
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
}

.friend-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.friend-request-name {
  font-weight: bold;
}

.friend-request-actions {
  display: flex;
  gap: 10px;
}

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

.accept-button:hover,
.deny-button:hover {
  color: #007bff;
}

.deny-button:hover {
  color: #dc3545;
}

.accept-button:hover {
  color: #28a745;
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
  .friend-requests-popover {
    bottom: 70px;
    right: 10px;
    width: calc(100% - 20px);
    max-width: 300px;
  }
}
</style>
