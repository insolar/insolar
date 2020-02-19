CREATE TABLE jet_trees (
    pulse_number bigint primary key,
    jet_tree bytea not null
);
/*******************************************************************************
 * Copyright 2020 Insolar Network Ltd.
 * All rights reserved.
 * This material is licensed under the Insolar License version 1.0,
 * available at https://github.com/insolar/insolar/blob/master/LICENSE.md.
 ******************************************************************************/

---- create above / drop below ----
DROP TABLE jet_trees;
