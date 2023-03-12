package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var db *sql.DB

// UsersTable handle the users table
func UsersTable(database *sql.DB) {
	db = database
}

func parseID(c *gin.Context) (int64, error) {
	id := c.Param("id")
	parsedID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "The user ID is not valid",
		})
		return parsedID, err
	}
	return parsedID, nil
}

// ListUsers list all users
func ListUsers(c *gin.Context) {
	var users []*User
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		fmt.Printf("Error getting users: %s \n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			fmt.Printf("Failed to scan user: %s \n", err)
		}

		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Rows error: %s \n", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}

// CreateUser creates a user in database
func CreateUser(c *gin.Context) {
	var user User
	c.BindJSON(&user)
	name, email, password := user.Name, user.Email, user.Password
	newUser := User{
		Name:      name,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := db.Exec("INSERT INTO users (name, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		name, email, password, newUser.CreatedAt, newUser.UpdatedAt)
	if err != nil {
		fmt.Printf("Failed to create user: %v \n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	lastInsertID, _ := result.LastInsertId()
	data := map[string]interface{}{
		"id":       lastInsertID,
		"name":     name,
		"email":    email,
		"password": password,
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"data":   data,
	})
}

// UpdateUser updates the user by Id
func UpdateUser(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}

	var user User
	c.BindJSON(&user)

	// Build update query
	var updateQuery []string
	if user.Name != "" {
		updateQuery = append(updateQuery, fmt.Sprintf("name='%s'", user.Name))
	}
	if user.Email != "" {
		updateQuery = append(updateQuery, fmt.Sprintf("email='%s'", user.Email))
	}
	if user.Name != "" {
		updateQuery = append(updateQuery, fmt.Sprintf("password='%s'", user.Password))
	}
	updateString := strings.Join(updateQuery, ",")

	// Build update statement
	updateStmt, err := db.Prepare("UPDATE users SET " + updateString + " WHERE id=?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Error preparing update statement",
		})
		return
	}

	// Execute update statement
	result, err := updateStmt.Exec(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Error updating user",
		})
		return
	}

	// checking
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Error checking rows affected",
		})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Could not find user",
		})
		return
	}

	var updatedUser User
	err = db.QueryRow("SELECT id, name, email, password FROM users WHERE id=?", id).Scan(&updatedUser.ID, &updatedUser.Name, &updatedUser.Email, &updatedUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Error getting updated user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   updatedUser,
	})
}

// DeleteUser deletes a user
func DeleteUser(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}

	result, err := db.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Error deleting user",
		})
		return
	}

	// checking...
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Error checking rows affected",
		})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Could not find user",
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
