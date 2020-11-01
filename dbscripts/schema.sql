CREATE DATABASE storeDB;

USE storeDB;

-- type Store struct {
-- 	BaseModel
-- 	StoreID  string
-- 	Name     string
-- 	AreaCode int64
-- }

CREATE TABLE storeDB.stores
(
    id INT(10) unsigned NOT NULL AUTO_INCREMENT,
    created_at DATETIME DEFAULT NULL,
    updated_at DATETIME DEFAULT NULL,
    store_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    area_code INT(10) unsigned NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE storeDB.requests
(
    id INT(10) unsigned NOT NULL AUTO_INCREMENT,
    created_at DATETIME DEFAULT NULL,
    updated_at DATETIME DEFAULT NULL,
    status VARCHAR(255) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE storeDB.images
(
    id INT(10) unsigned NOT NULL AUTO_INCREMENT,
    created_at DATETIME DEFAULT NULL,
    updated_at DATETIME DEFAULT NULL,
    request_tid INT(10) unsigned NOT NULL,
    image_url VARCHAR(255) NOT NULL,
    store_id  VARCHAR(255) NOT NULL,
    visit_time DATETIME DEFAULT NULL,
    perimeter FLOAT  DEFAULT NULL,
    status VARCHAR(255) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (request_tid) REFERENCES storeDB.requests(id)
);