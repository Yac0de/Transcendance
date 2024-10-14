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
          <label for="password">Password:</label>
          <input type="password" id="password" v-model="password" required />
        </div>
        <div>
          <label for="confirmPassword">Confirm Password:</label>
          <input type="password" id="confirmPassword" v-model="confirmPassword" required />
        </div>
        <div v-if="successMessage" class="alert alert-success">{{ successMessage }}</div>
        <div v-if="errorMessage" class="alert alert-error">{{ errorMessage }}</div>
        <button type="submit">Sign Up</button>
      </form>
      <button class="signin-button" @click="handleSignin">Back to sign in</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '../services/api'

const nickname = ref('')
const password = ref('')
const confirmPassword = ref('')
const errorMessage = ref('')
const successMessage = ref('')
const router = useRouter()

const handleSubmit = async () => {
  if (password.value !== confirmPassword.value) {
    errorMessage.value = "Passwords don't match!";
    return;
  }

  if (nickname.value.length < 3) {
    errorMessage.value = "Nickname must be at least 3 characters long!";
    return;
  }

  if (password.value.length < 6) {
    errorMessage.value = "Password must be at least 6 characters long!";
    return;
  }

  try {
    errorMessage.value = '';
    successMessage.value = '';
    const user = await api.auth.signup({
      nickname: nickname.value,
      password: password.value
    });
    console.log('Sign up successful', user);
    successMessage.value = `Sign up successful! Welcome, ${nickname.value}!`;
    router.push('/signin');
  } catch (err: any) {
    errorMessage.value = err.message;
    console.error('Sign up failed', err);
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

.alert {
  padding: 10px;
  margin-bottom: 10px;
  margin-top: 5px;
  border-radius: 5px;
}

.alert-success {
  background-color: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.alert-error {
  background-color: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}
</style>
