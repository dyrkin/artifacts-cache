create table partition
(
    id bigserial primary key,
    uuid varchar(36) not null,
    time timestamp not null
);

create table file
(
    id bigserial primary key,
    partition_id bigint not null,
    subset text not null,
    name text not null,
    "offset" bigint not null,
    size bigint not null,
    foreign key (partition_id) references partition (id)
);

create index on partition (uuid);
create index on file (subset, name);
create index on file (partition_id);