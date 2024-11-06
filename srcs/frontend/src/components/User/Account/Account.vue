<template>
  <div v-if="userExists" class="account-container">
    <div class="account-content">
      <h2>Account Details</h2>
      <AccountView v-if="!isViewingStats && !isEditing && !isDeleting" :user="userToDisplay"
        :isOwnProfile="isOwnProfile" @startEditing="startEditing" />
      <AccountEdit v-if="isEditing && !isDeleting" :user="userToDisplay" :errorMessage="errorMessage"
        @saveProfile="saveProfile" @cancelEdit="cancelEdit" @confirmDeleteAccount="confirmDeleteAccount"
        @updateErrorMessage="errorMessage = $event" />
      <DeleteAccountPrompt v-if="isDeleting" :deleted="deleted" @deleteAccount="deleteAccount"
        @cancelDelete="cancelDelete" />

      <AccountStats v-if="isViewingStats" />
      <div v-if="!isEditing && !isDeleting">
        <button @click="toggleStats" class="toggle-button">
          {{ isViewingStats ? 'Back to Account' : 'View Statistics' }}
        </button>
      </div>

      <div v-if="successMessage" class="alert alert-success">{{ successMessage }}</div>
      <div v-if="errorMessage" class="alert alert-error">{{ errorMessage }}</div>
    </div>
  </div>
  <NotFound v-else />
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import api from '../../../services/api'
import { useUserStore } from '../../../stores/user'
import NotFound from '../../General/NotFound.vue'
import AccountView from './AccountView.vue'
import AccountEdit from './AccountEdit.vue'
import DeleteAccountPrompt from './DeleteAccountPrompt.vue'
import AccountStats from './AccountStats.vue'
import { UserData } from '../../../types/models';

const isEditing = ref(false)
const isDeleting = ref(false)
const isViewingStats = ref(false)
const deleted = ref(false)
const userToDisplay = ref<UserData>({ id: '', nickname: '', displayname: '', avatar: '' })
const successMessage = ref('')
const errorMessage = ref('')
const router = useRouter()
const route = useRoute()
const isOwnProfile = ref(false)
const userStore = useUserStore()
const userExists = ref(true)

const resetMessages = () => {
  successMessage.value = ''
  errorMessage.value = ''
}

const checkOwnProfile = async () => {
  try {
    isOwnProfile.value = userStore.getNickname === route.params.nickname
  } catch (error) {
    console.error('Error checking if own profile:', error)
    isOwnProfile.value = false
  }
}

const fetchUserData = async (nickname: string) => {
  resetMessages()
  userExists.value = true

  try {
    let userData: UserData | null

    if (nickname === userStore.getNickname) {
      userData = {
        id: userStore.getId ?? '',
        nickname: userStore.getNickname ?? '',
        displayname: userStore.getDisplayName ?? '',
        avatar: userStore.getAvatarPath ?? ''
      };
    } else {
      userData = await api.user.getProfileData(nickname)
      if (!userData) {
        throw new Error("User not found!")
      }
    }
    userToDisplay.value = { ...userData }
  } catch (error) {
    userExists.value = false
  }
}

onMounted(async () => {
  await checkOwnProfile()
  await fetchUserData(route.params.nickname as string)
})

watch(
  () => route.params.nickname,
  async (newNickname) => {
    await checkOwnProfile()
    await fetchUserData(newNickname as string)
  }
)

const startEditing = () => {
  isEditing.value = true
  resetMessages()
}

const toggleStats = () => {
  isViewingStats.value = !isViewingStats.value
}

const saveProfile = async (updatedUser: UserData, newAvatarFile: File | null) => {
  resetMessages()

  if (updatedUser.nickname.length < 3) {
    errorMessage.value = "Nickname must be at least 3 characters long!";
    return;
  }

  if (updatedUser.displayname.length < 3) {
    errorMessage.value = "Displayname must be at least 3 characters long!";
    return;
  }

  try {
    await api.user.updateUserProfile(updatedUser, newAvatarFile)
    await userStore.fetchUser()
    userToDisplay.value = { ...updatedUser, avatar: userStore.getAvatarPath ?? '' }
    successMessage.value = 'Profile updated successfully'
  } catch (error: any) {
    const errorResponse = await error.response?.json()
    if (errorResponse && errorResponse.error) {
      errorMessage.value = errorResponse.error
    } else {
      errorMessage.value = 'Error updating profile: ' + (error as Error).message
    }
  }
}

const cancelEdit = () => {
  isEditing.value = false
  resetMessages()
}

const confirmDeleteAccount = () => {
  isDeleting.value = true
  resetMessages()
}

const deleteAccount = async (password: string) => {
  resetMessages()
  try {
    await api.user.deleteUserAccount(password)
    successMessage.value = 'Account deleted successfully'
    deleted.value = true
    isEditing.value = false
    await api.auth.signout()
    userStore.clearUser()

    setTimeout(() => {
      router.push('/signin')
    }, 1000)
  } catch (error) {
    errorMessage.value = 'Error deleting account: ' + (error as Error).message
  }
}

const cancelDelete = () => {
  isDeleting.value = false
  resetMessages()
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

.toggle-button {
  width: 100%;
  padding: 10px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
  margin-top: 20px;
}

.toggle-button:hover {
  background-color: #0056b3;
}
</style>
