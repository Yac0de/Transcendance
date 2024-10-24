<template>
	<div class="live-chat-container">
		<div class="chat-icon-container">
			<div class="chat-icon" @click="toggleChatInterface">
				<i class="fas fa-comment"></i>
			</div>
		</div>

		<div v-if="showChatInterface" class="chat-interface">
			<div class="chat-content">
				<div class="friend-list">
					<div class="friend-list-header">
						<h3>Friends</h3>
					</div>
					<div class="friend-list-content">
						<ul>
							<li v-for="friend in friends" :key="friend.id"
								:class="['friend-item', { 'active': currentFriendId === friend.id }]"
								@click="selectFriend(friend.id)">
								{{ friend.nickname }}
							</li>
						</ul>
					</div>
				</div>
				<div class="current-discussion">
					<template v-if="currentFriend">
						<h4>{{ currentFriend.nickname }}</h4>
						<div class="messages">
							<div v-for="message in currentFriendMessages" :key="message.id"
								:class="['message-wrapper', message.senderId === userStore.getId ? 'user-message' : 'receiver-message']">
								<div class="message-content">
									{{ message.content }}
								</div>
							</div>
						</div>
						<div class="message-input">
							<input v-model="newMessage" @keyup.enter="sendMessage"
								placeholder="Type a message..." />
							<button @click="sendMessage">Send</button>
						</div>
					</template>
					<template v-else>
						<p>Select a friend to start chatting</p>
					</template>
				</div>
			</div>
			<button @click="toggleChatInterface" class="close-button">
				<i class="fas fa-times"></i>
			</button>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { WebSocketService, ChatMessage } from '../../../services/chatService';
import { useUserStore } from '../../../stores/user'
import api from '../../../services/api';

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
const newMessage = ref('');
const userStore = useUserStore();
const friends = ref<Friend[]>([]);
const messagesByFriend = ref<{ [friendId: string]: Message[] }>({});

const webSocketService = new WebSocketService(userStore.getId);

const currentFriend = computed(() =>
	friends.value.find(f => f.id === currentFriendId.value)
);

const currentFriendMessages = computed(() =>
	currentFriendId.value
		? messagesByFriend.value[currentFriendId.value] || []
		: []
);

const handleWebSocketMessage = (message: Message) => {
	console.log("CALLBACK FROM THE SEVER: ", message);
	const friendId = message.senderId === userStore.getId
		? message.receiverId
		: message.senderId;
};

const toggleChatInterface = () => {
	showChatInterface.value = !showChatInterface.value;
};

const selectFriend = async (friendId: string) => {
	currentFriendId.value = friendId;

	await loadFriendDiscussion(friendId);
};

const loadFriendDiscussion = async (friendId: string) => {
	try {
		//const messages = await api.chat.getChatHistory(friendId);
		//discussions.value[friendId] = messages;
	} catch (error) {
		console.error('Failed to load discussion history', error);
	}
};

