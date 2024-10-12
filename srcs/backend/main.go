package main

import (
	"api/controllers"
	"api/database"
	"api/middleware"
	"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Compteur pour le nombre total de requêtes HTTP
	httpRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total", // Nom de la métrique
			Help: "Total number of HTTP request", // Description de la métrique
		},
		[]string{"method", "endpoint", "status"}, // Labels pour le type de méthode, l'endpoint et le statut
	)

	// Histogramme pour mesurer la durée des requêtes HTTP
	httpDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds", // Nom de la métrique
			Help: "Duration of HTTP requests", // Description de la métrique
		},
		[]string{"method", "endpoint"}, // Labels pour le type de méthode et l'endpoint
	)
)

// Middleware pour enregistrer les métriques Prometheus
func prometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now() // Horodatage du début de la requête
		c.Next() // Exécution du middleware suivant
		duration := time.Since(start) // Durée de la requête
		status := c.Writer.Status() // Récupération du statut de la réponse

		// Enregistrement des métriques
		httpRequests.WithLabelValues(c.Request.Method, c.FullPath(), string(rune(status))).Inc() // Incrémentation du compteur
		httpDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(duration.Seconds()) // Observation de la durée
	}
}

func main() {
	router := gin.Default() // Initialisation du routeur Gin
	database.New() // Connexion à la base de données

	// Configuration CORS pour autoriser certaines origines
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"}, // Origines autorisées
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Méthodes autorisées
		AllowHeaders:     []string{"Content-Type", "Authorization"}, // En-têtes autorisés
		AllowCredentials: true, // Autorise les credentials
		MaxAge:           12 * time.Hour, // Durée de validité des pré-requêtes CORS
	}
	router.Use(cors.New(config)) // Application de la configuration CORS
	router.Use(middleware.Token()) // Middleware pour la gestion des tokens
	router.Use(prometheusMiddleware()) // Ajout du middleware Prometheus

	// Route pour exposer les métriques Prometheus
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.Static("/users/avatar", "./avatars") // Route pour les avatars des utilisateurs

	// Groupes de routes pour les utilisateurs et l'authentification
	users := router.Group("/users")
	auth := router.Group("/auth")
	users.Use(middleware.AuthGuard()) // Middleware pour protéger les routes des utilisateurs
	controllers.Auth(auth) // Liaisons des contrôleurs d'authentification
	controllers.Users(users) // Liaisons des contrôleurs d'utilisateurs

	router.Run(":4000") // Démarrage du serveur sur le port 4000
}