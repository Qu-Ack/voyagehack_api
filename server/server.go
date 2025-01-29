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
	"github.com/Qu-Ack/voyagehack_api/services/messaging"
	"github.com/Qu-Ack/voyagehack_api/services/observers"
	"github.com/Qu-Ack/voyagehack_api/services/upload"
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

const (
	userContextKey = "user"
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
	messageRepo := messaging.NewMessageRepo(db)
	messageService := messaging.NewMessageService(messageRepo)

	uploadService, err := upload.NewUploadService()

	if err != nil {
		fmt.Println(err)
		panic("couldn't initialize upload service")
	}

	// GraphQL handler
	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			UserService:      userService,
			ObserverService:  observerService,
			MailService:      mailService,
			MessagingService: messageService,
			UploadService:    uploadService,
		},
	}))

	srv.AddTransport(&transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			EnableCompression: true,
		},
		InitFunc: func(ctx context.Context, initPayload transport.InitPayload) (context.Context, *transport.InitPayload, error) {
			// Try to get token from payload
			if auth, ok := initPayload["Authorization"].(string); ok {
				token := strings.TrimPrefix(auth, "Bearer ")
				fmt.Println(token)

				// Validate token
				parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
					return []byte("tryandbruteforcethisbitch"), nil
				})

				if err != nil || !parsedToken.Valid {
					return ctx, nil, fmt.Errorf("invalid token")
				}

				claims, ok := parsedToken.Claims.(jwt.MapClaims)
				if !ok {
					return ctx, nil, fmt.Errorf("invalid claims")
				}

				// Create authenticated user
				user := graph.AuthenticatedUser{
					ID:    claims["id"].(string),
					Email: claims["email"].(string),
					Role:  claims["role"].(string),
				}

				// Store user in context and return modified payload
				return context.WithValue(ctx, graph.UserContextKey, user), &initPayload, nil
			}

			// Handle test token if present
			if testToken, ok := initPayload["TestToken"].(string); ok {
				if testToken == "tryandbruteforcethisbitch" {
					return context.WithValue(ctx, graph.UserContextKey, "yes"), &initPayload, nil
				}
			}

			return ctx, &initPayload, nil
		},
	})

	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GRAPHQL{})

	router.GET("/", func(c *gin.Context) {
		playground.Handler("GraphQL Playground", "/query").ServeHTTP(c.Writer, c.Request)
	})
	router.POST("/query", func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	})
	router.GET("/query", func(ctx *gin.Context) {
		srv.ServeHTTP(ctx.Writer, ctx.Request)
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
	fmt.Println(uri)
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

// GinContextToContextMiddleware converts Gin context to standard context
func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "GinContextKey", c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func AuthTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user graph.AuthenticatedUser

		// Check for Authorization header
		authToken := c.Request.Header.Get("Authorization")
		if authToken == "" {
			// Handle test token
			testToken := c.Request.Header.Get("TestToken")
			if testToken == "tryandbruteforcethisbitch" {
				user := "yes"
				ctx := context.WithValue(c.Request.Context(), graph.UserContextKey, user)
				c.Request = c.Request.WithContext(ctx)
				c.Next()
				return
			} else if testToken == "" {
				c.Next()
				return
			}
		}

		// Process JWT token
		tokenString := strings.TrimPrefix(authToken, "Bearer ")
		if tokenString == authToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		// Validate token
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte("tryandbruteforcethisbitch"), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
			return
		}

		// Create authenticated user
		user = graph.AuthenticatedUser{
			ID:    claims["id"].(string),
			Email: claims["email"].(string),
			Role:  claims["role"].(string),
		}

		// Store in context using the custom key
		ctx := context.WithValue(c.Request.Context(), graph.UserContextKey, user)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func GetAuthenticatedUser(ctx context.Context) (*graph.AuthenticatedUser, error) {
	user, ok := ctx.Value(userContextKey).(graph.AuthenticatedUser)
	if !ok {
		return nil, fmt.Errorf("unauthorized: user not found in context")
	}
	return &user, nil

}
