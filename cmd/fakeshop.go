package main

import (
	"flag"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/golang/glog"
	"github.com/jsfan/fake-shop/internal/config"
	"github.com/jsfan/fake-shop/internal/graph"
	"github.com/jsfan/fake-shop/internal/graph/generated"
	"github.com/jsfan/fake-shop/internal/store"
	"log"
	"net/http"
	"os"
)

const defaultPort = "8888"
const stockFile = "config/stock.yaml"
const promotionsFile = "config/promotions.yaml"

func main() {
	stockFileOpt := flag.String("stock", stockFile, "Stock YAML file")
	promoFileOpt := flag.String("promotions", promotionsFile, "Promotions YAML file")

	flag.Parse()
	stock, err := config.ReadInventory(*stockFileOpt)
	if err != nil {
		glog.Fatalf("Could not read inventory: %+v", err)
	}
	promotions, err := config.ReadPromotions(*promoFileOpt)
	if err != nil {
		glog.Fatalf("Could not read promotions: %+v", err)
	}
	store.InitShop()
	if err := store.StockShop(stock); err != nil {
		glog.Fatalf("Inventory issue: %+v", err)
	}
	store.RegisterPromotions(promotions)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
