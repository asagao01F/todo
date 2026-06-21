package usecase

import (
    "context"
    "todo/base/model"
)

// Usecase側が「俺はこれらの機能が欲しい」とインターフェースを宣言する
type TodoRepository interface {
    Create(ctx context.Context, todo *model.Todo) error
}

type TodoUsecase struct {
    repo TodoRepository // Usecaseは自分が定義したインターフェースに依存する
}
