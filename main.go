package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Book struct{
	ID string `json:id`
	Title string `json:title`
	Description string `json:description`
	Author string `json:author`
}
type Data struct{
	Result string `json:result`
	Book Book
}
var db *sql.DB
var err error
func main() {
	db, err = sql.Open("mysql", "root:@/publicaciones")
if err != nil {
panic(err.Error())
}
	defer db.Close()
	router := gin.Default()// inicializamos el enrutador
	//Creacion de endpoints
	router.GET("/getAllBooks", getAllBooks)
	router.GET("/getBookById",getBookById)
	router.POST("/createBook", createBook)
	router.PUT("/updateBook", updateBook)
	router.DELETE("/deleteBook", deleteBook)


	router.Run(":8081")

}
func getAllBooks(c *gin.Context) {
	var books []Book
	result, err := db.Query("SELECT * from books")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var book Book
		err := result.Scan(&book.ID, &book.Title,&book.Description,&book.Author)
		if err != nil {
			panic(err.Error())
		}
		books = append(books, book)
	}
	if arrByte, err := json.Marshal(&books); err != nil {
		c.Error(err)
		c.Status(409)
	} else {
		arrByte= append(arrByte, '\n')
		c.Data(200, "application/json", arrByte)
	}
}
func getBookById(c *gin.Context) {
	params,found := c.GetQuery("id")
	if !found{
		println("no se encontro el parametro")
	}
	result, err := db.Query("SELECT * FROM books WHERE idBook = ?", params)
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var book Book
	for result.Next() {
		err := result.Scan(&book.ID, &book.Title,&book.Description,&book.Author)
		if err != nil {
			println(err.Error())
		}
	}
	if arrByte, err := json.Marshal(&book); err != nil {
		c.Error(err)
		c.Status(409)
	} else {
		arrByte = append(arrByte, '\n')
		c.Data(200, "application/json", arrByte)
	}
}
func createBook(c *gin.Context) {
	var book Book
	stmt, err := db.Prepare("INSERT INTO books(titleBook,descriptionBook,authorBook) VALUES(?,?,?)")
	if err != nil {
		println(err.Error())
	}
	err = c.ShouldBind(&book)
	if err != nil {
		println(err.Error())
	}

	_, err = stmt.Exec(book.Title,book.Description,book.Author)
	if err != nil {
		println(err.Error())
	}
	var result Data
	result.Result="Book added successfully"
	result.Book=book
	if arrByte, err := json.Marshal(&result); err != nil {
		c.Error(err)
		c.Status(409)
	} else {
		arrByte = append(arrByte, '\n')
		c.Data(200, "application/json", arrByte)
	}
}
func updateBook(c *gin.Context) {
	var book Book
	err = c.ShouldBind(&book)
	if err != nil {
		println("error en shouldBind", err.Error())
	}
	stmt, err := db.Prepare("UPDATE books SET titleBook = ?,descriptionBook=?,authorBook=? WHERE idBook = ?")
	if err != nil {
		println("Error al preparar query", err.Error())
	}
	_, e := stmt.Exec(book.Title, book.Description, book.Author, book.ID)
	if e != nil {
		println( e.Error())
	}
	var result Data
	result.Result = "Book updated successfully"
	result.Book = book
	if arrByte, err := json.Marshal(&result); err != nil {
		c.Error(err)
		c.Status(409)
	} else {
		arrByte = append(arrByte, '\n')
		c.Data(200, "application/json", arrByte)
	}
}
	func deleteBook(c *gin.Context) {
		param,found := c.GetQuery("id")
		if !found{
			println("no se encontro el parametro")
		}
		stmt, err := db.Prepare("DELETE FROM books WHERE idBook = ?")
		if err != nil {
			println(err.Error())
		}
		_, err = stmt.Exec(param)
		if err != nil {
			println(err.Error())
		}
		var result Data
		result.Result="Book deleted successfully"
		if arrByte, err := json.Marshal(&result); err != nil {
			c.Error(err)
			c.Status(409)
		} else {
			arrByte = append(arrByte, '\n')
			c.Data(200, "application/json", arrByte)
		}
	}


