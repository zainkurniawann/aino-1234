-- +goose Up
-- +goose StatementBegin
-- Table: qna

-- DROP TABLE IF EXISTS qna;

CREATE TABLE qna (
    qna_id bigint NOT NULL,
    qna_uuid uuid NOT NULL,
    question text NOT NULL,
    answer text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    created_by character varying NOT NULL,
    updated_at timestamp without time zone,
    updated_by character varying,
    deleted_at timestamp without time zone,
    deleted_by character varying
);


ALTER TABLE qna OWNER TO postgres;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS qna;
-- +goose StatementEnd