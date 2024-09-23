<template>
  <div class="signup-container">
    <div class="signup-form">
      <h2>Sign Up</h2>
      <form @submit.prevent="handleSubmit">
        <div>
          <label for="username">Username:</label>
          <input type="text" id="username" v-model="username" required>
        </div>
        <div>
          <label for="password">Password:</label>
          <input type="password" id="password" v-model="password" required>
        </div>
        <div>
          <label for="confirmPassword">Confirm Password:</label>
          <input type="password" id="confirmPassword" v-model="confirmPassword" required>
        </div>
        <button type="submit">Sign Up</button>
      </form>
      <button class="login-button" @click="handleLogin">Back to Login</button>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'

export default defineComponent({
  name: 'SignUpForm',
  setup(props, { emit }) {
    const username = ref('')
    const password = ref('')
    const confirmPassword = ref('')

    const handleSubmit = () => {
      if (password.value !== confirmPassword.value) {
        alert("Passwords don't match!")
        return
      }
      console.log('Sign up attempted with username:', username.value)
      alert('Sign up functionality not implemented yet')
    }

    const handleLogin = () => {
      emit('switch-to-login')
    }

    return {
      username,
      password,
      confirmPassword,
      handleSubmit,
      handleLogin
    }
  }
})
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
  box-shadow: 0 0 10px rgba(0,0,0,0.1);
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
.login-button {
  margin-top: 10px;
  background-color: #3498db;
}
.login-button:hover {
  background-color: #2980b9;
}
</style>
