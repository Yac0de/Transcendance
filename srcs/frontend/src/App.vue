<template>
  <div id="app">
    <nav class="sticky-nav">
      <div class="nav-content"> 
        <button @click="navigate('home')">Home</button>
        <button @click="navigate('login')">Login</button>
        <button @click="navigate('signup')">Sign Up</button>
        <button @click="navigate('pong')">Play Pong</button>
      </div>          
    </nav>
    <div class="content">
      <component :is="currentComponent" @navigate="navigate"></component>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, shallowRef } from 'vue'
import HomeView from './components/HomeView.vue'
import LoginForm from './components/LoginForm.vue'
import SignUpForm from './components/SignUpForm.vue'
import PongGame from './components/PongGame.vue'

type ComponentName = 'home' | 'login' | 'signup' | 'pong'

export default defineComponent({
  name: 'App',
  components: {
    HomeView,
    LoginForm,
    SignUpForm,
    PongGame
  },
  setup() {
    const currentComponent = shallowRef<typeof HomeView>(HomeView)

    const navigate = (page: ComponentName) => {
      switch(page) {
        case 'home':
          currentComponent.value = HomeView
          break
        case 'login':
          currentComponent.value = LoginForm
          break
        case 'signup':
          currentComponent.value = SignUpForm
          break
        case 'pong':
          currentComponent.value = PongGame
          break
      }
    }

    return {
      currentComponent,
      navigate
    }
  }
})
</script>

<style>
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
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  z-index: 1000;
}

.nav-content {
  display: flex;
  justify-content: flex-start;
  padding: 10px 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.nav-content button {
  margin-right: 10px;
  padding: 5px 10px;
  font-size: 16px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.nav-content button:hover {
  background-color: #0056b3;
}

.content {
  margin-top: 60px; /* Adjust this value based on the height of your nav bar */
  padding: 20px;
}
</style>
