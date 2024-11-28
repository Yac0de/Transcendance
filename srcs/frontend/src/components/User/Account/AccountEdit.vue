<template>
  <div class="account-edit">
    <div class="avatar-container">
      <div class="avatar-wrapper">
        <img :src="avatarUrl" alt="User Avatar" class="avatar-image" />
        <div class="avatar-edit-overlay" @click="triggerAvatarUpload">
          <span class="edit-icon">✎</span>
        </div>
      </div>
      <input type="file" @change="handleAvatarChange" accept="image/*" id="avatar-upload" class="avatar-upload"
        ref="avatarInput" />
    </div>

    <div class="edit-fields">
      <label for="edit-nickname">Nickname:</label>
      <input id="edit-nickname" v-model="editedUser.nickname" type="text" maxlength="20" />
      <label for="edit-displayname">Display Name:</label>
      <input id="edit-displayname" v-model="editedUser.displayname" type="text" maxlength="16" />

      <div class="password-toggle">
        <label for="change-password" class="toggle-label">Change password</label>
        <label class="switch">
          <input type="checkbox" v-model="changePassword" />
          <span class="slider round"></span>
        </label>
      </div>

      <div v-if="changePassword" class="change-password">
        <div class="current-password">
          <label for="current-password">Current Password:</label>
          <input id="current-password" v-model="currentPassword" type="password" placeholder="Enter current password"
            maxlength="50" />
        </div>
        <div class="new-password">
          <label for="new-password">New Password:</label>
          <input id="new-password" v-model="newPassword" type="password" placeholder="New password" maxlength="50" />

          <input id="confirm-password" v-model="confirmPassword" type="password" placeholder="Confirm new password"
            maxlength="50" />
        </div>
      </div>
    </div>

    <div class="account-actions">
      <button class="save-button" @click="saveChanges">Save Changes</button>
      <button class="cancel-button" @click="$emit('cancelEdit')">
        <i class="fas fa-arrow-left"></i>
      </button>
    </div>

    <button class="delete-button" @click="$emit('confirmDeleteAccount')">Delete account</button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import api from '../../../services/api';
import { UserData } from '../../../types/models';

interface Props {
  user: UserData;
}

const props = defineProps<Props>();
const emit = defineEmits(['saveProfile', 'cancelEdit', 'confirmDeleteAccount', 'updateErrorMessage']);

const editedUser = ref({ ...props.user });
const newAvatarFile = ref<File | null>(null);
const avatarInput = ref<HTMLInputElement | null>(null);
const changePassword = ref<boolean>(false);
const currentPassword = ref<string>('');
const newPassword = ref<string>('');
const confirmPassword = ref<string>('');


const avatarUrl = computed(() => {
  if (newAvatarFile.value) {
    return URL.createObjectURL(newAvatarFile.value);
  }
  return editedUser.value.avatar ? api.user.getAvatarUrl(editedUser.value.avatar) : api.user.getAvatarUrl('default.png');
});

const triggerAvatarUpload = () => {
  avatarInput.value?.click();
};

const handleAvatarChange = (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0];
  if (file) {
    newAvatarFile.value = file;
    editedUser.value.avatar = file.name; // This is temporary, just for display
  }
};

const saveChanges = async () => {
  // Reset any previous error messages
  emit('updateErrorMessage', '');

  // Validate nickname and display name fields first
  if (editedUser.value.nickname.length < 3) {
    emit('updateErrorMessage', 'Nickname must be at least 3 characters long.');
    return;
  }

  if (editedUser.value.displayname.length < 3) {
    emit('updateErrorMessage', 'Display name must be at least 3 characters long.');
    return;
  }

  // Validate password change if requested
  if (changePassword.value) {
    if (currentPassword.value.length === 0) {
      emit('updateErrorMessage', 'Please enter your current password.');
      return;
    }

    if (newPassword.value.length < 6) {
      emit('updateErrorMessage', 'New password must be at least 6 characters long.');
      return;
    }

    if (newPassword.value !== confirmPassword.value) {
      emit('updateErrorMessage', 'New passwords do not match.');
      return;
    }

    try {
      await api.user.changePassword(currentPassword.value, newPassword.value);
    } catch (error) {
      emit('updateErrorMessage', 'Failed to change password. Please verify your current password.');
      return;
    }
  }

  // If everything is valid, save the profile and optionally upload avatar
  emit('saveProfile', editedUser.value, newAvatarFile.value);
};

</script>

<style scoped>
.account-edit {
  display: flex;
  flex-direction: column;
  align-items: center;
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

.edit-fields {
  color: white;
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 20px;
}

.edit-fields label {
  font-weight: bold;
  text-shadow: 0.5px 0.5px 1px black;
}

.edit-fields input {
  padding: 5px;
  border: 1px solid #ccc;
  border-radius: 4px;
}

.account-actions {
  width: 100%;
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
}

button {
  width: 48%;
  padding: 10px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  color: white;
  font-size: 14px;
  transition: background-color 0.3s;
}

.save-button {
  background: linear-gradient(90deg, var(--secondary-dark-color), var(--secondary-bright-color));
}

.save-button:hover {
  opacity: 0.8;
  transform: scale(1.03);
}

.cancel-button {
  background-color: #95a5a6;
}

.cancel-button:hover {
  background-color: #7f8c8d;
  transform: scale(1.03);
}

.delete-button {
  width: 100%;
  background-color: #e74c3c;
}

.delete-button:hover {
  background-color: #c0392b;
  transform: scale(1.02);
}

.password-toggle {
  display: flex;
  flex-direction: column;
  margin-top: 10px;
  height: 45px;
  justify-content: space-between;
}

.change-password input {
  margin-top: 5px;
}

.change-password label {
  font-size: x-small;
}

.current-password,
.new-password {
  display: flex;
  flex-direction: column;
}

.new-password {
  margin-top: 15px;
}

.new-password input {
  margin-bottom: 7px;
}

.switch {
  position: relative;
  display: inline-block;
  width: 50px;
  height: 24px;
}

.switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #ccc;
  transition: 0.4s;
  border-radius: 24px;
}

.slider:before {
  position: absolute;
  content: "";
  height: 20px;
  width: 20px;
  left: 4px;
  bottom: 2px;
  background-color: white;
  transition: 0.4s;
  border-radius: 50%;
}

input:checked+.slider {
  background-color: #4CAF50;
}

input:checked+.slider:before {
  transform: translateX(26px);
}

.slider.round {
  border-radius: 24px;
}

.slider.round:before {
  border-radius: 50%;
}

.toggle-label {
  font-size: 12px;
  vertical-align: middle;
}

@media (max-height:850px) {
  .avatar-container {
  margin-bottom: 10px;
  }
  .avatar-wrapper {
  height: 6rem;
  width: 6rem;
  }
}

@media (max-height:780px) {
  .avatar-container {
  margin-bottom: 5px;
  }
  .avatar-wrapper {
  height: 4rem;
  width: 4rem;
  }
}
</style>
