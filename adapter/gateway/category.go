// package gateway は「ドメインと外部（DBなど）の橋渡し」を担う層。
package gateway

import (
	"gorm.io/gorm" // DBアクセスに使うORM

	"go-api-arch-clean-template/entity" // 扱うドメイン型（Category）
)

// CategoryRepository は「カテゴリの永続化」に必要な操作を定義したインターフェース（契約）。
// 使う側はこの抽象にだけ依存し、DBの具体実装を知らなくてよい。
type CategoryRepository interface {
	GetOrCreate(category *entity.Category) (*entity.Category, error) // あれば取得・なければ作成
}

// categoryRepository は CategoryRepository の実装（小文字＝パッケージ外に非公開）。
type categoryRepository struct {
	db *gorm.DB // 実際のDB接続を保持
}

// NewCategoryRepository はDB接続を受け取り、リポジトリ実装を生成して返すコンストラクタ。
// 戻り値の型はインターフェースなので、利用側は実装詳細に縛られない（差し替え可能）。
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

// GetOrCreate は渡されたカテゴリを「あれば取得・なければ新規作成」して返す。
func (c *categoryRepository) GetOrCreate(category *entity.Category) (*entity.Category, error) {
	var getOrCreatedCategory entity.Category // 結果を受け取る変数
	// category に一致する行を探し、無ければ作成する（FirstOrCreate）。
	tx := c.db.FirstOrCreate(&getOrCreatedCategory, category)
	if tx.Error != nil { // DBエラーがあれば
		return nil, tx.Error
	}
	return &getOrCreatedCategory, nil // 取得/作成できたカテゴリを返す
}
