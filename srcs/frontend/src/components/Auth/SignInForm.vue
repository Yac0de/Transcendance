<template>
  <AuthForm formTitle="Sign In" submitButtonLabel="Sign In" :fields="fields" :successMessage="successMessage"
    :errorMessage="errorMessage" @submit="handleSubmit" @secondaryAction="handleSignup"
    secondaryButtonLabel="Sign Up" />
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import api from '../../services/api';
import AuthForm from './AuthForm.vue';
import { useUserStore } from '../../stores/user';
import { Field } from '../../types/models';

const nickname = ref('');
const password = ref('');
const errorMessage = ref('');
const successMessage = ref('');
const router = useRouter();
const userStore = useUserStore();

// Form fields
const fields = ref<Field[]>([
  { label: 'Nickname', model: nickname, type: 'text', required: true, maxlength: 20 },
  { label: 'Password', model: password, type: 'password', required: true, maxlength: 50 },
]);

const handleSubmit = async () => {
  // Field validation
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
    await api.auth.signin({ nickname: nickname.value, password: password.value });
    await userStore.fetchUser();

    userStore.setWebSocketService(userStore.getId);
    console.log('Sign in successful', userStore.getNickname);

    successMessage.value = 'Sign in successful!';
    router.push('/');
  } catch (err: any) {
    console.log(err);
    errorMessage.value = err.error || 'An error occurred during sign in';
  }
};

const handleSignup = () => {
  router.push('/signup');
};
</script>

<style scoped></style>
