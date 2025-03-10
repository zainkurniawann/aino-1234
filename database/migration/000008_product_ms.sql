-- +goose Up
-- +goose StatementBegin
-- Table: product_ms

-- DROP TABLE IF EXISTS product_ms;

CREATE SEQUENCE public.product_order_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.product_order_seq OWNER TO postgres;


CREATE TABLE product_ms (
    product_id bigint NOT NULL,
    product_uuid character varying(128) NOT NULL,
    product_order integer DEFAULT nextval('product_order_seq'::regclass),
    product_name character varying(128) NOT NULL,
    product_owner character varying(128) NOT NULL,
    created_by character varying(100) NOT NULL,
    created_at timestamp(0) without time zone DEFAULT now() NOT NULL,
    updated_by character varying(100) DEFAULT ''::character varying,
    updated_at timestamp(0) without time zone,
    deleted_by character varying(100) DEFAULT ''::character varying,
    deleted_at timestamp(0) without time zone
);


ALTER TABLE product_ms OWNER TO postgres;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product_ms;
DROP SEQUENCE IF EXISTS public.product_order_seq;
-- +goose StatementEnd