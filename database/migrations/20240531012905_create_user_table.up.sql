use membership_system;

create table if not exists users
(
    id int auto_increment not null,
    name varchar(50) not null comment '使用者姓名',
    email varchar(128) not null comment '使用者信箱',
    password varchar(80) not null comment '使用者密碼',
    role tinyint unsigned not null comment '角色 0:admin/1:moderator/2:member',
    created_at int unsigned not null default (unix_timestamp()),
    updated_at int unsigned not null default (unix_timestamp()),
    primary key(id),
    unique (email)
) ENGINE=InnoDB default CHARSET=utf8mb4 collate=utf8mb4_unicode_ci comment '使用者';
