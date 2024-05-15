package services

import (
	"crypto/rand"
	"errors"
	"log"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"firebase.google.com/go/v4/auth"

	"github.com/jackc/pgx/v5"
	"github.com/rohan031/adgytec-api/firebase"
	"github.com/rohan031/adgytec-api/v1/custom"
	"github.com/rohan031/adgytec-api/v1/dbqueries"
	"github.com/rohan031/adgytec-api/v1/validation"
)

type UserCreationDetails struct {
	Name     string
	Email    string
	Password string
}

type User struct {
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Role      string    `json:"role" db:"role"`
	UserId    string    `json:"userId,omitempty" db:"user_id"`
	CreatedAt time.Time `json:"createdAt,omitempty" db:"created_at"`
}

func generateRandomPassword() (string, error) {
	// Define character sets for password generation
	upperChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowerChars := "abcdefghijklmnopqrstuvwxyz"
	digitChars := "0123456789"
	specialChars := "!@#$%^&*()-_=+[]{}|;:,.<>?~"

	// Concatenate all character sets
	allChars := upperChars + lowerChars + digitChars + specialChars

	var password strings.Builder
	for i := 0; i < 10; i++ {
		// Generate random index to select a character from allChars
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(allChars))))
		if err != nil {
			return "", err
		}

		// Append selected character to password
		password.WriteByte(allChars[randomIndex.Int64()])
	}

	return password.String(), nil
}

/*
if user exists in db return true
else delete user from firebase and re-add it
return false, err // internal server error
return true, nil // user exits
return false, nil // delete user from firebase and create new user
*/
func userExistsInDb(email string) (bool, error) {
	// fetching the user from db
	args := dbqueries.GetUserByEmailArgs(email)
	rows, err := db.Query(ctx, dbqueries.GetUserByEmail, args)
	if err != nil {
		log.Printf("Error fetching user from db: %v\n", err)
		return false, err
	}
	defer rows.Close()

	_, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// user doesn't exist in db
			u, err := firebase.FirebaseClient.GetUserByEmail(ctx, email)
			if err != nil {
				log.Printf("Error getting user data from firebase: %v\n", err)
				return false, err
			}

			err = firebase.FirebaseClient.DeleteUser(ctx, u.UID)
			if err != nil {
				log.Printf("Error deleting user from firebase: %v\n", err)
				return false, err
			}

			return false, nil
		}
		log.Printf("error reading rows: %v\n", err)
		return false, err
	}

	return true, nil
}

func (u *User) CreateUser() (string, error) {
	// creating random password
	password, err := generateRandomPassword()
	if err != nil {
		log.Printf("Error generating password: %v\n", err)
		return "", err
	}

	// creating user in firebase
	params := (&auth.UserToCreate{}).Email(u.Email).DisplayName(u.Name).Password(password)
	userRecord, err := firebase.FirebaseClient.CreateUser(ctx, params)
	if err != nil {
		if auth.IsEmailAlreadyExists(err) {
			// find user in db
			ispresent, err := userExistsInDb(u.Email)
			if err != nil {
				return "", err
			}

			// user already exists
			if ispresent {
				message := "The email address provided is already associated with an existing user account."
				return "", &custom.MalformedRequest{Status: http.StatusConflict, Message: message}
			}

			// create new user with given details
			return u.CreateUser()
		}

		log.Printf("Error creating user in firebase: %v\n", err)
		return "", err
	}

	// setting custom claims for newly created user
	uid := userRecord.UID
	claims := map[string]interface{}{"role": u.Role}
	err = firebase.FirebaseClient.SetCustomUserClaims(ctx, uid, claims)
	if err != nil {
		log.Printf("Error setting custom claims: %v\n", err)
		return "", err
	}

	// inserting into database user table
	args := dbqueries.CreateUserArgs(uid, u.Email, u.Name, u.Role)
	_, err = db.Exec(ctx, dbqueries.CreateUser, args)
	if err != nil {
		log.Printf("Error adding user in database: %v\n", err)
		return "", err
	}

	return password, nil
}

func (u *User) ValidateInput() bool {
	// validating email, role and name parameters
	return (validation.ValidateEmail(u.Email) &&
		validation.ValidateRole(u.Role) &&
		validation.ValidateName(u.Name))
}

// To do
// method to update name
// method to update role
// method to update name and role
// method to delete user
// method to get all users
// method to get a single user

func deleteUserFromFirebase(userId string, wg *sync.WaitGroup, errchan chan error) {
	defer wg.Done()

	err := firebase.FirebaseClient.DeleteUser(ctx, userId)
	if err != nil {
		log.Printf("Error deleting user from firebase: %v\n", err)
		errchan <- err
	}

	errchan <- nil
}

func deleteUserFromDatabase(userId string, wg *sync.WaitGroup, errchan chan error) {
	defer wg.Done()

	args := dbqueries.DeleteUserArgs(userId)
	_, err := db.Exec(ctx, dbqueries.DeleteUser, args)
	if err != nil {
		log.Printf("Error deleting user in database: %v\n", err)
		errchan <- err
	}

	errchan <- nil
}

// delete user
func (u *User) DeleteUser() error {
	errchan := make(chan error, 2)
	wg := new(sync.WaitGroup)

	wg.Add(2)
	go deleteUserFromFirebase(u.UserId, wg, errchan)
	go deleteUserFromDatabase(u.UserId, wg, errchan)

	wg.Wait()
	close(errchan)

	for err := range errchan {
		if err != nil {
			return err
		}
	}

	return nil
}
