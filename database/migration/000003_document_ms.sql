-- +goose Up
-- +goose StatementBegin
-- Table: document_ms

-- DROP TABLE IF EXISTS document_ms;

CREATE SEQUENCE public.document_order_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.document_order_seq OWNER TO postgres;

CREATE TABLE document_ms (
    document_id bigint NOT NULL,
    document_uuid character varying(128) NOT NULL,
    document_order integer DEFAULT nextval('document_order_seq'::regclass),
    document_code character varying(20) NOT NULL,
    document_name character varying(100) NOT NULL,
    created_by character varying(100) NOT NULL,
    created_at timestamp(0) without time zone DEFAULT now() NOT NULL,
    updated_by character varying(100) DEFAULT ''::character varying,
    updated_at timestamp(0) without time zone,
    deleted_by character varying(100) DEFAULT ''::character varying,
    deleted_at timestamp(0) without time zone
);

ALTER TABLE document_ms OWNER TO postgres;

-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO public.document_ms (
    document_id, document_uuid, document_order, document_code, document_name, 
    created_by, created_at, updated_by, updated_at, deleted_by, deleted_at
) VALUES
    (1720586841296101, '1f5dd571-d452-4c69-bb9e-75ba916df12b', 1, 'DA', 'Dampak Analisa', 
    'Adipiscingelite', '2024-07-10 11:26:53', '', NULL, '', NULL),

    (1723612477546683, '63c78d3e-2896-4e58-bc19-746db23179b1', 4, 'HA', 'Hak Akses', 
    'Adipiscingelite', '2024-08-14 11:23:03', '', NULL, '', NULL),

    (1720624585356006, 'f97b1dbe-0366-4ec5-aeff-bdf1f288e347', 3, 'BA', 'Berita Acara', 
    'Adipiscingelite', '2024-07-10 21:17:09', 'Adipiscingelite', '2024-07-16 09:04:20', '', NULL),

    (1720598392013780, '26b14f9e-c739-461c-b4c1-76b71afde5e3', 2, 'ITCM', 'IT Change Management', 
    'Adipiscingelite', '2024-07-10 13:51:38', 'Adipiscingelite', '2025-03-04 15:22:49', '', NULL);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS document_ms;
DROP SEQUENCE IF EXISTS document_order_seq;
-- +goose StatementEnd