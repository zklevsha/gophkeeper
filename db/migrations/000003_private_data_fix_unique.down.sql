ALTER TABLE private_data
DROP CONSTRAINT private_data_userid_name_unique;

ALTER TABLE private_data
ADD CONSTRAINT private_data_id_name_key UNIQUE (id, name);