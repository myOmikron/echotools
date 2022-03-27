package auth

import (
	"errors"
	"fmt"
	"github.com/myOmikron/echotools/middleware"
	"github.com/myOmikron/echotools/utilitymodels"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrAuthenticationFailed = errors.New("authentication failed")
	ErrUsernameNotFound     = errors.New("username not found")
	ErrUserNotFound         = errors.New("user not found")
	ErrHashError            = errors.New("hashing has failed")
)

//Authenticate Try to authenticate with the given credentials
func Authenticate(db *gorm.DB, username string, password string) (*utilitymodels.User, error) {
	var u utilitymodels.User
	var count int64

	db.Find(&u, "username = ?", username).Count(&count)
	if count == 0 {
		// Comparing static hash in order to deny username enumeration by looking at the time a request took
		bcrypt.CompareHashAndPassword(
			[]byte("$2b$12$KisigGoquLISbypB3kHB1eUOXZUWm7AwOZcwIIH9V9YejhxvIvlo6"),
			[]byte("Deny username enumeration"),
		)
		return nil, ErrUsernameNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		fmt.Println(err.Error())
		return nil, ErrAuthenticationFailed
	}

	return &u, nil
}

func SetNewPassword(db *gorm.DB, userID uint, newPassword string) error {
	var u utilitymodels.User
	var count int64

	if err := db.Find(&u, userID).Count(&count).Error; err != nil {
		return middleware.ErrDatabaseError
	}

	if count != 1 {
		return ErrUserNotFound
	}

	if hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12); err != nil {
		return ErrHashError
	} else {
		u.Password = string(hash)
	}

	if err := middleware.InvalidateSessions(db, userID); err != nil {
		return err
	}

	if err := db.Save(&u).Error; err != nil {
		fmt.Println("unable to update user")
		return middleware.ErrDatabaseError
	}
	return nil
}
