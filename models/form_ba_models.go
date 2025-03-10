package models

import (
	"database/sql"
	"time"
)

type BA struct {
	Judul          string `json:"judul"`
	Tanggal        string `json:"tanggal"`
	AppName        string `json:"nama_aplikasi" db:"nama_aplikasi"`
	NoDA           string `json:"no_da"`
	NoITCM         string `json:"no_itcm"`
	DilakukanOleh  string `json:"dilakukan_oleh"`
	DidampingiOleh string `json:"didampingi_oleh"`
	Keterangan     string `json:"keterangan"`
}

type FormsBA struct {
	FormUUID       string         `json:"form_uuid" db:"form_uuid"`
	FormNumber     string         `json:"form_number" db:"form_number"`
	FormTicket     string         `json:"form_ticket" db:"form_ticket"`
	FormStatus     string         `json:"form_status" db:"form_status"`
	DocumentName   string         `json:"document_name" db:"document_name"`
	ProjectName    string         `json:"project_name" db:"project_name"`
	CreatedBy      string         `json:"created_by" db:"created_by"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	UpdatedBy      sql.NullString `json:"updated_by" db:"updated_by"`
	UpdatedAt      sql.NullTime   `json:"updated_at" db:"updated_at"`
	DeletedBy      sql.NullString `json:"deleted_by" db:"deleted_by"`
	DeletedAt      sql.NullTime   `json:"deleted_at" db:"deleted_at"`
	Judul          string         `json:"judul" db:"judul"`
	Tanggal        string         `json:"tanggal" db:"tanggal"`
	AppName        string         `json:"nama_aplikasi" db:"nama_aplikasi"`
	NoDA           string         `json:"no_da" db:"no_da"`
	NoITCM         string         `json:"no_itcm" db:"no_itcm"`
	DilakukanOleh  string         `json:"dilakukan_oleh" db:"dilakukan_oleh"`
	DidampingiOleh string         `json:"didampingi_oleh" db:"didampingi_oleh"`
}

type FormsBAAll struct {
	FormUUID       string         `json:"form_uuid" db:"form_uuid"`
	FormNumber     string         `json:"form_number" db:"form_number"`
	FormTicket     string         `json:"form_ticket" db:"form_ticket"`
	FormStatus     string         `json:"form_status" db:"form_status"`
	DocumentName   string         `json:"document_name" db:"document_name"`
	ProjectName    string         `json:"project_name" db:"project_name"`
	CreatedBy      string         `json:"created_by" db:"created_by"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	UpdatedBy      sql.NullString `json:"updated_by" db:"updated_by"`
	UpdatedAt      sql.NullTime   `json:"updated_at" db:"updated_at"`
	DeletedBy      sql.NullString `json:"deleted_by" db:"deleted_by"`
	DeletedAt      sql.NullTime   `json:"deleted_at" db:"deleted_at"`
	Judul          string         `json:"judul" db:"judul"`
	Tanggal        string         `json:"tanggal" db:"tanggal"`
	AppName        string         `json:"nama_aplikasi" db:"nama_aplikasi"`
	NoDA           string         `json:"no_da" db:"no_da"`
	NoITCM         string         `json:"no_itcm" db:"no_itcm"`
	DilakukanOleh  string         `json:"dilakukan_oleh" db:"dilakukan_oleh"`
	DidampingiOleh string         `json:"didampingi_oleh" db:"didampingi_oleh"`
	Keterangan     string         `json:"keterangan" db:"keterangan"`
	UUID           string         `json:"sign_uuid" db:"sign_uuid"`
	Name           string         `json:"name" db:"name"`
	Position       string         `json:"position" db:"position"`
	Role           string         `json:"role_sign" db:"role_sign"`
	IsSign         bool           `json:"is_sign" db:"is_sign"`
}

type FormsBeritaAcara struct {
	FormUUID     string         `json:"form_uuid" db:"form_uuid"`
	FormNumber   string         `json:"form_number" db:"form_number"`
	FormStatus   string         `json:"form_status" db:"form_status"`
	DocumentName string         `json:"document_name" db:"document_name"`
	CreatedBy    string         `json:"created_by" db:"created_by"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
	UpdatedBy    sql.NullString `json:"updated_by" db:"updated_by"`
	UpdatedAt    sql.NullTime   `json:"updated_at" db:"updated_at"`
	DeletedBy    sql.NullString `json:"deleted_by" db:"deleted_by"`
	DeletedAt    sql.NullTime   `json:"deleted_at" db:"deleted_at"`
	BeritaAcara
}

type FormsBeritaAcaraAsset struct {
	FormUUID     string         `json:"form_uuid" db:"form_uuid"`
	FormNumber   string         `json:"form_number" db:"form_number"`
	FormStatus   string         `json:"form_status" db:"form_status"`
	DocumentName string         `json:"document_name" db:"document_name"`
	CreatedBy    string         `json:"created_by" db:"created_by"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
	UpdatedBy    sql.NullString `json:"updated_by" db:"updated_by"`
	UpdatedAt    sql.NullTime   `json:"updated_at" db:"updated_at"`
	DeletedBy    sql.NullString `json:"deleted_by" db:"deleted_by"`
	DeletedAt    sql.NullTime   `json:"deleted_at" db:"deleted_at"`
	NamaAsset    string         `json:"nama_asset" db:"asset_name"`
	SerialNumber string         `json:"serial_number" db:"serial_number"`
	Spesifikasi  string         `json:"spesifikasi" db:"asset_specification"`
	AssetType    string         `json:"asset_type" db:"asset_type"`
	Image        []string       `json:"image"`
	BeritaAcara
}

