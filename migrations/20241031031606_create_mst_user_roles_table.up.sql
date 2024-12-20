CREATE TABLE
    mst_user_roles (
        id_user INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_role INT REFERENCES mst_roles (id) ON UPDATE CASCADE ON DELETE RESTRICT
    );