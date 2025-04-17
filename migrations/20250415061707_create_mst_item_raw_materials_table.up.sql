CREATE TABLE
    mst_item_raw_materials (
        id SERIAL PRIMARY KEY,
        id_item INT REFERENCES mst_items (id) ON UPDATE CASCADE ON DELETE CASCADE,
        id_item_product_type INT REFERENCES mst_item_product_types (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_item_group_type INT REFERENCES mst_item_group_types (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_item_process INT REFERENCES mst_item_processes (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_item_surface INT REFERENCES mst_item_surfaces (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_item_source INT REFERENCES mst_item_sources (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        diameter_size VARCHAR NOT NULL,
        length_size VARCHAR NOT NULL,
        inner_diameter_size VARCHAR NOT NULL,
        id_createdby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_updatedby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ
    );