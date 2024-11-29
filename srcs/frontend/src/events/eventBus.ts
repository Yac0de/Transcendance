import mitt from 'mitt'
import type { GameEvent, GameStart, GameFinished } from '../types/game.ts';

import type { 
  LobbyInvitationToFriend,
  LobbyInvitationFromFriend,
  LobbyAcceptFromFriend,
  LobbyDenyFromFriend,
  LobbyCreated,
  LobbyPlayerStatus,
  LobbyTerminate,
  LobbyPregameRemainingTime,
  SpecialModeStatus,
} from '../types/lobby';

import type { 
   TournamentJoinWithCode,
   TournamentCreate,
   TournamentEvent,
   TournamentStart,
   TournamentTimer,
   TournamentGame,
   TournamentTreeState,
   TournamentError
} from '../types/tournament';

type Events = {
  'LOBBY_INVITATION_TO_FRIEND': LobbyInvitationToFriend;
  'LOBBY_INVITATION_FROM_FRIEND': LobbyInvitationFromFriend;
  'LOBBY_ACCEPT_FROM_FRIEND': LobbyAcceptFromFriend;
  'LOBBY_DENY_FROM_FRIEND': LobbyDenyFromFriend;
  'LOBBY_CREATED': LobbyCreated;
  'LOBBY_PLAYER_STATUS': LobbyPlayerStatus;
  'LOBBY_SPECIAL_MODE_TOGGLED': SpecialModeStatus;
  'LOBBY_TERMINATE': LobbyTerminate;
  'LOBBY_PREGAME_REMAINING_TIME': LobbyPregameRemainingTime;
  'LOBBY_DESTROYED': void;
  'GAME_EVENT' : GameEvent;
  'GAME_START': GameStart;
  'GAME_FINISHED': GameFinished;
  'TOURNAMENT_JOIN_WITH_CODE': TournamentJoinWithCode
  'TOURNAMENT_CREATE': TournamentCreate
  'TOURNAMENT_EVENT': TournamentEvent
  'TOURNAMENT_START': TournamentStart
  'TOURNAMENT_TERMINATE': void
  'TOURNAMENT_TIMER': TournamentTimer
  'TOURNAMENT_GAME': TournamentGame
  'TOURNAMENT_TREE_STATE': TournamentTreeState
  'TOURNAMENT_ERROR': TournamentError
}

export const eventBus = mitt<Events>();
