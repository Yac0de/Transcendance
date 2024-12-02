<template>
  <div id="app">
    <nav class="sticky-nav">
      <div class="nav-content">
        <div class="nav-left">
          <router-link to="/" class="nav-button home-button">HOME</router-link>
        </div>
        <button class="nav-button themes" @click="themeStore.nextTheme">SWITCH THEMES</button>
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
        <FriendList v-if="!isGameRoute && userStore.isSignedIn"/>
        <Chat v-if="!isGameRoute && userStore.isSignedIn"/>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, computed, watch } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useUserStore } from './stores/user';
import { useChatStore } from './stores/chatStore';
import { useThemeStore } from './stores/themeStore';
import api from './services/api';
import FriendList from './components/User/Friend/FriendMenu.vue';
import Chat from './components/User/Chat/Chat.vue';
import InvitePopUp from './components/Lobby/InvitePopUp.vue';

const userStore = useUserStore();
const chatStore = useChatStore();
const themeStore = useThemeStore();

const router = useRouter();
const route = useRoute();

const isGameRoute = computed(() => route.path.startsWith('/game'));

const checkAuth = async () => {
  await userStore.initializeStore();
};

const handleSignout = async () => {
  try {
    await api.auth.signout();
    userStore.clearUser();
    chatStore.resetUnreadMessage(0);
    router.push('/');
  } catch (error) {
    console.error('Error signing out:', error);
  }
};

onMounted(() => {
  checkAuth();
  themeStore.loadTheme(); // Charger le thème depuis localStorage
  themeStore.applyTheme(themeStore.currentTheme); // Appliquer le thème
});

watch(() => themeStore.currentTheme, (newTheme) => {
  themeStore.applyTheme(newTheme); // Appliquer le thème lorsque la valeur change
});
</script>

<style>
:root {
  font-family: "Audiowide", sans-serif;
}

.texturized-and-dynamic-theme {
  --main-color: #2f4454;
  --secondary-dark-color: #2e151b;
  --secondary-bright-color: #da7b93;
  --main-extra-color: #376e6f;
  --secondary-extra-color: #1c3334;
}

.metallic-chill-theme {
  --main-color: #3d52a0;
  --secondary-dark-color: #7091e6;
  --secondary-bright-color: #8697c4;
  --main-extra-color: #adbbda;
  --secondary-extra-color: #ede8f5;
}

.cool-and-collected-theme {
  --main-color: #003135;
  --secondary-dark-color: #024950;
  --secondary-bright-color: #964734;
  --main-extra-color: #0fa4af;
  --secondary-extra-color: #afdde5;
}

.erthy-and-serene-theme {
  --main-color: #3e362e;
  --secondary-dark-color: #865d36;
  --secondary-bright-color: #93785b;
  --main-extra-color: #ac8968;
  --secondary-extra-color: #a69080;
}

.mechanical-and-floaty-theme {
  --main-color: #141619;
  --secondary-dark-color: #2c2e3a;
  --secondary-bright-color: #050a44;
  --main-extra-color: #0a21c0;
  --secondary-extra-color: #b3b4bd;
}

.striking-and-simple-theme {
  --main-color: #0b0c10;
  --secondary-dark-color: #1f2833;
  --secondary-bright-color: #c5c6c7;
  --main-extra-color: #66fcf1;
  --secondary-extra-color: #45a29e;
}

.sleek-and-futuristic-theme {
  --main-color: #2c3531;
  --secondary-dark-color: #116466;
  --secondary-bright-color: #d9b08c;
  --main-extra-color: #ffcb9a;
  --secondary-extra-color: #d1e8e2;
}

.eye-catching-and-sleek-theme {
  --main-color: #501f3a;
  --secondary-dark-color: #cb2d6f;
  --secondary-bright-color: #cccccc;
  --main-extra-color: #14a098;
  --secondary-extra-color: #0f292f;
}

.impactful-and-striking-colors-theme {
  --main-color: #c34271;
  --secondary-dark-color: #f070a1;
  --secondary-bright-color: #16ffbd;
  --main-extra-color: #12c998;
  --secondary-extra-color: #439f76;
}

.vibrant-and-calming-theme {
  --main-color: #026670;
  --secondary-dark-color: #9fedd7;
  --secondary-bright-color: #f7f0a3;
  --main-extra-color: #fce181;
  --secondary-extra-color: #edeae5;
}

</style>

<style scoped>

#app {
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  display: flex;
  flex-direction: column;
  height: 100vh;
}

.sticky-nav {
  display: flex;
  align-items: center;
  justify-content: center;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 7vh; 
  display: flex;
  align-items: center;
  overflow: hidden;
  background: var(--main-color);
  box-shadow: 0 0px 10px 0px black;
}

.nav-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  
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
  margin-left: 10px;
}

.nav-button {
  margin-left: 10px;
  padding: 0.5vh 1vw;
  font-family: "Audiowide", sans-serif;
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
  white-space: nowrap;
}

.nav-button.themes {
  margin-left: 0;
}

.home-button {
  margin-right: 10px;
}

.nav-button:hover {
  background-color: var(--main-extra-color);
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
  background-color: var(--main-extra-color);
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

@media (max-width: 960px) {
  .sticky-nav {
  justify-content: space-between;
}
.nav-content {
  width: 100%;
}
.nav-button {
  margin-left: 0px;
}
.nav-right button,
.nav-right a {
  margin-left: 10px;
}
.nav-button.themes {
  margin-left: 10px;
}
.nav-left,
.nav-right {
  position: relative;
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
