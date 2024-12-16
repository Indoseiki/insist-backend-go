CREATE TABLE
    mst_role_permissions (
        id_role INT NOT NULL,
        id_menu INT NOT NULL,
        is_create BOOLEAN NOT NULL,
        is_update BOOLEAN NOT NULL,
        is_delete BOOLEAN NOT NULL
    );