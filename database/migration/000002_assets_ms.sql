-- +goose Up
-- +goose StatementBegin
-- Table: assets_ms

-- DROP TABLE IF EXISTS assets_ms;
CREATE TABLE assets_ms (
    asset_id bigint NOT NULL,
    asset_uuid character varying(128) NOT NULL,
    asset_code character varying(128) NOT NULL,
    asset_name character varying(128) NOT NULL,
    serial_number character varying(128) NOT NULL,
    asset_specification character varying NOT NULL,
    procurement_date date NOT NULL,
    price character varying(128) NOT NULL,
    asset_description character varying NOT NULL,
    system_classification character varying(128) NOT NULL,
    asset_location character varying(128) NOT NULL,
    asset_status character varying(128) NOT NULL,
    created_by character varying(128) NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_by character varying(128),
    updated_at timestamp without time zone,
    deleted_by character varying(128),
    deleted_at timestamp without time zone,
    merk character varying,
    asset_img jsonb,
    asset_type character varying NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS assets_ms;
-- +goose StatementEnd