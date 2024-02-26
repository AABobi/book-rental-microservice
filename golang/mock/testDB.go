package mock

import (
	"book-rental/data"
	"book-rental/db"
	"os"
	"path/filepath"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func CreateDBforTest(m *testing.M) {
	var err error

	rootDir := filepath.Join(".", "..")
	if err != nil {
		panic("failed to get current working directory")
	}

	dbName := filepath.Join(rootDir, "test.db")
	os.Remove(dbName)

	db.DB, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	user, books := testDataForDb()
	db.DB.AutoMigrate(&data.User{}, &data.Book{})
	fillDb(db.DB, user, books)

	os.Exit(m.Run())
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

	user := data.User{
		Email:    "test",
		Password: "test",
		Books:    []data.Book{},
	}

	return user, books
}
