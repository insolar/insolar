create table records (
    pulse_number bigint not null,
    position bigint not null,
    record_id bytea not null,
    object_id bytea not null,
    jet_id bytea not null,
    signature bytea not null,
    polymorph int not null,
    virtual bytea not null, -- serialized to protobuf
    primary key(pulse_number, position)
);

create index records_recod_id_idx on records(record_id);

create table records_last_position (
    pulse_number bigint primary key,
    position bigint not null
);

/*******************************************************************************
 * Copyright 2020 Insolar Network Ltd.
 * All rights reserved.
 * This material is licensed under the Insolar License version 1.0,
 * available at https://github.com/insolar/insolar/blob/master/LICENSE.md.
 ******************************************************************************/

---- create above / drop below ----
DROP TABLE records;
DROP TABLE records_last_position;
