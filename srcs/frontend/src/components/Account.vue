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
        <div v-if="!isEditing && !isDeleting">
          <p><strong>Nickname:</strong> {{ user?.nickname }}</p>
          <p><strong>Display Name:</strong> {{ user?.displayname }}</p>
        </div>
        <div v-if="isEditing && !isDeleting" class="edit-fields">
          <label for="edit-nickname">Nickname:</label>
          <input id="edit-nickname" v-model="editedUser.nickname" type="text" />
          <label for="edit-displayname">Display Name:</label>
          <input id="edit-displayname" v-model="editedUser.displayname" type="text" />
          <button v-if="isEditing" class="delete-button" @click="confirmDeleteAccount">Delete account</button>
        </div>

        <div v-if="isDeleting" class="confirm-delete-container">
          <p>Are you sure you want to delete your account?</p>
          <p>Please enter your password to confirm:</p>
          <input type="password" v-model="deletePassword" placeholder="Enter your password" />
          <div v-if="successMessage" class="alert alert-success">{{ successMessage }}</div>
          <div v-if="errorMessage" class="alert alert-error">{{ errorMessage }}</div>
          <div v-if="!deleted" class="delete-actions">
            <button class="confirm-delete-button" @click="deleteAccount">Confirm Delete</button>
            <button class="cancel-delete-button" @click="cancelDelete">Cancel</button>
          </div>
        </div>
      </div>

      <div v-if="successMessage && isEditing && !isDeleting" class="alert alert-success">{{ successMessage }}</div>
      <div v-if="errorMessage && isEditing && !isDeleting" class="alert alert-error">{{ errorMessage }}</div>

      <div v-if="!deleted" class="account-actions">
        <button v-if="!isEditing && isOwnProfile" class="edit-button" @click="startEditing">Edit Profile</button>
        <button v-if="isEditing && !isDeleting" class="save-button" @click="saveProfile">Save Changes</button>
        <button v-if="isEditing && !isDeleting" class="cancel-button" @click="cancelEdit">Cancel</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import api from '../services/api'

interface UserData {
  nickname: string;
  displayname: string;
  avatar: string;
}

const user = ref<UserData | null>(null)
const isEditing = ref(false)
const isDeleting = ref(false)
const deleted = ref(false)
const editedUser = ref<UserData>({ nickname: '', displayname: '', avatar: '' })
const avatarInput = ref<HTMLInputElement | null>(null)
const newAvatarFile = ref<File | null>(null)
const deletePassword = ref('')
const successMessage = ref('')
const errorMessage = ref('')
const router = useRouter()
const route = useRoute()
const isOwnProfile = ref(false)

const props = defineProps<{
  nickname: string
}>()

const resetMessages = () => {
  successMessage.value = ''
  errorMessage.value = ''
}


const checkOwnProfile = async () => {
  try {
    const currentUserData = await api.user.getUserData()
    const routeNickname = route.params.nickname as string
    isOwnProfile.value = currentUserData?.nickname === routeNickname
    console.log(isOwnProfile.value)
  } catch (error) {
    console.error('Error checking if own profile:', error)
    isOwnProfile.value = false
    console.log("not own profile")
  }
}

const avatarUrl = computed(() => {
  if (isEditing.value && newAvatarFile.value) {
    return URL.createObjectURL(newAvatarFile.value)
  }
  return editedUser.value.avatar ? api.user.getAvatarUrl(editedUser.value.avatar) : api.user.getAvatarUrl('default.png')
})

const fetchUserData = async () => {
  console.log("fetch user data")
  try {
    const userData = await api.user.getProfileData(props.nickname)
    if (userData) {
      user.value = userData
      editedUser.value = { ...userData }
    } else {
      user.value = null
    }
  } catch (error) {
    console.error('Error fetching user data:', error)
  }
}

onMounted(async () => {
  console.log("Component mounted")
  await checkOwnProfile()
  await fetchUserData()
})

const startEditing = () => {
  isEditing.value = true
  resetMessages()
}

const saveProfile = async () => {
  resetMessages()
  if (editedUser.value.nickname.length < 3) {
    errorMessage.value = "Nickname must be at least 3 characters long!";
    return;
  }

  if (editedUser.value.displayname.length < 3) {
    errorMessage.value = "Display Name must be at least 3 characters long!";
    return;
  }
  try {
    await api.user.updateUserProfile(editedUser.value, newAvatarFile.value)
    await fetchUserData()
    successMessage.value = 'Profile updated successfully'
  } catch (error: any) {
    const errorResponse = await error.response?.json();
    if (errorResponse && errorResponse.error) {
      errorMessage.value = errorResponse.error;
    } else {
      errorMessage.value = 'Error updating profile: ' + (error as Error).message;
    }
  }
}

const cancelEdit = () => {
  if (user.value) {
    editedUser.value = { ...user.value }
  }
  isEditing.value = false
  newAvatarFile.value = null
  resetMessages()
}

const cancelDelete = () => {
  deletePassword.value = ''
  isDeleting.value = false
}

const confirmDeleteAccount = () => {
  isDeleting.value = true
  resetMessages()
}

const deleteAccount = async () => {
  resetMessages()
  try {
    if (!deletePassword.value) {
      errorMessage.value = 'Password required to delete account'
      return
    }

    await api.user.deleteUserAccount(deletePassword.value)
    successMessage.value = 'Account deleted successfully'
    deleted.value = true
    isEditing.value = false
    await api.auth.signout()

    setTimeout(() => {
      router.push('/signin')
    }, 2000)
  } catch (error) {
    errorMessage.value = 'Error deleting account: ' + (error as Error).message
  }
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

.edit-fields input,
.confirm-delete-container input {
  padding: 5px;
  border: 1px solid #ccc;
  border-radius: 4px;
}

.delete-actions {
  margin-top: 20px;
}

.account-actions,
.delete-actions {
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

.delete-button,
.confirm-delete-button {
  background-color: #dc3545;
  border-color: #dc3545;
}

.cancel-button {
  color: #fff;
  background-color: #343a40;
  border-color: #343a40;
}

.delete-button:hover,
.confirm-delete-button:hover {
  background-color: #c82333;
  border-color: #bd2130;
}

.cancel-button:hover {
  background-color: #23272b;
  border-color: #1d2124;
}

.edit-button,
.save-button,
.cancel-delete-button {
  background-color: #3498db;
}

.edit-button:hover,
.save-button:hover,
.cancel-delete-button:hover {
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
