<template>
  <div class="signup-container">
    <div class="signup-form">
      <h2>Sign Up</h2>
      <form @submit.prevent="handleSubmit">
        <div>
          <label for="nickname">Nickname:</label>
          <input type="text" id="nickname" v-model="nickname" required />
        </div>
        <div>
          <label for="email">Email:</label>
          <input type="email" id="email" v-model="email" required />
        </div>
        <div>
          <label for="password">Password:</label>
          <input type="password" id="password" v-model="password" required />
        </div>
        <div>
          <label for="confirmPassword">Confirm Password:</label>
          <input type="password" id="confirmPassword" v-model="confirmPassword" required />
        </div>
        <button type="submit">Sign Up</button>
      </form>
      <button class="signin-button" @click="handleSignin">Back to sign in</button>
    </div>
    <div v-if="error" class="error-message">{{ error }}</div>
    <div v-if="successMessage" class="success-message">{{ successMessage }}</div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '../services/api'

const nickname = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const error = ref('')
const successMessage = ref('')
const router = useRouter()

const handleSubmit = async () => {
  if (password.value !== confirmPassword.value) {
    error.value = "Passwords don't match!"
    return
  }

  if (password.value.length < 6) {
    error.value = "Password must be at least 6 characters long!"
    return
  }

  try {
    error.value = ''
    successMessage.value = ''
    const user = await api.signup({
      nickname: nickname.value,
      email: email.value,
      password: password.value
    })
    console.log('Sign up successful', user)
    successMessage.value = `Sign up successful! Welcome, ${nickname.value}!`
    router.push('/signin')
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'An error occurred during sign up'
    console.error('Sign up failed', err)
  }
}

const handleSignin = () => {
  router.push('/signin')
}
</script>

<style scoped>
.signup-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}

.signup-form {
  width: 300px;
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 10px;
  background-color: #f9f9f9;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
}

.signup-form h2 {
  text-align: center;
  color: #333;
}

.signup-form div {
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

.signin-button {
  margin-top: 10px;
  background-color: #3498db;
}

.signin-button:hover {
  background-color: #2980b9;
}

.error-message {
  color: red;
  margin-top: 10px;
  text-align: center;
}

.success-message {
  color: green;
  margin-top: 10px;
  text-align: center;
}
</style>
