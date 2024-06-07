use membership_system;

create table if not exists oauth2_client
(
    id varchar(255) not null,
    client_name varchar(50) not null,
    client_secret varchar(80) not null,
    scope varchar(255) not null,
    grant_types json not null default ('[]'),
    audience json not null default ('[]'),
    post_logout_redirect_uris json not null default ('[]'),
    backchannel_logout_uri text not null default (''),
    redirect_uris json not null default ('[]'),
    created_at int unsigned not null default (unix_timestamp()),
    updated_at int unsigned not null default (unix_timestamp()),
    primary key(id)
) ENGINE=InnoDB default CHARSET=utf8mb4 collate=utf8mb4_unicode_ci;
