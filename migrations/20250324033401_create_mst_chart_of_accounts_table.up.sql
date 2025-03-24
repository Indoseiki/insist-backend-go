CREATE TABLE
    mst_chart_of_accounts (
        id SERIAL PRIMARY KEY,
        account INT NOT NULL,
        description VARCHAR NOT NULL,
        type VARCHAR NOT NULL,
        class VARCHAR NOT NULL,
        exchange_rate_type VARCHAR NOT NULL,
        id_createdby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_updatedby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ
    );