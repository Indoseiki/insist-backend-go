CREATE TABLE
    mst_item_processes (
        id SERIAL PRIMARY KEY,
        id_item_category INT REFERENCES mst_item_categories (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        code VARCHAR NOT NULL,
        description VARCHAR NOT NULL,
        remarks VARCHAR,
        id_createdby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_updatedby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ
    );