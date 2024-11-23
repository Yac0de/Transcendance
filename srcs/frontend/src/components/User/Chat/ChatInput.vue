<template>
	<div class="message-input">
		<input ref="inputRef" v-model="message" @keyup.enter="sendMessage" placeholder="Type a message..." />
		<button 
			class="action-button challenge-btn" 
			:class="{ 'disabled': !isOnline }"
			@click="handleChallengeClick"
			:disabled="!isOnline"
			title="Challenge to a match"
		>
			<i class="fas fa-gamepad"></i>
		</button>
		<button class="action-button" @click="sendMessage">Send</button>
	</div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useOnlineUsersStore } from '../../../stores/onlineUsers'
import { useRouter } from 'vue-router';

const router = useRouter();
const emit = defineEmits<{
	(e: 'send', message: string): void;
}>();

const props = defineProps<{
	currentFriendId: number;
}>();

const onlineUsersStore = useOnlineUsersStore();
const message = ref('');
const inputRef = ref<HTMLInputElement | null>(null);

const sendMessage = () => {
	if (message.value.trim()) {
		emit('send', message.value);
		message.value = '';
	}
};

const isOnline = computed(() => {
	return onlineUsersStore.isUserOnline(props.currentFriendId);
});

const handleChallengeClick = () => {
    if (!isOnline.value) return;
    router.push({        // Then navigate to lobby
        path: '/lobby',
        query: { friendId: props.currentFriendId.toString() }
    });
};

const focusInput = () => {
    if (inputRef.value)
        inputRef.value.focus();
};

defineExpose({
    focusInput
});
</script>

<style scoped>
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

.message-input .action-button {
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

.message-input .action-button.disabled {
	background-color: #ccc;
	cursor: not-allowed;
	transform: none;
}

.message-input .action-button:hover:not(.disabled) {
	background-color: #1557b0;
	transform: translateY(-1px);
}

.message-input .action-button:active:not(.disabled) {
	transform: translateY(0);
}

@media (max-width: 640px) {
	.message-input {
		padding: 12px;
	}
	.message-input input {
		font-size: 16px;
	}
}
</style>
