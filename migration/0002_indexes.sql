create table last_known_pulse_for_indexes
(
    object_id    bytea  not null,
    pulse_number bigint not null
);

create unique index last_known_pulse_for_indexes_unique on last_known_pulse_for_indexes (object_id, pulse_number);

create table indexes
(
    object_id             bytea   not null,
    pulse_number          bigint  not null,

    lifeline_last_used    bigint  not null,
    pending_records       bytea[] null,

    -- lifeline
    latest_state          bytea   null,
    state_id              integer not null,
    parent                bytea   not null,
    latest_request        bytea   null,
    earliest_open_request bigint  null,
    open_requests_count   bigint  not null,

    primary key (object_id, pulse_number)
);
