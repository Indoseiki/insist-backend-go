CREATE TABLE
    mst_menus (
        id SERIAL PRIMARY KEY,
        label VARCHAR NOT NULL,
        path VARCHAR NOT NULL,
        id_parent INT,
        icon VARCHAR NOT NULL,
        sort INT NOT NULL,
        id_createdby INT,
        id_updatedby INT,
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ
    );