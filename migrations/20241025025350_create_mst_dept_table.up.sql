CREATE TABLE
    mst_depts (
        id SERIAL PRIMARY KEY,
        code VARCHAR NOT NULL,
        description VARCHAR NOT NULL,
        remarks VARCHAR NOT NULL,
        id_createdby INT,
        id_updatedby INT,
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ
    );