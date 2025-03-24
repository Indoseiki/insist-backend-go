CREATE TABLE
    mst_tax_codes (
        id SERIAL PRIMARY KEY,
        name VARCHAR NOT NULL,
        description VARCHAR NOT NULL,
        type VARCHAR NOT NULL,
        rate FLOAT NOT NULL,
        include VARCHAR NOT NULL,
        id_account_ar INT REFERENCES mst_chart_of_accounts (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_account_ar_process INT REFERENCES mst_chart_of_accounts (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_account_ap INT REFERENCES mst_chart_of_accounts (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_createdby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_updatedby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ
    );