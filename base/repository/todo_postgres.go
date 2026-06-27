package repository

import (
    "context"
    "fmt"
    "todo/base/model"

    "gorm.io/gorm"
)

type PostgresTodoRepository struct {
    db *gorm.DB
}

func NewPostgresTodoRepository(db *gorm.DB) *PostgresTodoRepository {
    return &PostgresTodoRepository{db: db}
}

// 1. Create: TODOの作成
// GORMの .WithContext と .Create を使います。
// 実行後、自動採番されたIDや、データベース側で生成されたCreatedAt/UpdatedAtが自動的に todo 構造体に書き戻されます。
func (r *PostgresTodoRepository) Create(ctx context.Context, todo *model.Todo) error {
    result := r.db.WithContext(ctx).Create(todo)
    if result.Error != nil {
        return fmt.Errorf("failed to create todo: %w", result.Error)
    }
    return nil
}

// 2. FindByID: 1件取得
// .First を使うと、指定した主キー（ID）で検索します。
// データが見つからない場合は gorm.ErrRecordNotFound エラーが返るため、それをハンドリングします。
func (r *PostgresTodoRepository) FindByID(ctx context.Context, id uint) (*model.Todo, error) {
    var todo model.Todo
    result := r.db.WithContext(ctx).First(&todo, id)
    
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return nil, nil // 見つからない場合はエラーではなく nil を返す設計
        }
        return nil, fmt.Errorf("failed to find todo by id: %w", result.Error)
    }
    
    return &todo, nil
}

// 3. FindByAccountID: 特定ユーザーのTODO一覧取得
// .Where で条件を指定し、.Find でスライス（配列）に結果を詰め込みます。
// 最新のTODOが上にくるように .Order("created_at DESC") を挟んでいます。
func (r *PostgresTodoRepository) FindAll(ctx context.Context, accountId *int64) ([]*model.Todo, error) {
    var todos []*model.Todo
    query := r.db.WithContext(ctx).Order("created_at DESC")
    
    if accountId != nil {
        query = query.Where("account_id = ?", *accountId)
    }
    
    result := query.Find(&todos)
        
    if result.Error != nil {
        return nil, fmt.Errorf("failed to find todos by account id: %w", result.Error)
    }
    
    return todos, nil
}

// 4. Update: TODOの更新
// .Save は、主キー（ID）が含まれている構造体を渡すと、そのレコードを丸ごと更新（UPDATE）してくれます。
// 自動的に updated_at も現在時刻に更新されます。
func (r *PostgresTodoRepository) Update(ctx context.Context, todo *model.Todo) error {
    result := r.db.WithContext(ctx).Save(todo)
    if result.Error != nil {
        return fmt.Errorf("failed to update todo: %w", result.Error)
    }
    
    // もし「1件も更新されなかった（対象レコードが存在しなかった）」ことを検知したい場合
    if result.RowsAffected == 0 {
        return fmt.Errorf("todo not found for update")
    }
    
    return nil
}

// 5. Delete: TODOの削除
// .Delete に空のモデルとIDを渡すことで、該当するレコードを削除（DELETE）します。
func (r *PostgresTodoRepository) Delete(ctx context.Context, id uint) error {
    result := r.db.WithContext(ctx).Delete(&model.Todo{}, id)
    if result.Error != nil {
        return fmt.Errorf("failed to delete todo: %w", result.Error)
    }
    
    if result.RowsAffected == 0 {
        return fmt.Errorf("todo not found for delete")
    }
    
    return nil
}
