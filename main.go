package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Book struct {
	ID          int    `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
	Author      string `json:"Author"`
}

var books = []Book{
	{ID: 1, Title: "To Kill a Mockingbird", Description: "A young girl's experience with racial injustice in a small Alabama town.", Author: "Harper Lee"},
	{ID: 2, Title: "The Catcher in the Rye", Description: "A disillusioned teenager's struggles with the idea of growing up.", Author: "J.D. Salinger"},
	{ID: 3, Title: "1984", Description: "A dystopian novel depicting a totalitarian future society.", Author: "George Orwell"},
	{ID: 4, Title: "1990", Description: "A novel depicting a totalitarian future society.", Author: "George Orwell"},
	{ID: 5, Title: "The Great Gatsby", Description: "A young man's quest to win the heart of his beloved in the roaring twenties.", Author: "F. Scott Fitzgerald"},
	{ID: 6, Title: "The Count of Monte Cristo", Description: "A man's quest for revenge after being betrayed by his friends.", Author: "Alexandre Dumas"},
	{ID: 7, Title: "The Picture of Dorian Gray", Description: "A young man's descent into madness and sin after selling his soul for eternal youth.", Author: "Oscar Wilde"},
	{ID: 8, Title: "Alice in Wonderland", Description: "A young girl's adventures in a fantastical world.", Author: "Lewis Carroll"},
	{ID: 9, Title: "The Adventures of Sherlock Holmes", Description: "The stories of a brilliant detective and his trusty sidekick.", Author: "Sir Arthur Conan Doyle"},
	{ID: 10, Title: "The War of the Worlds", Description: "A Martian invasion of Earth.", Author: "H.G. Wells"},
}

// fetch all books (GET)
func allBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: allBooks")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// add new book (POST)
func addBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: addBook")

	var newBook Book
	json.NewDecoder(r.Body).Decode(&newBook)

	maxID := 0
	for _, book := range books {
		if book.ID > maxID {
			maxID = book.ID
		}
	}

	newBook.ID = maxID + 1

	books = append(books, newBook)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBook)
}

// update a book by id (PUT)
func updateBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: updateBook")

	vars := mux.Vars(r)
	idParam := vars["id"]

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedBook Book
	json.NewDecoder(r.Body).Decode(&updatedBook)

	for i, book := range books {
		if book.ID == id {
			books[i].Title = updatedBook.Title
			books[i].Description = updatedBook.Description
			books[i].Author = updatedBook.Author

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(books[i])
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// delete a book by id (DELETE)
func deleteBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: deleteBook")

	vars := mux.Vars(r)
	idParam := vars["id"]

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	bookFound := false

	for i, book := range books {
		if book.ID == id {
			books = append(books[:i], books[i+1:]...)
			bookFound = true
			break
		}
	}

	response := map[string]string{}

	if bookFound {
		response["message"] = "Book deleted successfully"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		errorResponse := map[string]string{"error": "Book not found"}
		json.NewEncoder(w).Encode(errorResponse)
	}
}

// update specific book field by id(PATCH)
func updateBookField(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: updateBookField")

	vars := mux.Vars(r)
	idParam := vars["id"]

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updateFields map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updateFields)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadGateway)
		return
	}

	var book *Book
	for i := range books {
		if books[i].ID == id {
			book = &books[i]
			break
		}
	}

	if book == nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	if title, ok := updateFields["Title"].(string); ok {
		book.Title = title
	}
	if description, ok := updateFields["Description"].(string); ok {
		book.Description = description
	}
	if author, ok := updateFields["Author"].(string); ok {
		book.Author = author
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// home page
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: homePage")

	fmt.Fprintf(w, "Home Page")
}

func handleRequests() {

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage).Methods("GET")
	myRouter.HandleFunc("/books", allBooks).Methods("GET")
	myRouter.HandleFunc("/books", addBook).Methods("POST")
	myRouter.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	myRouter.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")
	myRouter.HandleFunc("/books/{id}", updateBookField).Methods("PATCH")

	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	handleRequests()
}
