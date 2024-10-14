<template>
  <div id="app">
    <nav class="sticky-nav">
      <div class="nav-content">
        <div class="nav-left">
          <router-link to="/" class="nav-button home-button">Home</router-link>
        </div>
        <div class="nav-right">
          <template v-if="!isAuthenticated">
            <router-link to="/signin" class="nav-button">Sign In</router-link>
            <router-link to="/signup" class="nav-button">Sign Up</router-link>
          </template>
          <template v-else>
            <router-link to="/pong" class="nav-button">Play Pong</router-link>
            <router-link to="/account" class="nav-button">Account</router-link>
            <button @click="handleSignout" class="nav-button">Sign Out</button>
          </template>
        </div>
      </div>
    </nav>
    <div class="content">
      <router-view @signin-success="handleSigninSuccess"></router-view>
    </div>
    <FriendList v-if="isAuthenticated" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import api from './services/api';
import FriendList from './components/FriendList.vue';

const isAuthenticated = ref(false);
const router = useRouter();

const checkAuth = async () => {
  const userData = await api.getUserData();
  isAuthenticated.value = userData !== null;
};

const handleSigninSuccess = () => {
  isAuthenticated.value = true;
  checkAuth();
}

const handleSignout = async () => {
  await api.signout();
  isAuthenticated.value = false;
  router.push('/');
};

onMounted(checkAuth);
</script>

<style scoped>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}

.sticky-nav {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  background-color: #f8f9fa;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  z-index: 1000;
}

.nav-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 20px;
  max-width: 100%;
  margin: 0 auto;
}

.nav-left {
  flex: 0 0 auto;
}

.nav-right {
  flex: 0 0 auto;
  display: flex;
  justify-content: flex-end;
  align-items: center;
}

.nav-button {
  margin-left: 10px;
  padding: 5px 10px;
  font-size: 16px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  text-decoration: none;
  transition: background-color 0.3s;
}

.home-button {
  margin-left: 0;
}

.nav-button:hover {
  background-color: #0056b3;
}

.content {
  margin-top: 60px;
  padding: 20px;
}

@media (max-width: 600px) {
  .nav-content {
    flex-direction: column;
    align-items: stretch;
  }

  .nav-right {
    margin-top: 10px;
    justify-content: flex-start;
  }

  .nav-button {
    margin-left: 0;
    margin-right: 10px;
    margin-bottom: 5px;
  }

  .home-button {
    margin-bottom: 10px;
  }
}
</style>
