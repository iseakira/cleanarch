// package entity_test は entity パッケージ「外」からテストするための専用パッケージ。
// （_test を付けると、実際の利用者と同じく公開APIだけを使ってテストできる）
package entity_test

import (
	"testing" // Goの標準テスト機能（*testing.T など）

	"github.com/stretchr/testify/assert" // 検証を読みやすく書けるアサーションライブラリ

	"go-api-arch-clean-template/entity"      // テスト対象（Album, Category など）
	"go-api-arch-clean-template/pkg"         // 時刻変換などの共通ユーティリティ
	"go-api-arch-clean-template/pkg/tester"  // テスト補助（固定時刻の時計 NewMockClock）
)

// TestAlbum は「Album を組み立てると各フィールドが期待通りに入るか」を確認するテスト。
func TestAlbum(t *testing.T) {
	// --- Arrange（準備）: テスト用のデータを用意 ---
	category := entity.Category{ // 紐づけるカテゴリを作る
		ID:   1,
		Name: "sports",
	}

	now := pkg.Str2time("2023-01-01")     // 文字列を time.Time に変換して「基準の今日」を作る
	mockClock := tester.NewMockClock(now) // 常に 2023-01-01 を返す固定時計（周年計算に渡す）
	album := entity.Album{                // 検証対象の Album を組み立てる
		ID:          1,
		Title:       "album",
		ReleaseDate: now,
		CategoryID:  1,
		Category:    category,
	}

	// --- Assert（検証）: 各フィールドが期待通りかを1つずつ確認 ---
	assert.Equal(t, 1, album.ID)                                  // ID は 1
	assert.Equal(t, 0, album.Anniversary(mockClock))             // 今日=リリース日 なので 0 周年
	assert.Equal(t, "album", album.Title)                        // タイトル
	assert.Equal(t, now, album.ReleaseDate)                      // リリース日
	assert.Equal(t, 1, album.CategoryID)                         // 外部キー
	assert.Equal(t, 1, album.Category.ID)                        // 紐づくカテゴリのID
	assert.Equal(t, "sports", string(album.Category.Name))       // カテゴリ名（型変換して文字列比較）
}

// TestAlbumAnniversary は周年計算ロジックを、境界となる日付で集中的に検証するテスト。
func TestAlbumAnniversary(t *testing.T) {
	mockedClock := tester.NewMockClock(pkg.Str2time("2022-04-01")) // 「今日」を 2022-04-01 に固定

	// non-leap（うるう年をまたがないケース）
	album := entity.Album{ReleaseDate: pkg.Str2time("2022-04-01")} // 今日と同じ日
	assert.Equal(t, 0, album.Anniversary(mockedClock))             // → 0 周年

	album = entity.Album{ReleaseDate: pkg.Str2time("2021-04-02")}  // 記念日の前日（まだ丸1年経っていない）
	assert.Equal(t, 0, album.Anniversary(mockedClock))             // → 0 周年

	album = entity.Album{ReleaseDate: pkg.Str2time("2021-04-01")}  // ちょうど1年前
	assert.Equal(t, 1, album.Anniversary(mockedClock))             // → 1 周年

	// leap（うるう年 2020 をまたぐケース。補正が効くか確認）
	album = entity.Album{ReleaseDate: pkg.Str2time("2020-04-02")}  // 記念日前日相当だが…
	assert.Equal(t, 1, album.Anniversary(mockedClock))             // → 1 周年（うるう年補正が正しい）

	album = entity.Album{ReleaseDate: pkg.Str2time("2020-04-01")}  // ちょうど2年前
	assert.Equal(t, 2, album.Anniversary(mockedClock))             // → 2 周年
}
