package middleware

import (
	"document/models"
	"document/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	jose "github.com/dvsekhvalnov/jose2go"
	"github.com/fatih/color"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	UserId             int    `json:"user_id"`
	UserUUID           string `json:"user_uuid"`
	AppRoleId          int    `json:"application_role_id"`
	DivisionTitle      string `json:"division_title"`
	DivisionCode       string `json:"division_code"`
	RoleCode           string `json:"role_code"`
	Username           string `json:"user_name"`
	jwt.StandardClaims        // Embed the StandardClaims struct

}

func DecryptJWE(jweToken string, secretKey string) (string, error) {
	// Dekripsi token JWE
	decrypted, _, err := jose.Decode(jweToken, secretKey)
	if err != nil {
		return "", err
	}
	return decrypted, nil
}

// func DecryptJWE(jweToken string, secretKey string) (string, error) {
// 	// Dekripsi token JWE
// 	decrypted, _, err := jose.Decode(jweToken, secretKey)
// 	if err != nil {
// 		return "", err
// 	}
// 	return decrypted, nil
// }

func ExtractClaims(jwtToken string) (JwtCustomClaims, error) {
	claims := &JwtCustomClaims{}
	secretKey := "secretJwToken" // Ganti dengan kunci yang benar

	token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return JwtCustomClaims{}, err
	}

	return *claims, nil
}

func SuperAdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		secretKey := "secretJwToken" // Ganti dengan kunci yang benar
		_, exists := utils.InvalidTokens[tokenString]
		if exists {
			return c.JSON(http.StatusUnauthorized, &models.Response{
				Code:    401,
				Message: "Token tidak valid atau Anda telah logout",
				Status:  false,
			})
		}
		// Periksa apakah tokenString tidak kosong
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak ditemukan!",
				"status":  false,
			})
		}
		// _, exists := InvalidTokens[tokenString]
		// if exists {
		// 	return c.JSON(http.StatusUnauthorized, &models.Response{
		// 		Code:    401,
		// 		Message: "Token tidak valid atau Anda telah logout",
		// 		Status:  false,
		// 	})
		// }

		// Periksa apakah tokenString mengandung "Bearer "
		if !strings.HasPrefix(tokenString, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak valid!",
				"status":  false,
			})
		}

		// Hapus "Bearer " dari tokenString
		tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")

		// Langkah 1: Mendekripsi token JWE
		decrypted, err := DecryptJWE(tokenOnly, secretKey)
		if err != nil {
			fmt.Println("Gagal mendekripsi token:", err)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak valid!",
				"status":  false,
			})
		}

		fmt.Println("Token yang sudah dideskripsi:", decrypted)

		var claims JwtCustomClaims
		errJ := json.Unmarshal([]byte(decrypted), &claims)
		if errJ != nil {
			fmt.Println("Gagal mengurai klaim:", errJ)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak valid!",
				"status":  false,
			})
		}
		if time.Now().Unix() > claims.StandardClaims.ExpiresAt {
			// Token telah kedaluwarsa
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Sesi Anda sudah habis! Silahkan login kembali.",
				"status":  false,
			})
		}
		// Sekarang Anda memiliki data dalam struct JwtCustomClaims
		// Anda bisa mengakses UserId atau klaim lain sesuai kebutuhan
		// fmt.Println("UserID:", claims.UserId)

		userID := claims.UserId
		userUUID := claims.UserUUID // Mengakses UserID langsung
		userName := claims.Username
		roleID := claims.AppRoleId
		divisionTitle := claims.DivisionTitle
		roleCode := claims.RoleCode
		divisionCode := claims.DivisionCode
		if roleCode != "" {
			log.Print(roleCode)
		}

		fmt.Println("User ID:", userID)
		fmt.Println("User UUID:", userUUID)
		fmt.Println("User Name:", userName)
		fmt.Println("Role Code:", roleCode)
		fmt.Println("Division title:", divisionTitle)
		fmt.Println("Division Code : ", divisionCode)

		c.Set("user_id", userID)
		c.Set("user_name", userName)
		c.Set("division_code", divisionCode)
		c.Set("user_uuid", userUUID)
		c.Set("application_role_id", roleID)
		c.Set("division_title", divisionTitle)
		c.Set("role_code", roleCode)

		if roleCode != "SA" {
			log.Print(err)
			return c.JSON(http.StatusForbidden, &models.Response{
				Code:    403,
				Message: "Akses ditolak!",
				Status:  false,
			})
		}

		// Token JWE valid, Anda dapat melanjutkan dengan pengolahan berikutnya
		return next(c)
	}
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		secretKey := "secretJwToken" // Ganti dengan kunci yang benar

		// Periksa apakah tokenString tidak kosong
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak ditemukan!",
				"status":  false,
			})
		}

		// Periksa apakah tokenString mengandung "Bearer "
		if !strings.HasPrefix(tokenString, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak valid!",
				"status":  false,
			})
		}

		// Hapus "Bearer " dari tokenString
		tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")

		// Langkah 1: Mendekripsi token JWE
		decrypted, err := DecryptJWE(tokenOnly, secretKey)
		if err != nil {
			fmt.Println("Gagal mendekripsi token:", err)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak valid!",
				"status":  false,
			})
		}

		fmt.Println("Token yang sudah dideskripsi:", decrypted)

		var claims JwtCustomClaims
		errJ := json.Unmarshal([]byte(decrypted), &claims)
		if errJ != nil {
			fmt.Println("Gagal mengurai klaim:", errJ)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak valid!",
				"status":  false,
			})
		}

		// Langkah 3: Periksa apakah token sudah kedaluwarsa
		if time.Now().Unix() > claims.StandardClaims.ExpiresAt {
			// Token telah kedaluwarsa
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Sesi Anda sudah habis! Silahkan login kembali.",
				"status":  false,
			})
		}

		// Sekarang Anda memiliki data dalam struct JwtCustomClaims
		// Anda bisa mengakses UserId atau klaim lain sesuai kebutuhan
		// fmt.Println("UserID:", claims.UserId)

		userUUID := claims.UserUUID // Mengakses UserID langsung
		username := claims.Username
		userID := claims.UserId
		divisionCode := claims.DivisionCode
		// roleID := claims.AppRoleId
		// divisionTitle := claims.DivisionTitle
		// roleCode := claims.RoleCode
		// if roleCode != "" {
		// 	log.Print(roleCode)
		// }

		fmt.Println("User ID:", userID)
		fmt.Println("User UUID:", userUUID)
		fmt.Println("User Name:", username)
		fmt.Println("Division Code:", divisionCode)

		// fmt.Println("Role Code:", roleCode)

		c.Set("user_uuid", userUUID)
		c.Set("user_name", username)
		c.Set("user_id", userID)
		c.Set("division_code", divisionCode)
		// c.Set("application_role_id", roleID)
		// c.Set("division_title", divisionTitle)
		// c.Set("role_code", roleCode)

		// Token JWE valid, Anda dapat melanjutkan dengan pengolahan berikutnya
		return next(c)
	}
}

func GuestMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		secretKey := "secretJwToken" // Sesuaikan dengan kunci JWT

		// Jika token tidak dikirim, anggap sebagai guest
		if tokenString == "" {
			c.Set("is_guest", true) // Tandai sebagai guest
			return next(c)
		}

		// Jika token ada, cek apakah valid
		if !strings.HasPrefix(tokenString, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak valid!",
				"status":  false,
			})
		}

		tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")

		// Dekripsi token
		decrypted, err := DecryptJWE(tokenOnly, secretKey)
		if err != nil {
			fmt.Println("Gagal mendekripsi token:", err)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak valid!",
				"status":  false,
			})
		}

		var claims JwtCustomClaims
		errJ := json.Unmarshal([]byte(decrypted), &claims)
		if errJ != nil {
			fmt.Println("Gagal mengurai klaim:", errJ)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak valid!",
				"status":  false,
			})
		}

		// Periksa apakah token kedaluwarsa
		if time.Now().Unix() > claims.StandardClaims.ExpiresAt {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Sesi Anda sudah habis! Silahkan login kembali.",
				"status":  false,
			})
		}

		// Simpan data user ke context jika ada token
		c.Set("user_uuid", claims.UserUUID)
		c.Set("user_name", claims.Username)
		c.Set("user_id", claims.UserId)
		c.Set("division_code", claims.DivisionCode)
		c.Set("role_code", claims.RoleCode)
		c.Set("is_guest", false) // Tandai bukan guest

		return next(c)
	}
}

func AdminMemberMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		secretKey := "secretJwToken" // Ganti dengan kunci yang benar

		// Periksa apakah tokenString tidak kosong
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak ditemukan!",
				"status":  false,
			})
		}

		// Periksa apakah tokenString mengandung "Bearer "
		if !strings.HasPrefix(tokenString, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak valid!",
				"status":  false,
			})
		}

		// Hapus "Bearer " dari tokenString
		tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")

		// Langkah 1: Mendekripsi token JWE
		decrypted, err := DecryptJWE(tokenOnly, secretKey)
		if err != nil {
			fmt.Println("Gagal mendekripsi token:", err)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak valid!",
				"status":  false,
			})
		}

		fmt.Println("Token yang sudah dideskripsi:", decrypted)

		var claims JwtCustomClaims
		errJ := json.Unmarshal([]byte(decrypted), &claims)
		if errJ != nil {
			fmt.Println("Gagal mengurai klaim:", errJ)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak valid!",
				"status":  false,
			})
		}

		// Langkah 3: Periksa apakah token sudah kedaluwarsa
		if time.Now().Unix() > claims.StandardClaims.ExpiresAt {
			// Token telah kedaluwarsa
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Sesi Anda sudah habis! Silahkan login kembali.",
				"status":  false,
			})
		}

		// Sekarang Anda memiliki data dalam struct JwtCustomClaims
		// Anda bisa mengakses UserId atau klaim lain sesuai kebutuhan
		// fmt.Println("UserID:", claims.UserId)

		userUUID := claims.UserUUID // Mengakses UserID langsung
		username := claims.Username
		userID := claims.UserId
		divisionCode := claims.DivisionCode
		// roleID := claims.AppRoleId
		// divisionTitle := claims.DivisionTitle
		roleCode := claims.RoleCode
		if roleCode != "" {
			log.Print(roleCode)
		}

		fmt.Println("User ID:", userID)
		fmt.Println("User UUID:", userUUID)
		fmt.Println("User Name:", username)
		fmt.Println("Division Code:", divisionCode)
		fmt.Println("Role Code:", roleCode)

		c.Set("user_uuid", userUUID)
		c.Set("user_name", username)
		c.Set("user_id", userID)
		c.Set("division_code", divisionCode)
		c.Set("role_code", roleCode)

		if roleCode == "SA" {
			log.Print(err)
			// Jika role code adalah SA, kembalikan pesan Unauthorized
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Akses ditolak! Anda tidak memiliki izin untuk mengakses ini.",
				"status":  false,
			})
		}
		// c.Set("application_role_id", roleID)
		// c.Set("division_title", divisionTitle)
		// c.Set("role_code", roleCode)

		// Token JWE valid, Anda dapat melanjutkan dengan pengolahan berikutnya
		return next(c)
	}
}
func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		secretKey := "secretJwToken" // Ganti dengan kunci yang benar
		_, exists := utils.InvalidTokens[tokenString]
		if exists {
			return c.JSON(http.StatusUnauthorized, &models.Response{
				Code:    401,
				Message: "Token tidak valid atau Anda telah logout",
				Status:  false,
			})
		}
		// Periksa apakah tokenString tidak kosong
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak ditemukan!",
				"status":  false,
			})
		}

		// Periksa apakah tokenString mengandung "Bearer "
		if !strings.HasPrefix(tokenString, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak valid!",
				"status":  false,
			})
		}

		// Hapus "Bearer " dari tokenString
		tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")

		// Langkah 1: Mendekripsi token JWE
		decrypted, err := DecryptJWE(tokenOnly, secretKey)
		if err != nil {
			fmt.Println("Gagal mendekripsi token:", err)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak valid!",
				"status":  false,
			})
		}

		fmt.Println("Token yang sudah dideskripsi:", decrypted)

		var claims JwtCustomClaims
		errJ := json.Unmarshal([]byte(decrypted), &claims)
		if errJ != nil {
			fmt.Println("Gagal mengurai klaim:", errJ)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak valid!",
				"status":  false,
			})
		}

		if time.Now().Unix() > claims.StandardClaims.ExpiresAt {
			// Token telah kedaluwarsa
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Sesi Anda sudah habis! Silahkan login kembali.",
				"status":  false,
			})
		}
		// Sekarang Anda memiliki data dalam struct JwtCustomClaims
		// Anda bisa mengakses UserId atau klaim lain sesuai kebutuhan
		// fmt.Println("UserID:", claims.UserId)

		userID := claims.UserId
		userUUID := claims.UserUUID // Mengakses UserID langsung
		userName := claims.Username
		roleID := claims.AppRoleId
		divisionTitle := claims.DivisionTitle
		roleCode := claims.RoleCode
		if roleCode != "" {
			log.Print(roleCode)
		}

		fmt.Println("User ID:", userID)
		fmt.Println("User UUID:", userUUID)
		fmt.Println("User Name:", userName)
		fmt.Println("Role Code:", roleCode)
		fmt.Println("Division title:", divisionTitle)

		c.Set("user_id", userID)
		c.Set("user_uuid", userUUID)
		c.Set("user_name", userName)
		c.Set("application_role_id", roleID)
		c.Set("division_title", divisionTitle)
		c.Set("role_code", roleCode)
		if roleCode != "A" {
			return c.JSON(http.StatusForbidden, &models.Response{
				Code:    403,
				Message: "Akses ditolak!",
				Status:  false,
			})
		}

		// Token JWE valid, Anda dapat melanjutkan dengan pengolahan berikutnya
		return next(c)
	}
}

