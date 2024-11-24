<!-- TournamentTree.vue -->
<template>
  <div class="wrapper">
    <div class="item">
      <!-- Final -->
      <div class="item-parent">
        <p>Tournament Winner</p>
      </div>
      <div class="item-childrens">
        <!-- Semi-Final 1 -->
        <div class="item-child">
          <div class="item">
            <div class="item-parent">
              <p>Semi 1 Winner</p>
            </div>
            <div class="item-childrens">
              <div class="item-child">
                <p>Player 1</p>
              </div>
              <div class="item-child">
                <p>Player 2</p>
              </div>
            </div>
          </div>
        </div>
        <!-- Semi-Final 2 -->
        <div class="item-child">
          <div class="item">
            <div class="item-parent">
              <p>Semi 2 Winner</p>
            </div>
            <div class="item-childrens">
              <div class="item-child">
                <p>Player 3</p>
              </div>
              <div class="item-child">
                <p>Player 4</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// No script needed for now, just structure
</script>

<style scoped>
.wrapper {
  display: flex;
  height: 600px;
  justify-content: center;
}

.item {
  display: flex;
  flex-direction: row-reverse;
}

.item p {
  padding: 20px;
  margin: 0;
  background-color: #f5f5f5;
  border-radius: 4px;
  min-width: 120px;
  text-align: center;
}

.item-parent {
  position: relative;
  margin-left: 50px;
  display: flex;
  align-items: center;
}

.item-parent::after {
  position: absolute;
  content: '';
  width: 25px;
  height: 2px;
  left: 0;
  top: 50%;
  background-color: #e0e0e0;
  transform: translateX(-100%);
}

.item-childrens {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.item-child {
  display: flex;
  align-items: flex-start;
  justify-content: flex-end;
  margin-top: 10px;
  margin-bottom: 10px;
  position: relative;
}

.item-child::before {
  content: '';
  position: absolute;
  background-color: #e0e0e0;
  right: 0;
  top: 50%;
  transform: translateX(100%);
  width: 25px;
  height: 2px;
}

.item-child::after {
  content: '';
  position: absolute;
  background-color: #e0e0e0;
  right: -25px;
  height: calc(50% + 22px);
  width: 2px;
  top: 50%;
}

.item-child:last-child::after {
  transform: translateY(-100%);
}

.item-child:only-child::after {
  display: none;
}
</style>
