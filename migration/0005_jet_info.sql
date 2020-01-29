CREATE TABLE jets_info (
    pulse_number bigint primary key,
    info bytea not null
);

CREATE TABLE key_value (
    k varchar(256) primary key,
    v bytea not null
);

---- create above / drop below ----
DROP TABLE jets_info;
DROP TABLE key_value;
