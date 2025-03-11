CREATE TABLE
    mst_machine_statuses (
        id SERIAL PRIMARY KEY,
        id_machine INT REFERENCES mst_machines (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_reason INT REFERENCES mst_reasons (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        remarks VARCHAR,
        id_createdby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_updatedby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ
    );