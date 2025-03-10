-- +goose Up
-- +goose StatementBegin
-- Table: form_ms

-- DROP TABLE IF EXISTS form_ms;

CREATE TYPE public.status AS ENUM (
    'Draft',
    'Published'
);


ALTER TYPE public.status OWNER TO postgres;

CREATE TABLE form_ms (
    form_id bigint NOT NULL,
    form_uuid character varying(128) NOT NULL,
    document_id bigint NOT NULL,
    user_id bigint NOT NULL,
    project_id bigint,
    form_number character varying(100) NOT NULL,
    form_ticket character varying(100),
    form_status status NOT NULL,
    form_data jsonb NOT NULL,
    is_approve boolean,
    reason character varying(5000),
    created_by character varying(100) NOT NULL,
    created_at timestamp(0) without time zone DEFAULT now() NOT NULL,
    updated_by character varying(100) DEFAULT ''::character varying,
    updated_at timestamp(0) without time zone,
    deleted_by character varying(100) DEFAULT ''::character varying,
    deleted_at timestamp(0) without time zone,
    image_path character varying
);


ALTER TABLE form_ms OWNER TO postgres;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS form_ms;
DROP TYPE IF EXISTS public.status;
-- +goose StatementEnd