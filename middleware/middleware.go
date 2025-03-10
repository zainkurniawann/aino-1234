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

// Middleware Logger Warna
func ColoredLogger(next echo.HandlerFunc) echo.HandlerFunc {
	// Warna untuk method dan status code
	var (
		Blue    = color.New(color.FgBlue).SprintFunc()    // Method: GET, POST, UPDATE, DELETE
		Green   = color.New(color.FgGreen).SprintFunc()   // Status: 200 (Success)
		Yellow  = color.New(color.FgYellow).SprintFunc()  // Status: 400 (Bad Request), 401 (Unauthorized), 404 (Not Found)
		Red     = color.New(color.FgRed).SprintFunc()     // Status: 500 (Error)
		Magenta = color.New(color.FgMagenta).SprintFunc() // Status: 403 (Forbidden)
	)

	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		stop := time.Now()

		// Ambil informasi request
		method := c.Request().Method
		statusCode := c.Response().Status
		path := c.Request().URL.Path
		latency := stop.Sub(start)

		// Tentukan pesan berdasarkan status code
		var statusMessage string
		switch statusCode {
		case 200:
			statusMessage = Green("Kok Iso To")
		case 400:
			statusMessage = Yellow("Bad Request Kocak")
		case 401:
			statusMessage = Yellow("Tokenmu wir")
		case 403:
			statusMessage = Magenta("Forbidden Woi")
		case 404:
			statusMessage = Yellow("Not Found wir")
		case 500:
			statusMessage = Red("Error Woilah")
		default:
			statusMessage = Yellow("Jann Ngrepoti Tenan")
		}

		// Pilih warna berdasarkan method
		var methodColor func(a ...interface{}) string
		switch method {
		case "GET", "POST", "UPDATE", "DELETE":
			methodColor = Blue
		default:
			methodColor = Blue
		}

		// Pilih warna berdasarkan status code
		var statusColor func(a ...interface{}) string
		switch statusCode {
		case 200:
			statusColor = Green // Success
		case 400, 401, 404:
			statusColor = Yellow // Bad Request, Unauthorized, Not Found
		case 403:
			statusColor = Magenta // Forbidden
		case 500:
			statusColor = Red // Error
		default:
			statusColor = Yellow
		}

		// Print log dengan warna dan pesan status
		fmt.Printf("%s %s - %s [%s] (%s)\n", methodColor(method), statusColor(statusCode), statusMessage, path, latency)

		return err
	}
}
