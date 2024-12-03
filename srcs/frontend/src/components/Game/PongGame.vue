<template>
  <div class="pong-menu">
    <div class="pong-menu-container">
      <div class="menu-buttons">
        <button @click="startTournament" class="menu-button">{{ $t('tournament') }}</button>
        <button @click="startPlayerVsFriend" class="menu-button">{{ $t('playAgainstFriend') }}</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';

const router = useRouter();

const startTournament = () => {
  router.push('/tournament');
};

const startPlayerVsFriend = () => {
  router.push('/lobby');
};
</script>

<style scoped>
.pong-menu {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  color: #2c3e50;
}

.pong-menu-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: space-around;
  text-align: center;
  height: 40vh;
  min-height: 65px;
  min-width: 300px;
  padding: 2vh 5vw;
  border-radius: 20px;
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.5);
  overflow: hidden;
  background: var(--main-color);
}

.menu-title {
  font-size: 2.5rem;
  margin-bottom: 2rem;
}

.menu-buttons {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.menu-button {
  padding: 10px 20px;
  font-size: 1.2rem;
  font-weight: 600;
  background: linear-gradient(90deg, var(--secondary-dark-color), var(--secondary-bright-color));
  box-shadow: 0 0 2px rgba(0, 0, 0, 1);
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
  width: 200px;
}

.menu-button:hover {
  opacity: 0.85;
  transform: scale(1.02);
}

@media (max-width: 600px) {
  .menu-title {
    font-size: 2rem;
  }

  .menu-button {
    width: 100%;
  }
}
</style>
