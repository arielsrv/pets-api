create schema pets;

create table pets.apps_types
(
    id   bigint
        primary key,
    name varchar(255) not null,
    constraint name
        unique (name)
);

INSERT INTO pets.apps_types (id, name) VALUES (1, 'backend');
INSERT INTO pets.apps_types (id, name) VALUES (2, 'frontend');

create table pets.apps
(
    id          bigint auto_increment
        primary key,
    name        varchar(255) not null,
    project_id  bigint       not null,
    active      tinyint(1) default 1 not null,
    app_type_id bigint null,
    constraint name
        unique (name),
    constraint project_id
        unique (project_id),
    constraint apps_apps_types_apps
        foreign key (app_type_id) references pets.apps_types (id)
            on delete set null
);
