<template>
	<div class="live-chat-container">
		<ChatIcon @toggle="toggleChatInterface" />

		<div v-if="showChatInterface" class="chat-interface">
			<div class="chat-content">
				<FriendList :friends="friends" :currentFriendId="currentFriendId"
					@select-friend="selectFriend" />
				<ChatDiscussion :currentFriend="currentFriend" :messages="currentConversation"
					:userId="userStore.getId ?? 0" @send-message="sendMessage" />
			</div>
			<button @click="toggleChatInterface" class="close-button">
				<i class="fas fa-times"></i>
			</button>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { useUserStore } from '../../../stores/user';
import { useChatStore } from '../../../stores/chatStore';
import api from '../../../services/api';
import ChatIcon from './ChatIcon.vue';
import FriendList from './ChatFriendList.vue';
import ChatDiscussion from './ChatDiscussion.vue';
import { Friend, Message, ChatHistory } from '../../../types/models';
import { ChatMessage } from '../../../types/chat';

const showChatInterface = ref(false);
const currentFriendId = ref<number | null>(null);
const userStore = useUserStore();
const chatStore = useChatStore();
const friends = ref<Friend[]>([]);
const fetchedConversationsTracker = ref<Set<number>>(new Set());
const conversations = ref<{ [friendId: number]: Message[] }>({});

const currentFriend = computed(() =>
	friends.value.find((f: Friend) => f.id === currentFriendId.value)
);

const currentConversation = computed(() =>
	currentFriendId.value
		? conversations.value[currentFriendId.value] || []
		: []
);

const toggleChatInterface = () => {	
	showChatInterface.value = !showChatInterface.value;

	if (showChatInterface.value && currentFriendId.value)
    	chatStore.selectFriend(currentFriendId.value);

	if (userStore.getWebSocketService)
		setupChatMessageHandler();
};

const selectFriend = async (friendId: number) => {
	currentFriendId.value = friendId;
	chatStore.selectFriend(friendId);
	await loadFriendDiscussion(friendId);
};

const loadFriendDiscussion = async (friendId: number) => {
	if (fetchedConversationsTracker.value.has(friendId)) {
		return;
	}

	try {
		const messages: ChatHistory | null = await api.chat.getChatHistory(friendId);
		const formattedMessages = messages?.conversation?.map(msg => ({
			content: msg.content,
			senderId: msg.senderId,
			receiverId: msg.receiverId,
			createdAt: msg.createdAt
		})) || [];

		const conversationId = friendId;
		conversations.value[conversationId] = formattedMessages;
		fetchedConversationsTracker.value.add(friendId);
	} catch (error) {
		console.error('Failed to load discussion history', error);
	}
};

const sendMessage = (message: string) => {
	console.log("Sending message to friend ID:", currentFriendId.value);
	if (message.trim() && currentFriendId.value) {
		if (userStore.getWebSocketService?.isConnected()) {
			userStore.getWebSocketService?.sendChatMessage(
				message,
				userStore.getId ?? 0,
				currentFriendId.value
			);
		} else {
			console.error('WebSocket is not connected');
		}
	}
};

const setupChatMessageHandler = () => {
	if (!userStore.getWebSocketService) {
		return;
	}

	userStore.getWebSocketService.setMessageHandler('CHAT', (message: ChatMessage) => {
		const messageToPush: Message = {
			content: message.data,
			senderId: message.senderID,
			receiverId: message.receiverID,
			createdAt: new Date().toISOString()
		};

		const conversationId = messageToPush.senderId === userStore.getId
			? messageToPush.receiverId
			: messageToPush.senderId;

		if (!conversations.value[conversationId]) {
			conversations.value[conversationId] = [];
		}
		conversations.value[conversationId].push(messageToPush);
	});
};

const fetchFriendList = async () => {
	try {
		const fetchedFriends = await api.friendlist.getFriendList();
		if (fetchedFriends) {
			friends.value = fetchedFriends;
		}
	} catch (error) {
		console.error('Failed to fetch friend list', error);
	}
};

watch(() => chatStore.selectedFriendId, (newFriendId) => {
  if (newFriendId === null) {
    showChatInterface.value = false;
  } else {
    showChatInterface.value = true;
    selectFriend(newFriendId);

	if (userStore.getWebSocketService) {
      setupChatMessageHandler();
    }
  }
});

onMounted(() => {
	fetchFriendList();
});
</script>

<style scoped>
.live-chat-container {
	position: fixed;
	bottom: 20px;
	left: 20px;
	z-index: 1001;
}

.chat-interface {
	position: absolute;
	bottom: 60px;
	left: 0;
	width: 600px;
	height: 500px;
	background-color: white;
	border-radius: 12px;
	box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
	display: flex;
	flex-direction: column;
	overflow: hidden;
}

.chat-content {
	display: flex;
	height: 100%;
	overflow: hidden;
}

.close-button {
	position: absolute;
	top: 10px;
	right: 10px;
	background: none;
	border: none;
	cursor: pointer;
	font-size: 18px;
	color: #666;
	padding: 8px;
	border-radius: 50%;
	transition: background-color 0.2s ease;
}

.close-button:hover {
	background-color: #f0f0f0;
	color: #333;
}

@media (max-width: 640px) {
	.chat-interface {
		width: 100vw;
		height: 100vh;
		bottom: 0;
		left: 0;
		border-radius: 0;
	}

	.chat-content {
		flex-direction: column;
	}
}
</style>
