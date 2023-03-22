ALTER TABLE books ADD COLUMN IF NOT EXISTS library_id SERIAL;

ALTER TABLE books 
    ADD CONSTRAINT fk_library_id FOREIGN KEY (library_id) REFERENCES library (id);