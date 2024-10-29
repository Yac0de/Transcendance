<template>
	<div class="chat-icon-container">
		<div class="chat-icon" @click="$emit('toggle')">
			<i class="fas fa-comment"></i>
		</div>
	</div>
</template>

<script setup lang="ts">
defineEmits<{
	(e: 'toggle'): void;
}>();
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
