ALTER TABLE files RENAME COLUMN file_chroma to file_chroma_old;
ALTER TABLE files ADD COLUMN file_chroma integer DEFAULT -1;
UPDATE files SET file_chroma = file_chroma_old;
ALTER TABLE files DROP COLUMN file_chroma_old;
