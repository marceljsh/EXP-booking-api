package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	// router.GET("/books/:id", bookById)
	// router.GET("/books/:slug", bookBySlug)
	router.POST("/books", createBook)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.GET("/welcome", func(ctx *gin.Context) {
		firstname := ctx.DefaultQuery("firstname", "Guest")
		lastname := ctx.Query("lastname") // shortcut for ctx.Request.URL.Query().Get("lastname")

		ctx.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})
	router.Run("localhost:8080")
}

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
	Slug     string `json:"slug"`
}

func (b book) String() string {
	return fmt.Sprintf("%v|%v|%v|%v|%v", b.ID, b.Title, b.Author, b.Quantity, b.Slug)
}

var books = []book{
	{ID: "1", Title: "Mein Kampft", Author: "Adolf Hitler", Quantity: 69, Slug: "mein-kampft"},
	{ID: "2", Title: "Communist Manifesto", Author: "Friedrich Engels", Quantity: 420, Slug: "communist-manifesto"},
	{ID: "3", Title: "Das Kapital", Author: "Karl Marx", Quantity: 31, Slug: "das-kapital"},
}

func getBooks(ctx *gin.Context) {
	id, ok_id := ctx.GetQuery("id")
	slug, ok_slug := ctx.GetQuery("slug")

	fmt.Printf("==> id:%v | slug:%v\n", id, slug)

	if ok_id && ok_slug {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "ID and Slug cannot be used as Index simultaneously."})
	} else if !ok_id && !ok_slug {
		ctx.IndentedJSON(http.StatusOK, books)
	} else if ok_id {
		book, err := getBookById(id)
		if err != nil {
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": err})
			return
		}

		ctx.IndentedJSON(http.StatusOK, book)
	} else if ok_slug {
		book, err := getBookBySlug(slug)
		if err != nil {
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": err})
			return
		}

		ctx.IndentedJSON(http.StatusOK, book)
	}
}

/* ======================== using wild card, conflicted ======================== */
// func bookById(c *gin.Context) {
// 	id := ctx.Param("id")
// 	book, err := getBookById(id)

// 	if err != nil {
// 		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book with such id not found."})
// 		return
// 	}

// 	ctx.IndentedJSON(http.StatusOK, book)
// }

// func bookBySlug(c *gin.Context) {
// 	slug, ok := ctx.GetQuery("slug")
// 	if !ok {
// 		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing slug query parameter."})
// 		return
// 	}

// 	book, err := getBookBySlug(slug)

// 	if err != nil {
// 		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book with such slug not found."})
// 		return
// 	}

// 	ctx.IndentedJSON(http.StatusOK, book)
// }
/* ============================================================================= */

func checkoutBook(ctx *gin.Context) {
	id, ok := ctx.GetQuery("id")

	if !ok {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	book, err := getBookById(id)
	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	if book.Quantity <= 0 {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available"})
		return
	}

	book.Quantity--
	ctx.IndentedJSON(http.StatusOK, book)
}

func returnBook(ctx *gin.Context) {
	id, ok := ctx.GetQuery("id")

	if !ok {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	book, err := getBookById(id)
	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	book.Quantity++
	ctx.IndentedJSON(http.StatusOK, book)
}

func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}

	return nil, errors.New("book with such ID not found")
}

func getBookBySlug(slug string) (*book, error) {
	for i, b := range books {
		if b.Slug == slug {
			return &books[i], nil
		}
	}

	return nil, errors.New("book with such slug not found")
}

func createBook(ctx *gin.Context) {
	var newBook book

	if err := ctx.BindJSON(&newBook); err != nil {
		fmt.Println(err)
		return
	}

	books = append(books, newBook)
	ctx.IndentedJSON(http.StatusCreated, newBook)
}
