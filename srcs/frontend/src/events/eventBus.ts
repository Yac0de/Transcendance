import mitt from 'mitt'
import { GameEvent } from '../types/game.ts';

import type { 
  LobbyInvitationToFriend,
  LobbyInvitationFromFriend,
  LobbyAcceptFromFriend,
  LobbyDenyFromFriend,
  LobbyCreated,
  LobbyPlayerStatus,
  LobbyTerminate,
  LobbyPregameRemainingTime,
} from '../types/lobby';

import type { 
   TournamentJoinWithCode,
   TournamentCreate,
   TournamentEvent
} from '../types/tournament';

type Events = {
  'LOBBY_INVITATION_TO_FRIEND': LobbyInvitationToFriend;
  'LOBBY_INVITATION_FROM_FRIEND': LobbyInvitationFromFriend;
  'LOBBY_ACCEPT_FROM_FRIEND': LobbyAcceptFromFriend;
  'LOBBY_DENY_FROM_FRIEND': LobbyDenyFromFriend;
  'LOBBY_CREATED': LobbyCreated;
  'LOBBY_PLAYER_STATUS': LobbyPlayerStatus;
  'LOBBY_TERMINATE': LobbyTerminate;
  'LOBBY_PREGAME_REMAINING_TIME': LobbyPregameRemainingTime;
  'LOBBY_DESTROYED': void;
  'GAME_EVENT' : GameEvent;
  'TOURNAMENT_JOIN_WITH_CODE': TournamentJoinWithCode
  'TOURNAMENT_CREATE': TournamentCreate
  'TOURNAMENT_EVENT': TournamentEvent
}

export const eventBus = mitt<Events>();
