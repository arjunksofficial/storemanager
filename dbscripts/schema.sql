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

-- type Request struct {
-- 	BaseModel
-- 	StoreTID   int64
-- 	VisitTime  *time.Time
-- 	ImageCount int64
-- 	Status     string
-- }

CREATE TABLE storeDB.requests
(
    id INT(10) unsigned NOT NULL AUTO_INCREMENT,
    created_at DATETIME DEFAULT NULL,
    updated_at DATETIME DEFAULT NULL,
    status VARCHAR(255) NOT NULL,
    PRIMARY KEY (id)
);

-- type Image struct {
-- 	BaseModel
-- 	RequestID int64
-- 	ImageURL  string
-- 	Perimeter float64
-- 	Status    string
-- }

CREATE TABLE storeDB.images
(
    id INT(10) unsigned NOT NULL AUTO_INCREMENT,
    created_at DATETIME DEFAULT NULL,
    updated_at DATETIME DEFAULT NULL,
    request_tid INT(10) unsigned NOT NULL,
    image_url VARCHAR(255) NOT NULL,
    store_tid INT(10) unsigned NOT NULL,
    visit_time DATETIME DEFAULT NULL,
    perimeter FLOAT,
    status VARCHAR(255) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (request_tid) REFERENCES storeDB.requests(id),
    FOREIGN KEY (store_tid) REFERENCES storeDB.stores(id)
);