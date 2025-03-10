package models

import (
	"database/sql"
	"time"
)

type FormHA struct {
	UUID         string         `json:"form_uuid" db:"form_uuid"`
	DocumentUUID string         `json:"document_uuid" db:"document_uuid"`
	DocumentID   int64          `json:"document_id" db:"document_id"`
	UserID       int            `json:"user_id" db:"user_id" validate:"required"`
	FormNumber   string         `json:"form_number" db:"form_number"`
	FormTicket   string         `json:"form_ticket" db:"form_ticket"`
	FormStatus   string         `json:"form_status" db:"form_status"`
	Created_by   string         `json:"created_by" db:"created_by"`
	Created_at   time.Time      `json:"created_at" db:"created_at"`
	Updated_by   sql.NullString `json:"updated_by" db:"updated_by"`
	Updated_at   sql.NullTime   `json:"updated_at" db:"updated_at"`
	Deleted_by   sql.NullString `json:"deleted_by" db:"deleted_by"`
	Deleted_at   sql.NullTime   `json:"deleted_at" db:"deleted_at"`
}

type AddInfoHAReq struct {
	UUID string `json:"form_uuid" db:"form_uuid"`
	// Name     string `json:"name" db:"name"`
	NamaPengguna string `json:"nama_pengguna" db:"nama_pengguna"`
	RuangLingkup string `json:"ruang_lingkup" db:"ruang_lingkup"`
	JangkaWaktu  string `json:"jangka_waktu" db:"jangka_waktu"`
	// Password string `json:"password" db:"password"`
	// Scope    string `json:"scope" db:"scope"`
}

type AddInfoHA struct {
	UUID     string `json:"form_uuid" db:"form_uuid"`
	Name     string `json:"info_name" db:"name"`
	Instansi string `json:"instansi" db:"instansi"`
	Position string `json:"position" db:"position"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	Scope    string `json:"scope" db:"scope"`
}

type FormsHAReq struct {
	FormUUID       string         `json:"form_uuid" db:"form_uuid"`
	FormNumber     string         `json:"form_number" db:"form_number"`
	FormName       string         `json:"form_name" db:"form_name"`
	FormTicket     string         `json:"form_ticket" db:"form_ticket"`
	DocumentName   string         `json:"document_name" db:"document_name"`
	NamaTim        string         `json:"nama_tim" db:"nama_tim"`
	ProductManager string         `json:"product_manager" db:"product_manager"`
	NamaPengusul   string         `json:"nama_pengusul" db:"nama_pengusul"`
	TanggalUsul    string         `json:"tanggal_usul" db:"tanggal_usul"`
	FormType       string         `json:"form_type" db:"form_type"`
	FormStatus     string         `json:"form_status" db:"form_status"`
	CreatedBy      string         `json:"created_by" db:"created_by"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	UpdatedBy      sql.NullString `json:"updated_by" db:"updated_by"`
	UpdatedAt      sql.NullTime   `json:"updated_at" db:"updated_at"`
	DeletedBy      sql.NullString `json:"deleted_by" db:"deleted_by"`
	DeletedAt      sql.NullTime   `json:"deleted_at" db:"deleted_at"`
	ApprovalStatus string         `json:"approval_status" db:"approval_status"`
	Reason         string         `json:"reason" db:"reason"` // tambahkan field ini
}

