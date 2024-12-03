<template>
  <div class="game-header">
    <div class="scoreP1">
      <div class="player-info">
        <div class="avatar-and-name">
          <div class="avatar-container">
            <img :src="api.user.getAvatarUrl(player1?.avatar ?? null)" class="avatar" :alt="$t('player1AvatarAlt')" />
          </div>
          <span :style="nameStyle(player1)">{{ player1?.displayname || $t('player1') }}</span>
        </div>
        <span>{{ state.score.player1 }}</span>
      </div>
    </div>
    <div class="timer-container">
      <div class="timer-title">{{ $t('time') }}
        <div class="timer">{{ state.elapsedTime }}s</div>
      </div>
    </div>
    <div class="scoreP2">
      <div class="player-info">
        <div class="avatar-and-name">
          <div class="avatar-container">
            <img :src="api.user.getAvatarUrl(player2?.avatar ?? null)" class="avatar" :alt="$t('player2AvatarAlt')" />
          </div>
          <span :style="nameStyle(player2)">{{ player2?.displayname || $t('player2') }}</span>
        </div>
        <span>{{ state.score.player2 }}</span>
      </div>
    </div>
  </div>
</template>
  
<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { GameState } from '../../types/game';
import { UserData } from '../../types/models'
import { fetchUserById } from '../../utils/fetch'
import api from '../../services/api';

const props = defineProps<{
  state: GameState
  player1id: number | null
  player2id: number | null
}>()

const player1 = ref<UserData | null>(null)
const player2 = ref<UserData | null>(null)

const fetchPlayerData = async () => {
  if (props.player1id) {
    player1.value = await fetchUserById(props.player1id ?? null)
  }
  if (props.player2id) {
    player2.value = await fetchUserById(props.player2id ?? null)
  }
}

onMounted(fetchPlayerData)

watch([() => props.player1id, () => props.player2id], fetchPlayerData)

const nameStyle = (player: UserData | null) => {
  const displayname = player?.displayname || '';
  const minFontSize = 15;
  const maxFontSize = 24;
  const minLength = 3;

  const length = displayname.length;
  const fontSize = Math.max(minFontSize, Math.min(maxFontSize, maxFontSize - (length - minLength) * 1.5));

  return {
    fontSize: `${fontSize}px`,
  };
}
</script>

<style scoped>
.game-header {
  position: relative;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: linear-gradient(to right, var(--secondary-dark-color), color-mix(in srgb, var(--secondary-dark-color) 75%, white));
  color: black;
  width: 800px;
  height: 50px;
  font-size: x-large;
  font-weight: bolder;
}

.scoreP1, .scoreP2 {
  width: 40%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.player-info {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  color: white;
  font-size: 1.2rem;
  text-shadow: 0.5px 0.5px 1px black;
}

.avatar-and-name {
  display: flex;
  align-items: center;
  width: 250px;
}

.avatar-container {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  border: 2px solid var(--secondary-bright-color);
  background-color: var(--secondary-dark-color);
}

.avatar {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.timer-container {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: black;
  clip-path: polygon(0% 0, 100% 0, 80% 100%, 20% 100%);
  width: 25%;
  height: 98%;
  font-size: 1rem;
}

.timer-title {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--secondary-dark-color);
  text-shadow: 
    1px 1px 0 white,
    -1px 1px 0 white,
    1px -1px 0 white,
    -1px -1px 0 white;
  font-weight: bold;       
  letter-spacing: .2rem;
}

.timer {
  letter-spacing: normal;
  color: var(--secondary-bright-color);
  text-shadow: none;
}

.scoreP1 .player-info {
  flex-direction: row;
  margin-left: 20px;
  margin-right: 25px;
}

.scoreP1 .avatar-and-name span{
  margin-left: 15px;
}

.scoreP2 .player-info {
  flex-direction: row-reverse;
  margin-right: 20px;
  margin-left: 25px;
}

.scoreP2 .avatar-and-name {
  flex-direction: row-reverse;
}

.scoreP2 .avatar-and-name span{
  margin-right: 15px;
}

/* Optional: Add hover effect on avatar */
.avatar-container:hover {
  border-color: white;
  transform: scale(1.05);
  transition: all 0.2s ease-in-out;
}
</style>
