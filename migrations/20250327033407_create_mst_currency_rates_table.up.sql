CREATE TABLE
    mst_currency_rates (
        id SERIAL PRIMARY KEY,
        id_from_currency INT REFERENCES mst_currencies (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_to_currency INT REFERENCES mst_currencies (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        buy_rate FLOAT NOT NULL,
        sell_rate FLOAT NOT NULL,
        effective_date DATE NOT NULL,
        id_createdby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_updatedby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ
    );