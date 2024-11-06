<template>
	<div class="current-discussion">
		<template v-if="currentFriend">
			<h4>{{ currentFriend.nickname }}</h4>
			<div class="messages" ref="messageContainer">
				<div v-for="message in messages" :key="`${message.senderId} - ${message.createdAt}`"
					:class="['message-wrapper', message.senderId === userId ? 'user-message' : 'receiver-message']">
					<div class="message-content">
						{{ message.content }}
					</div>
				</div>
			</div>
			<ChatInput @send="handleSend" />
		</template>
		<template v-else>
			<p class="no-friend-selected"> Select a friend to start chatting</p>
		</template>
	</div>
</template>

<script setup lang="ts">
import ChatInput from './ChatInput.vue';
import { ref, watch, nextTick } from 'vue';
import { Friend, Message } from '../../../types/models';

const messageContainer = ref<HTMLElement | null>(null);

const scrollToBottom = () => {
	nextTick(() => {
		if (messageContainer.value) {
			messageContainer.value.scrollTop = messageContainer.value.scrollHeight;
		}
	});
};

const props = defineProps<{
	currentFriend: Friend | undefined;
	messages: Message[];
	userId: number;
}>();

const emit = defineEmits<{
	(e: 'send-message', message: string): void;
}>();

const handleSend = (message: string) => {
	emit('send-message', message);
};

watch(() => props.messages, () => {
	scrollToBottom();
}, { deep: true });

</script>

<style scoped>
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

.no-friend-selected {
	display: flex;
	justify-content: center;
	align-items: center;
	height: 100%;
	color: #666;
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

@media (max-width: 640px) {
	.current-discussion {
		height: 70%;
	}
}
</style>
