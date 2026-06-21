package usecase

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"todo/base/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockRepository(t *testing.T) (repository.PostgresTodoRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("failed to open a stub database connection: %v", err)
	}

	mock.ExpectPing()

	dialer := postgres.New(postgres.Config{
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gormDB, err := gorm.Open(dialer, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	repo := repository.NewPostgresTodoRepository(gormDB)
	return *repo, mock
}

// 1. CreateTodo のテスト
func TestTodoUsecase_CreateTodo(t *testing.T) {
	ctx := context.Background()

	t.Run("成功: バリデーションを通過し、TODOが正常に作成されること", func(t *testing.T) {
		repo, mock := setupMockRepository(t)
		uc := NewTodoUsecase(&repo)

		title := "テストタスク"
		description := "詳細な説明文"

		mock.ExpectBegin()
		// 【修正】テーブル名を "todo" に変更し、全7カラムの引数に対応
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "todo" ("account_id","title","description","is_completed","due_date","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).
			WithArgs(
				uint(1),              // account_id (デフォルト)
				title,                // title
				description,          // description
				false,                // is_completed
				nil,                  // due_date
				sqlmock.AnyArg(),     // created_at
				sqlmock.AnyArg(),     // updated_at
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint(1)))
		mock.ExpectCommit()

		todo, err := uc.CreateTodo(ctx, title, description, &[]int64{1}[0])

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if todo == nil {
			t.Fatal("expected todo not to be nil")
		}
		if todo.ID != 1 || todo.Title != title || todo.Description != description {
			t.Errorf("unexpected todo fields: %+v", todo)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("失敗: タイトルが空（空白のみ）の場合、バリデーションエラーになること", func(t *testing.T) {
		repo, _ := setupMockRepository(t)
		uc := NewTodoUsecase(&repo)

		todo, err := uc.CreateTodo(ctx, "   ", "説明文", &[]int64{1}[0])

		if err == nil {
			t.Error("expected error, got nil")
		}
		if err.Error() != "todo title cannot be empty" {
			t.Errorf("expected validation error message, got: %v", err)
		}
		if todo != nil {
			t.Errorf("expected todo to be nil, got %+v", todo)
		}
	})

	t.Run("失敗: リポジトリ側(DB)でエラーが発生した場合、エラーがそのまま返ること", func(t *testing.T) {
		repo, mock := setupMockRepository(t)
		uc := NewTodoUsecase(&repo)

		mock.ExpectBegin()
		// 【修正】ここも期待するテーブル名を "todo" に変更
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "todo"`)).
			WillReturnError(fmt.Errorf("db connection error"))
		mock.ExpectRollback()

		todo, err := uc.CreateTodo(ctx, "タイトル", "説明", &[]int64{1}[0])

		if err == nil {
			t.Error("expected error, got nil")
		}
		if todo != nil {
			t.Errorf("expected todo to be nil, got %+v", todo)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

// 2. GetTodoByID のテスト
func TestTodoUsecase_GetTodoByID(t *testing.T) {
	ctx := context.Background()

	t.Run("成功: レコードが存在する場合、データが返ること", func(t *testing.T) {
		repo, mock := setupMockRepository(t)
		uc := NewTodoUsecase(&repo)

		targetID := uint(10)
		now := time.Now()

		// 【修正】テーブル名を "todo" に変更
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "todo" WHERE "todo"."id" = $1 ORDER BY "todo"."id" LIMIT $2`)).
			WithArgs(targetID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "created_at"}).
				AddRow(targetID, "既存タスク", "既存の説明", now))

		todo, err := uc.GetTodoByID(ctx, targetID)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if todo == nil || todo.ID != targetID || todo.Title != "既存タスク" {
			t.Errorf("unexpected result: %+v", todo)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("失敗: レコードが存在しない場合（リポジトリがnilを返した場合）、todo not found エラーになること", func(t *testing.T) {
		repo, mock := setupMockRepository(t)
		uc := NewTodoUsecase(&repo)

		targetID := uint(999)

		// 【修正】テーブル名を "todo" に変更
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "todo" WHERE "todo"."id" = $1`)).
			WithArgs(targetID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		todo, err := uc.GetTodoByID(ctx, targetID)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if err.Error() != "todo not found" {
			t.Errorf("expected 'todo not found' error, got: %v", err)
		}
		if todo != nil {
			t.Errorf("expected todo to be nil, got %+v", todo)
		}
	})
}
