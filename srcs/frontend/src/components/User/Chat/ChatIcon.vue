<template>
	<div class="chat-icon-container">
		<div class="chat-icon" @click="$emit('toggle')">
			<i class="fas fa-comment"></i>
			<span v-if="unreadMessages > 0" class="notification-badge">{{ unreadMessages }}</span>
		</div>
	</div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useChatStore } from '../../../stores/chatStore';

const chatStore = useChatStore();

defineEmits<{
	(e: 'toggle'): void;
}>();

const unreadMessages = computed(() => chatStore.totalUnreadMessages);
</script>

<style scoped>
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

.notification-badge {
  position: absolute;
  top: -5px;
  right: -5px;
  background-color: red;
  color: white;
  font-size: 12px;
  font-weight: bold;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
  box-shadow: 0 0 5px rgba(0, 0, 0, 0.3);
}

@media (max-width: 640px) {
	.chat-icon {
		width: 40px;
		height: 40px;
	}

	.chat-icon i {
		font-size: 20px;
	}
}
</style>
