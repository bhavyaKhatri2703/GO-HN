package auth

import (
	"database/sql"
	"fmt"
	"net/http"
)

func registerHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	hash, err := encrypt(password)

	if err != nil {
		http.Error(w, "password hashing error", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO user VALUES (email , password_hash) , VALUES ($1,$2)", email, hash)

	if err != nil {
		http.Error(w, "Erorr registering", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "User registered successfully")

}

func loginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	var hash string
	var id int
	err := db.QueryRow("SELECT id, password_hash FROM users WHERE email=$1", email).Scan(&id, &hash)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if !CheckPasswordHash(password, hash) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "Welcome user %d!", id)
}
