CREATE TABLE
    mst_employees (
        number VARCHAR NOT NULL,
        name VARCHAR NOT NULL,
        division VARCHAR,
        department VARCHAR,
        position VARCHAR,
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        service VARCHAR,
        education VARCHAR,
        birthday DATE,
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ
    );