CREATE TABLE
    mst_approvals (
        id SERIAL PRIMARY KEY,
        id_menu INT REFERENCES mst_menus (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        status VARCHAR NOT NULL,
        action VARCHAR NOT NULL,
        count INT NOT NULL,
        level INT NOT NULL,
        id_createdby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_updatedby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ
    );