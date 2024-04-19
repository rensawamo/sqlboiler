-- +migrate Up
insert into pilots (id, name, age) values (1, 'taro yamada', 40);
insert into pilots (id, name, age) values (2, 'jiro tanaka', 50);

insert into jets (id, pilot_id, age, name, color) values (1, 1, 1, '727', 'green');
insert into jets (id, pilot_id, age, name, color) values (2, 1, 2, '707', 'red');
insert into jets (id, pilot_id, age, name, color) values (3, 2, 3, 'DC-8', 'white');
insert into jets (id, pilot_id, age, name, color) values (4, 2, 10, 'Comet', 'black');

insert into languages (id, language) values (1, 'english');
insert into languages (id, language) values (2, 'japanese');
insert into languages (id, language) values (3, 'esperanto');
insert into languages (id, language) values (4, 'klingon');

insert into pilot_languages (pilot_id, language_id) values (1, 1);
insert into pilot_languages (pilot_id, language_id) values (1, 2);
insert into pilot_languages (pilot_id, language_id) values (2, 1);
insert into pilot_languages (pilot_id, language_id) values (2, 3);

-- +migrate Down
delete from pilot_languages;
delete from languages;
delete from jets;
delete from pilots;

