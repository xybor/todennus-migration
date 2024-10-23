DROP TABLE oauth2_consents;

ALTER TABLE users ADD COLUMN allowed_scope VARCHAR;

UPDATE users
SET allowed_scope='*';
