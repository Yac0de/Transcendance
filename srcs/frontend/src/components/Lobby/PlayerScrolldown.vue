<template>
  <div class="scroll-down-container" v-click-outside="closeDropdown">
    <div class="selected-friend" @click="toggleDropdown">
      <span v-if="selectedFriend">{{ selectedFriend }}</span>
      <span v-else>Select a friend</span>
      <span class="arrow" :class="{ 'arrow-up': isOpen }">▼</span>
    </div>
    <div class="dropdown-list" v-if="isOpen">
      <div v-for="friend in friends" :key="friend" class="friend-item" @click="selectFriend(friend)"
        :class="{ 'selected': friend === selectedFriend }">
        {{ friend }}
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ClickableScrollDown',
  props: {
    friends: {
      type: Array,
      default: () => []
    }
  },
  data() {
    return {
      isOpen: false,
      selectedFriend: null
    }
  },
  directives: {
    'click-outside': {
      mounted(el, binding) {
        el.clickOutsideEvent = function (event) {
          if (!(el === event.target || el.contains(event.target))) {
            binding.value(event);
          }
        };
        document.addEventListener('click', el.clickOutsideEvent);
      },
      unmounted(el) {
        document.removeEventListener('click', el.clickOutsideEvent);
      }
    }
  },
  methods: {
    toggleDropdown() {
      this.isOpen = !this.isOpen;
    },
    closeDropdown() {
      this.isOpen = false;
    },
    selectFriend(friend) {
      this.selectedFriend = friend;
      this.isOpen = false;
      this.$emit('friend-selected', friend);
    }
  }
}
</script>

<style scoped>
.scroll-down-container {
  position: relative;
  width: 200px;
}

.selected-friend {
  padding: 10px 15px;
  border: 2px solid #2c3e50;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: white;
}

.arrow {
  transition: transform 0.2s ease;
}

.arrow-up {
  transform: rotate(180deg);
}

.dropdown-list {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background-color: white;
  border: 2px solid #2c3e50;
  border-top: none;
  border-radius: 0 0 8px 8px;
  max-height: 200px;
  overflow-y: auto;
  z-index: 1000;
}

.friend-item {
  padding: 10px 15px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.friend-item:hover {
  background-color: #f5f5f5;
}

.friend-item.selected {
  background-color: #e8e8e8;
}

/* Scrollbar styling */
.dropdown-list::-webkit-scrollbar {
  width: 8px;
}

.dropdown-list::-webkit-scrollbar-track {
  background: #f1f1f1;
}

.dropdown-list::-webkit-scrollbar-thumb {
  background: #888;
  border-radius: 4px;
}

.dropdown-list::-webkit-scrollbar-thumb:hover {
  background: #555;
}
</style>
