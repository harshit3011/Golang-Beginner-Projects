package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/thedevsaddam/renderer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var rnd *renderer.Render
var db *mongo.Database

const (
	dbName         = "demo_todo"
	collectionName = "todo"
	port           = ":9000"
)

type (
	todoModel struct {
		ID        primitive.ObjectID `bson:"_id,omitempty"`
		Title     string             `bson:"title"`
		Completed bool               `bson:"completed"`
		CreatedAt time.Time          `bson:"createAt"`
	}

	todo struct {
		ID        string    `json:"id"`
		Title     string    `json:"title"`
		Completed bool      `json:"completed"`
		CreatedAt time.Time `json:"created_at"`
	}
)

func init() {
	rnd = renderer.New()
	initMongoDB()
}

func initMongoDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	db = client.Database(dbName)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := rnd.Template(w, http.StatusOK, []string{"static/home.tpl"}, nil)
	if err != nil {
		log.Println(err)
		rnd.JSON(w, http.StatusInternalServerError, renderer.M{"message": "Internal Server Error"})
		return
	}
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var t todo

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		rnd.JSON(w, http.StatusProcessing, err)
		return
	}

	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The title field is required",
		})
		return
	}

	tm := todoModel{
		ID:        primitive.NewObjectID(),
		Title:     t.Title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	_, err := db.Collection(collectionName).InsertOne(context.Background(), tm)
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to save todo",
			"error":   err,
		})
		return
	}

	rnd.JSON(w, http.StatusCreated, renderer.M{
		"message": "Todo created successfully",
		"todo_id": tm.ID.Hex(),
	})
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))

	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The id is invalid",
		})
		return
	}

	var t todo

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		rnd.JSON(w, http.StatusProcessing, err)
		return
	}

	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The title field is required",
		})
		return
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"title": t.Title, "completed": t.Completed}}

	_, err := db.Collection(collectionName).UpdateOne(context.Background(), filter, update)
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to update todo",
			"error":   err,
		})
		return
	}

	rnd.JSON(w, http.StatusOK, renderer.M{
		"message": "Todo updated successfully",
	})
}

func fetchTodos(w http.ResponseWriter, r *http.Request) {
	cursor, err := db.Collection(collectionName).Find(context.Background(), bson.M{})
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to fetch todo",
			"error":   err,
		})
		return
	}
	defer cursor.Close(context.Background())

	var todos []todoModel
	if err := cursor.All(context.Background(), &todos); err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to fetch todo",
			"error":   err,
		})
		return
	}

	todoList := make([]todo, len(todos))
	for i, t := range todos {
		todoList[i] = todo{
			ID:        t.ID.Hex(),
			Title:     t.Title,
			Completed: t.Completed,
			CreatedAt: t.CreatedAt,
		}
	}

	rnd.JSON(w, http.StatusOK, renderer.M{"data": todoList})
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))

	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The id is invalid",
		})
		return
	}

	filter := bson.M{"_id": id}
	_, err := db.Collection(collectionName).DeleteOne(context.Background(), filter)
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to delete todo",
			"error":   err,
		})
		return
	}

	rnd.JSON(w, http.StatusOK, renderer.M{
		"message": "Todo deleted successfully",
	})
}

func main() {
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", homeHandler)

	r.Mount("/todo", todoHandlers())

	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("Listening on port", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	<-stopChan
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server gracefully stopped!")

	// Close MongoDB connection
	if err := db.Client().Disconnect(context.Background()); err != nil {
		log.Fatalf("Failed to disconnect from MongoDB: %v", err)
	}
}

func todoHandlers() http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
		r.Get("/", fetchTodos)
		r.Post("/", createTodo)
		r.Put("/{id}", updateTodo)
		r.Delete("/{id}", deleteTodo)
	})
	return rg
}
