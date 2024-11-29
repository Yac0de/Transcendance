<template>
	<div class="friend-list">
		<div class="friend-list-content">
			<ul>
				<li v-for="friend in friends" :key="friend.id"
					:class="['friend-item', { 'active': currentFriendId === friend.id }]"
					@click="$emit('select-friend', friend.id)">
					{{ friend.nickname }}
					<span v-if="getUnreadCount(friend.id) > 0" class="unread-badge">
                        {{ getUnreadCount(friend.id) }}
                    </span>
				</li>
			</ul>
		</div>
	</div>
</template>

<script setup lang="ts">
import { Friend } from '../../../types/models';
import { useChatStore } from '../../../stores/chatStore';

defineProps<{
	friends: Friend[];
	currentFriendId: number | null;
}>();

defineEmits<{
	(e: 'select-friend', friendId: number): void;
}>();

const chatStore = useChatStore();

const getUnreadCount = (friendId: number) => {
    return chatStore.getUnreadCountForFriend(friendId);
};
</script>

<style scoped>
.friend-list {
	width: 250px;
	flex-shrink: 0;
	border-right: 1px solid var(--main-extra-color);
	background-color: white;
	display: flex;
	flex-direction: column;
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
	position: relative;
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
	color: var(--main-extra-color);
}

.unread-badge {
    position: absolute;
    top: 50%;
    right: 10px;
    transform: translateY(-50%);
    background-color: red;
    color: white;
    font-size: 12px;
    font-weight: bold;
    padding: 2px 6px;
    border-radius: 12px;
    box-shadow: 0 0 5px rgba(0, 0, 0, 0.2);
    white-space: nowrap;
}

@media (max-width: 640px) {
	.friend-list {
		width: 100%;
		max-width: 100%;
		border-right: none;
		height: 30%;
		min-height: 200px;
	}
}
</style>
