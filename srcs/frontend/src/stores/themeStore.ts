import { defineStore } from 'pinia';

export const useThemeStore = defineStore('theme', {
  state: () => ({
    themes: ['main', 'dark', 'light', 'blue', 'green'], // Liste des thèmes disponibles
    currentThemeIndex: 0, // Index du thème actif
  }),
  getters: {
    currentTheme: (state) => state.themes[state.currentThemeIndex],
  },
  actions: {
    loadTheme() {
      const savedTheme = localStorage.getItem('theme');
      if (savedTheme && this.themes.includes(savedTheme)) {
        this.currentThemeIndex = this.themes.indexOf(savedTheme);
      } else {
        this.currentThemeIndex = 0; // Défaut au premier thème
      }
    },
    applyTheme(theme: string) {
      const rootElement = document.documentElement;
      const themes = ['main', 'dark', 'light', 'blue', 'green'];

      // Supprime toutes les classes de thèmes existantes
      themes.forEach((t) => rootElement.classList.remove(`${t}-theme`));

      // Ajoute la classe du thème sélectionné
      rootElement.classList.add(`${theme}-theme`);
    },
    nextTheme() {
      this.currentThemeIndex = (this.currentThemeIndex + 1) % this.themes.length;
      const currentTheme = this.currentTheme;
      this.applyTheme(currentTheme);
      localStorage.setItem('theme', currentTheme); // Sauvegarde dans localStorage
    },
  },
});
