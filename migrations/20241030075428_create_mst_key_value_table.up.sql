CREATE TABLE
    mst_key_values (
        id SERIAL PRIMARY KEY,
        key VARCHAR NOT NULL,
        value VARCHAR NOT NULL,
        remarks VARCHAR,
        id_createdby INT,
        id_updatedby INT,
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ
    );