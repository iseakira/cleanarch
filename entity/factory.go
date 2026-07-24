// package entity はアプリのドメイン（業務データ）を表す型をまとめたパッケージ。
package entity

// NewDomains は全ドメインモデルの一覧を返す。
// テストの SetupSuite などがこれを回して AutoMigrate（テーブル自動生成）に使う。
// 新しいエンティティを追加したら、ここにも足すとマイグレーション対象に含まれる。
func NewDomains() []interface{} {
	return []interface{}{&Category{}, &Album{}} // Category と Album のポインタを列挙して返す
}
