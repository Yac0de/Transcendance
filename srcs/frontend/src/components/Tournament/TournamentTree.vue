<!-- TournamentTree.vue -->
<template>
  <div class="tournament-container">
    <div class="bracket">
      <!-- Final -->
      <div class="match-winner">
        <p>Tournament Winner</p>
      </div>
      <div class="match-connections">
        <!-- Semi-Final 1 -->
        <div class="match-branch">
          <div class="bracket">
            <div class="match-winner">
              <p>Semi 1 Winner</p>
            </div>
            <div class="match-connections">
              <div class="match-branch">
                <p>Player 1:  {{ usersgame1[0]?.displayname }} </p>
              </div>
              <div class="match-branch">
                <p>Player 2:  {{ usersgame1[1]?.displayname }} </p>
              </div>
            </div>
          </div>
        </div>
        <!-- Semi-Final 2 -->
        <div class="match-branch">
          <div class="bracket">
            <div class="match-winner">
              <p>Semi 2 Winner</p>
            </div>
            <div class="match-connections">
              <div class="match-branch">
                <p>Player 3:  {{ usersgame2[0]?.displayname }} </p>
              </div>
              <div class="match-branch">
                <p>Player 4:  {{ usersgame2[1]?.displayname }} </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { UserData } from '../../types/models'
import { fetchMultipleUsers } from '../../utils/fetch'

const usersgame1 = ref<(UserData | null)[]>([null, null]); 
const usersgame2 = ref<(UserData | null)[]>([null, null]); 

const props = defineProps<{
  game1array: number[]; 
  game2array: number[]; 
}>();

onMounted(async () => {
  usersgame1.value = await fetchMultipleUsers(props.game1array); 
  usersgame2.value = await fetchMultipleUsers(props.game2array); 
  console.log(usersgame1.value[0].displayname)
})

</script>

<style scoped>
.tournament-container {
  display: flex;
  height: 600px;
  justify-content: center;
}

.bracket {
  display: flex;
  flex-direction: row-reverse;
}

.bracket p {
  padding: 20px;
  margin: 0;
  background-color: #f5f5f5;
  border-radius: 4px;
  min-width: 120px;
  text-align: center;
}

.match-winner {
  position: relative;
  margin-left: 50px;
  display: flex;
  align-items: center;
}

.match-winner::after {
  position: absolute;
  content: '';
  width: 25px;
  height: 2px;
  left: 0;
  top: 50%;
  background-color: #e0e0e0;
  transform: translateX(-100%);
}

.match-connections {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.match-branch {
  display: flex;
  align-items: flex-start;
  justify-content: flex-end;
  margin-top: 10px;
  margin-bottom: 10px;
  position: relative;
}

.match-branch::before {
  content: '';
  position: absolute;
  background-color: #e0e0e0;
  right: 0;
  top: 50%;
  transform: translateX(100%);
  width: 25px;
  height: 2px;
}

.match-branch::after {
  content: '';
  position: absolute;
  background-color: #e0e0e0;
  right: -25px;
  height: calc(50% + 22px);
  width: 2px;
  top: 50%;
}

.match-branch:last-child::after {
  transform: translateY(-100%);
}

.match-branch:only-child::after {
  display: none;
}
</style>
