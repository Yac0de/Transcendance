<template>
  <div class="auth-container">
    <div class="auth-form">
      <h2>{{ formTitle }}</h2>
      <form @submit.prevent="onSubmit">
        <div v-for="(field, index) in fields" :key="index">
          <label :for="field.label">{{ field.label }}</label>
          <input :type="field.type" :id="field.label" v-model="field.model.value" :required="field.required"
            :maxlength="field.maxlength" />
        </div>
        <div v-if="successMessage" class="alert alert-success">{{ successMessage }}</div>
        <div v-if="errorMessage" class="alert alert-error">{{ errorMessage }}</div>
        <button type="submit">{{ submitButtonLabel }}</button>
      </form>
      <button class="secondary-button" @click="onSecondaryAction">{{ secondaryButtonLabel }}</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from 'vue';
import { Field } from '../../types/models';

defineProps<{
  formTitle: string,
  fields: Field[],
  submitButtonLabel: string,
  successMessage: string,
  errorMessage: string,
  secondaryButtonLabel: string,
}>();

const emit = defineEmits(['submit', 'secondaryAction']);

const onSubmit = () => emit('submit');
const onSecondaryAction = () => emit('secondaryAction');
</script>

<style scoped>
.auth-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}

.auth-form {
  width: 300px;
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 10px;
  background-color: #f9f9f9;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
}

.auth-form h2 {
  text-align: center;
  color: #333;
}

.auth-form div {
  margin-bottom: 15px;
}

label {
  display: block;
  margin-bottom: 5px;
  color: #666;
}

input {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  box-sizing: border-box;
}

button {
  width: 100%;
  padding: 10px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
}

button:hover {
  background-color: #45a049;
}

.secondary-button {
  margin-top: 10px;
  background-color: #3498db;
}

.secondary-button:hover {
  background-color: #2980b9;
}

.alert {
  padding: 10px;
  margin-bottom: 10px;
  margin-top: 5px;
  border-radius: 5px;
}

.alert-success {
  background-color: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.alert-error {
  background-color: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}
</style>
