CREATE TABLE
    mst_role_menus (
        id_role INT REFERENCES mst_roles (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_menu INT REFERENCES mst_menus (id) ON UPDATE CASCADE ON DELETE RESTRICT
    );