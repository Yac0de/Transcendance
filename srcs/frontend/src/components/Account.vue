<template>
  <div class="account-container">
    <div class="account-content">
      <h2>Account Details</h2>
      <div class="avatar-container">
        <img :src="avatarUrl" alt="User Avatar" class="avatar-image" />
        <input type="file" @change="uploadAvatar" accept="image/*" id="avatar-upload" class="avatar-upload" />
        <label for="avatar-upload" class="avatar-upload-label">Change Avatar</label>
      </div>
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
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../services/api'

const user = ref({
  nickname: '',
  email: '',
  avatar: ''
})

const router = useRouter()

const avatarUrl = computed(() => {
  if (user.value.avatar) {
    return api.getAvatarUrl(user.value.avatar)
  }
    return api.getAvatarUrl("default.png")
})

const fetchUserData = async () => {
  try {
    const userData = await api.getUserData()
    user.value = userData
  } catch (error) {
    console.error('Error fetching user data:', error)
  }
}

onMounted(() => {
  fetchUserData()
})

const editProfile = () => {
  console.log('Edit profile clicked')
}

const logout = () => {
  console.log('User logged out')
  router.push('/')
}

const uploadAvatar = async (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (file) {
    try {
      const avatarPath = await api.uploadAvatar(file)
      user.value.avatar = avatarPath
      console.log('Avatar updated')
    } catch (error) {
      console.error('Error uploading avatar:', error)
    }
  }
}
</script>

<style scoped>
/* Styles remain the same as in the previous version */
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

.avatar-container {
  text-align: center;
  margin-bottom: 20px;
}

.avatar-image {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  object-fit: cover;
}

.avatar-upload {
  display: none;
}

.avatar-upload-label {
  display: inline-block;
  padding: 8px 12px;
  background-color: #3498db;
  color: white;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  margin-top: 10px;
}

.avatar-upload-label:hover {
  background-color: #2980b9;
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
