CREATE TABLE
    mst_fcs_buildings (
        id_fcs INT REFERENCES mst_fcs (id) ON UPDATE CASCADE ON DELETE RESTRICT,
        id_building INT REFERENCES mst_buildings (id) ON UPDATE CASCADE ON DELETE RESTRICT
    );