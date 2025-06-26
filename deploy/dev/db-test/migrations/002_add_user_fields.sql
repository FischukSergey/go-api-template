-- Добавляем новые поля в таблицу users
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS last_name varchar(100),
ADD COLUMN IF NOT EXISTS phone varchar(20),
ADD COLUMN IF NOT EXISTS sex varchar(20),
ADD COLUMN IF NOT EXISTS city varchar(50),
ADD COLUMN IF NOT EXISTS country varchar(50);
