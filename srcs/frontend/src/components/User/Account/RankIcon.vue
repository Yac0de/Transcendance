<!-- RankBadge.vue -->
<template>
	<div class="rank-badge">
	  <img 
		:src="getRankImage" 
		:alt="getRankName"
		class="rank-image"
	  />
	</div>
  </template>
  
  <script setup lang="ts">
  import { computed } from 'vue';
  
  // Importez vos images
  import rank1Image from '../../../assets/ranks/rank1.png';
  import rank2Image from '../../../assets/ranks/rank2.png';
  import rank3Image from '../../../assets/ranks/rank3.png';
  import rank4Image from '../../../assets/ranks/rank4.png';
  
  const props = defineProps<{
	winrate: number
  }>();
  
  const getRankImage = computed(() => {
	if (props.winrate <= 20) return rank4Image;
	if (props.winrate <= 50) return rank3Image;
	if (props.winrate <= 75) return rank2Image;
	return rank1Image;
  });
  
  const getRankName = computed(() => {
	if (props.winrate < 20) return 'Rank IV';
	if (props.winrate < 50) return 'Rank III';
	if (props.winrate < 75) return 'Rank II';
	return 'Rank I';
  });
  </script>
  
  <style scoped>
  .rank-badge {
	display: inline-flex;
	align-items: center;
	margin-left: 10px;
  }
  
  .rank-image {
	width: 8rem;
	height: 8rem;
	object-fit: contain;
  }
  </style>