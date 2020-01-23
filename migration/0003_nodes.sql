CREATE TABLE nodes (
    pulse_number bigint not null,
    node_num int not null,
    polymorph bigint not null,
    node_id bytea not null, -- consider moving this row to a separate table
    role bigint not null,
    primary key(pulse_number, node_num)
);
---- create above / drop below ----
DROP TABLE nodes;
