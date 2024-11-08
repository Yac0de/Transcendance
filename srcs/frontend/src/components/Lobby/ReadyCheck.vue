<template>
  <div class="ready-check">
    <button class="ready-button" :class="{ 'ready': isReady }" @click="toggleReady">
      {{ isReady ? 'Ready!' : 'Not Ready' }}
    </button>
  </div>
</template>

<script>
export default {
  name: 'ReadyCheck',
  data() {
    return {
      isReady: false
    }
  },
  methods: {
    toggleReady() {
      this.isReady = !this.isReady
      this.$emit('ready-changed', this.isReady)
    }
  }
}
</script>

<style scoped>
.ready-check {
  width: 120px;
  margin: 10px;
}

.ready-button {
  width: 100%;
  padding: 10px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: bold;
  transition: all 0.3s ease;
}

.ready-button:not(.ready) {
  background-color: #dc3545;
  color: white;
}

.ready-button.ready {
  background-color: #28a745;
  color: white;
}
</style>
