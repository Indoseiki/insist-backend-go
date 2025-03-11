CREATE TABLE
    approval_histories (
        id SERIAL PRIMARY KEY,
        id_approval INT REFERENCES mst_approvals (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        ref_table VARCHAR NOT NULL,
        ref_id INT NOT NULL,
        key VARCHAR NOT NULL,
        message VARCHAR NOT NULL,
        id_createdby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        created_at TIMESTAMPTZ
    );