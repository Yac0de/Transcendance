<template>
  <div class="homepage-container">
    <div class="homepage-content">
      <h1 class="logo">2001</h1>
      <div class="button-container">
        <router-link to="/pong" class="pong-button">Play Pong</router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
.homepage-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: white;
}

.homepage-content {
  text-align: center;
  padding: 80px;
  border-radius: 20px;
  background-color: white;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
}

.logo {
  font-size: 120px;
  color: #333;
  margin-bottom: 50px;
}

.button-container {
  display: flex;
  justify-content: center;
  gap: 20px;
  flex-wrap: wrap;
}

.button-container a {
  display: inline-block;
  width: 120px;
  padding: 10px;
  color: white;
  background-color: #4CAF50;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
  text-align: center;
  text-decoration: none;
  transition: background-color 0.3s;
}

.button-container a.pong-button {
  background-color: #dfba1c;
}

.button-container a:hover {
  opacity: 0.8;
}
</style>
