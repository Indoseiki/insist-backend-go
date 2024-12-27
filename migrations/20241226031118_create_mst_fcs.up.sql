CREATE TABLE
    mst_fcs (
        id SERIAL PRIMARY KEY,
        id_building INT REFERENCES mst_buildings (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        code VARCHAR NOT NULL,
        description VARCHAR NOT NULL,
        remarks VARCHAR,
        id_createdby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_updatedby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ
    );