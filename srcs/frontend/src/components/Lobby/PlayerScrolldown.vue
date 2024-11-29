<template>
  <div class="scroll-down-container" ref="containerRef">
    <div v-if="isWaiting" class="selected-friend">
      <span class="select-friend-text">Waiting for friend's answer...</span>
    </div>
    <div v-else>
      <div class="selected-friend" @click="toggleDropdown">
        <span v-if="selectedFriend">{{ selectedFriend }}</span>
        <span class="select-friend-text" v-else>Select a friend</span>
        <span class="arrow" :class="{ 'arrow-up': isOpen }">▼</span>
      </div>
      <div class="dropdown-list" v-if="isOpen">
        <div v-if="onlineFriends.length === 0" class="friend-item no-friends">
          No friends online
        </div>
        <div v-for="friend in onlineFriends" :key="friend.id" class="friend-item" @click="selectFriend(friend.id)">
          {{ friend.nickname }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useOnlineUsersStore } from '../../stores/onlineUsers'
import { Friend } from '../../types/models';
import api from '../../services/api';

let friendList: Friend[] = [];
const onlineUsersStore = useOnlineUsersStore();
const onlineUsers = computed(() => onlineUsersStore.getOnlineUsers);
const onlineFriends = computed(() =>
  friendList.filter(friend => onlineUsers.value.includes(friend.id))
);

const isOpen = ref(false);
const containerRef = ref<HTMLElement | null>(null);

interface Props {
  selectedFriend?: number | null
  isWaiting?: boolean
}

withDefaults(defineProps<Props>(), {
  selectedFriend: null,
  isWaiting: false
})

const emit = defineEmits<{
  (e: 'friend-selected', userId: number): void
}>()

const toggleDropdown = () => {
  isOpen.value = !isOpen.value
}

const selectFriend = (userId: number) => {
  emit('friend-selected', userId)
  isOpen.value = false
}

const handleClickOutside = (event: MouseEvent) => {
  if (containerRef.value && !containerRef.value.contains(event.target as Node)) {
    isOpen.value = false
  }
}

onMounted(async () => {
  friendList = await api.friendlist.getFriendList();
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.scroll-down-container {
  position: relative;
  width: 220px;
}

.selected-friend {
  background-color: #f8f9fa;
  color: white;
  padding: 12px 16px;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  justify-content: space-between;
  align-items: center;
  transition: background-color 0.2s;
}

.select-friend-text {
  font-weight: bold;
  color: #2c3e50;
}

.selected-friend:hover {
  background-color: #f8f8fa;
}

.arrow {
  display: inline-block;
  transition: transform 0.2s ease;
  color: #ffffff;
  font-size: 0.8em;
  margin-left: 8px;
}

.arrow-up {
  transform: rotate(180deg);
}

.dropdown-list {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  right: 0;
  background-color: #ffffff;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  z-index: 1000;
  max-height: 200px;
  overflow-y: auto;
}

.friend-item {
  padding: 12px 16px;
  cursor: pointer;
  transition: all 0.2s ease;
  color: #2c3e50;
  border-bottom: 1px solid #eee;
}

.friend-item:last-child {
  border-bottom: none;
}

.friend-item:hover {
  background-color: #f8f9fa;
  color: #2c3e50;
}

/* Scrollbar styling */
.dropdown-list::-webkit-scrollbar {
  width: 6px;
}

.dropdown-list::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

.dropdown-list::-webkit-scrollbar-thumb {
  background: #cbd5e0;
  border-radius: 4px;
}

.dropdown-list::-webkit-scrollbar-thumb:hover {
  background: #718096;
}

/* Smooth transitions */
.dropdown-list {
  transition: opacity 0.2s ease, transform 0.2s ease;
}

@media (max-width: 768px) {
  .scroll-down-container {
    width: 100%;
    max-width: 220px;
  }
}
</style>
