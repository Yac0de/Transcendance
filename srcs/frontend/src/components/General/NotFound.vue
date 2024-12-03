<template>
	<div class="not-found-container">
	  <h1>{{ $t('notFoundTitle') }}</h1>
	  <p>{{ $t('notFoundMessage') }}</p>
	  <router-link to="/">{{ $t('goBackHome') }}</router-link>
	</div>
  </template>  

  <script setup lang="ts">
  </script>

  <style scoped>
  .not-found-container {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	height: 100vh;
	text-align: center;
  }

  h1 {
	font-size: 4rem;
	margin-bottom: 1rem;
	color: #ff4757;
  }

  p {
	font-size: 1.5rem;
	margin-bottom: 1.5rem;
	color: white;
  }

  a {
	font-size: 1.25rem;
	color: var(--secondary-bright-color);
	text-decoration: underline;
  }
  </style>