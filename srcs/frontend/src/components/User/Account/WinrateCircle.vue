<template>
	<div class="relative inline-flex items-center justify-center">
	  <svg class="w-32 h-32 -rotate-90">
		<circle
		  class="text-gray-200"
		  stroke-width="8"
		  stroke="currentColor"
		  fill="transparent"
		  :r="radius"
		  cx="64"
		  cy="64"
		/>
		<circle
		  class="text-blue-600"
		  stroke-width="8"
		  stroke-linecap="round"
		  stroke="currentColor"
		  fill="transparent"
		  :r="radius"
		  cx="64"
		  cy="64"
		  :style="{
			strokeDasharray: circumference,
			strokeDashoffset: strokeDashoffset
		  }"
		/>
	  </svg>
	  <span class="absolute text-2xl font-bold text-white">
		{{ Math.round(percentage) }}%
	  </span>
	</div>
  </template>
  
  <script setup lang="ts">
  import { computed } from 'vue'
  
  const props = defineProps<{
	percentage: number
  }>()
  
  const radius = 50
  const circumference = computed(() => 2 * Math.PI * radius)
  const strokeDashoffset = computed(() => 
	circumference.value - (props.percentage / 100) * circumference.value
  )
  </script>
  
  <style scoped>
  .text-gray-200 {
	color: #BC4749;
  }
  
  .text-blue-600 {
	color: #80ED99;
  }
  
  .w-32 {
	width: 8rem;
  }
  
  .h-32 {
	height: 8rem;
  }
  
  .-rotate-90 {
	transform: rotate(-90deg);
  }
  
  .relative {
	position: relative;
  }
  
  .absolute {
	position: absolute;
  }
  
  .inline-flex {
	display: inline-flex;
  }
  
  .items-center {
	align-items: center;
  }
  
  .justify-center {
	justify-content: center;
  }
  
  .text-2xl {
	font-size: 1.5rem;
  }
  
  .font-bold {
	font-weight: 700;
  }
  </style>