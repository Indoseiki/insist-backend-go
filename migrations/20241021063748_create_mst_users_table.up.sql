CREATE TABLE
    mst_users (
        id SERIAL PRIMARY KEY,
        id_dept INT NOT NULL,
        name VARCHAR NOT NULL,
        email VARCHAR UNIQUE,
        username VARCHAR NOT NULL UNIQUE,
        password VARCHAR NOT NULL,
        refresh_token TEXT,
        otp_url VARCHAR,
        otp_key VARCHAR,
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        is_two_fa BOOLEAN NOT NULL DEFAULT FALSE,
        id_createdby INT,
        id_updatedby INT,
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ
    );