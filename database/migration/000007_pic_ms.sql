-- +goose Up
-- +goose StatementBegin
-- Table: pic_ms

-- DROP TABLE IF EXISTS pic_ms;

CREATE TABLE pic_ms (
    asset_id bigint NOT NULL,
    pic_uuid character varying(128) NOT NULL,
    pic_name character varying(128) NOT NULL,
    pic_description character varying(128) NOT NULL,
    created_by character varying NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_by character varying,
    updated_at timestamp without time zone,
    deleted_by character varying,
    deleted_at timestamp without time zone,
    start_at timestamp without time zone NOT NULL,
    ended_at timestamp without time zone
);


ALTER TABLE pic_ms OWNER TO postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pic_ms;
-- +goose StatementEnd