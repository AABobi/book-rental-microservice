package db

import (
	"book-rental/data"
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitGDB() {
	dbName := "api2.db"

	_, existFileErr := os.Stat(dbName)
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	if !os.IsNotExist(existFileErr) {
		fmt.Printf("Database file %s exists\n", dbName)
	} else {
		user, books := testDataForDb()
		db.AutoMigrate(&data.User{}, &data.Book{})
		fillDb(db, user, books)
	}
	DB = db
}

func fillDb(db *gorm.DB, user data.User, books []data.Book) {
	db.Create(&user)
	for _, book := range books {
		result := db.Create(&book)
		if result.Error != nil {
			panic(result.Error)
		}
	}
}

func testDataForDb() (data.User, []data.Book) {

	var books []data.Book

	bookNames := []string{
		"The Alchemist", "To Kill a Mockingbird", "1984", "Brave New World", "The Catcher in the Rye",
		"The Great Gatsby", "Moby-Dick", "Pride and Prejudice", "The Hobbit", "Harry Potter and the Sorcerer's Stone",
		"The Lord of the Rings", "The Shining", "One Hundred Years of Solitude", "The Odyssey", "The Hunger Games",
		"The Girl with the Dragon Tattoo", "Sapiens: A Brief History of Humankind", "The Road", "A Game of Thrones",
		"The Da Vinci Code", "The Fault in Our Stars", "The Hitchhiker's Guide to the Galaxy", "The Martian",
		"The Silence of the Lambs", "The Girl on the Train", "Jurassic Park", "The Chronicles of Narnia",
		"The Kite Runner", "The Grapes of Wrath", "The Color Purple", "Dune", "Fahrenheit 451",
		"The Godfather", "The Outsiders", "The Picture of Dorian Gray", "The Secret Garden", "Wuthering Heights",
	}

	for _, name := range bookNames {
		book := data.Book{
			BookID:          0,
			Name:            name,
			DateToReturn:    nil,
			UserID:          0,
			ShippingAddress: "",
		}
		books = append(books, book)
	}

	// Print the result for verification
	for _, book := range books {
		fmt.Printf("BookID: %d, Name: %s, DateToReturn: %v, UserID: %d, ShippingAddress: %s\n",
			book.BookID, book.Name, book.DateToReturn, book.UserID, book.ShippingAddress)
	}

	user := data.User{
		Email:    "test",
		Password: "test",
		Books:    []data.Book{}, // Initialize Books as an empty slice
	}

	return user, books
}
