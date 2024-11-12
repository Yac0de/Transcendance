import mitt from 'mitt'

import type { 
  LobbyInvitationToFriend,
  LobbyInvitationFromFriend,
  LobbyAcceptFromFriend,
  LobbyDenyFromFriend,
  LobbyTerminate
} from '../types/lobby';

type Events = {
  'LOBBY_INVITATION_TO_FRIEND': LobbyInvitationToFriend;
  'LOBBY_INVITATION_FROM_FRIEND': LobbyInvitationFromFriend;
  'LOBBY_ACCEPT_FROM_FRIEND': LobbyAcceptFromFriend;
  'LOBBY_DENY_FROM_FRIEND': LobbyDenyFromFriend;
  'LOBBY_TERMINATE': LobbyTerminate;
}

export const eventBus = mitt<Events>();
