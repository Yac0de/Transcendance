import { defineStore } from 'pinia';

export const useGameSettingsStore = defineStore("gameSettings", {
    state: () => ({
      gameMode: false,
    }),
    actions: {
      setGameMode(mode: boolean) {
        this.gameMode = mode;
      },
    },
  });
  