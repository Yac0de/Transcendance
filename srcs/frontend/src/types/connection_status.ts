export interface OnlineUsersMessage {
  type: 'ONLINE_USERS';
  usersOnline: number[];
}

export interface UserStatusMessage {
  type: 'USER_DISCONNECTED | NEW_CONNECTION';
  user: number;
}
