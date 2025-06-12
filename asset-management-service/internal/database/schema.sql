CREATE DATABASE asset_management;

CREATE TABLE assets(
    ip_address VARCHAR(64) PRIMARY KEY ,
    mac_address VARCHAR(64),
    hostname VARCHAR(255),
    asset_type VARCHAR(50),
    location VARCHAR(100),
    operating_system VARCHAR(100),
    first_seen TIMESTAMP WITH TIME ZONE NOT NULL,
    last_seen TIMESTAMP WITH TIME ZONE NOT NULL ,
    is_managed BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE services(
    id SERIAL PRIMARY KEY UNIQUE ,
    asset_ip VARCHAR(64) REFERENCES assets(ip_address) on DELETE CASCADE,
    name VARCHAR(100),
    port INTEGER,
    protocol VARCHAR(10),
    description TEXT,
    is_managed BOOLEAN DEFAULT FALSE
);