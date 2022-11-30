-- eKYC schema

CREATE TABLE plan(
    id INT PRIMARY KEY,
    name VARCHAR NOT NULL,
    daily_base_cost NUMERIC NOT NULL,
    api_face_match_cost NUMERIC NOT NULL,
    api_ocr_cost NUMERIC NOT NULL,
    image_upload_cost NUMERIC NOT NULL
)

CREATE TABLE client(
    id INT PRIMARY KEY,
    plan_id INT REFERENCES plan(id),
    name VARCHAR NOT NULL,
    access_key VARCHAR NOT NULL,
    secret_key VARCHAR NOT NULL,
    email_id VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
)

CREATE TABLE client_plans(
    id BIGINT PRIMARY KEY,
    client_id INT REFERENCES client(id),
    plan_id INT REFERENCES plan(id),
    timestamp TIMESTAMP NOT NULL
)

CREATE TABLE file_type(
    id INT PRIMARY KEY,
    name VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
)

CREATE TABLE file(
    id BIGINT PRIMARY KEY,
    client_id INT REFERENCES client(id),
    file_type_id INT REFERENCES file_type(id),
    extension VARCHAR NOT NULL,
    size NUMERIC NOT NULL,
    uuid VARCHAR UNIQUE NOT NULL,
    timestamp TIMESTAMP NOT NULL
)

CREATE TABLE face_match(
    id BIGINT PRIMARY KEY,
    client_id INT REFERENCES client(id),
    image1 BIGINT REFERENCES file(id),
    image2 BIGINT REFERENCES file(id),
    score NUMERIC NOT NULL,
    timestamp TIMESTAMP NOT NULL
)

CREATE TABLE ocr_info(
    id BIGINT PRIMARY KEY,
    client_id INT REFERENCES client(id), 
    id_card BIGINT REFERENCES file(id),
    name VARCHAR NOT NULL,
    gender VARCHAR NOT NULL,
    id_number VARCHAR NOT NULL,
    address_line_1 VARCHAR NOT NULL,
    address_line_2 VARCHAR NOT NULL,
    pincode VARCHAR NOT NULL,
    timestamp TIMESTAMP NOT NULL
)