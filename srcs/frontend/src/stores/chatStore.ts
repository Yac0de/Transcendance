import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useChatStore = defineStore('chat', () => {
  const selectedFriendId = ref<number | null>(null);

  const selectFriend = (friendId: number) => {
    selectedFriendId.value = friendId;
  };

  const toggleFriend = (friendId: number) => {
    if (selectedFriendId.value === friendId) {
      selectedFriendId.value = null;
    } else {
      selectedFriendId.value = friendId;
    }
  };

  return {
    selectedFriendId,
    selectFriend,
    toggleFriend
  };
});
