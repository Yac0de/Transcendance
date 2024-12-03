<template>
  <AuthForm v-if="!show2fa"
    formTitle="formTitleSignIn"
    submitButtonLabel="submitButtonLabelSignIn"
    :fields="fields"
    :successMessage="successMessage"
    :errorMessage="errorMessage"
    @submit="handleSubmit"
    @secondaryAction="handleSignup"
    secondaryButtonLabel="secondaryButtonLabelSignUp" 
  />
  <div v-else class="confirmTwoFA">
    <label for="confirmationCode">Please enter the confirmation code:</label>
      <form
        class="container-input-confirm-2fa"
        @submit.prevent="confirm2FA"
      >
        <input
          v-model="confirmationCode"
          type="text"
          id="confirmationCode"
          required
          placeholder="Enter 2FA code"
        />
        <input type="submit" value="Send" />
      </form>
      <div
        v-if="errorMessage || successMessage"
        :class="{
          'error-message': errorMessage,
          'success-message': successMessage
        }"
      >
        {{ errorMessage || successMessage }}
      </div>
  </div>
</template>


<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import AuthForm from './AuthForm.vue';
import { useUserStore } from '../../stores/user';
import { Field } from '../../types/models';
import { useI18n } from 'vue-i18n';
import { API_BASE_URL } from '../../services/apiUtils';

const nickname = ref('');
const password = ref('');
const errorMessage = ref('');
const successMessage = ref('');
const router = useRouter();
const userStore = useUserStore();

const { t } = useI18n();

const show2fa = ref<boolean>(false);
const confirmationCode = ref('');

const fields: Field[] = [
  { label: 'nicknameField', model: nickname, type: 'text', required: true, maxlength: 16 },
  { label: 'passwordField', model: password, type: 'password', required: true, maxlength: 50 },
];

const handleSubmit = async () => {
  if (nickname.value.length < 3) {
    errorMessage.value = t('errorMessageNicknameTooShort');
    return;
  }

  if (password.value.length < 6) {
    errorMessage.value = t('errorMessagePasswordTooShort');
    return;
  }
  try {
    errorMessage.value = '';
    successMessage.value = '';
    const response = await fetch(`${API_BASE_URL}/auth/signin`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: "include",
            body: JSON.stringify({ nickname: nickname.value, password: password.value }),
        });
    let result = await response.json()
    console.log("RESULT", result)

    if (response.status === 202) {
      successMessage.value = 'Two-factor authentication is required. Please enter your 2FA code.';

      show2fa.value = true;
      return;
    }
    
    await userStore.fetchUser();

    const userId: number | null = userStore.getId;
    if (userId) {
      userStore.setWebSocketService(userId);
    }
    console.log('Sign in successful', userStore.getNickname);

    successMessage.value = t('successMessageSignIn');
    router.push('/');
  } catch (err: any) {
    console.log(err);
    errorMessage.value = err.error || t('errorUnexpected');
  }
};

const handleSignup = () => {
  router.push('/signup');
};

const resetMessages = () => {
  errorMessage.value = '';
  successMessage.value = '';
};

const confirm2FA = async () => {
  resetMessages();
  try {
    const response = await fetch(`${API_BASE_URL}/auth/check2FA`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ 
        code: confirmationCode.value,
        nickname: nickname.value
       }),
    });

    if (!response.ok) {
      throw new Error('Failed to confirm 2FA code');
    }

    const result = await response.json();
    successMessage.value = `2FA Confirmation successful: ${result.message}`;
    await userStore.fetchUser();

    const userId: number | null = userStore.getId;
    if (userId) {
      userStore.setWebSocketService(userId);
    }
    console.log('Sign in successful', userStore.getNickname);

    successMessage.value = 'Sign in successful!';
    router.push('/');
    
  } catch (error) {
    console.error('Error confirming 2FA:', error);
    errorMessage.value = "Failed to confirm 2FA. Please try again.";
  }
};
</script>

<style scoped>
.confirmTwoFA {
  display: flex;
  flex-direction: column;
  color: white;
}

.container-input-confirm-2fa {
  display: flex;
  justify-content: space-between;
}
.container-input-confirm-2fa input[type='text'] {
  width: 100%;
}
.container-input-confirm-2fa input[type='submit'] {
  width: 25%;
  padding: 10px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  color: white;
  font-size: 14px;
  transition: background-color 0.3s;
  background: linear-gradient(
    to right,
    var(--secondary-dark-color),
    color-mix(in srgb, var(--secondary-dark-color) 75%, white)
  );
}

.error-message,
.success-message {
  position: static;
  margin: 16px 0;
  padding: 8px;
  border-radius: 4px;
  font-size: 14px;
  text-align: center;
  word-wrap: break-word;
}


.error-message {
  background-color: #ffebee;
  color: #d32f2f;
  border: 1px solid #ef9a9a;
}

.success-message {
  background-color: #e8f5e9;
  color: #388e3c;
  border: 1px solid #a5d6a7;
}
</style>
