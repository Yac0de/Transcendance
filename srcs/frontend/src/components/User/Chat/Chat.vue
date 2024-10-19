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
								:class="['message', message.sender === 'user' ? 'user-message' : 'receiver-message']">
								{{ message.content }}
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
import { ref, computed } from 'vue';

const showChatInterface = ref(false);
const currentDiscussionId = ref(null);
const newMessage = ref('');

// Mock data for discussions
const discussions = ref([
	{ id: 1, name: 'General Chat', messages: [{ id: 1, content: 'Welcome to the general chat!', sender: 'system' }] },
	{ id: 2, name: 'Support', messages: [{ id: 1, content: 'How can we help you?', sender: 'system' }] },
]);

const currentDiscussion = computed(() =>
	discussions.value.find(d => d.id === currentDiscussionId.value)
);

const toggleChatInterface = () => {
	showChatInterface.value = !showChatInterface.value;
};

const selectDiscussion = (id) => {
	currentDiscussionId.value = id;
};

const sendMessage = () => {
	if (newMessage.value.trim() && currentDiscussion.value) {
		currentDiscussion.value.messages.push({
			id: currentDiscussion.value.messages.length + 1,
			content: newMessage.value,
			sender: 'user'
		});
		newMessage.value = '';
	}
};
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
}

.discussion-list {
	width: 200px;
	border-right: 1px solid #dee2e6;
	overflow-y: auto;
}

.discussion-list ul {
	list-style-type: none;
	padding: 0;
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
	display: flex;
	flex-direction: column;
	padding: 10px;
}

.messages {
	flex-grow: 1;
	overflow-y: auto;
}

.message {
	margin-bottom: 10px;
	padding: 5px 10px;
	border-radius: 4px;
}

.user-message {
	background-color: #007bff;
	color: white;
	align-self: flex-end;
}

.receiver-message {
	background-color: #f1f3f5;
	color: black;
	align-self: flex-start;
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
