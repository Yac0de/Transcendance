<template>
  <div class="auth-container">
    <div class="auth-form">
      <h1>{{ formTitle }}</h1>
      <form @submit.prevent="onSubmit">
        <div v-for="(field, index) in fields" :key="index">
          <label :for="field.label">{{ field.label }}</label>
          <input :type="field.type" :id="field.label" v-model="field.model.value" :required="field.required"
            :maxlength="field.maxlength" />
        </div>
        <div v-if="successMessage" class="alert alert-success">{{ successMessage }}</div>
        <div v-if="errorMessage" class="alert alert-error">{{ errorMessage }}</div>
        <button class="primary-button" type="submit">{{ submitButtonLabel }}</button>
        <button class="secondary-button" @click="onSecondaryAction">{{ secondaryButtonLabel }}</button>
      </form>
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
  font-family: "Audiowide", sans-serif;
}

.auth-form {
  display: flex;
  flex-direction: column;
  justify-content: space-around;
  height: 60vh;
  min-height: 320px;
  width: 50vw;
  min-width: 210px;
  max-width: 500px;
  padding: 1vh 1.5vw;
  border-radius: 20px;
  background-color: var(--main-color);
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.5);
  overflow: hidden;
}

.auth-form h1 {
  text-align: center;
  color: white;
  font-size: 2rem;
  font-weight: 650;
  text-shadow: 1px 1px 2px black;
}

.auth-form div {
  margin-bottom: 35px;
}

label {
  display: block;
  margin-bottom: 5px;
  color: white;
  font-size: 1rem;
  font-weight: 500;
  text-shadow: 0.5px 0.5px 1px black;
}

input {
  width: 100%;
  padding: 8px;
  border: 1px solid #dddddd;
  border-radius: 4px;
  box-sizing: border-box;
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.3);
}

button {
  width: 100%;
  padding: 10px;
  transition: background 0.3s ease;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1.2rem;
}

.primary-button {
  background: linear-gradient(to right, var(--secondary-dark-color), #502d36);
}

.primary-button:hover {
  background: linear-gradient(to right, #502d36, var(--secondary-dark-color));
}

.secondary-button {
  margin-top: 10px;
  background: linear-gradient(to right, var(--secondary-bright-color), #f5a3b4);
}

.secondary-button:hover {
  background: linear-gradient(to right, #e86c89, #ffbccd);
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

@media (max-height: 900px) {
  .auth-form div {
  margin-bottom: 25px;
  }
  input{
    padding: 6px;
    box-shadow: 0 0 15px rgba(0, 0, 0, 0.3);
  }
  .auth-form h1 {
  font-size: 1.8rem;
  }
  label{
    font-size: 0.9rem;
    font-weight: 500;
    text-shadow: 0.4px 0.4px 0.9px black;
  }
  button {
  padding: 5px;
  }
  .alert {
  font-size: 0.9rem;
  padding: 8px;
  margin-bottom: 8px;
  margin-top: 4px;
  }
}

@media (max-height: 725px) {
  .auth-form div {
  margin-bottom: 20px;
  }
  .auth-form h1 {
  font-size: 1.6rem;
  }
  label{
    font-size: 0.8rem;
    font-weight: 400;
    text-shadow: 0.3px 0.3px 0.7px black;
  }
  input{
    padding: 4px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.3);
  }
  button {
  padding: 2.5px;
  }
  .alert {
  font-size: 0.8rem;
  padding: 6px;
  margin-bottom: 6px;
  margin-top: 3px;
  }
}

@media (max-height: 630px) {
  .auth-form div {
  margin-bottom: 15px;
  }
  .auth-form h1 {
  font-size: 1.4rem;
  }
  label{
    font-size: 0.7rem;
  }
  input{
    padding: 2px;
    box-shadow: 0 0 5px rgba(0, 0, 0, 0.3);
  }
  button {
  padding: 1.25px;
  }
  .alert {
  font-size: 0.7rem;
  padding: 4px;
  margin-bottom: 4px;
  margin-top: 2px;
  }
}
</style>
