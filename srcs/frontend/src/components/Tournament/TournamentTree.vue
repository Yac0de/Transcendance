<template>
  <div class="p-8 w-full max-w-4xl mx-auto">
    <div class="relative">
      <!-- Semi-finals -->
      <div class="flex justify-between mb-32">
        <!-- Left side -->
        <div class="flex flex-col">
          <div class="p-4 w-48 mb-4 bg-white border rounded shadow">
            <div class="font-semibold">{{ players[0].name }}</div>
          </div>
          <div class="p-4 w-48 bg-white border rounded shadow">
            <div class="font-semibold">{{ players[1].name }}</div>
          </div>
        </div>
        
        <!-- Right side -->
        <div class="flex flex-col">
          <div class="p-4 w-48 mb-4 bg-white border rounded shadow">
            <div class="font-semibold">{{ players[2].name }}</div>
          </div>
          <div class="p-4 w-48 bg-white border rounded shadow">
            <div class="font-semibold">{{ players[3].name }}</div>
          </div>
        </div>
      </div>

      <!-- Final -->
      <div class="absolute left-1/2 bottom-0 transform -translate-x-1/2">
        <div class="p-4 w-48 bg-white border rounded shadow">
          <div class="font-semibold text-center">Final Winner</div>
        </div>
      </div>

      <!-- Connecting lines -->
      <svg class="absolute top-0 left-0 w-full h-full -z-10">
        <!-- Left bracket -->
        <path
          d="M 190 32 L 240 32 L 240 96 L 290 96"
          stroke="#888"
          fill="none"
          stroke-width="2"
        />
        <path
          d="M 190 96 L 240 96"
          stroke="#888"
          fill="none"
          stroke-width="2"
        />
        
        <!-- Right bracket -->
        <path
          d="M 430 32 L 380 32 L 380 96 L 330 96"
          stroke="#888"
          fill="none"
          stroke-width="2"
        />
        <path
          d="M 430 96 L 380 96"
          stroke="#888"
          fill="none"
          stroke-width="2"
        />

        <!-- Final lines -->
        <path
          d="M 290 96 L 310 96 L 310 160"
          stroke="#888"
          fill="none"
          stroke-width="2"
        />
        <path
          d="M 330 96 L 310 96"
          stroke="#888"
          fill="none"
          stroke-width="2"
        />
      </svg>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const players = ref([
  { id: 1, name: "Player 1" },
  { id: 2, name: "Player 2" },
  { id: 3, name: "Player 3" },
  { id: 4, name: "Player 4" }
])
</script>
