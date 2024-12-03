<template>
  <div class="add-friend-popover">
    <div class="add-friend-header">
      <h3>{{ $t('addFriendTitle') }}</h3>
      <button @click="toggleAddFriend" class="close-button">&times;</button>
    </div>
    <div class="add-friend-content">
      <input 
        v-model="newFriendNickname" 
        type="text" 
        :placeholder="$t('enterFriendNamePlaceholder')" 
        class="friend-input" 
      />
      <div v-if="errorMessage" class="error-message">{{ errorMessage }}</div>
      <div v-if="successMessage" class="success-message">{{ successMessage }}</div>
      <button @click="addFriend" class="add-friend-btn" :disabled="loadingAddFriend">
        <span v-if="loadingAddFriend">{{ $t('addingFriend') }}</span>
        <span v-else>{{ $t('addFriendButton') }}</span>
      </button>
    </div>
  </div>
</template>

    
<script setup lang="ts">
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
import api from '../../../services/api';

const { t } = useI18n();

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
    errorMessage.value = t('errorEnterFriendName');
    return;
  }

  if (newFriendNickname.value.length < 3) {
    errorMessage.value = t('errorNicknameTooShort');
    return;
  }

  loadingAddFriend.value = true;
  try {
    await api.friendlist.sendFriendRequest(newFriendNickname.value.trim());
    newFriendNickname.value = '';
    props.fetchFriendRequests();
    successMessage.value = t('friendRequestSuccess');
    errorMessage.value = '';
  } catch (error: any) {
    errorMessage.value = error.error || t('errorUnexpected');
  } finally {
    loadingAddFriend.value = false;
  }
};

</script>
  
<style scoped>
.add-friend-popover {
  position: fixed;
  bottom: 80px;
  right: 20px;
  width: 300px;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  z-index: 1000;
}

.add-friend-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 15px;
  border-bottom: 1px solid var(--main-extra-color);
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
  display: flex;
  flex-direction: column;
  padding: 15px;
  max-height: 400px;
  overflow-y: auto;
  overflow-x: hidden;
}

.friend-input {
  width: 93%;
  padding: 8px;
  margin-bottom: 10px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 14px;
}

.add-friend-btn {
  width: 100%;
  padding: 8px;
  background-color: var(--main-extra-color);
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.3s;
}

.add-friend-btn:hover {
	background-color: color-mix(in srgb, var(--main-extra-color) 85%, white);
	transform: scale(1.025);
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
  