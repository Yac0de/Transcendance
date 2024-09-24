# Étape 1 : Utiliser une image Go pour compiler le projet
FROM golang:1.22.2-alpine AS build-stage

# Définir le répertoire de travail dans le conteneur
WORKDIR /app

# Copier seulement les fichiers nécessaires pour les dépendances
COPY go.mod go.sum ./
RUN go mod download

# Copier tout le code source dans le conteneur
COPY . .

# Compiler l'application Go
RUN go build -o backend .

# Étape 2 : Utiliser une image plus légère pour l'exécution
FROM alpine:latest

# Installer les dépendances requises
RUN apk --no-cache add ca-certificates

# Définir le répertoire de travail pour l'exécution
WORKDIR /app

# Copier le binaire compilé depuis l'étape de build
COPY --from=build-stage /app/backend .

# Exposer le port 4000 pour l'application
EXPOSE 4000

# Démarrer l'application
CMD ["./backend"]
