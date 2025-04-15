-- Fungsi untuk memeriksa apakah item digunakan sebagai foreign key di tabel lain
CREATE OR REPLACE FUNCTION "public"."check_mst_item_used_as_fk"("item_id" int4)
  RETURNS "pg_catalog"."bool" AS $BODY$
DECLARE
    ref_table TEXT;
    ref_column TEXT;
    fk_count INTEGER;
BEGIN
    -- Melakukan iterasi pada tabel-tabel yang memiliki foreign key ke mst_items.id
    FOR ref_table, ref_column IN (
        SELECT 
            kcu.table_name, 
            kcu.column_name
        FROM 
            information_schema.table_constraints AS tc
        JOIN 
            information_schema.key_column_usage AS kcu
            ON tc.constraint_name = kcu.constraint_name
            AND tc.constraint_schema = kcu.constraint_schema
        JOIN 
            information_schema.constraint_column_usage AS ccu
            ON ccu.constraint_name = tc.constraint_name
            AND ccu.constraint_schema = tc.constraint_schema
        WHERE 
            tc.constraint_type = 'FOREIGN KEY'
            AND ccu.table_name = 'mst_items'  -- Menyaring hanya tabel yang merujuk ke mst_items
            AND ccu.column_name = 'id'       -- Menyaring hanya kolom id yang dijadikan foreign key
    ) LOOP
        -- Mengeksekusi query untuk menghitung jumlah penggunaan item_id di kolom foreign key
        EXECUTE format(
            'SELECT COUNT(*) FROM %I WHERE %I = $1',
            ref_table, ref_column
        )
        INTO fk_count
        USING item_id;

        -- Jika ditemukan lebih dari 0, berarti item_id digunakan sebagai foreign key
        IF fk_count > 0 THEN
            RETURN TRUE;
        END IF;
    END LOOP;

    -- Jika tidak ditemukan referensi, kembalikan FALSE
    RETURN FALSE;
END;
$BODY$
  LANGUAGE plpgsql VOLATILE
  COST 100;


-- Fungsi Trigger untuk mencegah update pada mst_items jika item digunakan sebagai foreign key
CREATE OR REPLACE FUNCTION prevent_update_if_used_as_fk()
RETURNS TRIGGER AS $$
BEGIN
    -- Mengecek apakah item yang akan diupdate digunakan sebagai foreign key
    IF check_mst_item_used_as_fk(OLD.id) THEN
        -- Jika digunakan, maka raise exception dan batalkan update
        RAISE EXCEPTION 'Tidak dapat memperbarui data mst_item dengan ID % karena sedang digunakan sebagai kunci asing di tabel lain.', OLD.id;
    END IF;
    -- Jika tidak digunakan sebagai foreign key, lanjutkan proses update
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- Trigger untuk memanggil fungsi prevent_update_if_used_as_fk sebelum update pada mst_items
CREATE TRIGGER prevent_mst_item_update
BEFORE UPDATE ON mst_items
FOR EACH ROW
EXECUTE FUNCTION prevent_update_if_used_as_fk();
