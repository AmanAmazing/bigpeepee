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
('john@example.com', 'johndoe', '$2a$14$9gXg5n7LWwSjY/LeCdzKU.V1nFxdmebzgZfuz.h65JVE3bFBzEg6a'),
('jane@example.com', 'janesmith', '$2a$14$9gXg5n7LWwSjY/LeCdzKU.V1nFxdmebzgZfuz.h65JVE3bFBzEg6a'),
('mike@example.com', 'mikejohnson', '$2a$14$9gXg5n7LWwSjY/LeCdzKU.V1nFxdmebzgZfuz.h65JVE3bFBzEg6a'),
('emily@example.com', 'emilybrown', '$2a$14$9gXg5n7LWwSjY/LeCdzKU.V1nFxdmebzgZfuz.h65JVE3bFBzEg6a'),
('david@example.com', 'davidmiller', '$2a$14$9gXg5n7LWwSjY/LeCdzKU.V1nFxdmebzgZfuz.h65JVE3bFBzEg6a'),
('sarah@example.com', 'sarahwilson', '$2a$14$9gXg5n7LWwSjY/LeCdzKU.V1nFxdmebzgZfuz.h65JVE3bFBzEg6a'),
('kevin@example.com', 'kevinjones', '$2a$14$9gXg5n7LWwSjY/LeCdzKU.V1nFxdmebzgZfuz.h65JVE3bFBzEg6a'),
('amy@example.com', 'amytaylor', '$2a$14$9gXg5n7LWwSjY/LeCdzKU.V1nFxdmebzgZfuz.h65JVE3bFBzEg6a'),
('james@example.com', 'jameslee', '$2a$14$9gXg5n7LWwSjY/LeCdzKU.V1nFxdmebzgZfuz.h65JVE3bFBzEg6a'),
('olivia@example.com', 'oliviaharris', '$2a$14$9gXg5n7LWwSjY/LeCdzKU.V1nFxdmebzgZfuz.h65JVE3bFBzEg6a');

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

insert into user_roles (user_id, role_id)
select
(select id from users where email = 'david@example.com'),
(select id from roles where name = 'user');

insert into user_roles (user_id, role_id)
select
(select id from users where email = 'sarah@example.com'),
(select id from roles where name = 'manager');

insert into user_roles (user_id, role_id)
select
(select id from users where email = 'kevin@example.com'),
(select id from roles where name = 'user');

insert into user_roles (user_id, role_id)
select
(select id from users where email = 'amy@example.com'),
(select id from roles where name = 'manager');

insert into user_roles (user_id, role_id)
select
(select id from users where email = 'james@example.com'),
(select id from roles where name = 'user');

insert into user_roles (user_id, role_id)
select
(select id from users where email = 'olivia@example.com'),
(select id from roles where name = 'manager');

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

insert into user_departments (user_id, department_id)
select
(select id from users where email = 'david@example.com'),
(select id from departments where name = 'sales');

insert into user_departments (user_id, department_id)
select
(select id from users where email = 'sarah@example.com'),
(select id from departments where name = 'engineering');

insert into user_departments (user_id, department_id)
select
(select id from users where email = 'kevin@example.com'),
(select id from departments where name = 'purchasing');

insert into user_departments (user_id, department_id)
select
(select id from users where email = 'amy@example.com'),
(select id from departments where name = 'finance');

insert into user_departments (user_id, department_id)
select
(select id from users where email = 'james@example.com'),
(select id from departments where name = 'purchasing');

insert into user_departments (user_id, department_id)
select
(select id from users where email = 'olivia@example.com'),
(select id from departments where name = 'sales');

-- Insert test suppliers
insert into suppliers (name) values
('argos'),
('amazon'),
('zebra'),
('microsoft'),
('golak');

-- Insert test nominals. not sure what that even means lol
insert into nominals (name) values
('computer costs - gsit'),
('office - printer'),
('agency staff'),
('advertising');

-- Insert test products.
insert into products (name) values
('desk appliances'),
('desktops'),
('printer - subscription'),
('office furniture');

-- Insert purchase order approvers
insert into purchase_order_approvers (user_id)
select id from users where email in ('john@example.com', 'emily@example.com', 'amy@example.com');

-- Insert purchase orders at various stages
-- Pending purchase order
insert into purchase_orders (user_id, department_id, title, description, status, priority)
select
(select id from users where email = 'mike@example.com'),
(select id from departments where name = 'engineering'),
'Pending PO 1',
'This is a pending purchase order',
'pending',
'medium';

-- Manager approved purchase order
insert into purchase_orders (user_id, department_id, title, description, status, priority)
select
(select id from users where email = 'jane@example.com'),
(select id from departments where name = 'marketing'),
'Manager Approved PO 1',
'This is a manager approved purchase order',
'manager_approved',
'high';

-- Partially approved purchase order
insert into purchase_orders (user_id, department_id, title, description, status, priority)
select
(select id from users where email = 'david@example.com'),
(select id from departments where name = 'sales'),
'Partially Approved PO 1',
'This is a partially approved purchase order',
'partial_approved',
'low';

-- Fully approved purchase order
insert into purchase_orders (user_id, department_id, title, description, status, priority)
select
(select id from users where email = 'kevin@example.com'),
(select id from departments where name = 'purchasing'),
'Fully Approved PO 1',
'This is a fully approved purchase order',
'approved',
'high';

-- Insert purchase order items for the above purchase orders
insert into purchase_order_items (purchase_order_id, item_name, quantity, unit_price, total_price, link, status)
values
((select id from purchase_orders where title = 'Pending PO 1'), 'Item 1', 10, 50.00, 500.00, 'https://example.com/item1', 'pending'),
((select id from purchase_orders where title = 'Pending PO 1'), 'Item 2', 5, 100.00, 500.00, 'https://example.com/item2', 'pending'),
((select id from purchase_orders where title = 'Manager Approved PO 1'), 'Item 3', 8, 75.00, 600.00, 'https://example.com/item3', 'manager_approved'),
((select id from purchase_orders where title = 'Manager Approved PO 1'), 'Item 4', 3, 200.00, 600.00, 'https://example.com/item4', 'manager_approved'),
((select id from purchase_orders where title = 'Partially Approved PO 1'), 'Item 5', 12, 30.00, 360.00, 'https://example.com/item5', 'approved'),
((select id from purchase_orders where title = 'Partially Approved PO 1'), 'Item 6', 6, 80.00, 480.00, 'https://example.com/item6', 'pending'),
((select id from purchase_orders where title = 'Fully Approved PO 1'), 'Item 7', 20, 25.00, 500.00, 'https://example.com/item7', 'approved'),
((select id from purchase_orders where title = 'Fully Approved PO 1'), 'Item 8', 10, 120.00, 1200.00, 'https://example.com/item8', 'approved');

COMMIT;
