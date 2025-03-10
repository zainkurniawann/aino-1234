-- +goose Up
-- +goose StatementBegin
-- Table: hak_akses_info

-- DROP TABLE IF EXISTS hak_akses_info;

CREATE TABLE hak_akses_info (
    form_id bigint NOT NULL,
    info_uuid character varying(128) NOT NULL,
    host character varying(128),
    name character varying(128) NOT NULL,
    instansi character varying(128) NOT NULL,
    "position" character varying(128) NOT NULL,
    username character varying(128) NOT NULL,
    password character varying(128) NOT NULL,
    scope character varying(128) NOT NULL,
    type character varying(128),
    matched boolean,
    description text,
    created_by character varying(100) NOT NULL,
    created_at timestamp(0) without time zone DEFAULT now() NOT NULL,
    updated_by character varying(100) DEFAULT ''::character varying,
    updated_at timestamp(0) without time zone,
    deleted_by character varying(100) DEFAULT ''::character varying,
    deleted_at timestamp(0) without time zone
);


ALTER TABLE hak_akses_info OWNER TO postgres;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS hak_akses_info;
-- +goose StatementEnd