type FormsHA struct {
	FormUUID       string         `json:"form_uuid" db:"form_uuid"`
	FormNumber     string         `json:"form_number" db:"form_number"`
	FormTicket     string         `json:"form_ticket" db:"form_ticket"`
	DocumentName   string         `json:"document_name" db:"document_name"`
	FormName       string         `json:"form_name" db:"form_name"`
	FormStatus     string         `json:"form_status" db:"form_status"`
	CreatedBy      string         `json:"created_by" db:"created_by"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	UpdatedBy      sql.NullString `json:"updated_by" db:"updated_by"`
	UpdatedAt      sql.NullTime   `json:"updated_at" db:"updated_at"`
	DeletedBy      sql.NullString `json:"deleted_by" db:"deleted_by"`
	DeletedAt      sql.NullTime   `json:"deleted_at" db:"deleted_at"`
	ApprovalStatus string         `json:"approval_status" db:"approval_status"`
	Reason         string         `json:"reason" db:"reason"` // tambahkan field ini
}

type HAReq struct {
	FormName       string `json:"form_name" db:"form_name"`
	NamaTim        string `json:"nama_tim" db:"nama_tim"`
	ProductManager string `json:"product_manager" db:"product_manager"`
	NamaPengusul   string `json:"nama_pengusul" db:"nama_pengusul"`
	TanggalUsul    string `json:"tanggal_usul" db:"tanggal_usul"`
	FormType       string `json:"form_type" db:"form_type"`
}

type HA struct {
	FormName string `json:"form_name" db:"form_name"`
	FormType string `json:"form_type" db:"form_type"`
}
type FormsHAAll struct {
	FormUUID      string         `json:"form_uuid" db:"form_uuid"`
	FormStatus    string         `json:"form_status" db:"form_status"`
	DocumentName  string         `json:"document_name" db:"document_name"`
	CreatedBy     string         `json:"created_by" db:"created_by"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
	UpdatedBy     sql.NullString `json:"updated_by" db:"updated_by"`
	UpdatedAt     sql.NullTime   `json:"updated_at" db:"updated_at"`
	DeletedBy     sql.NullString `json:"deleted_by" db:"deleted_by"`
	DeletedAt     sql.NullTime   `json:"deleted_at" db:"deleted_at"`
	FormName      string         `json:"form_name" db:"form_name"`
	InfoUUID      string         `json:"info_uuid" db:"info_uuid"`
	InfoName      string         `json:"info_name" db:"info_name"` // Ubah nama field dari "name" menjadi "info_name"
	Instansi      string         `json:"instansi" db:"instansi"`
	InfoPosition  string         `json:"position" db:"position"`
	Username      string         `json:"username" db:"username"`
	Password      string         `json:"password" db:"password"`
	Scope         string         `json:"scope" db:"scope"`
	UUID          string         `json:"sign_uuid" db:"sign_uuid"`
	SignatoryName string         `json:"signatory_name" db:"signatory_name"`         // Ubah nama field dari "name" menjadi "signatory_name"
	Position      string         `json:"signatory_position" db:"signatory_position"` // Ubah nama field dari "position" menjadi "signatory_position"
	Role          string         `json:"role_sign" db:"role_sign"`
	IsSign        bool           `json:"is_sign" db:"is_sign"`
}

type HakAksesInfo struct {
	InfoUUID string `json:"info_uuid" db:"info_uuid"`
	InfoName string `json:"info_name" db:"info_name"`
	Instansi string `json:"instansi" db:"instansi"`
	Position string `json:"position" db:"position"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	Scope    string `json:"scope" db:"scope"`
}

type HakAksesRequest struct {
	HAUUID       string `json:"ha_uuid" db:"ha_uuid"`
	NamaPengguna string `json:"nama_pengguna" db:"nama_pengguna"`
	RuangLingkup string `json:"ruang_lingkup" db:"ruang_lingkup"`
	JangkaWaktu  string `json:"jangka_waktu" db:"jangka_waktu"`
}

type HakAksesUpdateRequest struct {
	Added   []HakAksesRequest `json:"added"`
	Updated []HakAksesRequest `json:"updated"`
	InfoHA  []HakAksesRequest `json:"info_ha"`
	Deleted []HakAksesRequest `json:"deleted"`
}

type FormRequest struct {
	FormData  HAReq                 `json:"formData"`
	InfoHA    HakAksesUpdateRequest `json:"hakAksesInfoData"`
	Signatory []Signatory           `json:"signatories"`
}

////////////////////////////////////////////

type HakAksesReviewUpdateRequest struct {
	Added        []HakAksesInfo `json:"added"`
	Updated      []HakAksesInfo `json:"updated"`
	InfoHAReview []HakAksesInfo `json:"info_ha_review"`
	Deleted      []HakAksesInfo `json:"deleted"`
}

type FormRequestReview struct {
	FormData     HA                          `json:"formData"`
	Signatory    []Signatory                 `json:"signatories"`
	InfoHAReview HakAksesReviewUpdateRequest `json:"hak_akses_info_data"`
}
