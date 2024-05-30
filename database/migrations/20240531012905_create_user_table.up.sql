use membership_system;

create table if not exists users
(
    id int auto_increment not null,
    name varchar(50) not null comment '使用者姓名',
    created_at timestamp default current_timestamp() not null,
    updated_at timestamp default current_timestamp() on update current_timestamp() not null,
    primary key(id),
    index (name)
) ENGINE=InnoDB default CHARSET=utf8mb4 collate=utf8mb4_unicode_ci comment '使用者';
