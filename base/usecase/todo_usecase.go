package usecase

import (
	"context"
	"errors"
	"strings"
	"todo/base/model"
    "todo/base/repository"
	"time"
)

type TodoUsecaseInterface interface {
	CreateTodo(ctx context.Context, title string, description string, accountId *int64, dueDate time.Time) (*model.Todo, error)
	GetTodoByID(ctx context.Context, id uint) (*model.Todo, error)
	FetchTodoList(ctx context.Context, accountId *int64) ([]*model.Todo, error)
}

type TodoUsecase struct {
	todoRepo *repository.PostgresTodoRepository // 以前定義したリポジトリ
}

func NewTodoUsecase(todoRepo *repository.PostgresTodoRepository) *TodoUsecase {
	return &TodoUsecase{todoRepo: todoRepo}
}

// 1. CreateTodo: ビジネスロジックを伴うTODO作成
func (u *TodoUsecase) CreateTodo(ctx context.Context, title string, description string, accountId *int64, dueDate time.Time) (*model.Todo, error) {
	if strings.TrimSpace(title) == "" {
		return nil, errors.New("todo title cannot be empty")
	}

	todo := &model.Todo{
		Title:       title,
		Description: description,
		AccountId:   accountId,
		DueDate:     dueDate,
		// time.Now() などの生成やID自動採番は、UsecaseやDB（GORM）の責務にします
	}

	if err := u.todoRepo.Create(ctx, todo); err != nil {
		return nil, err
	}

	return todo, nil
}

// 2. GetTodoByID: 1件取得（スタブから本番用に切り替え可能に）
func (u *TodoUsecase) GetTodoByID(ctx context.Context, id uint) (*model.Todo, error) {
	todo, err := u.todoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if todo == nil {
		return nil, errors.New("todo not found")
	}
	return todo, nil
}

// 3. FetchTodoList: 一覧取得（新規追加）
func (u *TodoUsecase) FetchTodoList(ctx context.Context, accountId *int64) ([]*model.Todo, error) {
	return u.todoRepo.FindAll(ctx, accountId)
}
