<template>
  <div class="login-container">
    <div class="login-form">
      <h2>Login</h2>
      <form @submit.prevent="handleSubmit">
        <div>
          <label for="nickname">Nickname:</label>
          <input type="text" id="nickname" v-model="nickname" required />
        </div>
        <div>
          <label for="password">Password:</label>
          <input type="password" id="password" v-model="password" required />
        </div>
        <button type="submit">Login</button>
      </form>
      <button class="signup-button" @click="handleSignup">Sign Up</button>
    </div>
    <div v-if="error" class="error-message">{{ error }}</div>
    <div v-if="user" class="success-message">Welcome, {{ user.name }}!</div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '../services/api'

const nickname = ref('')
const password = ref('')
const error = ref('')
const user = ref<any>(null)
const router = useRouter()

const emit = defineEmits(['login-success']);

const handleSubmit = async () => {
  try {
    error.value = ''
    user.value = await api.login({ nickname: nickname.value, password: password.value })
    console.log('Login successful', user.value)
    emit('login-success')
    router.push('/pong')
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'An error occurred during login'
    console.error('Login failed', err)
  }
}

const handleSignup = () => {
  router.push('/signup')
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}

.login-form {
  width: 300px;
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 10px;
  background-color: #f9f9f9;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
}

.login-form h2 {
  text-align: center;
  color: #333;
}

.login-form div {
  margin-bottom: 15px;
}

label {
  display: block;
  margin-bottom: 5px;
  color: #666;
}

input {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  box-sizing: border-box;
}

button {
  width: 100%;
  padding: 10px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
}

button:hover {
  background-color: #45a049;
}

.signup-button {
  margin-top: 10px;
  background-color: #3498db;
}

.signup-button:hover {
  background-color: #2980b9;
}

.error-message {
  color: red;
  margin-top: 10px;
}

.success-message {
  color: green;
  margin-top: 10px;
}
</style>
