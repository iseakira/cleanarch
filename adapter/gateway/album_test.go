// package gateway_test は gateway パッケージを「外部から」テストする専用パッケージ。
package gateway_test

import (
	"errors"  // テスト用のエラー値を作る
	"regexp"  // 期待SQLを正規表現用にエスケープ（QuoteMeta）
	"strings" // 文字列の部分一致判定
	"testing" // Goの標準テスト機能
	"time"    // 時刻（time.Now など）

	"github.com/DATA-DOG/go-sqlmock"    // DBを使わずSQLをシミュレートするモック
	"github.com/stretchr/testify/suite" // 前処理/後処理付きのテストスイート

	"go-api-arch-clean-template/adapter/gateway" // テスト対象（AlbumRepository）
	"go-api-arch-clean-template/entity"          // ドメイン型（Album, Category）
	"go-api-arch-clean-template/pkg"             // 時刻変換などの共通ユーティリティ
	"go-api-arch-clean-template/pkg/tester"      // テスト補助（SQLiteスイート・MockDB）
)

// AlbumRepositorySuite はリポジトリのテストをまとめるスイート。
// tester.DBSQLiteSuite を埋め込むことで、本物のSQLiteを使う準備（SetupSuite等）を引き継ぐ。
type AlbumRepositorySuite struct {
	tester.DBSQLiteSuite                          // SQLite接続やマイグレーションの機能を継承
	repository           gateway.AlbumRepository  // テスト対象のリポジトリ（契約型で保持）
}

// TestAlbumRepositorySuite は go test の入口。suite.Run でスイート内の各テストを実行する。
func TestAlbumRepositorySuite(t *testing.T) {
	suite.Run(t, new(AlbumRepositorySuite))
}

// SetupSuite はスイート開始前に一度だけ呼ばれる。
func (suite *AlbumRepositorySuite) SetupSuite() {
	suite.DBSQLiteSuite.SetupSuite()                          // 親のSetupSuiteでSQLiteを用意＆マイグレーション
	suite.repository = gateway.NewAlbumRepository(suite.DB)   // 本物のSQLite接続でリポジトリを生成
}

// MockDB はリポジトリを「モックDB版」に差し替え、期待値設定用のオブジェクトを返す。
// エラー系のテストで使う（本物DBではエラーを再現しにくいため）。
func (suite *AlbumRepositorySuite) MockDB() sqlmock.Sqlmock {
	mock, mockGormDB := tester.MockDB()                        // モックDBを生成
	suite.repository = gateway.NewAlbumRepository(mockGormDB)  // リポジトリをモック接続で作り直す
	return mock                                               // 「このSQLでこう返す」を仕込む側を返す
}

// AfterTest は各テストの後に呼ばれる。MockDBで差し替えた状態を元の本物DBに戻す。
func (suite *AlbumRepositorySuite) AfterTest(suiteName, testName string) {
	suite.repository = gateway.NewAlbumRepository(suite.DB)
}

