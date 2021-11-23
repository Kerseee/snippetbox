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
	DB *sql.DB
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
	// Retrive id and hashed password with given email.
	var id int
	var hashedPassword []byte
	stmt := "SELECT id, hashed_password FROM users WHERE email = ? AND active = TRUE"
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// Validate password.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// Return user id.
	return id, nil
}

// Get return the user detail based on given user id.
func (m *UserModel) Get(id int) (*models.User, error) {
	u := &models.User{}

	stmt := `SELECT id, name, email, created, active FROM users WHERE id = ?`
	err := m.DB.QueryRow(stmt, id).Scan(&u.ID, &u.Name, &u.Email, &u.Created, &u.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return u, nil
}
