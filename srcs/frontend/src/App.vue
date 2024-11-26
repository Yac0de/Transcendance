<template>
  <div id="app">
    <nav class="sticky-nav">
      <div class="nav-content">
        <div class="nav-left">
          <router-link to="/" class="nav-button home-button">Home</router-link>
        </div>
        <div class="nav-right">
          <template v-if="!userStore.isSignedIn">
            <router-link to="/signin" class="nav-button">Sign In</router-link>
            <router-link to="/signup" class="nav-button">Sign Up</router-link>
          </template>
          <template v-else>
            <router-link to="/pong" class="nav-button">Play Pong</router-link>
            <router-link :to="`/${userStore.nickname}`" class="nav-button">Account</router-link>
            <button @click="handleSignout" class="nav-button">Sign Out</button>
          </template>
        </div>
      </div>
    </nav>
    <div class="content">
      <router-view></router-view>
      <InvitePopUp />
      <FriendList v-if="userStore.isSignedIn" />
      <Chat v-if="userStore.isSignedIn" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from './stores/user';
import api from './services/api';
import FriendList from './components/User/Friend/FriendMenu.vue';
import Chat from './components/User/Chat/Chat.vue';
import InvitePopUp from './components/Lobby/InvitePopUp.vue';

const userStore = useUserStore();
const router = useRouter();

const checkAuth = async () => {
  await userStore.fetchUser();
};

const handleSignout = async () => {
  try {
    await api.auth.signout();
    userStore.clearUser();
    router.push('/');
  } catch (error) {
    console.error('Error signing out:', error);
  }
};

onMounted(checkAuth);
</script>

<style scoped>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  background: #2f4454;
}

.sticky-nav {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  background-color: #2f4454;
}

.nav-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1vh 2vw;
}

.nav-left {
  flex: 0 0 auto;
}

.nav-right {
  display: flex;
  justify-content: flex-end;
  align-items: center;
}

.nav-button {
  margin-left: 10px;
  padding: 0.5vh 1vw;
  font-size: 1.2rem;
  background: #376e6f;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  text-decoration: none;
  transition: background-color 0.3s;
}

.home-button {
}

.nav-button:hover {
  background-color: #2b5758;
}

.content {
}
</style>
