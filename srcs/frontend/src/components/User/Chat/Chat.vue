<template>
	<div class="live-chat-container">
		<ChatIcon @toggle="toggleChatInterface" />

		<div v-if="showChatInterface" class="chat-interface">
			<div class="chat-content">
				<FriendList :friends="friends" :currentFriendId="currentFriendId"
					@select-friend="selectFriend" />
				<ChatDiscussion :currentFriend="currentFriend" :messages="currentConversation"
					:userId="userStore.getId ?? 0" @send-message="sendMessage" />
			</div>
			<button @click="toggleChatInterface" class="close-button">
				<i class="fas fa-times"></i>
			</button>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue';
import { useUserStore } from '../../../stores/user';
import { useChatStore } from '../../../stores/chatStore';
import { eventBus } from '../../../events/eventBus'
import api from '../../../services/api';
import ChatIcon from './ChatIcon.vue';
import FriendList from './ChatFriendList.vue';
import ChatDiscussion from './ChatDiscussion.vue';
import { Friend, Message } from '../../../types/models';
import { ChatMessage } from '../../../types/chat';

const TOURNAMENT_MASTER: Friend = {
    id: 0,
    nickname: "Tournament Master",
    displayname: "Tournament Master",
    avatar: '',
    isOnline: true
};

const showChatInterface = ref(false);
const currentFriendId = ref<number | null>(null);
const userStore = useUserStore();
const chatStore = useChatStore();
const friends = ref<Friend[]>([]);
const fetchedConversationsTracker = ref<Set<number>>(new Set());
const conversations = ref<{ [friendId: number]: Message[] }>({});

const currentFriend = computed(() =>
	friends.value.find((f: Friend) => f.id === currentFriendId.value)
);

const currentConversation = computed(() =>
	currentFriendId.value !== null
		? conversations.value[currentFriendId.value] || []
		: []
);

const toggleChatInterface = () => {
    showChatInterface.value = !showChatInterface.value;

    if (!showChatInterface.value) {
        // close chat
        currentFriendId.value = null;
        chatStore.selectFriend(-1);
    } else {
        // open chat
        if (currentFriendId.value) {
            chatStore.selectFriend(currentFriendId.value);
        }
        userStore.getWebSocketService?.isConnected() && setupChatMessageHandler();
    }
};

const selectFriend = async (friendId: number) => {
    if (friendId !== -1) {
        currentFriendId.value = friendId;
        chatStore.selectFriend(friendId);
        await loadFriendDiscussion(friendId);
    }
};

const loadFriendDiscussion = async (friendId: number) => {
    if (friendId === -1 || friendId === 0 || fetchedConversationsTracker.value.has(friendId)) return;

    try {
        const messages = await api.chat.getChatHistory(friendId);
        conversations.value[friendId] = messages?.conversation?.map(msg => ({
            content: msg.content,
            senderId: msg.senderId,
            receiverId: msg.receiverId,
            createdAt: msg.createdAt,
        })) || [];
        fetchedConversationsTracker.value.add(friendId);
    } catch (error) {
        console.error('Failed to load discussion history', error);
    }
};

const sendMessage = (message: string) => {
    if (message.trim() && currentFriendId.value) {

        if (!conversations.value[currentFriendId.value]) {
            conversations.value[currentFriendId.value] = [];
        }

        if (userStore.getWebSocketService?.isConnected()) {
            userStore.getWebSocketService?.sendChatMessage(
                message,
                userStore.getId ?? 0,
                currentFriendId.value
            );
        } else {
            console.error('WebSocket is not connected');
        }
    }
};

const setupChatMessageHandler = () => {
    const webSocketService = userStore.getWebSocketService;
    if (!webSocketService?.isConnected()) {
        console.log("WebSocket not ready.");
        return;
    }

    webSocketService.setMessageHandler('CHAT', (message: ChatMessage) => {
        const conversationId = message.senderID === userStore.getId
            ? message.receiverID
            : message.senderID;

        if (!conversationId) {
            console.warn('Unable to determine conversation ID for message:', message);
            return;
        }

        const newMessage: Message = {
            content: message.data,
            senderId: message.senderID,
            receiverId: message.receiverID,
            createdAt: new Date().toISOString(),
        };

        if (!conversations.value[conversationId]) {
            conversations.value[conversationId] = [];
        }
        conversations.value[conversationId].push(newMessage);

        if (conversationId !== chatStore.selectedFriendId) {
            chatStore.addUnreadMessage(conversationId);
        }
    });
    console.log("WebSocket handlers set up.");
};

const fetchFriendList = async () => {
	try {
		const fetchedFriends = await api.friendlist.getFriendList();
		if (fetchedFriends) {
			friends.value = [
                            TOURNAMENT_MASTER, 
                            ...fetchedFriends
                        ];
		}
	} catch (error) {
		console.error('Failed to fetch friend list', error);
	}
};

watch(() => chatStore.selectedFriendId, async (newFriendId) => {
    if (newFriendId === -1) {
        showChatInterface.value = false;
        currentFriendId.value = null;
    } else {
        chatStore.unreadMessagesCount[newFriendId] = 0;

        showChatInterface.value = true;
        currentFriendId.value = newFriendId;
        await loadFriendDiscussion(newFriendId);
    }
});

onMounted(() => {
    eventBus.on('CHAT_FROM_TOURNAMENT_MASTER', (message: string) => {
        if (!conversations.value[0] || message === "You just started a tournament, good luck ..") {
            conversations.value[0] = [];
            chatStore.resetUnreadMessage(0);
        }
        const formattedMessage: Message = {
            content: message,
            senderId: 0,  
            receiverId: userStore.getId ?? 0,
            createdAt: new Date().toISOString()
        };
        conversations.value[0].push(formattedMessage);

        if (currentFriendId.value !== 0) {
            chatStore.addUnreadMessage(0);
        }
    })
    fetchFriendList();
});

onUnmounted(() => {
    if (!conversations?.value[0]) {
        chatStore.resetUnreadMessage(0);
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
	color: white;
	padding: 4px 8px;
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
