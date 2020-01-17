create table pulses(
  pulse_number bigint primary key,
  prev_pn bigint not null,
  next_pn bigint not null,
  tstamp timestamp not null,
  epoch bigint not null,
  origin_id bytea not null,
  entropy bytea not null
);

create table pulse_signs(
  pulse_number bigint not null references pulses(pulse_number),
  chosen_public_key text not null,
  entropy bytea not null,
  signature bytea not null,
  primary key(pulse_number, chosen_public_key)
);
---- create above / drop below ----
DROP TABLE pulse_signs;
DROP TABLE pulses;
