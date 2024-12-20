CREATE TABLE
    mst_menus (
        id SERIAL PRIMARY KEY,
        label VARCHAR NOT NULL,
        path VARCHAR NOT NULL,
        id_parent INT,
        icon VARCHAR NOT NULL,
        sort INT NOT NULL,
        id_createdby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_updatedby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ
    );