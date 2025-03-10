package controller

import (
	"database/sql"
	"document/database"
	"document/models"
	"document/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	jose "github.com/dvsekhvalnov/jose2go"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

var db = database.Connection()

func DecryptJWE(jweToken string, secretKey string) (string, error) {
	decrypted, _, err := jose.Decode(jweToken, secretKey)
	if err != nil {
		return "", err
	}
	return decrypted, nil
}

type JwtCustomClaims struct {
	UserID   int    `json:"user_id"`
	UserUUID string `json:"user_uuid"`
	RoleCode string `json:"role_code"`
	UserName string `json:"user_name"`
	DivisionTitle      string `json:"division_title"`
	DivisionCode      string `json:"division_code"`
	jwt.StandardClaims       
}

func AddDocument(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken"

	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

	if !strings.HasPrefix(tokenString, "Bearer ") {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")

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
	userName := c.Get("user_name").(string)  

	fmt.Println("Token yang sudah dideskripsi:", decrypted)

	fmt.Println("User name:", userName)

	if userName == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Invalid token atau token tidak ditemukan!",
			"status":  false,
		})
	}
	var addDocument models.Document
	if err := c.Bind(&addDocument); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak valid!",
			Status:  false,
		})
	}

	whitespace := regexp.MustCompile(`^\s`)
	if whitespace.MatchString(addDocument.Code) {
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Code tidak boleh dimulai dengan spasi!",
			Status:  false,
		})
	}

	if whitespace.MatchString(addDocument.Name) {
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Name tidak boleh dimulai dengan spasi!",
			Status:  false,
		})
	}

	errVal := c.Validate(&addDocument)

	if errVal == nil {
		var existingDocumentID int
		err := db.QueryRow("SELECT document_id FROM document_ms WHERE (document_code = $1 OR document_name = $2) AND deleted_at IS NULL", addDocument.Code, addDocument.Name).Scan(&existingDocumentID)

		if err == nil {
			return c.JSON(http.StatusBadRequest, &models.Response{
				Code:    400,
				Message: "Gagal menambahkan document. Document sudah ada!",
				Status:  false,
			})
		} else {
			addroleErr := service.AddDocument(addDocument, userName)
			if addroleErr != nil {
				return c.JSON(http.StatusInternalServerError, &models.Response{
					Code:    500,
					Message: "Terjadi kesalahan internal pada server. Coba beberapa saat lagi",
					Status:  false,
				})
			}

			return c.JSON(http.StatusCreated, &models.Response{
				Code:    201,
				Message: "Berhasil menambahkan document!",
				Status:  true,
			})
		}
	} else {
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Data tidak boleh kosong!",
			Status:  false,
		})
	}

}

func GetAllDoc(c echo.Context) error {
	documents, err := service.GetAllDoc()
	if err != nil {
		response := models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal server. Mohon coba beberapa saat lagi",
			Status:  false,
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	return c.JSON(http.StatusOK, documents)

}

func ShowDocById(c echo.Context) error {
	id := c.Param("id")

	var getDoc models.Document

	getDoc, err := service.ShowDocById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			response := models.Response{
				Code:    404,
				Message: "Document tidak ditemukan!",
				Status:  false,
			}
			return c.JSON(http.StatusNotFound, response)
		} else {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
				Status:  false,
			})
		}
	}

	return c.JSON(http.StatusOK, getDoc)
}

func UpdateDocument(c echo.Context) error {
	id := c.Param("id")

	perviousContent, errGet := service.ShowDocById(id)
	log.Println("perviousContent", perviousContent)
	if errGet != nil {
		return c.JSON(http.StatusNotFound, &models.Response{
			Code:    404,
			Message: "Gagal mengupdate document. Document tidak ditemukan!",
			Status:  false,
		})
	}

	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken"

	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

	if !strings.HasPrefix(tokenString, "Bearer ") {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")

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
	userName := c.Get("user_name").(string)

	fmt.Println("Token yang sudah dideskripsi:", decrypted)

	fmt.Println("User name:", userName)

	if userName == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Invalid token atau token tidak ditemukan!",
			"status":  false,
		})
	}

	var editDoc models.Document
	if err := c.Bind(&editDoc); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data invalid!",
			Status:  false,
		})
	}
	whitespace := regexp.MustCompile(`^\s`)
	if whitespace.MatchString(editDoc.Code) {
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Code tidak boleh dimulai dengan spasi!",
			Status:  false,
		})
	}

	if whitespace.MatchString(editDoc.Name) {
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Name tidak boleh dimulai dengan spasi!",
			Status:  false,
		})
	}
	err = c.Validate(&editDoc)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    422,
			Message: "Data tidak boleh kosong!",
			Status:  false,
		})
	}
	if err == nil {
		existingDoc, err := service.GetDocCodeName(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Terjadi kesalahan internal pada server.",
				Status:  false,
			})
		}

		if editDoc.Code != "" && editDoc.Code != existingDoc.Code {
			existingDocID, err := service.GetDocumentIDByCode(editDoc.Code)
			if err == nil && strconv.Itoa(existingDocID) != id {
				return c.JSON(http.StatusBadRequest, &models.Response{
					Code:    400,
					Message: "Document dengan code tersebut sudah ada! Document tidak boleh sama!",
					Status:  false,
				})
			}
		}

		if editDoc.Name != "" && editDoc.Name != existingDoc.Name {
			existingDocID, err := service.GetDocumentIDByName(editDoc.Name)
			if err == nil && strconv.Itoa(existingDocID) != id {
				return c.JSON(http.StatusBadRequest, &models.Response{
					Code:    400,
					Message: "Document dengan name tersebut sudah ada! Document tidak boleh sama!",
					Status:  false,
				})
			}
		}

		_, errService := service.UpdateDocument(editDoc, id, userName)
		if errService != nil {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
				Status:  false,
			})
		}

		return c.JSON(http.StatusOK, &models.Response{
			Code:    200,
			Message: "Document berhasil diperbarui!",
			Status:  true,
		})
	} else {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
			Status:  false,
		})
	}
}

func DeleteDoc(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken"

	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

	if !strings.HasPrefix(tokenString, "Bearer ") {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")

	decrypted, err := DecryptJWE(tokenOnly, secretKey)
	if err != nil {
		fmt.Println("Gagal mendekripsi token:", err)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}
	id := c.Param("id")
	_, errGet := service.ShowDocById(id)
	if errGet != nil {
		return c.JSON(http.StatusNotFound, &models.Response{
			Code:    404,
			Message: "Gagal menghapus document. Document tidak ditemukan!",
			Status:  false,
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

	userName := c.Get("user_name").(string)

	errService := service.DeleteDoc(id, userName)
	if errService != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
			Status:  false,
		})

	}

	return c.JSON(http.StatusOK, &models.Response{
		Code:    200,
		Message: "Document berhasil dihapus!",
		Status:  true,
	})

}
