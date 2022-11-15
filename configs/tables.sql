CREATE TABLE if not exists account (
    account_id int PRIMARY KEY,
    balance decimal
);

CREATE TABLE if not exists reservation (
    reservation_id serial PRIMARY KEY,
    account_id int REFERENCES account(account_id),
    service_id int NOT NULL,
    order_id int NOT NULL UNIQUE,
    amount decimal NOT NULL,
    reserve_date date NOT NULL,
    reservation_status varchar(20) NOT NULL
);

CREATE TABLE if not exists accounting_report (
    report_id serial PRIMARY KEY,
    service_id int NOT NULL,
    order_id int NOT NULL UNIQUE,
    amount decimal NOT NULL,
    account_id int,
    report_date date,
    status varchar(20) NOT NULL
)