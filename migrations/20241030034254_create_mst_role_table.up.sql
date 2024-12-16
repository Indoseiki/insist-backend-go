CREATE TABLE
    mst_roles (
        id SERIAL PRIMARY KEY,
        name VARCHAR NOT NULL,
        id_createdby INT,
        id_updatedby INT,
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ
    );