-- +goose Up
-- +goose StatementBegin
-- Table: project_ms

-- DROP TABLE IF EXISTS project_ms;

CREATE SEQUENCE public.project_order_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.project_order_seq OWNER TO postgres;

CREATE TABLE project_ms (
    project_id bigint NOT NULL,
    project_uuid character varying(128) NOT NULL,
    product_id bigint NOT NULL,
    project_order integer DEFAULT nextval('project_order_seq'::regclass),
    project_name character varying(128) NOT NULL,
    project_code character varying(20) NOT NULL,
    project_manager character varying(128) NOT NULL,
    created_by character varying(100) NOT NULL,
    created_at timestamp(0) without time zone DEFAULT now() NOT NULL,
    updated_by character varying(100) DEFAULT ''::character varying,
    updated_at timestamp(0) without time zone,
    deleted_by character varying(100) DEFAULT ''::character varying,
    deleted_at timestamp(0) without time zone
);


ALTER TABLE project_ms OWNER TO postgres;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS project_ms;
DROP SEQUENCE IF EXISTS public.project_order_seq;
-- +goose StatementEnd