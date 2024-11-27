<template>
  <div class="friend-request-item">
    <div class="friend-request-name">{{ request.nickname }}</div>
    <div class="friend-request-actions">
      <button @click="acceptRequest" class="accept-button" :disabled="loadingAccept">
        <i class="fas fa-check"></i>
      </button>
      <button @click="denyRequest" class="deny-button" :disabled="loadingDeny">
        <i class="fas fa-times"></i>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { Friend } from '../../../types/models';

const props = defineProps<{
  request: Friend;
  acceptFriend: (id: number) => Promise<void>;
  denyFriend: (id: number) => Promise<void>;
}>();

const loadingAccept = ref(false);
const loadingDeny = ref(false);

const acceptRequest = async () => {
  loadingAccept.value = true;
  try {
    await props.acceptFriend(props.request.id);
  } catch (error) {
    console.error('Failed to accept friend request', error);
  } finally {
    loadingAccept.value = false;
  }
};

const denyRequest = async () => {
  loadingDeny.value = true;
  try {
    await props.denyFriend(props.request.id);
  } catch (error) {
    console.error('Failed to deny friend request', error);
  } finally {
    loadingDeny.value = false;
  }
};
</script>

<style scoped>
.friend-request-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid var(--main-extra-color);
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
  transition: color 0.3s;
}

.accept-button:hover {
  color: #28a745;
}

.deny-button:hover {
  color: #dc3545;
}
</style>
