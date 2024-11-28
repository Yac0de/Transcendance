<template>
  <div class="pong-menu">
    <div class="menu-buttons">
      <button @click="startTournament" class="menu-button">Tournament</button>
      <button @click="startPlayerVsRobot" class="menu-button">Player VS Robot</button>
      <button @click="startPlayerVsFriend" class="menu-button">Play against your friend</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';

const router = useRouter();

const startTournament = () => {
  // TODO: Implement tournament logic 
  console.log('Starting tournament');
};

const startPlayerVsRobot = () => {
  // TODO: Implement player vs robot logic 
  console.log('Starting Player VS Robot game');
};

const startPlayerVsFriend = () => {
  // TODO: Implement player vs friend logic 
  console.log('Starting Play VS Friend game');
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
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
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
