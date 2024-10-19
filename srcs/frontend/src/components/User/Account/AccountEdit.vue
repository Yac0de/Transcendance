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
      <input id="edit-nickname" v-model="editedUser.nickname" type="text" />
      <label for="edit-displayname">Display Name:</label>
      <input id="edit-displayname" v-model="editedUser.displayname" type="text" />
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

interface UserData {
  nickname: string;
  displayname: string;
  avatar: string;
}

interface Props {
  user: UserData;
}

const props = defineProps<Props>();
const emit = defineEmits(['saveProfile', 'cancelEdit', 'confirmDeleteAccount']);

const editedUser = ref({ ...props.user });
const newAvatarFile = ref<File | null>(null);
const avatarInput = ref<HTMLInputElement | null>(null);

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

const saveChanges = () => {
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
  width: 100%;
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
  background-color: #3498db;
}

.save-button:hover {
  background-color: #2980b9;
}

.cancel-button {
  background-color: #95a5a6;
}

.cancel-button:hover {
  background-color: #7f8c8d;
}

.delete-button {
  width: 100%;
  background-color: #e74c3c;
}

.delete-button:hover {
  background-color: #c0392b;
}
</style>
