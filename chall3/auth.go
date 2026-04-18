package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate = validator.New()
var jwtSecretKey = []byte("301af9181c7594747ec6498fedfaa185c1919fb1a15fe5f3e1aa4cea731cae1f")

// LoginRequest is cast directly from strict JSON input to enforce types
type LoginRequest struct {
	Username string `json:"username" validate:"required,alphanum,min=4,max=30"`
	Password string `json:"password" validate:"required,alphanum,min=8,max=30"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,alphanum,min=4,max=30"`
	Password string `json:"password" validate:"required,alphanum,min=8,max=30"`
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user User
	collection := DB.Collection("users")
	fmt.Printf("%v", req)
	var query bson.M
	err = json.Unmarshal(body, &query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("%v", query)
	err = collection.FindOne(r.Context(), query).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tokenString, err := GenerateJWT(user.ID.Hex(), user.Role)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
		"token":   tokenString,
	})
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	collection := DB.Collection("users")

	count, err := collection.CountDocuments(r.Context(), bson.M{"username": req.Username})
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	newUser := User{
		Username:  req.Username,
		Password:  req.Password,
		Role:      1,
		CreatedAt: time.Now(),
	}

	_, err = collection.InsertOne(r.Context(), newUser)
	if err != nil {
		http.Error(w, "Failed to register user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]string{"message": "Registration successful"}
	respBytes, _ := json.Marshal(response)
	w.Write(respBytes)
}
func GenerateJWT(userID string, role uint8) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey)
}
