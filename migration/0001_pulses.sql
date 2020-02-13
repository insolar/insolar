CREATE TABLE key_value (
    k varchar(256) primary key,
    v bytea not null
);

create table pulses(
  pulse_number bigint primary key,
  prev_pn bigint not null,
  next_pn bigint not null,
  tstamp bigint not null,
  epoch bigint not null,
  origin_id bytea not null,
  entropy bytea not null
);

create table pulse_signs(
  pulse_number bigint not null references pulses(pulse_number) on delete cascade, -- on delete part is only for tests!
  chosen_public_key text not null,
  entropy bytea not null,
  signature bytea not null,
  primary key(pulse_number, chosen_public_key)
);
/*******************************************************************************
 * Copyright 2020 Insolar Network Ltd.
 * All rights reserved.
 * This material is licensed under the Insolar License version 1.0,
 * available at https://github.com/insolar/insolar/blob/master/LICENSE.md.
 ******************************************************************************/

---- create above / drop below ----
DROP TABLE key_value;
DROP TABLE pulse_signs;
DROP TABLE pulses;