const sendMessage = () => {
	if (newMessage.value.trim() && currentFriendId.value) {
		if (webSocketService.isConnected()) {
			webSocketService.sendMessage(
				newMessage.value,
				userStore.getId,
				currentFriendId.value
			);

			newMessage.value = '';
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
	console.log("MSG RECEIEVED IN THE HANDLER");
	if (message.Type === 'CHAT') {
		const messageToPush = {
			content: message.Data,
			senderId: message.SenderID,
			receiverId: message.ReceiverID,
			timestamp: new Date().toISOString()
		};

		//This computes the id of the friend who you are discussing with (even when the message
		//is your own message, returned by the server.
		const conversationId = messageToPush.senderId === userStore.getId
			? messageToPush.receiverId
			: messageToPush.senderId;

		console.log("C ID = ", conversationId);

		if (!messagesByFriend.value[conversationId]) {
			messagesByFriend.value[conversationId] = [];
		}

		messagesByFriend.value[conversationId].push(messageToPush);
		console.log("MTP = ", messageToPush.content);
		console.log("All messages = ", messagesByFriend.value[conversationId].map(msg => msg.content));
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

.chat-icon-container {
	position: relative;
}

.chat-icon {
	width: 50px;
	height: 50px;
	background-color: #28a745;
	border-radius: 50%;
	display: flex;
	justify-content: center;
	align-items: center;
	cursor: pointer;
	box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
	transition: background-color 0.3s ease, transform 0.2s ease;
}

.chat-icon:hover {
	background-color: #218838;
	transform: scale(1.05);
}

.chat-icon i {
	color: white;
	font-size: 24px;
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

.chat-content {
	display: flex;
	height: 100%;
	overflow: hidden;
}

.friend-list {
	width: 250px;
	flex-shrink: 0;
	border-right: 1px solid #e0e0e0;
	background-color: white;
	display: flex;
	flex-direction: column;
}

.friend-list-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 16px;
	border-bottom: 1px solid #e0e0e0;
	background-color: #f8f9fa;
}

.friend-list-header h3 {
	margin: 0;
	font-size: 18px;
	color: #333;
	font-weight: 600;
}

.friend-list-content {
	padding: 12px;
	flex-grow: 1;
	overflow-y: auto;
}

.friend-list-content::-webkit-scrollbar {
	width: 6px;
}

.friend-list-content::-webkit-scrollbar-track {
	background: #f1f1f1;
}

.friend-list-content::-webkit-scrollbar-thumb {
	background: #888;
	border-radius: 3px;
}

.friend-list-content ul {
	list-style-type: none;
	padding: 0;
	margin: 0;
}

.friend-item {
	padding: 12px 16px;
	margin-bottom: 8px;
	background-color: #f8f9fa;
	border-radius: 8px;
	cursor: pointer;
	transition: all 0.2s ease;
	font-size: 14px;
	color: #333;
}

.friend-item:hover {
	background-color: #e9ecef;
	transform: translateX(2px);
}

.friend-item.active {
	background-color: #e2e8f0;
	font-weight: 500;
	color: #1a73e8;
}

.current-discussion {
	flex-grow: 1;
	display: flex;
	flex-direction: column;
	background-color: #fff;
}

.current-discussion h4 {
	margin: 0;
	padding: 16px;
	font-size: 16px;
	font-weight: 600;
	color: #333;
	border-bottom: 1px solid #e0e0e0;
	background-color: #f8f9fa;
}

.messages {
	flex-grow: 1;
	overflow-y: auto;
	padding: 16px;
	display: flex;
	flex-direction: column;
	gap: 8px;
}

.messages::-webkit-scrollbar {
	width: 6px;
}

.messages::-webkit-scrollbar-track {
	background: #f1f1f1;
}

.messages::-webkit-scrollbar-thumb {
	background: #888;
	border-radius: 3px;
}

.message-wrapper {
	display: flex;
	margin-bottom: 4px;
	width: 100%;
	min-width: 0;
}

.user-message {
	justify-content: flex-end;
}

.receiver-message {
	justify-content: flex-start;
}

.message-content {
	max-width: 50%;
	min-width: 0;
	padding: 8px 16px;
	border-radius: 16px;
	word-wrap: break-word;
	overflow-wrap: break-word;
	word-break: break-word;
	hyphens: auto;
	font-size: 14px;
	line-height: 1.4;
	white-space: pre-wrap;
}

.user-message .message-content {
	background-color: #1a73e8;
	color: white;
	border-bottom-right-radius: 4px;
}

.receiver-message .message-content {
	background-color: #f1f3f4;
	color: #333;
	border-bottom-left-radius: 4px;
}

.message-input {
	padding: 16px;
	border-top: 1px solid #e0e0e0;
	display: flex;
	gap: 8px;
	background-color: #fff;
}

.message-input input {
	flex-grow: 1;
	padding: 10px 16px;
	border: 1px solid #e0e0e0;
	border-radius: 24px;
	font-size: 14px;
	transition: all 0.2s ease;
	background-color: #f8f9fa;
}

.message-input input:focus {
	outline: none;
	border-color: #1a73e8;
	background-color: #fff;
	box-shadow: 0 0 0 2px rgba(26, 115, 232, 0.1);
}

.message-input button {
	padding: 10px 20px;
	background-color: #1a73e8;
	color: white;
	border: none;
	border-radius: 24px;
	cursor: pointer;
	font-size: 14px;
	font-weight: 500;
	transition: all 0.2s ease;
}

.message-input button:hover {
	background-color: #1557b0;
	transform: translateY(-1px);
}

.message-input button:active {
	transform: translateY(0);
}

@media (max-width: 640px) {
	.chat-interface {
		width: 100vw;
		height: 100vh;
		bottom: 0;
		left: 0;
		border-radius: 0;
	}

	.chat-icon {
		width: 40px;
		height: 40px;
	}

	.chat-icon i {
		font-size: 20px;
	}

	.friend-list {
		width: 100%;
		max-width: 100%;
		border-right: none;
	}

	.chat-content {
		flex-direction: column;
	}

	.friend-list {
		height: 30%;
		min-height: 200px;
	}

	.current-discussion {
		height: 70%;
	}

	.message-input {
		padding: 12px;
	}

	.message-input input {
		font-size: 16px;
	}
}
</style>
