CREATE TABLE nodes (
    pulse_number bigint not null,
    node_num int not null,
    polymorph bigint not null,
    node_id bytea not null, -- consider moving this row to a separate table
    role bigint not null,
    primary key(pulse_number, node_num)
);
/*******************************************************************************
 * Copyright 2020 Insolar Network Ltd.
 * All rights reserved.
 * This material is licensed under the Insolar License version 1.0,
 * available at https://github.com/insolar/insolar/blob/master/LICENSE.md.
 ******************************************************************************/

---- create above / drop below ----
DROP TABLE nodes;
