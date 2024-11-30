<template>
	<div class="not-found-container">
	  <h1>404 - Page Not Found</h1>
	  <p>The page you are looking for does not exist.</p>
	  <router-link to="/">Go back to Home</router-link>
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