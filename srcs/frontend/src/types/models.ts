import { Ref } from 'vue';
import { WebSocketService } from '../services/websocketService';

export interface UserData {
  id: string;
  nickname: string;
  displayname: string;
  avatar: string;
}

export interface UserState {
  id: string | null;
  nickname: string | null;
  displayname: string | null;
  avatar: string | null;
  webSocketService: InstanceType<typeof WebSocketService> | null, 
}

export interface FriendRequest {
  id: string;
  nickname: string;
}

export interface Credentials {
    nickname: string;
    password: string;
}

export interface Friend {
    id: string;
    displayname: string;
    nickname: string;
    avatar: string;
    isOnline: boolean;
}

export interface Message {
	content: string;
	senderId: string;
	receiverId: string;
	createdAt: string;
}

export interface Field {
	label: string;
	model: Ref<string>;
	type: string;
	required: boolean;
	maxlength: number;
}

export interface ChatHistory {
	conversation: Message[]; 
}