type BeritaAcara struct {
	AssetUUID           string `json:"asset_uuid" db:"asset_uuid"`
	PicUUID             string `json:"pic_uuid" db:"pic_uuid"`
	PihakPertama        string `json:"pihak_pertama" db:"pihak_pertama"`
	JabatanPihakPertama string `json:"jabatan_pihak_pertama" db:"jabatan_pihak_pertama"`
	NamaPIC             string `json:"nama_pic" db:"nama_pic"`
	JabatanPIC          string `json:"jabatan_pic" db:"jabatan_pic"`
	Jenis               string `json:"jenis" db:"jenis"`
	KodeAsset           string `json:"kode_asset" db:"kode_asset"`
	Aksesoris           string `json:"aksesoris"`
	Reason              string `json:"reason"`
	Kondisi             string `json:"kondisi"`
	Merk                string `json:"merk"`
	Model               string `json:"model"`
	Start               string `json:"start_at" db:"start_at"`
	Ended               string `json:"ended_at" db:"ended_at"`
	Image               string `json:"image" db:"image_path"`
}

type UpdateRequest struct {
	Asset Asset           `json:"asset"`
	Pic   UpdatePicStruct `json:"pic"`
}
type UpdateImageRequest struct {
	Asset ImageAssetUpdateReq `json:"asset"`
	Pic   UpdatePicStruct     `json:"pic"`
}

type UpdatePicStruct struct {
	Added   []Pic `json:"added"`
	Updated []Pic `json:"updated"`
	Deleted []Pic `json:"deleted"`
	Pics    []Pic `json:"pics"`
}

type Pic struct {
	PicUUID    string `json:"pic_uuid" db:"pic_uuid"`
	NamaPic    string `json:"nama_pic" db:"pic_name"`
	Keterangan string `json:"keterangan" db:"pic_description"`
	Start      string `json:"start_at" db:"start_at"`
	Ended      string `json:"ended_at" db:"ended_at"`
}

type Asset struct {
	AssetUUID    string         `json:"asset_uuid" db:"asset_uuid"`
	Kode         string         `json:"kode_asset" db:"asset_code"`
	NamaAsset    string         `json:"nama_asset" db:"asset_name"`
	Merk         string         `json:"merk"`
	Model        string         `json:"model"`
	SerialNumber string         `json:"serial_number" db:"serial_number"`
	Spesifikasi  string         `json:"spesifikasi" db:"asset_specification"`
	TglPengadaan string         `json:"tgl_pengadaan" db:"procurement_date"`
	Harga        string         `json:"harga" db:"price"`
	Deskripsi    string         `json:"deskripsi" db:"asset_description"`
	Klasifikasi  string         `json:"klasifikasi" db:"system_classification"`
	Lokasi       string         `json:"lokasi" db:"asset_location"`
	Status       string         `json:"status" db:"asset_status"`
	AssetID      string         `json:"-" db:"asset_id"`
	AssetImg     []string       `json:"image" db:"asset_img"`
	AssetType    string         `json:"asset_type" db:"asset_type"`
	CreatedBy    string         `json:"created_by" db:"created_by"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
	UpdatedBy    sql.NullString `json:"updated_by" db:"updated_by"`
	UpdatedAt    sql.NullTime   `json:"updated_at" db:"updated_at"`
	DeletedBy    sql.NullString `json:"deleted_by" db:"deleted_by"`
	DeletedAt    sql.NullTime   `json:"deleted_at" db:"deleted_at"`
	Pic          []Pic          `json:"pic"`
}

type ImageAssetUpdateReq struct {
	AssetUUID    string         `json:"asset_uuid" db:"asset_uuid"`
	Kode         string         `json:"kode_asset" db:"asset_code"`
	NamaAsset    string         `json:"nama_asset" db:"asset_name"`
	SerialNumber string         `json:"serial_number" db:"serial_number"`
	Spesifikasi  string         `json:"spesifikasi" db:"asset_specification"`
	TglPengadaan string         `json:"tgl_pengadaan" db:"procurement_date"`
	Harga        string         `json:"harga" db:"price"`
	Deskripsi    string         `json:"deskripsi" db:"asset_description"`
	Klasifikasi  string         `json:"klasifikasi" db:"system_classification"`
	Lokasi       string         `json:"lokasi" db:"asset_location"`
	Status       string         `json:"status" db:"asset_status"`
	AssetType    string         `json:"asset_type" db:"asset_type"`
	CreatedBy    string         `json:"created_by" db:"created_by"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
	UpdatedBy    sql.NullString `json:"updated_by" db:"updated_by"`
	UpdatedAt    sql.NullTime   `json:"updated_at" db:"updated_at"`
	DeletedBy    sql.NullString `json:"deleted_by" db:"deleted_by"`
	DeletedAt    sql.NullTime   `json:"deleted_at" db:"deleted_at"`
	Pic          []Pic          `json:"pic"`
	Image        ImageUpdate    `json:"image"`
}

type ImageUpdate struct {
	Added   []string `json:"added"`
	Deleted []string `json:"deleted"`
	Images  []string `json:"images"`
}
