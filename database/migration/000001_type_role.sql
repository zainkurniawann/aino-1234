-- -- +goose Up
-- -- +goose StatementBegin
-- -- Table: assets_ms

-- -- DROP TABLE IF EXISTS assets_ms;

-- CREATE TYPE role AS ENUM (
--     'Pemohon',
--     'Atasan Pemohon',
--     'Penerima',
--     'Atasan Penerima',
--     'Diketahui oleh',
--     'Disusun oleh',
--     'Disahkan oleh',
--     'Direview oleh',
--     'Pengaju',
--     'Atasan Pengaju',
--     'Pihak Pertama',
--     'Pihak Kedua'
-- );


-- ALTER TYPE role OWNER TO postgres;

-- -- +goose StatementEnd

-- -- +goose Down
-- -- +goose StatementBegin
-- DROP type IF EXISTS role;
-- -- +goose StatementEnd