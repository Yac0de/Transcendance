<template>
  <div class="homepage-container">
    <div class="homepage-content">
      <h1 class="logo">43</h1>
      <div class="button-container">
        <button class="login-button" @click="navigateTo('login')">Login</button>
        <button class="signup-button" @click="navigateTo('signup')">Sign Up</button>
        <button class="pong-button" @click="navigateTo('pong')">Play Pong</button>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

type NavigationDestination = 'login' | 'signup' | 'pong'

export default defineComponent({
  name: 'HomeView',
  methods: {
    navigateTo(page: NavigationDestination): void {
      this.$emit('navigate', page)
    }
  }
})
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
  padding: 20px;
  border-radius: 10px;
  background-color: white;
  box-shadow: 0 0 10px rgba(0,0,0,0.1);
}
.logo {
  font-size: 120px;
  color: #333;
  margin-bottom: 30px;
}
.button-container {
  display: flex;
  justify-content: center;
  gap: 20px;
  flex-wrap: wrap;
}
button {
  width: 120px;
  padding: 10px;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
  transition: background-color 0.3s;
}
.login-button {
  background-color: #4CAF50;
}
.login-button:hover {
  background-color: #45a049;
}
.signup-button {
  background-color: #3498db;
}
.signup-button:hover {
  background-color: #2980b9;
}
.pong-button {
  background-color: #e74c3c;
}
.pong-button:hover {
  background-color: #c0392b;
}
</style>
