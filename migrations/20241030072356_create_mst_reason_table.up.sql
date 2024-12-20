CREATE TABLE
    mst_reasons (
        id SERIAL PRIMARY KEY,
        id_menu INT REFERENCES mst_menus (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        key VARCHAR NOT NULL,
        code VARCHAR,
        description VARCHAR NOT NULL,
        remarks VARCHAR,
        id_createdby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_updatedby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ
    );