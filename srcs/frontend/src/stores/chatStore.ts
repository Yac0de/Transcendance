import { defineStore } from 'pinia';
import { ref, computed } from 'vue';

export const useChatStore = defineStore('chat', () => {
  const selectedFriendId = ref<number | null>(null);
  const unreadMessagesCount = ref<{ [friendId: number]: number }>({});

  const selectFriend = (friendId: number) => {
    selectedFriendId.value = friendId;
    if (unreadMessagesCount.value[friendId]) {
      unreadMessagesCount.value[friendId] = 0;
    }
  };

  const addUnreadMessage = (friendId: number) => {
    if (!unreadMessagesCount.value[friendId]) {
      unreadMessagesCount.value[friendId] = 0;
    }
    unreadMessagesCount.value[friendId]++;
  };

  const totalUnreadMessages = computed(() =>
    Object.values(unreadMessagesCount.value).reduce((a, b) => a + b, 0)
  );

  const toggleFriend = (friendId: number) => {
    if (selectedFriendId.value === friendId) {
      selectedFriendId.value = null;
    } else {
      selectedFriendId.value = friendId;
    }
  };

  return {
    selectedFriendId,
    unreadMessagesCount,
    selectFriend,
    addUnreadMessage,
    totalUnreadMessages,
    toggleFriend,
  };
});
