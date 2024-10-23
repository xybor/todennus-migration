CREATE TABLE oauth2_consents (
    user_id BIGINT REFERENCES users(id),
    client_id BIGINT REFERENCES oauth2_clients(id),
    scope VARCHAR,
    created_at timestamp,
    updated_at timestamp,
    expires_at timestamp,
    PRIMARY KEY (user_id, client_id)
);

ALTER TABLE users DROP COLUMN allowed_scope;
