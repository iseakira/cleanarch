// package tester はテスト用のヘルパーをまとめたパッケージ。
package tester

import (
	"fmt" // 文字列フォーマット（DBファイル名の組み立てに使用）
	"os"  // 環境変数の操作に使用

	"github.com/stretchr/testify/suite" // testify のテストスイート機能（前処理付きのテスト群）
	"gorm.io/gorm"                       // ORM 本体（テスト対象に渡す DB ハンドル）

	"go-api-arch-clean-template/entity"                  // ドメインのエンティティ定義（マイグレーション対象）
	"go-api-arch-clean-template/infrastructure/database" // DB 接続設定・ファクトリ
)

// DBSQLiteSuite は「SQLite を使って行うテスト」用のスイート。
// MySQL コンテナ版（DBMySQLSuite）と違い、Docker 不要で軽量に実行できる。
type DBSQLiteSuite struct {
	suite.Suite         // testify スイートの基本機能を継承
	DB     *gorm.DB     // テストから使う DB ハンドル
	DBName string       // 使用する SQLite の DB（ファイル）名
}

// SetupSuite はスイート全体の開始前に一度だけ呼ばれる（testify のフック）。
func (suite *DBSQLiteSuite) SetupSuite() {

	// テスト名を元に一意な SQLite ファイル名を作る（テスト間でDBを分離するため）。
	suite.DBName = fmt.Sprintf("%s.unittest.sqlite", suite.T().Name())
	os.Setenv("DB_NAME", suite.DBName)                                // ファクトリがこの環境変数を読んで接続先を決める
	db, err := database.NewDatabaseSQLFactory(database.InstanceSQLite) // SQLite インスタンスへ接続
	suite.Assert().Nil(err)                                           // 接続でエラーが出ていないことを確認
	suite.DB = db                                                     // 以降のテストで使えるよう保持

	for _, model := range entity.NewDomains() { // 全ドメインモデルについて
		err := suite.DB.AutoMigrate(model) // テーブルを自動生成（マイグレーション）
		suite.Assert().Nil(err)
	}

}
