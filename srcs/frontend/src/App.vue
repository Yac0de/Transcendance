<template>
  <div id="app">
    <nav class="sticky-nav">
      <div class="nav-content">
        <div class="nav-left">
          <router-link to="/" class="nav-button home-button">HOME</router-link>
        </div>
        <div class="nav-right">
          <template v-if="!userStore.isSignedIn">
            <router-link to="/signin" class="nav-button">SIGN IN</router-link>
            <router-link to="/signup" class="nav-button">SIGN UP</router-link>
          </template>
          <template v-else>
            <router-link to="/pong" class="nav-button">Play Pong</router-link>
            <router-link :to="`/${userStore.nickname}`" class="nav-button">Account</router-link>
            <button @click="handleSignout" class="nav-button">Sign Out</button>
          </template>
        </div>
      </div>
    </nav>
    <div class="gradient_backgroud">
      <div class="content">
        <router-view></router-view>
        <InvitePopUp />
        <FriendList v-if="userStore.isSignedIn" />
        <Chat v-if="userStore.isSignedIn" />
      </div>
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

<style>
:root {
  --main-color: #376e6f;
  --secondary-dark-color: #2e151b;
  --secondary-bright-color: #da7b93;
  --main-extra-color: #2f4454;
  --secondary-extra-color: #1c3334;
}
</style>

<style scoped>

#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  display: flex;
  flex-direction: column;
  height: 100vh;
}

.sticky-nav {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 7vh; 
  display: flex;
  align-items: center;
  overflow: hidden;
  background: var(--main-color);
}

.nav-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1vh 2vw;
}

.nav-left {
  position: absolute;
  left: 2vw;
}

.nav-right {
  position: absolute;
  right: 2.5vw;
  display: flex;
  justify-content: flex-end;
  align-items: center;
}

.nav-button {
  margin-left: 10px;
  padding: 0.5vh 1vw;
  font-size: 1.2rem;
  font-weight: 600;
  text-shadow: 0.5px 0.5px 1px black;
  background: var(--main-color);
  border: 1px solid white;
  color: white;
  /* border: none; */
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

.gradient_backgroud {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(90deg, var(--secondary-dark-color), var(--secondary-bright-color));
  width: 100%;
  overflow: hidden;
  margin-top: 7vh;
}

.content {
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--main-color);
  width: 95vw;
  height: 85vh;
  box-shadow: 0 0 40px rgba(0, 0, 0, 1);
  
}

@media (max-width: 1320px), (max-height: 655px) {
  .nav-button {
    font-size: 1.1rem;
  }
}

@media (max-width: 800px), (max-height: 430px) {
  .nav-button {
    font-size: 1rem;
  }
}

@media (max-width: 768px), (max-height: 430px) {
  .nav-button {
    font-size: 0.9rem;
  }
}

@media (max-width: 520px), (max-height: 300px)  {
  .nav-button {
    font-size: 0.8rem;
  }
}

@media (max-width: 370px), (max-height: 250px) {
  .nav-button {
    font-size: 0.7rem;
  }
}

@media (max-height: 200px) {
  .nav-button {
    font-size: 0.6rem;
  }
}
</style>
