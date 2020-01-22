create type local_ref as (pulse_and_scope bigint, hash bytea);

create table records (
    id local_ref not null,
    position bigint not null,
    object_id local_ref not null,
    jet_id local_ref not null,
    signature bytea not null,
    polymorph int not null,
    virtual bytea not null, -- serialized to protobuf
    primary key(id, position)
);

create table records_last_position (
    pulse_number bigint primary key,
    position bigint not null
);

---- create above / drop below ----
DROP TABLE records;
DROP TABLE records_last_position;
