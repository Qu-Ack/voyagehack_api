package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Qu-Ack/voyagehack_api/api/graph"
	"github.com/Qu-Ack/voyagehack_api/services/mail"
	"github.com/Qu-Ack/voyagehack_api/services/observers"
	"github.com/Qu-Ack/voyagehack_api/services/user"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	*http.Server
	mongoClient *mongo.Client
}

type contextKey string

const (
	userContextKey = contextKey("user")
)

func New() (*http.Server, error) {
	// MongoDB connection
	client, err := connectMongoDB()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Initialize services with MongoDB collections
	db := client.Database(os.Getenv("MONGO_DB"))

	// Gin router setup
	router := gin.Default()
	router.Use(GinContextToContextMiddleware())
	router.Use(AuthTokenMiddleware())

	userRepo := user.NewUserRepo(db)
	userService := user.NewUserService(userRepo)

	observerService := observers.NewObserversService()
	mailRepo := mail.NewMailRepo(db)
	mailService := mail.NewMailService(mailRepo)

	// GraphQL handler
	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			UserService:     userService,
			ObserverService: observerService,
			MailService:     mailService,
		},
	}))

	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})
	srv.AddTransport(transport.GRAPHQL{})

	router.GET("/", func(c *gin.Context) {
		playground.Handler("GraphQL Playground", "/query").ServeHTTP(c.Writer, c.Request)
	})
	router.POST("/query", func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	})

	return &http.Server{
		Addr:    ":8080",
		Handler: router,
	}, nil
}

func connectMongoDB() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Successfully connected to MongoDB")
	return client, nil
}

type AuthenticatedUser struct {
	ID    string
	Email string
	Role  string
}

// GinContextToContextMiddleware converts Gin context to standard context
func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "GinContextKey", c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func AuthTokenMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		AuthToken := ctx.Request.Header.Get("Authorization")

		if AuthToken == "" {
			fmt.Println("here")
			testToken := ctx.Request.Header.Get("TestToken")
			fmt.Println(testToken)

			if testToken == "tryandbruteforcethisbitch" {
				c := context.WithValue(ctx.Request.Context(), "testuser", "yes")
				ctx.Request = ctx.Request.WithContext(c)
				ctx.Next()
				return
			} else if testToken == "" {
				ctx.Next()
				return
			}
		}

		tokenString := strings.TrimPrefix(AuthToken, "Bearer ")

		if tokenString == AuthToken {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "invalid token"})
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return "tryandbruteforcethisbitch", nil
		})

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "invalid token"})
			return
		}

		if !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
			return
		}

		user := AuthenticatedUser{
			ID:    claims["id"].(string),
			Email: claims["email"].(string),
			Role:  claims["role"].(string),
		}

		c := context.WithValue(ctx.Request.Context(), userContextKey, user)
		ctx.Request = ctx.Request.WithContext(c)

		ctx.Next()

	}
}

func GetAuthenticatedUser(ctx context.Context) (*AuthenticatedUser, error) {
	user, ok := ctx.Value(userContextKey).(AuthenticatedUser)
	if !ok {
		return nil, fmt.Errorf("unauthorized: user not found in context")
	}
	return &user, nil
}
