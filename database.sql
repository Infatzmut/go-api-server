create database inventoryñ
use inventory;

create table products(
    id int NOT NULL AUTO_INCREMENT,
    name varchar(255) NOT NULL,
    quantity int,
    price float(10,7),
    PRIMARY KEY(id)
);

insert into products values(1, "chair", 100, 200.00);
insert into products values(1, "desk", 800, 600.00);
