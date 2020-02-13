CREATE TABLE jet_trees (
    pulse_number bigint primary key,
    jet_tree bytea not null
);
---- create above / drop below ----
DROP TABLE jet_trees;
