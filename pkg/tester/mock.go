// package tester はテスト用のモック（偽物）ヘルパーをまとめたパッケージ。
package tester

import (
	"time" // 時刻を扱う標準ライブラリ（mockClock で使用）

	"github.com/DATA-DOG/go-sqlmock" // 実DBに接続せずSQLをシミュレートするモックライブラリ
	"gorm.io/driver/mysql"           // GORM 用の MySQL ドライバ（モックDBを MySQL として見せる）
	"gorm.io/gorm"                   // ORM 本体

	"go-api-arch-clean-template/pkg/logger" // プロジェクト共通のロガー（致命的エラー出力に使用）
)

// MockDB はテスト用のモックDBを生成して返す。
// mock       : 「このSQLが来たらこの結果を返す」と期待値を設定する操作用オブジェクト
// mockGormDB : テスト対象コードに渡す *gorm.DB（実DBの代わり）
func MockDB() (mock sqlmock.Sqlmock, mockGormDB *gorm.DB) {
	// sqlmock でモックDBを生成する。
	// QueryMatcherRegexp を指定し、期待SQLを「正規表現の部分マッチ」で照合できるようにする
	// （GORM が生成する複雑なSQLと完全一致させるのは困難なため）。
	mockDB, mock, err := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))

	if err != nil { // モックDB生成に失敗したら
		logger.Fatal(err.Error()) // ログを出してプロセスを停止
	}

	// 生成したモックDBを GORM から使えるようにラップする。
	mockGormDB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       "mock_db", // 接続文字列（モックなのでダミー値）
		DriverName:                "mysql",   // ドライバ名
		Conn:                      mockDB,    // 実接続の代わりにモックDBを差し込む
		SkipInitializeWithVersion: true,      // 起動時の SELECT VERSION() を抑止（モックでは取得不可のため必須）
	}), &gorm.Config{})

	if err != nil { // GORM 初期化に失敗したら
		logger.Fatal(err.Error()) // ログを出してプロセスを停止
	}

	return mock, mockGormDB // 期待値設定用の mock と、テスト対象に渡す DB を返す
}

// mockClock は「常に固定時刻を返す時計」。
// time.Now() に依存する処理をテストで再現可能にするために使う。
type mockClock struct {
	t time.Time // 返し続ける固定の時刻
}

// NewMockClock は指定した時刻 t に固定された mockClock を生成する。
func NewMockClock(t time.Time) mockClock {
	return mockClock{t}
}

// Now は現在時刻の代わりに、生成時に固定した時刻を常に返す。
func (m mockClock) Now() time.Time {
	return m.t
}
