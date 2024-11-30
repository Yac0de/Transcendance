import { Ref } from 'vue';
import { WebSocketService } from '../services/websocketService';

export interface UserData {
  id: number;
  nickname: string;
  displayname: string;
  avatar: string;
}

export interface UserState {
  id: number | null;
  nickname: string | null;
  displayname: string | null;
  avatar: string | null;
  webSocketService: InstanceType<typeof WebSocketService> | null; 
  isRedirectPending: boolean;
}

export interface GameHistory {
  id: number
  player1_id: number
  player2_id: number
  CreatedAt: string 
  UpdatedAt: string
  DeletedAt: string | null
  winner_id: number
  score1: number
  score2: number
  is_winner: boolean
  player1: UserData
  player2: UserData
}

export interface Credentials {
    nickname: string;
    password: string;
}

export interface Friend {
    id: number;
    displayname: string;
    nickname: string;
    avatar: string;
    isOnline: boolean;
}

export interface Message {
	content: string;
	senderId: number;
	receiverId: number;
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




