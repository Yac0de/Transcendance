export interface UserData {
  nickname: string;
  displayname: string;
  avatar: string;
}

export interface FriendRequest {
  id: string;
  nickname: string;
}

interface Credentials {
    nickname: string;
    password: string;
}

interface Friend {
    id: string;
    displayname: string;
    nickname: string;
    avatar: string;
    isOnline: boolean;
}

interface Message {
	content: string;
	senderId: string;
	receiverId: string;
	timestamp: string;
}




