CREATE TABLE
    mst_materials (
        id SERIAL PRIMARY KEY,
        code VARCHAR REFERENCES mst_items (code) ON UPDATE RESTRICT ON DELETE RESTRICT,
        id_createdby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_updatedby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ
    );