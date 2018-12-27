DROP DATABASE elibrarygo;
CREATE DATABASE elibrarygo;

DROP DATABASE elibrary;
CREATE DATABASE elibrary;

-- View data
select * from bukus;
select * from peminjamans;
select * from users;

-- add data to buku
insert into bukus(judul,tahunterbit,pengarang) values ('testjudul2',1998,'kasnadi');
insert into bukus(judul,tahunterbit,pengarang) values ('testjudul3',1998,'kasnadi');
insert into bukus(judul,tahunterbit,pengarang) values ('testjudul4',1998,'kasnadi');
insert into bukus(judul,tahunterbit,pengarang) values ('testjudul5',1998,'kasnadi');


-- add data to user
insert into users(nama,email,password) values ('kasnadi1','kasnadi1@gmail.com','kasnadi1');
insert into users(nama,email,password) values ('kasnadi2','kasnadi2@gmail.com','kasnadi2');