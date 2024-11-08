<template>
  <div class="leave-button-container">
    <button class="nav-button leave-lobby-button" @click="leaveLobby">Leave Lobby</button>
  </div>
</template>

<script>
export default {
  name: 'LeaveLobbyButton',
  methods: {
    leaveLobby() {
      this.$emit('leave-lobby')
    }
  }
}
</script>

<style scoped>
.leave-button-container {
  position: absolute;
  top: 80px;
  /* Adjusted to be below navbar */
  left: 20px;
  z-index: 100;
}

.leave-lobby-button {
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  text-decoration: none;
  transition: background-color 0.3s;
  padding: 5px 10px;
  font-size: 16px;
}

.leave-lobby-button:hover {
  background-color: #0056b3;
}

@media (max-width: 600px) {
  .leave-button-container {
    top: 120px;
    /* Additional space for collapsed navbar on mobile */
  }
}
</style>
