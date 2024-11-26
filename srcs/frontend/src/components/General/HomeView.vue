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
/* @font-face {
  font-family: 'Chailcy Noggin';
  src: url("../../fonts/ChailceNogginRegular.ttf");
} */

.homepage-container {
  font-family: "Audiowide", sans-serif;
  font-weight: 400;
  font-style: normal;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  min-height: 250px;
}

.homepage-content {
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
}

.welcome {
  color: white;
  text-shadow: 1px 1px 2px black;
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
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.2);
  transition: background-color 0.3s;
}

.pong-button:hover {
  opacity: 0.8;
}

@media (max-width: 1320px), (max-height: 655px) {
  .welcome h1:nth-of-type(1) {
  font-size: 2rem;
  }
  .welcome h1:nth-of-type(2) {
  font-size: 3rem;
  }
  .pong-button {
    font-size: 2rem;
  }
}

@media (max-width: 800px), (max-height: 550px) {
  .welcome h1:nth-of-type(1) {
  font-size: 1.5rem;
  }
  .welcome h1:nth-of-type(2) {
  font-size: 2.5rem;
  }
  .pong-button {
    font-size: 1.5rem;
  }
}

@media (max-width: 768px), (max-height: 430px) {
  .welcome h1:nth-of-type(1) {
  font-size: 1.4rem;
  }
  .welcome h1:nth-of-type(2) {
  font-size: 2rem;
  }
  /* .pong-button {
    font-size: 1.5rem;
  } */
}

@media (max-width: 520px), (max-height: 400px)  {
  .welcome h1:nth-of-type(1) {
  font-size: 1rem;
  }
  .welcome h1:nth-of-type(2) {
  font-size: 1.4rem;
  }
  .pong-button {
    font-size: 1.2rem;
  }
}

@media (max-width: 370px), (max-height: 250px) {
  .welcome h1:nth-of-type(1) {
  font-size: 0.8rem;
  }
  .welcome h1:nth-of-type(2) {
  font-size: 1rem;
  }
  .pong-button {
    font-size: 1rem;
  }
}

@media (max-height: 200px) {
  .welcome h1:nth-of-type(1) {
  font-size: 0.6rem;
  }
  .welcome h1:nth-of-type(2) {
  font-size: 1rem;
  }
  .pong-button {
    font-size: 0.8rem;
  }
}
</style>


