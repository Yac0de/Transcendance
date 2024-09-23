<template>
  <div class="login-container">
    <div class="login-form">
      <h2>Login</h2>
      <form @submit.prevent="handleSubmit">
        <div>
          <label for="username">Username:</label>
          <input type="text" id="username" v-model="username" required>
        </div>
        <div>
          <label for="password">Password:</label>
          <input type="password" id="password" v-model="password" required>
        </div>
        <button type="submit">Login</button>
      </form>
      <button class="signup-button" @click="handleSignup">Sign Up</button>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'

export default defineComponent({
  name: 'LoginForm',
  setup(props, { emit }) {
    const username = ref('')
    const password = ref('')

    const handleSubmit = () => {
      console.log('Login attempted with username:', username.value)
      alert('Login functionality not implemented yet')
    }

    const handleSignup = () => {
      emit('switch-to-signup')
    }

    return {
      username,
      password,
      handleSubmit,
      handleSignup
    }
  }
})
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
  box-shadow: 0 0 10px rgba(0,0,0,0.1);
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
</style>
