import { defineStore } from 'pinia';

export const useThemeStore = defineStore('theme', {
  state: () => ({
    themes: ['texturized-and-dynamic', 'metallic-chill', 'cool-and-collected', 'erthy-and-serene', 'mechanical-and-floaty', 'striking-and-simple', 'sleek-and-futuristic', 'eye-catching-and-sleek', 'impactful-and-striking-colors', 'vibrant-and-calming'],
    currentThemeIndex: 0,
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
        this.currentThemeIndex = 0;
      }
    },
    applyTheme(theme: string) {
      const rootElement = document.documentElement;
      const themes = ['texturized-and-dynamic', 'metallic-chill', 'cool-and-collected', 'erthy-and-serene', 'mechanical-and-floaty', 'striking-and-simple', 'sleek-and-futuristic', 'eye-catching-and-sleek', 'impactful-and-striking-colors', 'vibrant-and-calming'];

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
