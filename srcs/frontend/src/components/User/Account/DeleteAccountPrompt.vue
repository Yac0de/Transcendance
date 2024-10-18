<template>
  <div class="delete-account-prompt">
    <h3>Delete Account</h3>
    <p>Are you sure you want to delete your account?</p>
    <p>This action cannot be undone.</p>
    <div v-if="!deleted" class="password-input-container">
      <div class="password-input">
        <label for="delete-password">Please enter your password to confirm:</label>
        <input type="password" id="delete-password" v-model="password" placeholder="Enter your password" />
      </div>
    </div>
    <div v-if="!deleted" class="delete-actions">
      <button class="confirm-delete-button" @click="confirmDelete">Confirm Delete</button>
      <button class="cancel-delete-button" @click="$emit('cancelDelete')">Cancel</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

interface Props {
  deleted: boolean;
}

const props = defineProps<Props>();
const emit = defineEmits(['deleteAccount', 'cancelDelete']);
const password = ref('');

const confirmDelete = () => {
  if (password.value.trim() !== '') {
    emit('deleteAccount', password.value);
  }
};
</script>

<style scoped>
.delete-account-prompt {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  max-width: 400px;
  margin: 0 auto;
  padding: 20px;
}

h3 {
  color: #e74c3c;
  margin-bottom: 15px;
}

p {
  margin-bottom: 10px;
}

.password-input-container {
  width: 100%;
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
}

.password-input {
  width: 100%;
  max-width: 300px;
}

.password-input label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
  text-align: left;
}

.password-input input {
  width: 100%;
  padding: 8px;
  border: 1px solid #ccc;
  border-radius: 4px;
}

.delete-actions {
  width: 100%;
  display: flex;
  justify-content: center;
  gap: 10px;
}

button {
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  color: white;
  font-size: 14px;
  transition: background-color 0.3s;
}

.confirm-delete-button {
  background-color: #e74c3c;
}

.confirm-delete-button:hover {
  background-color: #c0392b;
}

.cancel-delete-button {
  background-color: #95a5a6;
}

.cancel-delete-button:hover {
  background-color: #7f8c8d;
}
</style>
