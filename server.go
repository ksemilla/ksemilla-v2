package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/ksemilla/ksemilla-v2/graph/generated"
	"github.com/ksemilla/ksemilla-v2/graph/model"
	graph "github.com/ksemilla/ksemilla-v2/graph/resolvers"
	"github.com/ksemilla/ksemilla-v2/middleware"
	"github.com/ksemilla/ksemilla-v2/utils"

	"github.com/joho/godotenv"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	godotenv.Load(filepath.Join(".", ".env"))

	r := chi.NewRouter()
	// r.Use(middleware.Logger)
	r.Use(middleware.GetCorsHandler())
	r.Use(middleware.UserContextBody)

	c := generated.Config{Resolvers: &graph.Resolver{}}
	c.Directives.HasRole = func(ctx context.Context, obj interface{}, next graphql.Resolver, role []model.Role) (interface{}, error) {
		ctxUserVal := ctx.Value(middleware.GetUserCtx())
		if ctxUserVal == nil {
			return nil, fmt.Errorf("access denied")
		}
		user := ctx.Value(middleware.GetUserCtx()).(model.User)
		var roles []string
		for _, val := range role {
			roles = append(roles, string(val))
		}
		if !utils.Contains(roles, user.Role) {
			return nil, fmt.Errorf("access denied")
		}
		return next(ctx)
	}

	// srv := handler.NewDefaultServer(generated.NewExecutableSchema(c))

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(c))

	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	// log.Fatal(http.ListenAndServe(":"+port, nil))
	log.Fatal(http.ListenAndServe(":"+port, r))
}
