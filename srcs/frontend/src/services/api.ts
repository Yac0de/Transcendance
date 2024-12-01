import authService from './authService';
import userService from './userService';
import friendlistService from './friendlistService';
import chatService from './chatService';
import gameHistoryService from './gameHistoryService';


export default {
    auth: authService,
    user: userService,
    friendlist: friendlistService,
    chat: chatService,
    gameHistory: gameHistoryService
};
