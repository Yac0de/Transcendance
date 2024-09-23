# Étape 1 : Utiliser une image Go pour compiler le projet
FROM golang:1.22.2-alpine AS build-stage

# Définir le répertoire de travail dans le conteneur
WORKDIR /app

# Copier seulement go.mod
COPY go.mod ./
RUN go mod download

# Copier tout le code source dans le conteneur
COPY . .

# Compiler l'application Go
RUN go build -o main .

# Étape 2 : Utiliser une image plus légère pour l'exécution
FROM alpine:latest

# Définir le répertoire de travail pour l'exécution
WORKDIR /app

# Copier le binaire compilé depuis l'étape de build
COPY --from=build-stage /app/main .

# Exposer le port 5000 pour l'application
EXPOSE 5000

# Démarrer l'application
CMD ["./main"]
