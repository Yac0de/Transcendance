<template>
  <div class="homepage-container">
    <div class="homepage-content">
      <div class="welcome">
        <h1>WELCOME TO</h1>
        <h1>TRANSCENDENCE PONG</h1>
      </div>
      <router-link to="/pong" class="pong-button">PLAY NOW</router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
@font-face {
  font-family: 'Chailcy Noggin';
  src: url("../../fonts/ChailceNogginRegular.ttf");
}

.homepage-container {
  font-family: 'Chailcy Noggin', 'Courier New', Courier, monospace;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}

.homepage-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: space-around;
  text-align: center;
  height: 40vh;
  padding: 2vh 5vw;
  border-radius: 20px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
}

.welcome {
  color: white;
}

.welcome h1:nth-of-type(1) {
  font-size: 3rem;
}

.welcome h1:nth-of-type(2) {
  font-size: 5rem;
}

.pong-button {
  display: inline-block;
  padding: 1vh 1vw;
  width: auto;
  color: white;
  background: linear-gradient(90deg, #2e151b, #da7b93);
  border-radius: 50px;
  cursor: pointer;
  font-size: 3rem;
  letter-spacing: 0.12rem;
  font-weight: 650;
  text-align: center;
  text-decoration: none;
  transition: background-color 0.3s;
}

.pong-button:hover {
  opacity: 0.8;
}
</style>
