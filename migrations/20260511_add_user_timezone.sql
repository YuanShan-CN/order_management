ALTER TABLE users 
ADD COLUMN timezone VARCHAR(50) NOT NULL DEFAULT 'Asia/Shanghai' 
AFTER password;
