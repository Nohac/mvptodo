CREATE TABLE todos (    
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    title varchar(100),
    description TEXT,
    status boolean
);
