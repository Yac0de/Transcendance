<template>
    <div class="account-stats">
      <h3>Game Statistics</h3>
  
      <div class="dropdown">
        <select v-model="selectedCategory" class="dropdown-select">
          <option value="general">General Statistics</option>
          <option value="tournament">Tournament Mode</option>
          <option value="vs-ai">VS AI Mode</option>
          <option value="vs-friends">VS Friends Mode</option>
        </select>
      </div>
  
      <div v-if="selectedCategory === 'general'" class="stats-section">
        <div class="stats-item">
          <span class="stat-label">Total Games Played:</span>
          <span class="stat-value">{{ totalGamesPlayed }}</span>
        </div>
        <div class="stats-item">
          <span class="stat-label">Wins:</span>
          <span class="stat-value">{{ wins }}</span>
        </div>
        <div class="stats-item">
          <span class="stat-label">Losses:</span>
          <span class="stat-value">{{ losses }}</span>
        </div>
        <div class="stats-item">
          <span class="stat-label">Win Rate:</span>
          <span class="stat-value">{{ winRate }}%</span>
        </div>
      </div>

      <div v-if="selectedCategory === 'tournament'" class="stats-section">
        <div class="stats-item">
          <span class="stat-label">Best Score:</span>
          <span class="stat-value">{{ bestTournamentScore }}</span>
        </div>
        <div class="stats-item">
          <span class="stat-label">Total Games:</span>
          <span class="stat-value">{{ totalTournamentsPlayed }}</span>
        </div>
        <div class="stats-item">
          <span class="stat-label">Win Rate:</span>
          <span class="stat-value">{{ tournamentWinRate }}%</span>
        </div>
      </div>
  
      <div v-if="selectedCategory === 'vs-ai'" class="stats-section">
        <div class="stats-item">
          <span class="stat-label">Best Score:</span>
          <span class="stat-value">{{ bestAIScore }}</span>
        </div>
        <div class="stats-item">
          <span class="stat-label">Total Games:</span>
          <span class="stat-value">{{ totalGamesAgainstAI }}</span>
        </div>
        <div class="stats-item">
          <span class="stat-label">Win Rate:</span>
          <span class="stat-value">{{ winRateAgainstAI }}%</span>
        </div>
      </div>
  
      <div v-if="selectedCategory === 'vs-friends'" class="stats-section">
        <div class="stats-item">
          <span class="stat-label">Best Score:</span>
          <span class="stat-value">{{ bestFriendScore }}</span>
        </div>
        <div class="stats-item">
          <span class="stat-label">Total Games:</span>
          <span class="stat-value">{{ totalGamesAgainstFriends }}</span>
        </div>
        <div class="stats-item">
          <span class="stat-label">Win Rate:</span>
          <span class="stat-value">{{ winRateAgainstFriends }}%</span>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref } from 'vue';
  
  const totalGamesPlayed = 50;
  const wins = 30;
  const losses = 20;
  const winRate = ((wins / totalGamesPlayed) * 100).toFixed(2);
  
  const selectedCategory = ref('general');
  
  const bestTournamentScore = 120;
  const totalTournamentsPlayed = 10;
  const tournamentWinRate = ((5 / totalTournamentsPlayed) * 100).toFixed(2);
  
  const bestAIScore = 150;
  const totalGamesAgainstAI = 25;
  const winRateAgainstAI = ((18 / totalGamesAgainstAI) * 100).toFixed(2);
  
  const bestFriendScore = 200;
  const totalGamesAgainstFriends = 15;
  const winRateAgainstFriends = ((9 / totalGamesAgainstFriends) * 100).toFixed(2);
  </script>
  
  <style scoped>
  .account-stats {
    margin-top: 20px;
    padding: 10px;
    background-color: #f0f0f0;
    border-radius: 10px;
  }
  
  .dropdown {
    margin-bottom: 20px;
  }
  
  .dropdown-select {
    padding: 10px;
    font-size: 16px;
    border: 1px solid #ccc;
    border-radius: 4px;
  }
  
  .stats-section {
    margin-top: 20px;
  }
  
  .stats-item {
    display: flex;
    justify-content: space-between;
    margin-bottom: 10px;
  }
  
  .stat-label {
    font-weight: bold;
  }
  
  .stat-value {
    color: #333;
  }
  </style>
