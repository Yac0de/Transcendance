<template>
  <AuthForm
    formTitle="Sign In"
    submitButtonLabel="Sign In"
    :fields="fields"
    :successMessage="successMessage"
    :errorMessage="errorMessage"
    @submit="handleSubmit"
    @secondaryAction="handleSignup"
    secondaryButtonLabel="Sign Up"
  />
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import api from '../../services/api';
import AuthForm from './AuthForm.vue';

const nickname = ref('');
const password = ref('');
const errorMessage = ref('');
const successMessage = ref('');
const user = ref<any>(null);
const router = useRouter();

// Form fields
const fields = ref([
  { label: 'Nickname', model: nickname, type: 'text', required: true },
  { label: 'Password', model: password, type: 'password', required: true },
]);

const emit = defineEmits(['signin-success']);

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
    user.value = await api.auth.signin({ nickname: nickname.value, password: password.value });

    console.log('Sign in successful', user.value);
    successMessage.value = 'Sign in successful!';
    emit('signin-success');
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

<style scoped>

</style>
