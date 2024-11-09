import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useChatStore = defineStore('chat', () => {
  const selectedFriendId = ref<number | null>(null);

  const selectFriend = (friendId: number) => {
    selectedFriendId.value = friendId;
  };

  return {
    selectedFriendId,
    selectFriend,
  };
});
