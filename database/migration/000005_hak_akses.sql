-- +goose Up
-- +goose StatementBegin
-- Table: hak_akses

-- DROP TABLE IF EXISTS hak_akses;

CREATE TABLE hak_akses (
    form_id bigint NOT NULL,
    ha_uuid character varying(128) NOT NULL,
    nama_pengguna character varying(128) NOT NULL,
    ruang_lingkup character varying(128) NOT NULL,
    jangka_waktu character varying(128) NOT NULL,
    created_by character varying(128) NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_by character varying(128) DEFAULT ''::character varying,
    updated_at timestamp without time zone,
    deleted_by character varying(128) DEFAULT ''::character varying,
    deleted_at timestamp without time zone
);


ALTER TABLE hak_akses OWNER TO postgres;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS hak_akses;
-- +goose StatementEnd