// Middleware untuk logging request
func ColoredLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		err := next(c)

		method := c.Request().Method
		statusCode := c.Response().Status
		path := c.Request().URL.Path
		latency := time.Since(start)

		methodColor := color.New(color.FgBlue).SprintFunc()
		statusColor := getStatusColor(statusCode)
		
		// Format status code + status message dengan warna
		coloredStatus := statusColor(fmt.Sprintf("%d %s", statusCode, getStatusMessage(statusCode)))

		fmt.Printf("%s %s - %s (%v)\n", methodColor(method), path, coloredStatus, latency)

		return err
	}
}

// Fungsi untuk menentukan warna berdasarkan status code
func getStatusColor(statusCode int) func(a ...interface{}) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return green
	case statusCode >= 400 && statusCode < 500:
		return yellow
	case statusCode >= 500:
		return red
	default:
		return white
	}
}

// Fungsi untuk menentukan status message berdasarkan kode status
func getStatusMessage(statusCode int) string {
	switch statusCode {
	case 200:
		return "Success"
	case 201:
		return "Created"
	case 204:
		return "No Content (No Changes)"
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized"
	case 403:
		return "Forbidden"
	case 404:
		return "Not Found"
	case 405:
		return "Method Not Allowed"
	case 409:
		return "Conflict"
	case 422:
		return "Unprocessable Entity"
	case 500:
		return "Internal Server Error"
	default:
		return "Unknown Status"
	}
}

// Warna untuk status code
var (
	green  = color.New(color.FgGreen).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	white  = color.New(color.FgWhite).SprintFunc()
)
