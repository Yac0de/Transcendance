<template>
  <AuthForm 
    formTitle="formTitleSignUp"
    submitButtonLabel="submitButtonLabelSignUp"
    :fields="fields"
    :successMessage="successMessage"
    :errorMessage="errorMessage"
    @submit="handleSubmit"
    @secondaryAction="handleSignin"
    secondaryButtonLabel="secondaryButtonLabelSignIn" 
  />
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import api from '../../services/api';
import AuthForm from './AuthForm.vue';
import { Field } from '../../types/models';
import { useI18n } from 'vue-i18n';

const nickname = ref('');
const password = ref('');
const confirmPassword = ref('');
const errorMessage = ref('');
const successMessage = ref('');
const router = useRouter();

const { t } = useI18n();

const fields: Field[] = [
  { label: 'nicknameField', model: nickname, type: 'text', required: true, maxlength: 16 },
  { label: 'passwordField', model: password, type: 'password', required: true, maxlength: 50 },
  { label: 'confirmPasswordField', model: confirmPassword, type: 'password', required: true, maxlength: 50 },
];

const handleSubmit = async () => {
  if (password.value !== confirmPassword.value) {
    errorMessage.value = t('errorMessagePasswordsDontMatch');
    return;
  }

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
    await api.auth.signup({ nickname: nickname.value, password: password.value });
    successMessage.value = t('successMessageSignUp', { nickname: nickname.value });
    router.push('/signin');
  } catch (err: any) {
    errorMessage.value = err.error || t('errorUnexpected');
  }
};

const handleSignin = () => {
  router.push('/signin');
};
</script>


<style scoped></style>
