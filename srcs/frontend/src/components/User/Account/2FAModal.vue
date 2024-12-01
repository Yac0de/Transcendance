<template>
  <div v-if="show" class="popup-container">
    <div class="popup">
      <h2>Two-Factor Authentication</h2>
      <div v-if="currentStep === 'verify'">
        <p>Please verify your password to continue</p>
        <input 
          type="password" 
          v-model="password" 
          placeholder="Enter your password"
          class="input-field"
        />
        <button class="confirm-button" @click="handlePasswordVerify">Verify</button>
      </div>
      <div v-if="currentStep === 'code'">
        <p>Enter the verification code sent to your email</p>
        <input 
          type="text" 
          v-model="verificationCode" 
          placeholder="Enter 6-digit code"
          class="input-field"
        />
        <button class="confirm-button" @click="handleCodeVerify">Verify Code</button>
      </div>
      <button class="close-button" @click="handleClose">Cancel</button>
    </div>
  </div>
</template>

<script setup lang="ts">
// Script remains the same as it handles the logic, not the presentation
import { ref, watch } from 'vue';

interface Props {
  show: boolean;
}

const emit = defineEmits(['close', 'passwordVerify', 'codeVerify']);
const props = defineProps<Props>();

const currentStep = ref<'verify' | 'code'>('verify');
const password = ref('');
const verificationCode = ref('');

watch(() => props.show, (newValue) => {
  if (!newValue) {
    resetState();
  }
});

const handlePasswordVerify = () => {
  emit('passwordVerify', password.value);
  currentStep.value = 'code';
};

const handleCodeVerify = () => {
  emit('codeVerify', verificationCode.value);
};

const handleClose = () => {
  emit('close');
  resetState();
};

const resetState = () => {
  currentStep.value = 'verify';
  password.value = '';
  verificationCode.value = '';
};
</script>

<style scoped>
.popup-container {
  /* Creates a positioning context that's fixed relative to the viewport */
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  /* Centers the popup and creates a semi-transparent backdrop */
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 1000;
}

.popup {
  /* Controls the size of the popup */
  width: 30rem;
  max-width: 90vw;
  min-height: 60vh;
  /* Creates space for content */
  padding: 2rem;
  /* Matches your app's theme */
  background-color: var(--secondary-dark-color);
  border-radius: 8px;
  /* Adds subtle depth */
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

h2 {
  color: white;
  margin-bottom: 1.5rem;
  font-size: 1.5rem;
}

.input-field {
  width: 100%;
  padding: 0.75rem;
  margin: 1rem 0;
  border: 1px solid var(--primary-color);
  border-radius: 4px;
  background-color: var(--background-light-color);
  color: white;
  font-size: 1rem;
}

.confirm-button, .close-button {
  width: 100%;
  padding: 0.75rem;
  margin-top: 1rem;
  border: none;
  border-radius: 4px;
  color: white;
  cursor: pointer;
  font-size: 1rem;
  transition: opacity 0.2s;
}

.confirm-button {
  background: var(--primary-color);
}

.close-button {
  background: var(--secondary-dark-color);
  border: 1px solid var(--primary-color);
}

.confirm-button:hover, .close-button:hover {
  opacity: 0.9;
}
</style>
