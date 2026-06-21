package repository

import (
    "context"
    "regexp"
    "testing"
    "time"
    "todo/base/model"
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/assert" // go get github.com/stretchr/testify をしておくと便利です（標準の t.Errorf でも可）
)

// Create メソッドの単体テスト
func TestPostgresTodoRepository_Create(t *testing.T) {
    // モックDBの初期化
    db, mock := SetupTestDB()
    repo := NewPostgresTodoRepository(db)

    ctx := context.Background()
    now := time.Now()
    
    // テスト対象のデータ
    todo := &model.Todo{
        AccountID:   1,
        Title:       "テストタスク",
        Description: "これは単体テストです",
        IsCompleted: false,
    }

    // 【期待する挙動の設定】
    // GORMが内部で発行する INSERT クエリと、その結果返ってくるダミーのID（1）を定義します
    // ※GORMのINSERTは「INSERT INTO "todo" ...」となるため、正規表現でマッチさせます
    mock.ExpectBegin() // トランザクションの開始を期待
    mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "todo"`)).
        WithArgs(todo.AccountID, todo.Title, todo.Description, todo.IsCompleted, todo.DueDate, sqlmock.AnyArg(), sqlmock.AnyArg()).
        WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(uint(1), now, now))
    mock.ExpectCommit() // トランザクションのコミットを期待

    // 【実行】
    err := repo.Create(ctx, todo)

    // 【検証】
    assert.NoError(t, err)
    assert.Equal(t, uint(1), todo.ID) // 自動採番されたID（1）が構造体にセットされているか
    
    // すべての期待したクエリが予定通り実行されたかチェック
    err = mock.ExpectationsWereMet()
    assert.NoError(t, err)
}

// FindByID メソッドの単体テスト（正常系と異常系）
func TestPostgresTodoRepository_FindByID(t *testing.T) {
    db, mock := SetupTestDB()
    repo := NewPostgresTodoRepository(db)
    ctx := context.Background()

    t.Run("正常にTODOが取得できる場合", func(t *testing.T) {
        expectedID := uint(1)
        
        // SELECTクエリが走り、1行のデータが返ってくることを期待
        rows := sqlmock.NewRows([]string{"id", "account_id", "title", "description", "is_completed"}).
            AddRow(expectedID, uint(100), "テストタイトル", "詳細", false)
            
        mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "todo" WHERE "todo"."id" = $1`)).
            WithArgs(expectedID, 1). // Firstを使用するためLIMIT 1の1が含まれる
            WillReturnRows(rows)

        // 実行
        todo, err := repo.FindByID(ctx, expectedID)

        // 検証
        assert.NoError(t, err)
        assert.NotNil(t, todo)
        assert.Equal(t, expectedID, todo.ID)
        assert.Equal(t, "テストタイトル", todo.Title)
    })

    t.Run("データが存在しない場合（nilを返す設計の検証）", func(t *testing.T) {
        targetID := uint(999)
        
        // データが空（0行）で返ってくることを期待
        mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "todo" WHERE "todo"."id" = $1`)).
            WithArgs(targetID, 1).
            WillReturnRows(sqlmock.NewRows([]string{}))

        // 実行
        todo, err := repo.FindByID(ctx, targetID)

        // 検証（リポジトリの設計通り、エラーにならず nil が返ってくるか）
        assert.NoError(t, err)
        assert.Nil(t, todo)
    })

    // すべての期待値の検証
    err := mock.ExpectationsWereMet()
    assert.NoError(t, err)
}
