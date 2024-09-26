<template>
    <div class="account-container">
      <div class="account-content">
        <h2>Account Details</h2>
        <div class="account-info">
          <p><strong>Nickname:</strong> {{ user.nickname }}</p>
          <p><strong>Email:</strong> {{ user.email }}</p>
        </div>
        <div class="account-actions">
          <button class="edit-button" @click="editProfile">Edit Profile</button>
          <button class="logout-button" @click="logout">Logout</button>
        </div>
      </div>
    </div>
  </template>

  <script setup lang="ts">
  import { ref } from 'vue'
  import { useRouter } from 'vue-router'

  // Simulated user data (replace with actual API data or store in a real scenario)
  const user = ref({
    nickname: 'Test',
    email: 'exemple@example.com'
  })

  const router = useRouter()

  const editProfile = () => {
    // Logic to edit the profile (could open a modal or navigate to a profile editing route)
    console.log('Edit profile clicked')
  }

  const logout = () => {
    // Example logout logic: Clear session or token
    console.log('User logged out')

    // Redirect to the home page after logging out
    router.push('/')
  }
  </script>

  <style scoped>
  .account-container {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;
  }

  .account-content {
    padding: 20px;
    border: 1px solid #ccc;
    border-radius: 10px;
    background-color: #f9f9f9;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
    width: 300px;
  }

  h2 {
    text-align: center;
    margin-bottom: 20px;
  }

  .account-info p {
    font-size: 16px;
    margin-bottom: 10px;
  }

  .account-actions {
    display: flex;
    justify-content: space-between;
  }

  button {
    width: 45%;
    padding: 10px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    color: white;
    font-size: 14px;
    transition: background-color 0.3s;
  }

  .edit-button {
    background-color: #3498db;
  }

  .edit-button:hover {
    background-color: #2980b9;
  }

  .logout-button {
    background-color: #e74c3c;
  }

  .logout-button:hover {
    background-color: #c0392b;
  }
  </style>