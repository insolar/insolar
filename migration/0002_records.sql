create type local_ref as (pulse_and_scope bigint, hash bytea);

create table records (
    id local_ref primary key,
    position bigint not null,
    object_id local_ref not null,
    jet_id local_ref not null,
    signature bytea not null,
    polymorph int not null,
    virtual bytea not null -- serialized to protobuf
);

---- create above / drop below ----
DROP TABLE records;
