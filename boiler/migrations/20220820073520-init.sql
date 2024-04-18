-- +migrate Up
create table resources (
    -- resource id
    id BINARY(16) primary key,
    -- status
    status enum('STATUS_ACTIVE', 'STATUS_DELETED') not null,
    -- date
    created_at datetime not null,
    updated_at datetime not null,
    deleted_at datetime
);

create table properties (
    -- pk
    id int auto_increment primary key,
    -- kf
    resource_id binary(16) not null,
    -- resource properties
    first_name text not null,
    last_name text not null,
    -- constraint
    constraint fk_resource_id foreign key (resource_id) references resources (id),
    constraint unique_resource_id unique (resource_id)
);

-- +migrate Down
drop table properties;

drop table resources;