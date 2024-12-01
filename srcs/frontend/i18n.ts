// src/i18n.ts
import { createI18n } from 'vue-i18n';

// Définitions des messages traduits
const messages = {
  en: {
    home: "Home",
    signin: "Sign In",
    signup: "Sign Up",
    playPong: "Play Pong",
    account: "Account",
    signout: "Sign Out",
    welcomeTo: "WELCOME TO",
    transcendencePong: "TRANSCENDENCE PONG",
    playNow: "PLAY NOW",
  },
  fr: {
    home: "Accueil",
    signin: "Connexion",
    signup: "Inscription",
    playPong: "Jouer à Pong",
    account: "Compte",
    signout: "Déconnexion",
    welcomeTo: "BIENVENUE À",
    transcendencePong: "TRANSCENDANCE PONG",
    playNow: "JOUER MAINTENANT",
  },
  es: {
    home: "Inicio",
    signin: "Iniciar sesión",
    signup: "Registrarse",
    playPong: "Jugar al Pong",
    account: "Cuenta",
    signout: "Cerrar sesión",
    welcomeTo: "BIENVENIDO A",
    transcendencePong: "TRANSCENDENCIA PONG",
    playNow: "JUGAR AHORA",
  },
};

// Initialisation de Vue I18n
const i18n = createI18n({
  locale: 'en', // Langue par défaut  
  fallbackLocale: 'en', // Langue de secours
  messages, // Messages traduits
});

export default i18n;
