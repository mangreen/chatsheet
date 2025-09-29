-- Up Migration: 創建資料表

-- 1. 創建 'users' 資料表
CREATE TABLE users (
    -- Email 作為主鍵
    email VARCHAR(255) PRIMARY KEY NOT NULL,
    
    -- Password 欄位用於存儲密碼哈希值
    password VARCHAR(255) NOT NULL,
    
    -- 創建和更新時間戳
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 創建一個函式來自動更新 updated_at 欄位
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW(); 
   RETURN NEW;
END;
$$ language 'plpgsql';

-- 為 users 表格創建觸發器
CREATE TRIGGER update_user_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


-- 2. 創建 'unipile_accounts' 資料表
CREATE TABLE unipile_accounts (
    -- ID 作為主鍵，使用 PostgreSQL 的 UUID 類型和 gen_random_uuid() 函數生成
    id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,

    -- 關聯到 User.Email
    user_email VARCHAR(255) NOT NULL,
    
    -- 服務提供者 (例如: linkedin)
    provider VARCHAR(50) NOT NULL,
    
    -- Unipile 服務返回的帳號 ID，必須是唯一的
    account_id VARCHAR(255) UNIQUE NOT NULL,

    -- 創建和更新時間戳
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- 設定外鍵約束，確保 user_email 必須存在於 users.email 中
    CONSTRAINT fk_user_email
        FOREIGN KEY(user_email) 
        REFERENCES users(email)
        ON DELETE CASCADE -- 當 user 被刪除時，其所有帳號也一併刪除
);

-- 為 unipile_accounts 表格創建觸發器
CREATE TRIGGER update_unipile_account_updated_at
BEFORE UPDATE ON unipile_accounts
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


-- Down Migration: 刪除資料表 (用於回滾)
-- ❗ 注意：如果您使用 migration 工具 (例如 Goose 或 Migrate)，您會將這些分成兩個檔案。

/*
DROP TRIGGER IF EXISTS update_unipile_account_updated_at ON unipile_accounts;
DROP TABLE IF EXISTS unipile_accounts;
DROP TRIGGER IF EXISTS update_user_updated_at ON users;
DROP TABLE IF EXISTS users;
DROP FUNCTION IF EXISTS update_updated_at_column;
*/