package data

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Book struct {
	BookID          uint       `json:"book_id" gorm:"primaryKey;autoIncrement;not null"`
	Name            string     `json:"name" gorm:"not null"`
	DateToReturn    *time.Time `json:"date_to_return"`
	UserID          uint       `json:"user_id" gorm:"foreignKey:UserID"`
	ShippingAddress string     `json:"shipping_address"`
}

// in progress
// sent
func FindBooks(db *gorm.DB, condition string) []Book {
	var books []Book
	db.Where(condition).Find(&books)
	return books
}

func FindAllBooks(db *gorm.DB, books *[]Book) {
	db.Find(&books)
}

func GetRentedBooksFromDb(db *gorm.DB, userId *uint, books *[]Book) {
	db.Where("user_id = ?", *userId).Find(&books)
}

func RentBook(db *gorm.DB, userId *uint, books []Book) error {
	tx := db.Begin()

	if tx.Error != nil {
		return tx.Error
	}
	currentTime := time.Now()
	twoWeeksLater := currentTime.Add(2 * 7 * 24 * time.Hour)

	for _, book := range books {
		shippingAddress := book.ShippingAddress
		err := tx.Where("name = ?", book.Name).First(&book).Error
		if err != nil {
			errorHandler("User", err, tx)
		}

		updateBook(tx, &twoWeeksLater, *userId, shippingAddress, book)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func RemoveUserIdFromBooks(db *gorm.DB, userId uint, books *[]Book) error {
	var date *time.Time = nil
	address := ""
	var id uint = 0

	tx := db.Begin()

	if tx.Error != nil {
		return tx.Error
	}

	for _, book := range *books {
		bookName := book.Name
		err := tx.Where("user_id = ?", userId).Where("name = ?", bookName).First(&book).Error

		if err != nil {
			errorHandler("User", err, tx)
		}

		err = updateBook(tx, date, id, address, book)

		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func updateBook(tx *gorm.DB, date *time.Time, userId uint, shippingAddress string, book Book) error {
	err := tx.Model(&book).Update("date_to_return", date).Error

	if err != nil {
		return err
	}

	err = tx.Model(&book).Update("user_id", userId).Error
	if err != nil {
		return err
	}

	err = tx.Model(&book).Update("shipping_address", shippingAddress).Error
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func errorHandler(wantedRecord string, err error, tx *gorm.DB) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Printf("%s not found.", wantedRecord)
	} else {
		fmt.Println("Error:", err)
	}

	tx.Rollback()
	return err
}
