-- auto-generated definition
create table "ShortUrl"
(
    id     serial
        constraint shorturl_pk
            primary key,
    domain varchar(256)  not null,
    path   varchar(2048) not null
);

alter table "ShortUrl"
    owner to postgres;

create unique index shorturl_id_uindex
    on "ShortUrl" (id);

create index domain_path
    on "ShortUrl" (domain, path);

