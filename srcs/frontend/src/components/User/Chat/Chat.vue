<template>
	<div class="live-chat-container">
		<ChatIcon @toggle="toggleChatInterface" />

		<div v-if="showChatInterface" class="chat-interface">
			<div class="chat-content">
				<FriendList :friends="friends" :currentFriendId="currentFriendId"
					@select-friend="selectFriend" />
				<ChatDiscussion :currentFriend="currentFriend" :messages="currentConversation"
					:userId="userStore.getId" @send-message="sendMessage" />
			</div>
			<button @click="toggleChatInterface" class="close-button">
				<i class="fas fa-times"></i>
			</button>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { WebSocketService } from '../../../services/chatService';
import { useUserStore } from '../../../stores/user';
import api from '../../../services/api';
import ChatIcon from './ChatIcon.vue';
import FriendList from './ChatFriendList.vue';
import ChatDiscussion from './ChatDiscussion.vue';

interface Friend {
	id: string;
	avatar: string;
	nickname: string;
}

interface Message {
	content: string;
	senderId: string;
	receiverId: string;
	timestamp: string;
}

const showChatInterface = ref(false);
const currentFriendId = ref<string | null>(null);
const userStore = useUserStore();
const friends = ref<Friend[]>([]);
const fetchedConversationsTracker = ref<Set<string>>(new Set());
const conversations = ref<{ [friendId: string]: Message[] }>({});

const webSocketService = new WebSocketService(userStore.getId);

const currentFriend = computed(() =>
	friends.value.find(f => f.id === currentFriendId.value)
);

const currentConversation = computed(() =>
	currentFriendId.value
		? conversations.value[currentFriendId.value] || []
		: []
);

const toggleChatInterface = () => {
	showChatInterface.value = !showChatInterface.value;
};

const selectFriend = async (friendId: string) => {
	currentFriendId.value = friendId;
	await loadFriendDiscussion(friendId);
};

const loadFriendDiscussion = async (friendId: string) => {
	if (fetchedConversationsTracker.value.has(friendId)) {
		return;
	}

	try {
		const messages = await api.chat.getChatHistory(friendId);
		const formattedMessages = messages.conversation.map(msg => ({
			content: msg.content,
			senderId: msg.senderId,
			receiverId: msg.receiverId,
			timestamp: msg.createdAt
		}));

		const conversationId = friendId;
		conversations.value[conversationId] = formattedMessages;
		fetchedConversationsTracker.value.add(friendId);
	} catch (error) {
		console.error('Failed to load discussion history', error);
	}
};

const sendMessage = (message: string) => {
	if (message.trim() && currentFriendId.value) {
		if (webSocketService.isConnected()) {
			webSocketService.sendMessage(
				message,
				userStore.getId,
				currentFriendId.value
			);
		} else {
			console.error('WebSocket is not connected');
		}
	}
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

webSocketService.setMessageHandler((message) => {
	if (message.Type === 'CHAT') {
		const messageToPush = {
			content: message.Data,
			senderId: message.SenderID,
			receiverId: message.ReceiverID,
		};

		const conversationId = messageToPush.senderId === userStore.getId
			? messageToPush.receiverId
			: messageToPush.senderId;

		if (!conversations.value[conversationId]) {
			conversations.value[conversationId] = [];
		}

		conversations.value[conversationId].push(messageToPush);
	}
});

onMounted(() => {
	webSocketService.connect();
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
