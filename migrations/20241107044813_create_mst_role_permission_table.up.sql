CREATE TABLE
    mst_role_permissions (
        id_role INT REFERENCES mst_roles (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_menu INT REFERENCES mst_menus (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        is_create BOOLEAN NOT NULL,
        is_update BOOLEAN NOT NULL,
        is_delete BOOLEAN NOT NULL
    );