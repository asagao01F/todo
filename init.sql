CREATE TABLE IF NOT EXISTS account (
    id SERIAL PRIMARY KEY,                        -- GORMのuint型 + primaryKey
    username VARCHAR(255) NOT NULL UNIQUE,         -- not null;unique
    email VARCHAR(255) NOT NULL UNIQUE,            -- not null;unique (types.Emailは文字列として格納)
    password VARCHAR(255) NOT NULL,                -- not null
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- time.Time
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP  -- time.Time
);

CREATE TABLE IF NOT EXISTS todo (
    id SERIAL PRIMARY KEY,                        -- GORMのuint型 + primaryKey
    account_id BIGINT,                             -- *int64 に対応（NULLを許容）
    title VARCHAR(255) NOT NULL,                   -- タイトル
    description TEXT,                              -- 説明文（長文に対応できるようTEXT型にしています）
    is_completed BOOLEAN DEFAULT FALSE,            -- bool型。デフォルトは未完了(FALSE)
    due_date TIMESTAMP WITH TIME ZONE,             -- 期限日
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- 💡 Accountテーブルとのリレーション（外部キー制約）を設定
    CONSTRAINT fk_account FOREIGN KEY (account_id) 
        REFERENCES account(id) 
        ON DELETE SET NULL                         -- アカウントが消えたら、紐づくTodoのaccount_idをNULLにする
);
