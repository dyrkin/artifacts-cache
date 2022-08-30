create table partition
(
    id integer primary key,
    uuid text not null,
    time text not null
);

create table file
(
    id integer primary key,
    partition_id integer not null,
    subset text not null,
    path text not null,
    "offset" integer not null,
    size integer not null,
    foreign key (partition_id) references partition (id)
);

create index partition___uuid on partition (uuid);
create index file___subset__path on file (subset, path);
create index file___partition_id on file (partition_id);