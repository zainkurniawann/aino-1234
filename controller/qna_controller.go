package controller

import (
	"document/models"
	"document/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// GetAllQnA mengambil semua data QnA
func GetAllQnA(c echo.Context) error {
	qnaList, err := service.GetAllQnA()
	if err != nil {
		log.Print(err)
		response := models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal server. Mohon coba beberapa saat lagi",
			Status:  false,
			Data:    nil,
		}
		return c.JSON(http.StatusInternalServerError, response)
	}

	response := models.Response{
		Code:    200,
		Message: "Data QnA berhasil diambil",
		Status:  true,
		Data:    qnaList,
	}

	return c.JSON(http.StatusOK, response)
}

func GetSpecQnA(c echo.Context) error {
	id := c.Param("id")
	
	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken"

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

	// Dekripsi token JWE
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
	if err := json.Unmarshal([]byte(decrypted), &claims); err != nil {
		fmt.Println("Gagal mengurai klaim:", err)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	qna, err := service.GetSpecQnA(id)
	if err != nil {
		log.Print(err)
		response := models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal server. Mohon coba beberapa saat lagi",
			Status:  false,
			Data:    nil,
		}
		return c.JSON(http.StatusInternalServerError, response)
	}

	if qna == nil {
		response := models.Response{
			Code:    404,
			Message: "Data QnA tidak ditemukan",
			Status:  false,
			Data:    nil,
		}
		return c.JSON(http.StatusNotFound, response)
	}

	response := models.Response{
		Code:    200,
		Message: "Data QnA berhasil diambil",
		Status:  true,
		Data:    qna,
	}

	return c.JSON(http.StatusOK, response)
}

func AddQnA(c echo.Context) error {
	var req models.QnARequest

	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken"

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

	// Dekripsi token JWE
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
	if err := json.Unmarshal([]byte(decrypted), &claims); err != nil {
		fmt.Println("Gagal mengurai klaim:", err)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	// Ambil username dari claims JWT
	createdBy := c.Get("user_name").(string)
	if createdBy == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Invalid token atau token tidak ditemukan!",
			"status":  false,
		})
	}

	// Bind request body ke struct
	if err := c.Bind(&req); err != nil {
		response := models.Response{
			Code:    400,
			Message: "Format data tidak valid",
			Status:  false,
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Validasi data wajib
	if req.Question == "" || req.Answer == "" {
		response := models.Response{
			Code:    400,
			Message: "Question dan Answer harus diisi",
			Status:  false,
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Memanggil service untuk menambahkan QnA
	if err := service.AddQnA(req, createdBy); err != nil {
		log.Print(err)
		response := models.Response{
			Code:    500,
			Message: "Gagal menambahkan data QnA",
			Status:  false,
		}
		return c.JSON(http.StatusInternalServerError, response)
	}

	response := models.Response{
		Code:    201,
		Message: "Data QnA berhasil ditambahkan",
		Status:  true,
	}
	return c.JSON(http.StatusCreated, response)
}

func UpdateQnA(c echo.Context) error {
	var req models.QnARequest
	id := c.Param("id")

	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken"

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

	// Dekripsi token JWE
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
	if err := json.Unmarshal([]byte(decrypted), &claims); err != nil {
		fmt.Println("Gagal mengurai klaim:", err)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	// Ambil username dari claims JWT
	updatedBy := c.Get("user_name").(string)
	if updatedBy == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Invalid token atau token tidak ditemukan!",
			"status":  false,
		})
	}

	// Bind request body ke struct
	if err := c.Bind(&req); err != nil {
		response := models.Response{
			Code:    400,
			Message: "Format data tidak valid",
			Status:  false,
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Validasi data wajib
	if req.Question == "" || req.Answer == "" {
		response := models.Response{
			Code:    400,
			Message: "Question dan Answer harus diisi",
			Status:  false,
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Memanggil service untuk menambahkan QnA
	if err := service.UpdateQnA(id, req, updatedBy); err != nil {
		log.Print(err)
		response := models.Response{
			Code:    500,
			Message: "Gagal mengedit data QnA",
			Status:  false,
		}
		return c.JSON(http.StatusInternalServerError, response)
	}

	response := models.Response{
		Code:    201,
		Message: "Data QnA berhasil diedit",
		Status:  true,
	}
	return c.JSON(http.StatusCreated, response)
}

func DeleteQnA(c echo.Context) error {
	id := c.Param("id")

	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken"

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

	// Dekripsi token JWE
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
	if err := json.Unmarshal([]byte(decrypted), &claims); err != nil {
		fmt.Println("Gagal mengurai klaim:", err)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	// Ambil username dari claims JWT
	deletedBy := c.Get("user_name").(string)
	if deletedBy == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Invalid token atau token tidak ditemukan!",
			"status":  false,
		})
	}

	// Memanggil service untuk menambahkan QnA
	if err := service.DeleteQnA(id, deletedBy); err != nil {
		log.Print(err)
		response := models.Response{
			Code:    500,
			Message: "Gagal menghapus data QnA",
			Status:  false,
		}
		return c.JSON(http.StatusInternalServerError, response)
	}

	response := models.Response{
		Code:    201,
		Message: "Data QnA berhasil dihapus",
		Status:  true,
	}
	return c.JSON(http.StatusCreated, response)
}
