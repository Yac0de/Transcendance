<template>
	<div class="message-input">
		<input 
			ref="inputRef" 
			maxlength="800" 
			v-model="message" 
			@keyup.enter="sendMessage" 
			:placeholder="$t('typeMessagePlaceholder')" 
		/>
		<button 
			class="action-button challenge-btn" 
			:class="{ 'disabled': !isOnline }"
			@click="handleChallengeClick"
			:disabled="!isOnline"
			:title="$t('challengeMatchTitle')"
		>
			<i class="fas fa-gamepad"></i>
		</button>
		<button class="action-button" @click="sendMessage">
			{{ $t('sendButton') }}
		</button>
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
    router.push({
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
	padding: 12px 16px;
	border-top: 1px solid var(--main-extra-color);
	display: flex;
	gap: 8px;
	background-color: #fff;
	width: 100%;
	box-sizing: border-box;
}

.message-input input {
	flex: 1;
	min-width: 0; /* Prevents input from overflowing */
	padding: 10px 16px;
	border: 1px solid var(--main-extra-color);
	border-radius: 24px;
	font-size: 14px;
	transition: all 0.2s ease;
	background-color: #f8f9fa;
}

.message-input input:focus {
	outline: none;
	border-color: var(--main-extra-color);
	background-color: #fff;
	box-shadow: 0 0 0 2px rgba(26, 115, 232, 0.1);
}

.action-buttons {
	display: flex;
	gap: 8px;
	flex-shrink: 0; /* Prevents buttons from shrinking */
}

.message-input .action-button {
	padding: 10px;
	background-color: var(--main-extra-color);
	color: white;
	border: none;
	border-radius: 24px;
	cursor: pointer;
	font-size: 14px;
	font-weight: 500;
	transition: all 0.2s ease;
	flex-shrink: 0;
}

.message-input .challenge-btn {
	width: 40px;
	padding: 10px 0;
}

.message-input .send-btn {
	padding: 10px 20px;
}

.message-input .action-button.disabled {
	background-color: var(--secondary-extra-color);
	cursor: not-allowed;
	transform: none;
}

.message-input .action-button:hover:not(.disabled) {
	background-color: color-mix(in srgb, var(--main-extra-color) 85%, white);
	transform: translateY(-1px);
}

.message-input .action-button:active:not(.disabled) {
	transform: translateY(0);
}

@media (max-width: 640px) {
	.message-input {
		padding: 8px;
	}
	
	.message-input input {
		font-size: 16px;
	}

	.message-input .send-btn {
		padding: 10px 16px;
	}
}
</style>
