CREATE TABLE
    password_resets (
        id SERIAL PRIMARY KEY,
        id_user INT NOT NULL,
        token VARCHAR NOT NULL UNIQUE,
        is_used BOOLEAN DEFAULT FALSE,
        expired_at TIMESTAMP NOT NULL,
        id_createdby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_updatedby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ
    );