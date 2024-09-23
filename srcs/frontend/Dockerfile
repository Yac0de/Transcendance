# Utilisation de l'image node pour le développement
FROM node:20-alpine

# Définir le répertoire de travail
WORKDIR /app

# Copier les fichiers package.json et package-lock.json et installer les dépendances
COPY package*.json ./
RUN npm install

# Mettre à jour npm à la dernière version
RUN npm install -g npm@latest

# Copier tout le code source dans le conteneur
COPY . .

# Exposer le port utilisé par Vite (3000 par défaut)
EXPOSE 3000

# Démarrer le serveur de développement de Vite
CMD ["npm", "run", "dev"]
