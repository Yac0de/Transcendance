export interface ChatMessage {
  type: 'CHAT';
  data: string;
  senderID: number;
  receiverID: number;
}
