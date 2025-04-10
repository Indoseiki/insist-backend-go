CREATE TABLE
    mst_billing_terms (
        id SERIAL PRIMARY KEY,
        code VARCHAR NOT NULL,
        description VARCHAR NOT NULL,
        due_days INT,
        discount_days INT,
        is_cash_only BOOLEAN,
        prox_due_day INT,
        prox_discount_day INT,
        prox_months_forward INT,
        prox_discount_months_forward INT,
        cutoff_day INT,
        discount_percent DECIMAL(5, 3),
        holiday_offset_method VARCHAR,
        is_advanced_terms BOOLEAN,
        prox_code INT,
        id_createdby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_updatedby INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ
    );