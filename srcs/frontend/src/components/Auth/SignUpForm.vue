<template>
  <AuthForm formTitle="SIGN UP" submitButtonLabel="Sign Up" :fields="fields" :successMessage="successMessage"
    :errorMessage="errorMessage" @submit="handleSubmit" @secondaryAction="handleSignin"
    secondaryButtonLabel="Back to sign in" />
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import api from '../../services/api';
import AuthForm from './AuthForm.vue';
import { Field } from '../../types/models';

const nickname = ref('');
const password = ref('');
const confirmPassword = ref('');
const errorMessage = ref('');
const successMessage = ref('');
const router = useRouter();

// Form fields
const fields: Field[] = [
  { label: 'Nickname', model: nickname, type: 'text', required: true, maxlength: 16 },
  { label: 'Password', model: password, type: 'password', required: true, maxlength: 50 },
  { label: 'Confirm Password', model: confirmPassword, type: 'password', required: true, maxlength: 50 },
];

const handleSubmit = async () => {
  // Field validation
  if (password.value !== confirmPassword.value) {
    errorMessage.value = "Passwords don't match!";
    return;
  }

  if (nickname.value.length < 3) {
    errorMessage.value = "Nickname must be at least 3 characters long!";
    return;
  }

  if (password.value.length < 6) {
    errorMessage.value = "Password must be at least 6 characters long!";
    return;
  }

  // API call for authentication
  try {
    errorMessage.value = '';
    successMessage.value = '';
    await api.auth.signup({ nickname: nickname.value, password: password.value });
    successMessage.value = `Sign up successful! Welcome, ${nickname.value}!`;
    router.push('/signin');
  } catch (err: any) {
    errorMessage.value = err.error || 'An error occurred during sign up';
  }
};

const handleSignin = () => {
  router.push('/signin');
};
</script>


<style scoped></style>
