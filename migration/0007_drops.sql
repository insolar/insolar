create table drops(
    pulse_number bigint not null,
    id_prefix bytea not null,
    jet_id bytea not null,
    split_threshold_exceeded int not null,
    split boolean not null,

    primary key (pulse_number, id_prefix)
);
---- create above / drop below ----
DROP TABLE drops;