<template>
	<div class="live-chat-container">
		<div class="chat-icon-container">
			<div class="chat-icon" @click="toggleChatInterface">
				<i class="fas fa-comment"></i>
			</div>
		</div>

		<div v-if="showChatInterface" class="chat-interface">
			<div class="chat-content">
				<div class="discussion-list">
					<ul>
						<li v-for="discussion in discussions" :key="discussion.id"
							@click="selectDiscussion(discussion.id)"
							:class="{ 'active': currentDiscussionId === discussion.id }">
							{{ discussion.name }}
						</li>
					</ul>
				</div>
				<div class="current-discussion">
					<template v-if="currentDiscussion">
						<h4>{{ currentDiscussion.name }}</h4>
						<div class="messages">
							<div v-for="message in currentDiscussion.messages"
								:key="message.id"
								:class="['message-wrapper', message.sender === 'user' ? 'user-message' : 'receiver-message']">
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
						<p>Select a discussion to start chatting</p>
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
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { WebSocketService, ChatMessage } from '../../../services/chatService';
import { API_BASE_URL } from './apiUtils';
import { useUserStore } from '../../../stores/user'

const showChatInterface = ref(false);
const currentDiscussionId = ref(null);
const newMessage = ref('');
const userStore = useUserStore();

const webSocketService = new WebSocketService(userStore.getId);

// MODIFIED: Updated mock data structure to match ChatMessage interface
const discussions = ref([
	{ id: 1, name: 'General Chat', messages: [] as ChatMessage[] },
	{ id: 2, name: 'Support', messages: [] as ChatMessage[] },
]);

const currentDiscussion = computed(() =>
	discussions.value.find(d => d.id === currentDiscussionId.value)
);

const handleWebSocketMessage = (message: ChatMessage) => {
	if (currentDiscussionId.value) {
		const discussion = discussion.value.find(d => d.id === currentDiscussionId.value);
		if (discussion) {
			discussion.messages.push(message);
		}
	}
}

const toggleChatInterface = () => {
	showChatInterface.value = !showChatInterface.value;
};

const selectDiscussion = (id: number) => {
	currentDiscussionId.value = id;
};

const sendMessage = () => {
	if (newMessage.value.trim() && currentDiscussionId.value) {
		const message: ChatMessage = {
			sender: 'User', //security here ?
			content: newMessage.value,
			timestamps: new Date().toISOString()
		};

		if (webSocketService.isConnected()) {
			webSocketService.sendMessage(newMessage.value);
			const discussion = discussions.value.find(
				d => d.id === currentDiscussionId.value
			);
			if (discussion) {
				discussion.message.push(message);
			}

			newMessage.value = '';
		} else {
			console.error('Websocket is not connected');
		}
	}
};

onMounted(() => {
	webSocketService.connect();

	if (webSocketService['ws']) {
		webSocketService['ws'].onMessage = (event) => {
			try {
				const message = JSON.parse(event.data);
				handleWebSocketMessage(message);
			} catch (error) {
				console.error('Error parsing the mess', error);
			}
		};
	}
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
	transition: background-color 0.3s;
}

.chat-icon:hover {
	background-color: #218838;
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
	height: 400px;
	background-color: white;
	border-radius: 8px;
	box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
	display: flex;
	flex-direction: column;
}

.close-button {
	position: absolute;
	top: 10px;
	right: 10px;
	background: none;
	border: none;
	cursor: pointer;
	font-size: 18px;
}

.chat-content {
	display: flex;
	height: 100%;
	overflow: hidden;
	/* Prevent content from expanding beyond container */
}

.discussion-list {
	width: 200px;
	/* Fixed width */
	flex-shrink: 0;
	/* Prevent shrinking */
	border-right: 1px solid #dee2e6;
	overflow-y: auto;
}

.discussion-list ul {
	list-style-type: none;
	padding: 0;
	margin: 0;
}

.discussion-list li {
	padding: 10px;
	cursor: pointer;
}

.discussion-list li.active {
	background-color: #e9ecef;
}

.current-discussion {
	flex-grow: 1;
	min-width: 0;
	/* Allow shrinking if necessary */
	display: flex;
	flex-direction: column;
	padding: 10px;
	overflow-x: hidden;
}

.current-discussion h4 {
	margin-top: 0;
	margin-bottom: 10px;
}

.messages {
	flex-grow: 1;
	overflow-y: auto;
	display: flex;
	flex-direction: column;
	padding: 10px;
}

.message-wrapper {
	display: flex;
	margin-bottom: 10px;
	width: 100%;
}

.user-message {
	justify-content: flex-end;
}

.receiver-message {
	justify-content: flex-start;
}

.message-content {
	max-width: calc(100% - 20px);
	/* Adjust for padding */
	padding: 5px 10px;
	border-radius: 4px;
	word-wrap: break-word;
	overflow-wrap: break-word;
	hyphens: auto;
}

.user-message .message-content {
	background-color: #007bff;
	color: white;
}

.receiver-message .message-content {
	background-color: #f1f3f5;
	color: black;
}

.message-input {
	display: flex;
	margin-top: 10px;
}

.message-input input {
	flex-grow: 1;
	padding: 5px;
	border: 1px solid #ced4da;
	border-radius: 4px 0 0 4px;
}

.message-input button {
	padding: 5px 10px;
	background-color: #28a745;
	color: white;
	border: none;
	border-radius: 0 4px 4px 0;
	cursor: pointer;
}

@media (max-width: 600px) {
	.chat-icon {
		width: 40px;
		height: 40px;
	}

	.chat-icon i {
		font-size: 20px;
	}

	.chat-interface {
		width: 100%;
		height: calc(100% - 80px);
		position: fixed;
		bottom: 70px;
		left: 0;
		right: 0;
		border-radius: 0;
	}
}
</style>
