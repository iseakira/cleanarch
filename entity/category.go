// package entity はアプリのドメイン（業務データ）を表す型をまとめたパッケージ。
package entity

import (
	"errors" // エラー値を作るための標準ライブラリ
)

// 取りうるカテゴリ名を定数として定義する（OpenAPI の enum: food/music/sports に対応）。
const (
	Food   CategoryName = "food"   // 食べ物
	Music  CategoryName = "music"  // 音楽
	Sports CategoryName = "sports" // スポーツ
)

// CategoryName はカテゴリ名を表す独自の文字列型。
// ただの string ではなく専用型にすることで、決められた値以外が入るのを防ぎやすくする。
type CategoryName string

// NewCategoryName は文字列を検証しつつ CategoryName を生成するコンストラクタ。
// 不正な値なら error を返し、正しければ *CategoryName を返す。
func NewCategoryName(value string) (*CategoryName, error) {
	var categoryName CategoryName            // 空の CategoryName を用意
	if err := categoryName.Set(value); err != nil { // 検証しつつ値をセット
		return nil, err // 不正ならエラーを返す
	}
	return &categoryName, nil // 正常なら生成した値のポインタを返す
}

// IsValid は現在の値が許可された3つのいずれかかどうかを判定する。
func (c *CategoryName) IsValid() bool {
	return *c == Food || *c == Music || *c == Sports // 定義済み定数のどれかに一致すればtrue
}

// Set は与えられた文字列を検証し、妥当ならレシーバに書き込む。
func (c *CategoryName) Set(value string) error {
	newCategoryName := CategoryName(value)  // 文字列を CategoryName 型に変換
	if !newCategoryName.IsValid() {         // 許可された値でなければ
		return errors.New("Invalid value for CategoryName") // エラーを返す
	}
	*c = newCategoryName // 妥当なのでレシーバの値を更新
	return nil
}

// Category はアルバムの分類を表すエンティティ。
type Category struct {
	ID   int          // カテゴリの識別子（DBの主キーに対応）
	Name CategoryName // カテゴリ名（検証済みの CategoryName が入る）
}

// NewCategory は名前の文字列から Category を生成するコンストラクタ。
// 名前の妥当性チェックは NewCategoryName に委ねる。
func NewCategory(name string) (*Category, error) {
	categoryName, err := NewCategoryName(name) // 名前を検証して CategoryName を作る
	if err != nil {
		return nil, err // 不正な名前ならエラーを返す
	}
	return &Category{
		Name: *categoryName, // 検証済みの名前を持つ Category を返す（ID は未設定）
	}, nil
}
