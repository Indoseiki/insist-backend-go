CREATE TABLE
    mst_approvals (
        id SERIAL PRIMARY KEY,
        id_menu INT NOT NULL,
        status VARCHAR NOT NULL,
        action VARCHAR NOT NULL,
        count INT NOT NULL,
        level INT NOT NULL,
        id_createdby INT,
        id_updatedby INT,
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ
    );