package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"back-end/config"
	"back-end/utils"

	"github.com/golang-jwt/jwt/v5"
)

type AuthController struct{}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

var JwtKey = []byte("your_secret_key")

// POST /api/auth/register
func Register(w http.ResponseWriter, r *http.Request) {
	var input RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, "error", "Invalid JSON", nil)
		return
	}

	if input.Name == "" || input.Password == "" {
		utils.RespondJSON(w, http.StatusBadRequest, "error", "All fields are required", nil)
		return
	}

	var exists bool
	err := config.DB.QueryRow(
		"SELECT EXISTS (SELECT 1 FROM users WHERE name = $1)",
		input.Name,
	).Scan(&exists)
	if err != nil {
		log.Printf("Register: failed to check duplicate name: %v", err)
		utils.RespondJSON(w, http.StatusInternalServerError, "error", "DB error", nil)
		return
	}
	if exists {
		utils.RespondJSON(w, http.StatusConflict, "error", "Name already exists", nil)
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, "error", "Could not hash password", nil)
		return
	}

	now := time.Now()
	_, err = config.DB.Exec(
		"INSERT INTO users (name, password, created_at, updated_at) VALUES ($1, $2, $3, $3)",
		input.Name,
		hashedPassword,
		now,
	)
	if err != nil {
		log.Printf("Register: failed to insert user: %v", err)
		utils.RespondJSON(w, http.StatusInternalServerError, "error", "Could not register user", nil)
		return
	}

	utils.RespondJSON(w, http.StatusCreated, "success", "User registered", nil)
}

// POST /api/auth/login

func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, "error", "Invalid JSON", nil)
		return
	}

	if req.UserName == "" || req.Password == "" {
		utils.RespondJSON(w, http.StatusBadRequest, "error", "Username and password are required", nil)
		return
	}

	var (
		id             int
		hashedPassword string
		name           string
	)

	err := config.DB.QueryRow("SELECT  id, name,  password FROM users WHERE name = $1", req.UserName).Scan(&id, &name, &hashedPassword)
	if err == sql.ErrNoRows {
		utils.RespondJSON(w, http.StatusUnauthorized, "error", "Invalid credentials", nil)
		return
	}
	if err != nil {
		log.Printf("Login: failed to query user: %v", err)
		utils.RespondJSON(w, http.StatusInternalServerError, "error", "DB error", nil)
		return
	}
	if !utils.CheckPasswordHash(req.Password, hashedPassword) {
		utils.RespondJSON(w, http.StatusUnauthorized, "error", "Invalid credentials Password", nil)
		return
	}

	// สร้าง JWT (access token)
	claims := jwt.MapClaims{
		"user_id":  id,
		"username": req.UserName,
		"exp":      time.Now().Add(15 * time.Minute).Unix(), // access token อายุ 15 นาที
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(JwtKey)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, "error", "Token error", nil)
		return
	}

	// สร้าง refresh token
	refreshToken, err := utils.GenerateRandomToken(32)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, "error", "Could not create refresh token", nil)
		return
	}

	// เก็บ refresh token ลง DB (7 วัน)
	expiredAt := time.Now().Add(7 * 24 * time.Hour)
	_, err = config.DB.Exec("INSERT INTO refresh_tokens (user_id, token, expired_at) VALUES ($1, $2, $3)", id, refreshToken, expiredAt)
	if err != nil {
		log.Printf("Login: failed to insert refresh token for user %d: %v", id, err)
		utils.RespondJSON(w, http.StatusInternalServerError, "error", "DB error", nil)
		return
	}

	// ส่งทั้งสอง token กลับ
	utils.RespondJSON(w, http.StatusOK, "success", "Login successful", map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"datas": map[string]interface{}{
			"name": name,
		},
	})
}

// POST /api/auth/logout
func Logout(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		RefreshToken string `json:"refresh_token"`
	}

	var input Req
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.RefreshToken == "" {
		utils.RespondJSON(w, http.StatusBadRequest, "error", "Refresh token required", nil)
		return
	}

	// ลบ refresh token จาก DB
	_, err = config.DB.Exec("DELETE FROM refresh_tokens WHERE token = $1", input.RefreshToken)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, "error", "DB error", nil)
		return
	}
	// Logout ลบ refresh token ใน DB
	_, err = config.DB.Exec("DELETE FROM refresh_tokens WHERE token = $1", input.RefreshToken)
	utils.RespondJSON(w, http.StatusOK, "success", "Logged out", nil)
}
