package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	City     string `json:"city" bson:"city"`
}

type Book struct {
	Author string `json:"author" bson:"author"`
	Title  string `json:"title" bson:"title"`
	Email  string `json:"email" bson:"email"`
	City   string `json:"city" bson:"city"`
	Name   string `json:"name" bson:"name"`
}

var client *mongo.Client

func homePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func accountPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/account.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func myListPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/mylist.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func allBooksPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/allbooks.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		collection := client.Database("yourdatabase").Collection("users")
		_, err = collection.InsertOne(context.Background(), user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/account/", http.StatusSeeOther)
	}
}

func getUserInfo(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	log.Println("Received email:", email)
	var user User

	collection := client.Database("yourdatabase").Collection("users")
	cursor, err := collection.Find(context.Background(), map[string]interface{}{"email": email})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	if cursor.Next(context.Background()) {
		if err := cursor.Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func addBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Println("Received a POST request to add a book")

		var book Book
		err := json.NewDecoder(r.Body).Decode(&book)
		if err != nil {
			fmt.Println("Error decoding request body:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Printf("Decoded book: %+v\n", book)

		collection := client.Database("yourdatabase").Collection("books")
		_, err = collection.InsertOne(context.Background(), book)
		if err != nil {
			fmt.Println("Error inserting book into database:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		filter := bson.M{"email": book.Email}

		cursor, err := collection.Find(context.Background(), filter)
		if err != nil {
			fmt.Println("Error querying database:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer cursor.Close(context.Background())

		var books []Book
		err = cursor.All(context.Background(), &books)
		if err != nil {
			fmt.Println("Error decoding query result:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(books)
		if err != nil {
			fmt.Println("Error encoding response JSON:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("Books successfully added and response sent")
	} else {
		fmt.Println("Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	var user User
	userCollection := client.Database("yourdatabase").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	city := user.City

	bookCollection := client.Database("yourdatabase").Collection("books")
	cursor, err := bookCollection.Find(ctx, bson.M{"city": city})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var books []Book
	if err = cursor.All(ctx, &books); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func searchBooks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	fmt.Println("Received search query:", query)

	filter := map[string]interface{}{
		"$or": []map[string]interface{}{
			{"title": primitive.Regex{Pattern: query, Options: "i"}},
			{"author": primitive.Regex{Pattern: query, Options: "i"}},
		},
	}

	collection := client.Database("yourdatabase").Collection("books")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error querying database:", err)
		return
	}
	defer cursor.Close(ctx)

	fmt.Println("Successfully queried database")

	var books []Book
	err = cursor.All(ctx, &books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error decoding query result:", err)
		return
	}

	fmt.Println("Successfully decoded query result:", books)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error encoding JSON response:", err)
		return
	}

	fmt.Println("Successfully encoded and sent response")
}

func getOwnerInfo(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	filter := bson.M{"Email": email}

	collection := client.Database("yourdatabase").Collection("users")
	var user User
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("Connected to MongoDB")

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	imgFs := http.FileServer(http.Dir("static/images"))
	http.Handle("/images/", http.StripPrefix("/images/", imgFs))

	http.HandleFunc("/", homePage)
	http.HandleFunc("/register_user", registerUser)
	http.HandleFunc("/account/", accountPage)
	http.HandleFunc("/mylist/", myListPage)
	http.HandleFunc("/allbooks/", allBooksPage)
	http.HandleFunc("/api/userinfo", getUserInfo)
	http.HandleFunc("/api/add_book", addBook)
	http.HandleFunc("/books", getBooks)
	http.HandleFunc("/search", searchBooks)
	http.HandleFunc("/owner", getOwnerInfo)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
