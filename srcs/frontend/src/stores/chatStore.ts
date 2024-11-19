import { defineStore } from 'pinia';
import { ref, computed } from 'vue';

export const useChatStore = defineStore('chat', () => {
  const selectedFriendId = ref<number>(-1);
  const unreadMessagesCount = ref<{ [friendId: number]: number }>({});

  const selectFriend = (friendId: number | -1) => {
    selectedFriendId.value = friendId;
    if (friendId !== -1) unreadMessagesCount.value[friendId] = 0;
  };

  const addUnreadMessage = (friendId: number) => {
    unreadMessagesCount.value[friendId] = unreadMessagesCount.value[friendId] || 0;
    unreadMessagesCount.value[friendId]++;
  };

  const getUnreadCountForFriend = (friendId: number) => unreadMessagesCount.value[friendId] || 0;

  const totalUnreadMessages = computed(() =>
    Object.values(unreadMessagesCount.value).reduce((a, b) => a + b, 0)
  );

  const toggleFriend = (friendId: number) => {
    if (selectedFriendId.value === friendId) {
      selectedFriendId.value = -1;
    } else {
      selectedFriendId.value = friendId;
    }
  };

  return {
    selectedFriendId,
    unreadMessagesCount,
    selectFriend,
    addUnreadMessage,
    getUnreadCountForFriend,
    totalUnreadMessages,
    toggleFriend,
  };
});
