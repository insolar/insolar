create table drops(
    pulse_number bigint not null,
    id_prefix bytea not null,
    jet_id bytea not null,
    split_threshold_exceeded int not null,
    split boolean not null,

    primary key (pulse_number, id_prefix)
);
/*******************************************************************************
 * Copyright 2020 Insolar Network Ltd.
 * All rights reserved.
 * This material is licensed under the Insolar License version 1.0,
 * available at https://github.com/insolar/insolar/blob/master/LICENSE.md.
 ******************************************************************************/

---- create above / drop below ----
DROP TABLE drops;
