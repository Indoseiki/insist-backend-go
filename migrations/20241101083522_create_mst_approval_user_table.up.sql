CREATE TABLE
    mst_approval_users (
        id_approval INT REFERENCES mst_approvals (id) ON UPDATE CASCADE ON DELETE CASCADE,
        id_user INT REFERENCES mst_users (id) ON UPDATE CASCADE ON DELETE CASCADE
    );