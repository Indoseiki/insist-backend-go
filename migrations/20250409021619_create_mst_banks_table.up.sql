CREATE TABLE
    mst_banks (
        id SERIAL PRIMARY KEY,
        code VARCHAR NOT NULL,
        name VARCHAR NOT NULL,
        account_num VARCHAR NOT NULL,
        id_account INT REFERENCES mst_chart_of_accounts (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_currency INT REFERENCES mst_currencies (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        bic VARCHAR,
        country VARCHAR NOT NULL,
        state VARCHAR NOT NULL,
        city VARCHAR NOT NULL,
        address VARCHAR NOT NULL,
        zip_code VARCHAR NOT NULL,
        remarks VARCHAR,
        id_createdby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_updatedby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ
    );