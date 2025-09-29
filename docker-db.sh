#!/bin/bash
set -e

# --- 0. 設置必要的目錄和權限 ---
echo "Creating necessary directories..."
mkdir -p /var/log/postgresql/
chown -R postgres:postgres /var/log/postgresql/
chmod 755 /var/log/postgresql/
# PGDATA 是 postgres 官方映像檔用於儲存資料的路徑，通常是 /var/lib/postgresql/data
PGDATA="/var/lib/postgresql/data" 

# --- 1. 啟動 PostgreSQL ---
echo "Starting PostgreSQL initialization and server..."

# 確保 PGDATA 存在且權限正確
mkdir -p "$PGDATA"
chown -R postgres:postgres "$PGDATA"

# 檢查資料庫是否已初始化
if [ ! -s "$PGDATA/PG_VERSION" ]; then
    echo "Initializing PostgreSQL data directory..."
    # 💡 [關鍵修正 1]: 
    # 為了讓 initdb 使用 POSTGRES_USER/PASSWORD/DB 環境變數，
    # 這些變數必須在 gosu postgres 執行時可見。
    # 此外，我們需要確保 initdb 確實創建了用戶。
    
    # 這是 postgres 映像檔中 initdb 的標準執行方式
    export PGPASSWORD="${POSTGRES_PASSWORD}"
    gosu postgres initdb
    
    # 💡 [關鍵修正 2]: 在初始化後，創建應用程式需要的角色和資料庫
    echo "Creating user and database..."
    # 啟動臨時服務來執行 SQL
    gosu postgres pg_ctl -D "$PGDATA" -w start -t 5 
    
    # 使用 psql 創建角色和資料庫
    # -v ON_ERROR_STOP=1 是為了確保任何 SQL 錯誤都會導致腳本停止
    # 這裡我們使用 -U postgres (預設超級用戶) 來執行創建
    gosu postgres psql --username postgres -d postgres <<-EOSQL
        -- 確保角色存在並設定密碼
        CREATE ROLE "${POSTGRES_USER}" WITH LOGIN PASSWORD '${POSTGRES_PASSWORD}';
        -- 創建資料庫並將所有權賦予新創建的角色
        CREATE DATABASE "${POSTGRES_DB}" OWNER "${POSTGRES_USER}";
        -- 撤銷 PUBLIC 寫入權限 (標準安全做法)
        REVOKE ALL ON DATABASE "${POSTGRES_DB}" FROM PUBLIC;
        GRANT ALL PRIVILEGES ON DATABASE "${POSTGRES_DB}" TO "${POSTGRES_USER}";
EOSQL
    
    # 停止臨時服務
    gosu postgres pg_ctl -D "$PGDATA" stop
    echo "User and database created."

fi

# --- 2. 啟動最終 PostgreSQL 服務 ---
echo "Starting final PostgreSQL server in background..."
# 這裡我們使用 -c 來設定監聽位址，確保應用程式可以通過 localhost 連線
gosu postgres pg_ctl -D "$PGDATA" -l /var/log/postgresql/server.log start &

# --- 3. 等待 PostgreSQL 準備就緒 (連線測試) ---
echo "Waiting for PostgreSQL to be ready for application connections..."
# 連線測試應該使用新創建的角色
until pg_isready -h localhost -p 5432 -U "${POSTGRES_USER}" > /dev/null 2>&1; do
  echo "PostgreSQL not ready, waiting..."
  sleep 1
done
echo "PostgreSQL is ready! (User: ${POSTGRES_USER})"

# --- 4. 啟動 Go 應用程式 ---
echo "Starting Go application..."
/chatsheet/myapp "$@" # 運行 Go 應用程式，並傳遞 CMD 參數