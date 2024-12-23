CREATE TABLE
    activity_logs (
        id SERIAL PRIMARY KEY,
        id_user INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        ip_address VARCHAR NOT NULL,
        action VARCHAR NOT NULL,
        is_success BOOLEAN NOT NULL DEFAULT FALSE,
        message VARCHAR,
        user_agent VARCHAR,
        os VARCHAR,
        created_at TIMESTAMPTZ
    );