-- +migrate Up
create table pilots
(
    id int  not null primary key,
    name text not null,
    age      int  not null
);

create table jets
(
    id       int  not null  primary key,
    pilot_id int  not null,
    age      int  not null,
    name     text not null,
    color    text not null,
    
    constraint jets_pilots_fk foreign key (pilot_id) references pilots (id)
);

create table languages
(
    id int not null primary key,
    language text not null
);


-- Join table
create table pilot_languages
(
    pilot_id    int not null,
    language_id int not null,
    primary key (pilot_id, language_id), 

    constraint pilot_languages_languages_fk foreign key (language_id) references languages (id),
    constraint pilot_languages_pilots_fk    foreign key (pilot_id) references pilots (id)
);

-- Composite primary key
-- +migrate Down
drop table pilot_languages;
drop table languages;
drop table jets;
drop table pilots;