CREATE TABLE
    mst_items (
        id SERIAL PRIMARY KEY,
        id_item_category INT REFERENCES mst_item_categories (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_uom INT REFERENCES mst_uoms (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        code VARCHAR NOT NULL UNIQUE,
        description VARCHAR NOT NULL,
        infor_code VARCHAR NOT NULL,
        infor_description VARCHAR NOT NULL,
        remarks VARCHAR,
        id_createdby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_updatedby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ
    );