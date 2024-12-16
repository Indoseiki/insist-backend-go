CREATE TABLE
    mst_reasons (
        id SERIAL PRIMARY KEY,
        id_menu INT NOT NULL,
        key VARCHAR NOT NULL,
        code VARCHAR,
        description VARCHAR NOT NULL,
        remarks VARCHAR,
        id_createdby INT,
        id_updatedby INT,
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ
    );