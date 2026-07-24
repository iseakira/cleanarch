// package tester はテスト用のヘルパーをまとめたパッケージ。
package tester

import (
	"context" // キャンセルやタイムアウトを伝搬するための context
	"fmt"     // 文字列フォーマット（ポート指定の組み立てに使用）
	"time"    // タイムアウト時間の指定に使用

	"gorm.io/gorm" // ORM 本体（テスト対象に渡す DB ハンドル）

	"github.com/stretchr/testify/suite"                // testify のテストスイート機能（前処理/後処理付きのテスト群）
	"github.com/testcontainers/testcontainers-go"      // Docker コンテナをテストから起動するライブラリ
	"github.com/testcontainers/testcontainers-go/wait" // コンテナ起動完了を待つ条件を定義

	"go-api-arch-clean-template/entity"                 // ドメインのエンティティ定義（マイグレーション対象）
	"go-api-arch-clean-template/infrastructure/database" // DB 接続設定・ファクトリ
	"go-api-arch-clean-template/pkg"                     // 共通ユーティリティ（ポート待機など）
)

// DBMySQLSuite は「実際の MySQL コンテナを立てて行う結合テスト」用のスイート。
// testify の suite.Suite を埋め込むことで、各テスト前後のフックを持てる。
type DBMySQLSuite struct {
	suite.Suite                              // testify スイートの基本機能を継承
	mySQLContainer testcontainers.Container  // 起動した MySQL コンテナへの参照（終了処理で使用）
	ctx            context.Context           // コンテナ操作に使う context
	DB             *gorm.DB                  // テストから使う DB ハンドル
}

// SetupTestContainers は MySQL の Docker コンテナを起動する。
func (suite *DBMySQLSuite) SetupTestContainers() (err error) {
	configs := database.NewConfigMySQL()                             // DB 接続設定（ホスト/ポート/ユーザーなど）を取得
	pkg.WaitForPort(configs.Database, configs.Port, 10*time.Second)  // 指定ポートが使えるようになるまで最大10秒待機
	suite.ctx = context.Background()                                 // 空の context を用意

	// 起動するコンテナの仕様を定義する。
	req := testcontainers.ContainerRequest{
		Image: "mysql:8.2", // 使用する MySQL イメージ
		Env: map[string]string{ // コンテナに渡す環境変数
			"MYSQL_DATABASE":             configs.Database, // 初期作成するDB名
			"MYSQL_USER":                 configs.User,     // 作成するユーザー
			"MYSQL_PASSWORD":             configs.Password, // そのパスワード
			"MYSQL_ALLOW_EMPTY_PASSWORD": "yes",            // root の空パスワードを許可
		},
		ExposedPorts: []string{fmt.Sprintf("%s:3306/tcp", configs.Port)}, // ホスト側ポート→コンテナ3306 のマッピング
		WaitingFor:   wait.ForLog("port: 3306  MySQL Community Server"),   // このログが出たら起動完了とみなす
	}

	// 仕様に従ってコンテナを生成・起動する。
	suite.mySQLContainer, err = testcontainers.GenericContainer(suite.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true, // 生成と同時に起動する
	})

	suite.Assert().Nil(err) // 起動でエラーが出ていないことを確認
	return nil
}

// SetupSuite はスイート全体の開始前に一度だけ呼ばれる（testify のフック）。
func (suite *DBMySQLSuite) SetupSuite() {
	err := suite.SetupTestContainers() // MySQL コンテナを起動
	suite.Assert().Nil(err)            // 失敗していないことを確認

	db, err := database.NewDatabaseSQLFactory(database.InstanceMySQL) // 起動したMySQLへ接続
	suite.Assert().Nil(err)
	suite.DB = db // 以降のテストで使えるよう保持
	for _, model := range entity.NewDomains() { // 全ドメインモデルについて
		err = suite.DB.AutoMigrate(model) // テーブルを自動生成（マイグレーション）
		suite.Assert().Nil(err)
	}
}

// TearDownSuite はスイート全体の終了後に一度だけ呼ばれ、後片付けを行う。
func (suite *DBMySQLSuite) TearDownSuite() {
	if suite.mySQLContainer == nil { // コンテナが起動していなければ
		return // 何もしない
	}
	err := suite.mySQLContainer.Terminate(suite.ctx) // コンテナを停止・破棄
	suite.Assert().Nil(err)
}
