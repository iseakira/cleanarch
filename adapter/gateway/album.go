// package gateway は「ドメインと外部（DBなど）の橋渡し」を担う層。
// ここではDBアクセス（リポジトリ）の実装を置く。
package gateway

import (
	"github.com/jinzhu/copier" // 構造体から構造体へフィールドをコピーするライブラリ
	"gorm.io/gorm"             // DBアクセスに使うORM

	"go-api-arch-clean-template/entity" // 扱うドメイン型（Album, Category）
)

// AlbumRepository は「Albumの永続化」に必要な操作を定義したインターフェース（契約）。
// 使う側はこの抽象にだけ依存し、DBの具体実装を知らなくてよい。
type AlbumRepository interface {
	Create(album *entity.Album) (*entity.Album, error) // 新規作成
	Get(ID int) (*entity.Album, error)                 // ID で1件取得
	Save(*entity.Album) (*entity.Album, error)         // 更新（保存）
	Delete(ID int) error                               // 削除
}

// albumRepository は AlbumRepository の実装（小文字＝パッケージ外に非公開）。
type albumRepository struct {
	db *gorm.DB // 実際のDB接続を保持
}

// NewAlbumRepository はDB接続を受け取り、リポジトリ実装を生成して返すコンストラクタ。
// 戻り値の型はインターフェースなので、利用側は実装詳細に縛られない。
func NewAlbumRepository(db *gorm.DB) AlbumRepository {
	return &albumRepository{db: db}
}

// GetOrCreateCategory はアルバムのカテゴリを「あれば取得・なければ作成」し、
// そのIDと実体をアルバムに紐づける補助メソッド。
func (a *albumRepository) GetOrCreateCategory(album *entity.Album) error {
	var category entity.Category // 結果を受け取る変数
	// 同名のカテゴリを探し、無ければ作成する（FirstOrCreate）。
	tx := a.db.FirstOrCreate(&category, entity.Category{Name: album.Category.Name})
	if tx.Error != nil { // DBエラーがあれば
		return tx.Error
	}
	album.CategoryID = category.ID // 取得/作成したカテゴリのIDを外部キーにセット
	album.Category = category      // 実体もセット
	return nil
}

// Create は新しいアルバムをDBに作成する。
func (a *albumRepository) Create(album *entity.Album) (*entity.Album, error) {
	if err := a.GetOrCreateCategory(album); err != nil { // 先にカテゴリを用意
		return nil, err
	}
	if err := a.db.Create(album).Error; err != nil { // アルバム本体を挿入
		return nil, err
	}
	return album, nil // 作成できたアルバムを返す
}

// Get は ID を指定してアルバムを1件取得する。
func (a *albumRepository) Get(ID int) (*entity.Album, error) {
	var album = entity.Album{} // 受け取る器
	// Preload("Category") で紐づくカテゴリも同時に読み込む（N+1を避ける）。
	if err := a.db.Preload("Category").First(&album, ID).Error; err != nil {
		return nil, err // 見つからない/エラーなら nil
	}
	return &album, nil
}

// Save は既存のアルバムを更新する。
func (a *albumRepository) Save(album *entity.Album) (*entity.Album, error) {
	selectedAlbum, err := a.Get(album.ID) // まず現在のレコードを取得
	if err != nil {
		return nil, err
	}

	if err := a.GetOrCreateCategory(album); err != nil { // カテゴリを用意
		return nil, err
	}

	// time.Time や types.Date 以外のフィールドをコピーする。
	// IgnoreEmpty: 空の項目は上書きしない / DeepCopy: 深いコピー。
	if err := copier.CopyWithOption(selectedAlbum, album, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}
	if err := a.db.Save(&selectedAlbum).Error; err != nil { // 変更を保存
		return nil, err
	}

	return selectedAlbum, nil
}

// Delete は ID を指定してアルバムを削除する。
func (a *albumRepository) Delete(ID int) error {
	album := entity.Album{ID: ID} // 削除対象を表す値
	if err := a.db.Where("id = ?", &album.ID).Delete(&album).Error; err != nil {
		return err
	}
	return nil
}