// TestAlbumRepositoryCRUD は本物のSQLiteを使い、作成→取得→更新→削除の一連を検証する。
func (suite *AlbumRepositorySuite) TestAlbumRepositoryCRUD() {
	now := pkg.Str2time("2023-01-01") // 基準となる日付
	album := &entity.Album{           // 登録するアルバムを用意
		Title:       "test",
		ReleaseDate: now,
		Category:    entity.Category{Name: entity.CategoryName("sports")},
	}
	// --- Create（作成）---
	album, err := suite.repository.Create(album)
	suite.Assert().Nil(err)                                    // エラーなし
	suite.Assert().NotZero(album.ID)                           // IDが採番されている
	suite.Assert().Equal("test", album.Title)                 // タイトル一致
	suite.Assert().Equal(now, album.ReleaseDate)              // 日付一致
	suite.Assert().NotZero(album.Category.ID)                 // カテゴリも作成されID付与
	suite.Assert().Equal("sports", string(album.Category.Name))

	// --- Get（取得）---
	getAlbum, err := suite.repository.Get(album.ID)
	suite.Assert().Nil(err)
	suite.Assert().Equal("test", getAlbum.Title)
	suite.Assert().Equal(now, album.ReleaseDate)
	suite.Assert().Equal(album.Category.ID, getAlbum.Category.ID) // Preloadでカテゴリも取れている
	suite.Assert().Equal("sports", string(getAlbum.Category.Name))

	// --- Save（更新）---
	getAlbum.Title = "updated"                                    // タイトルを変更
	updatedAlbum, err := suite.repository.Save(getAlbum)
	suite.Assert().Nil(err)
	suite.Assert().Equal("updated", updatedAlbum.Title)          // 更新が反映されている
	suite.Assert().NotNil(updatedAlbum.ReleaseDate)
	suite.Assert().NotNil(updatedAlbum.Category.ID)
	suite.Assert().Equal("sports", string(updatedAlbum.Category.Name))

	// --- Delete（削除）---
	err = suite.repository.Delete(updatedAlbum.ID)
	suite.Assert().Nil(err)
	deletedAlbum, err := suite.repository.Get(updatedAlbum.ID)   // 削除後に取得を試みる
	suite.Assert().Nil(deletedAlbum)                            // もう取れない
	suite.Assert().True(strings.Contains("record not found", err.Error())) // 「見つからない」エラー
}

// TestAlbumCreateFailure は Create 時のDBエラーを、モックDBで再現して検証する。
func (suite *AlbumRepositorySuite) TestAlbumCreateFailure() {
	mockDB := suite.MockDB() // リポジトリをモックに差し替え
	// カテゴリ検索クエリが来たら、わざとエラーを返すよう仕込む（QuoteMetaで正規表現エスケープ）。
	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` WHERE `categories`.`name` = ? ORDER BY `categories`.`id` LIMIT 1")).WithArgs("sports").WillReturnError(errors.New("create error"))

	album := &entity.Album{
		Title:       "test",
		ReleaseDate: time.Now(),
		Category:    entity.Category{Name: entity.CategoryName("sports")},
	}

	createdAlbum, err := suite.repository.Create(album)
	suite.Assert().Nil(createdAlbum)                    // 失敗したので結果はnil
	suite.Assert().NotNil(err)                          // エラーが返る
	suite.Assert().Equal("create error", err.Error())  // 仕込んだエラーが伝播している
}

// TestAlbumDeleteFailure は Delete 時のDBエラーを、モックDBで再現して検証する。
func (suite *AlbumRepositorySuite) TestAlbumDeleteFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectBegin() // トランザクション開始を期待
	// DELETE文が来たらエラーを返すよう仕込む。
	mockDB.ExpectExec(regexp.QuoteMeta("DELETE FROM `albums` WHERE id = ? AND `albums`.`id` = ?")).WithArgs(1, 1).WillReturnError(errors.New("delete error"))
	mockDB.ExpectRollback() // エラーなのでロールバックを期待
	mockDB.ExpectCommit()

	err := suite.repository.Delete(1)
	suite.Assert().NotNil(err)                          // エラーが返る
	suite.Assert().Equal("delete error", err.Error())
}

// TestAlbumSaveFailure は Save の最初の取得（SELECT）でのDBエラーを、モックDBで再現して検証する。
func (suite *AlbumRepositorySuite) TestAlbumSaveFailure() {
	mockDB := suite.MockDB()
	// Saveは内部でまず対象を取得する。そのSELECTでエラーを返すよう仕込む。
	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `albums` WHERE `albums`.`id` = ? ORDER BY `albums`.`id` LIMIT 1")).WithArgs(1).WillReturnError(errors.New("save error"))

	album := &entity.Album{
		ID:       1,
		Title:    "test",
		Category: entity.Category{Name: entity.CategoryName("sports")},
	}

	album, err := suite.repository.Save(album)
	suite.Assert().Nil(album)                        // 失敗したのでnil
	suite.Assert().NotNil(err)                       // エラーが返る
	suite.Assert().Equal("save error", err.Error())
}
