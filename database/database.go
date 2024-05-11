package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB() (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("%s://%s:%s@localhost:%s/%s", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	return pgxpool.New(context.Background(), dsn)
}

func TestDB() {
	conn, err := pgx.Connect(context.Background(), fmt.Sprintf("%s://%s:%s@localhost:%s/%s", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME")))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())
	sqlStatements := `
	BEGIN;

	-- Insert test departments
	insert into departments (name) values
	('sales'),
	('marketing'),
	('engineering'),
	('finance'),
	('purchasing');

	-- Insert test roles
	insert into roles (name) values
	('admin'),
	('manager'),
	('user');

	-- Insert test users
	insert into users (email, username, password) values
	('john@example.com', 'johndoe', '$2a$14$xrr8iyht40fgq75ijylko.zktdbdadmtxi2.zarpxkjhlky8onbqq'),
	('jane@example.com', 'janesmith', '$2a$14$xrr8iyht40fgq75ijylko.zktdbdadmtxi2.zarpxkjhlky8onbqq'),
	('mike@example.com', 'mikejohnson', '$2a$14$xrr8iyht40fgq75ijylko.zktdbdadmtxi2.zarpxkjhlky8onbqq'),
	('emily@example.com', 'emilybrown', '$2a$14$xrr8iyht40fgq75ijylko.zktdbdadmtxi2.zarpxkjhlky8onbqq');

	-- Insert user roles
	insert into user_roles (user_id, role_id)
	select 
	(select id from users where email = 'john@example.com'),
	(select id from roles where name = 'admin');

	insert into user_roles (user_id, role_id)
	select
	(select id from users where email = 'jane@example.com'),
	(select id from roles where name = 'manager');

	insert into user_roles (user_id, role_id)
	select
	(select id from users where email = 'mike@example.com'),
	(select id from roles where name = 'user');

	insert into user_roles (user_id, role_id)
	select
	(select id from users where email = 'emily@example.com'),
	(select id from roles where name = 'manager');

	insert into user_roles (user_id, role_id)
	select
	(select id from users where email = 'emily@example.com'),
	(select id from roles where name = 'user');

	-- Insert user departments
	insert into user_departments (user_id, department_id)
	select
	(select id from users where email = 'john@example.com'),
	(select id from departments where name = 'sales');

	insert into user_departments (user_id, department_id)
	select
	(select id from users where email = 'john@example.com'),
	(select id from departments where name = 'marketing');

	insert into user_departments (user_id, department_id)
	select
	(select id from users where email = 'jane@example.com'),
	(select id from departments where name = 'marketing');

	insert into user_departments (user_id, department_id)
	select
	(select id from users where email = 'mike@example.com'),
	(select id from departments where name = 'engineering');

	insert into user_departments (user_id, department_id)
	select
	(select id from users where email = 'emily@example.com'),
	(select id from departments where name = 'finance');

	COMMIT;
	`

	_, err = conn.Exec(context.Background(), sqlStatements)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("test data inserted successfully")
}
