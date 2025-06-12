CREATE DATABASE log_analysis;
\c log_analysis
-- Raw logs storage
CREATE TABLE IF NOT EXISTS raw_logs (
    id SERIAL PRIMARY KEY,
    source_type VARCHAR(50) NOT NULL, -- dhcp, dns, firewall, router, switch
    log_time TIMESTAMP NOT NULL,
    raw_content TEXT NOT NULL,
    processed BOOLEAN DEFAULT false,
    ingested_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- Parsed logs with structured data
CREATE TABLE IF NOT EXISTS parsed_logs (
    id SERIAL PRIMARY KEY,
    raw_log_id INTEGER REFERENCES raw_logs(id),
    log_type VARCHAR(50) NOT NULL, -- dhcp_lease, dns_query, firewall_accept, etc.
    log_time TIMESTAMP NOT NULL,
    source_ip VARCHAR(45),
    destination_ip VARCHAR(45),
    source_port INTEGER,
    destination_port INTEGER,
    protocol VARCHAR(20),
    action VARCHAR(50),
    details JSONB -- Additional structured data specific to log type
    );

-- Discovered assets
CREATE TABLE IF NOT EXISTS assets (
    id SERIAL PRIMARY KEY,
    ip_address VARCHAR(45),
    mac_address VARCHAR(20),
    hostname VARCHAR(255),
    first_seen TIMESTAMP NOT NULL,
    last_seen TIMESTAMP NOT NULL,
    asset_type VARCHAR(50), -- server, workstation, network, etc.
    metadata JSONB -- Additional information
    );
-- Services for each assets
CREATE TABLE services(
    id SERIAL PRIMARY KEY,
    asset_ip  VARCHAR(45),
    name VARCHAR(100) NOT NULL,
    port INTEGER,
    protocol VARCHAR(10)
);

CREATE DATABASE asset_management;
\c asset_management
CREATE TABLE if not exists assets(
    ip_address VARCHAR(64) PRIMARY KEY,
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
    id SERIAL PRIMARY KEY,
    asset_ip VARCHAR(64) REFERENCES assets(ip_address) on DELETE CASCADE,
    name VARCHAR(100),
    port INTEGER,
    protocol VARCHAR(10),
    description TEXT,
    is_managed BOOLEAN DEFAULT FALSE
);