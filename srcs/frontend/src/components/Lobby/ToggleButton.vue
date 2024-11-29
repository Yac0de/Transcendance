<template>
    <button
      class="toggle-button"
      :class="{ 'active': isGameMode }"
      @click="toggle"
      :disabled="disabled"
    >
      {{ isGameMode ? activeLabel : inactiveLabel }}
    </button>
  </template>
  
  <script setup lang="ts">
  interface Props {
    activeLabel: string;
    inactiveLabel: string;
    isGameMode: boolean;
    disabled?: boolean;
  }
  
  const props = withDefaults(defineProps<Props>(), {
    disabled: false,
  });
  
  const emit = defineEmits<{
    (event: 'toggle', value: boolean): void;
  }>();
  
  const toggle = () => {
    emit('toggle', !props.isGameMode);
  };
  </script>
  
  <style scoped>
  .toggle-button {
    width: 120px;
    padding: 10px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-weight: bold;
    transition: all 0.3s ease;
  }
  
  .toggle-button {
    background: linear-gradient(to right, var(--secondary-dark-color), color-mix(in srgb, var(--secondary-dark-color) 75%, white));
    color: white;
  }
  
  .toggle-button.active {
    background: linear-gradient(to right, var(--secondary-bright-color), color-mix(in srgb, var(--secondary-bright-color) 75%, white));
    color: white;
  }
  </style>
  
