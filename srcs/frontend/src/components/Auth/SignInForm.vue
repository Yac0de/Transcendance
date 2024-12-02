<template>
  <AuthForm 
    formTitle="formTitleSignIn"
    submitButtonLabel="submitButtonLabelSignIn"
    :fields="fields"
    :successMessage="successMessage"
    :errorMessage="errorMessage"
    @submit="handleSubmit"
    @secondaryAction="handleSignup"
    secondaryButtonLabel="secondaryButtonLabelSignUp" 
  />
</template>


<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import api from '../../services/api';
import AuthForm from './AuthForm.vue';
import { useUserStore } from '../../stores/user';
import { Field } from '../../types/models';
import { useI18n } from 'vue-i18n';

const nickname = ref('');
const password = ref('');
const errorMessage = ref('');
const successMessage = ref('');
const router = useRouter();
const userStore = useUserStore();

const { t } = useI18n();

// Form fields
const fields: Field[] = [
  { label: t('nicknameField'), model: nickname, type: 'text', required: true, maxlength: 16 },
  { label: t('passwordField'), model: password, type: 'password', required: true, maxlength: 50 },
];

const handleSubmit = async () => {
  // Field validation
  if (nickname.value.length < 3) {
    errorMessage.value = t('errorMessageNicknameTooShort');
    return;
  }

  if (password.value.length < 6) {
    errorMessage.value = t('errorMessagePasswordTooShort');
    return;
  }

  // API call for authentication
  try {
    errorMessage.value = '';
    successMessage.value = '';
    await api.auth.signin({ nickname: nickname.value, password: password.value });
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
</script>

<style scoped></style>
