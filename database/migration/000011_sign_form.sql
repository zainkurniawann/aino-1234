-- +goose Up
-- +goose StatementBegin
-- Table: sign_form

-- DROP TABLE IF EXISTS sign_form;

CREATE TYPE role AS ENUM (
    'Pemohon',
    'Atasan Pemohon',
    'Penerima',
    'Atasan Penerima',
    'Diketahui oleh',
    'Disusun oleh',
    'Disahkan oleh',
    'Direview oleh',
    'Pengaju',
    'Atasan Pengaju',
    'Pihak Pertama',
    'Pihak Kedua',
    'Pelaksana',
    'Mengetahui1',
    'Mengetahui2'
);


ALTER TYPE role OWNER TO postgres;

CREATE TABLE sign_form (
    user_id bigint NOT NULL,
    sign_uuid character varying(128) NOT NULL,
    form_id bigint NOT NULL,
    name character varying(128) NOT NULL,
    "position" character varying(128) NOT NULL,
    role_sign role NOT NULL,
    is_sign boolean DEFAULT false,
    created_by character varying(100) NOT NULL,
    created_at timestamp(0) without time zone DEFAULT now() NOT NULL,
    updated_by character varying(100) DEFAULT ''::character varying,
    updated_at timestamp(0) without time zone,
    deleted_by character varying(100) DEFAULT ''::character varying,
    deleted_at timestamp(0) without time zone,
    sign_date timestamp without time zone,
    sign_img character varying,
    is_guest boolean DEFAULT false
);


ALTER TABLE sign_form OWNER TO postgres;


ALTER TABLE qna OWNER TO postgres;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sign_form;
DROP TYPE IF EXISTS role;
-- +goose StatementEnd