<template>
  <div id="app">
    <nav class="sticky-nav">
      <div class="nav-content">
        <router-link to="/" class="nav-button">Home</router-link>
        <router-link to="/login" class="nav-button">Login</router-link>
        <router-link to="/signup" class="nav-button">Sign Up</router-link>
        <router-link to="/pong" class="nav-button">Play Pong</router-link>
        <router-link to="/account" class="nav-button">Account</router-link>
      </div>
    </nav>
    <div class="content">
      <router-view></router-view> <!-- Affiche le composant correspondant à la route -->
    </div>
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}

.sticky-nav {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  background-color: #f8f9fa;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  z-index: 1000;
}

.nav-content {
  display: flex;
  justify-content: flex-start;
  padding: 10px 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.nav-button {
  margin-right: 10px;
  padding: 5px 10px;
  font-size: 16px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  text-decoration: none;
  transition: background-color 0.3s;
}

.nav-button:hover {
  background-color: #0056b3;
}

.content {
  margin-top: 60px;
  padding: 20px;
}
</style>
