<template>
  <div class="account-container">
    <div class="account-content">
      <h2>Account Details</h2>
      <div class="avatar-container">
        <div class="avatar-wrapper">
          <img :src="avatarUrl" alt="User Avatar" class="avatar-image" />
          <div v-if="isEditing" class="avatar-edit-overlay" @click="triggerAvatarUpload">
            <span class="edit-icon">✎</span>
          </div>
        </div>
        <input type="file" @change="handleAvatarChange" accept="image/*" id="avatar-upload" class="avatar-upload"
          ref="avatarInput" />
      </div>
      <div class="account-info">
        <div v-if="!isEditing">
          <p><strong>Nickname:</strong> {{ user.nickname }}</p>
          <p><strong>Email:</strong> {{ user.email }}</p>
        </div>
        <div v-if="isEditing" class="edit-fields">
          <label for="edit-nickname">Nickname:</label>
          <input id="edit-nickname" v-model="editedUser.nickname" type="text" />
          <label for="edit-email">Email:</label>
          <input id="edit-email" v-model="editedUser.email" type="email" />
        </div>
      </div>
      <div class="account-actions">
        <button v-if="!isEditing" class="edit-button" @click="editProfile">Edit Profile</button>
        <button v-if="isEditing" class="save-button" @click="saveProfile">Save Changes</button>
        <button v-if="isEditing" class="cancel-button" @click="cancelEdit">Cancel</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import api from '../services/api'

interface UserData {
  nickname: string;
  email: string;
  avatar: string;
}

const user = ref<UserData>({
  nickname: '',
  email: '',
  avatar: ''
})

const isEditing = ref(false)
const editedUser = ref<UserData>({
  nickname: '',
  email: '',
  avatar: ''
})
const avatarInput = ref<HTMLInputElement | null>(null)
const newAvatarFile = ref<File | null>(null)

const avatarUrl = computed(() => {
  if (isEditing.value && newAvatarFile.value) {
    return URL.createObjectURL(newAvatarFile.value)
  }
  if (editedUser.value.avatar) {
    return api.getAvatarUrl(editedUser.value.avatar)
  }
  return api.getAvatarUrl("default.avif")
})

const fetchUserData = async () => {
  try {
    const userData = await api.getUserData()
    user.value = userData
    editedUser.value = { ...userData }
  } catch (error) {
    console.error('Error fetching user data:', error)
  }
}

onMounted(() => {
  fetchUserData()
})

const editProfile = () => {
  isEditing.value = true
}

const saveProfile = async () => {
  try {
    const updatedUserData = await api.updateUserProfile(editedUser.value, newAvatarFile.value)
    user.value = updatedUserData
    isEditing.value = false
    newAvatarFile.value = null
    console.log('Profile updated successfully')
  } catch (error) {
    console.error('Error updating profile:', error)
    isEditing.value = false
  }
}

const cancelEdit = () => {
  editedUser.value = { ...user.value }
  isEditing.value = false
  newAvatarFile.value = null
}

const handleAvatarChange = (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (file) {
    newAvatarFile.value = file
    editedUser.value.avatar = file.name // This is temporary, just for display
  }
}

const triggerAvatarUpload = () => {
  avatarInput.value?.click()
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

.avatar-container {
  text-align: center;
  margin-bottom: 20px;
}

.avatar-wrapper {
  position: relative;
  display: inline-block;
  height: 8rem;
  width: 8rem;
}

.avatar-image {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
}

.avatar-edit-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  border-radius: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  opacity: 0;
  transition: opacity 0.3s;
}

.avatar-wrapper:hover .avatar-edit-overlay {
  opacity: 1;
}

.edit-icon {
  color: white;
  font-size: 24px;
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

.edit-fields {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 20px;
}

.edit-fields label {
  font-weight: bold;
}

.edit-fields input {
  padding: 5px;
  border: 1px solid #ccc;
  border-radius: 4px;
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

.edit-button,
.save-button {
  background-color: #3498db;
}

.edit-button:hover,
.save-button:hover {
  background-color: #2980b9;
}

.cancel-button {
  background-color: #95a5a6;
}

.cancel-button:hover {
  background-color: #7f8c8d;
}
</style>
