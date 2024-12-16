CREATE TABLE
    mst_approval_users (
        id_approval INT NOT NULL,
        id_user INT NOT NULL,
        CONSTRAINT fk_mst_approval FOREIGN KEY (id_approval) REFERENCES mst_approval (id) ON DELETE CASCADE ON UPDATE CASCADE,
        CONSTRAINT fk_mst_users FOREIGN KEY (id_user) REFERENCES users (id) ON DELETE RESTRICT ON UPDATE CASCADE
    );