import { ref } from 'vue';

export interface UserData {
  nickname: string;
  displayname: string;
  avatar: string;
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
	required: bool;
	maxlength: number;
}

export interface ChatHistory {
	conversation: Message[]; 
}




