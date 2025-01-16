CREATE TABLE
    mst_locations (
        id SERIAL PRIMARY KEY,
        id_warehouse INT REFERENCES mst_warehouses (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        location VARCHAR NOT NULL,
        remarks VARCHAR,
        id_createdby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_updatedby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ
    );