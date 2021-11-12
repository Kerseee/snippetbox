package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"kerseeeHuang.com/snippetbox/pkg/models"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// UserModel is a wrapper of sql.db connection pool toward the users table in db.
type UserModel struct {
	DB	*sql.DB
}

// Insert insert an user into db if given user info are all valid.
func (m *UserModel) Insert(name, email, password string) error {
	// Hash the password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	// Prepare the statement.
	stmt := `INSERT INTO users (name, email, hashed_password, created)
		VALUES(?, ?, ?, UTC_TIMESTAMP())`

	// Execute the statement and handle errors if any.
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

// Authenticate authenticates the email addres and password, and return id
// if it pass the verification.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Get return the user detail based on given user id.
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}