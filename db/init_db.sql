/*

инициализация БД

*/

-- =============================
-- under user postgres:

CREATE USER appuser WITH PASSWORD 'appuser';

CREATE DATABASE app_db
TEMPLATE = 'template0'
ENCODING = 'utf-8'
LC_COLLATE = 'C.UTF-8'
LC_CTYPE = 'C.UTF-8';

REVOKE ALL ON DATABASE app_db FROM PUBLIC;
ALTER DATABASE app_db OWNER TO appuser;

-- =============================
-- connect to database:

\c app_db

-- =============================
-- under user appuser:

SET ROLE appuser;

DROP TABLE IF EXISTS users;
CREATE TABLE users (
uuid CHAR(36) UNIQUE NOT NULL,
user_name VARCHAR(100) NOT NULL,
email VARCHAR(254) UNIQUE NOT NULL,

PRIMARY KEY(uuid)
);

DROP TABLE IF EXISTS groups;
CREATE TABLE groups (
uuid CHAR(36) UNIQUE NOT NULL,
group_name VARCHAR(50) UNIQUE NOT NULL,
group_type VARCHAR(50) NOT NULL,
descr VARCHAR(254) NOT NULL,

PRIMARY KEY(uuid)
);


DROP TABLE IF EXISTS membership;
CREATE TABLE membership (
group_uuid CHAR(36) REFERENCES groups (uuid),
user_uuid  CHAR(36) REFERENCES users (uuid),

constraint membership_fk_group_uuid FOREIGN KEY (group_uuid) references groups (uuid) on delete restrict,
constraint membership_fk_user_uuid FOREIGN KEY (user_uuid) references users (uuid) on delete restrict,
constraint membership_uniq_gid_uid_pair UNIQUE (group_uuid, user_uuid)
);



INSERT INTO users (uuid, user_name, email) VALUES
('fd0250b7-fed0-426d-991b-45bc9e927b14','Laron','feest.margaretta@example.net'),
('e707a24d-79c1-4d56-8d01-f79d30054a09','Maribel','guiseppe31@example.org'),
('206c5c37-41c5-49fa-b443-47e97839c0b8','Seth','deangelo.kilback@example.net'),
('19dfc336-3172-4be9-be06-5670e3610123','Gerry','xschneider@example.net'),
('8fb59c30-5263-45c2-a13d-8e07fde474c3','Dorris','johns.fay@example.com'),
('49bdd317-6986-4df9-b14c-db2dea7e7c0d','Coy','hcummings@example.org'),
('5841aa24-4f74-49a1-973c-837f06cae09a','Georgette','brenden69@example.org'),
('913a6629-44ec-4512-ab19-c1eab86f2082','Pearline','ryleigh.terry@example.net'),
('d37f409a-62d3-4572-b2c8-2f73500b6241','Josefa','elmore.orn@example.com'),
('f166aaba-69bf-4cc9-b1de-92116eba9685','Ulises','lonie.o''kon@example.org');


INSERT INTO groups (uuid, group_name, group_type, descr) VALUES
('c3b5be32-8510-49de-ba1a-7dc5102a95ab','Dolorem accusantium','project','Laboriosam asperiores a omnis dolores fuga. Provident vitae qui reprehenderit quo voluptas.'),
('855ce14d-fb0b-47e0-9617-c8fa1ecafa8b','Tempora ab totam','community','Culpa aut commodi vero modi commodi aliquid. Enim dolores qui sapiente beatae doloremque.'),
('10b29442-945d-433d-be56-2a9dc4094472','Delectus et quae','organization','Tempore laborum deleniti et officia ab et omnis. Possimus perferendis maxime itaque in.'),
('4b5dd85f-c53b-466f-b9a4-3eed66325442','Nobis nesciunt','project','Odit hic nostrum neque porro eum maiores harum. Harum fuga id tempore eum reprehenderit fugit error.');


INSERT INTO membership (group_uuid, user_uuid) VALUES
('c3b5be32-8510-49de-ba1a-7dc5102a95ab','fd0250b7-fed0-426d-991b-45bc9e927b14'),
('c3b5be32-8510-49de-ba1a-7dc5102a95ab','e707a24d-79c1-4d56-8d01-f79d30054a09'),
('c3b5be32-8510-49de-ba1a-7dc5102a95ab','206c5c37-41c5-49fa-b443-47e97839c0b8'),
('855ce14d-fb0b-47e0-9617-c8fa1ecafa8b','fd0250b7-fed0-426d-991b-45bc9e927b14'),
('855ce14d-fb0b-47e0-9617-c8fa1ecafa8b','206c5c37-41c5-49fa-b443-47e97839c0b8'),
('855ce14d-fb0b-47e0-9617-c8fa1ecafa8b','8fb59c30-5263-45c2-a13d-8e07fde474c3'),
('855ce14d-fb0b-47e0-9617-c8fa1ecafa8b','49bdd317-6986-4df9-b14c-db2dea7e7c0d'),
('855ce14d-fb0b-47e0-9617-c8fa1ecafa8b','5841aa24-4f74-49a1-973c-837f06cae09a'),
('855ce14d-fb0b-47e0-9617-c8fa1ecafa8b','913a6629-44ec-4512-ab19-c1eab86f2082'),
('855ce14d-fb0b-47e0-9617-c8fa1ecafa8b','d37f409a-62d3-4572-b2c8-2f73500b6241'),
('10b29442-945d-433d-be56-2a9dc4094472','fd0250b7-fed0-426d-991b-45bc9e927b14'),
('10b29442-945d-433d-be56-2a9dc4094472','e707a24d-79c1-4d56-8d01-f79d30054a09'),
('10b29442-945d-433d-be56-2a9dc4094472','206c5c37-41c5-49fa-b443-47e97839c0b8'),
('10b29442-945d-433d-be56-2a9dc4094472','19dfc336-3172-4be9-be06-5670e3610123'),
('10b29442-945d-433d-be56-2a9dc4094472','8fb59c30-5263-45c2-a13d-8e07fde474c3'),
('10b29442-945d-433d-be56-2a9dc4094472','49bdd317-6986-4df9-b14c-db2dea7e7c0d'),
('10b29442-945d-433d-be56-2a9dc4094472','5841aa24-4f74-49a1-973c-837f06cae09a'),
('10b29442-945d-433d-be56-2a9dc4094472','913a6629-44ec-4512-ab19-c1eab86f2082'),
('10b29442-945d-433d-be56-2a9dc4094472','d37f409a-62d3-4572-b2c8-2f73500b6241'),
('10b29442-945d-433d-be56-2a9dc4094472','f166aaba-69bf-4cc9-b1de-92116eba9685');
