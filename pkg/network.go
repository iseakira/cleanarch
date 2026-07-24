// package pkg は共通ユーティリティをまとめたパッケージ。
package pkg

import (
	"net"     // TCP 接続でポートの空き状況を調べる
	"net/url" // URL の解析・結合
	"os"      // 環境変数の読み取り
	"time"    // タイムアウト・待機時間の扱い
)

// CheckPort は host:port に TCP 接続を試み、「ポートが空いているか」を返す。
// 戻り値 true  … 接続できない＝まだ誰も使っていない（空き）
// 戻り値 false … 接続できた＝すでに何かが待ち受けている（使用中）
func CheckPort(host string, port string) bool {
	conn, err := net.Dial("tcp", net.JoinHostPort(host, port)) // 実際に接続を試す（IPv6も安全に扱える）

	if conn != nil { // 接続できた＝誰かが待ち受けている
		conn.Close()  // 開いた接続は閉じる
		return false  // 使用中なので「空きではない」
	}

	if err != nil { // 接続に失敗した＝待ち受けがいない
		return true // 「空き」とみなす
	}
	return false
}

// WaitForPort は指定ポートが空くまで、timeout を上限に1秒間隔で待ち続ける。
// 空きになれば true、時間切れなら false を返す。
func WaitForPort(host string, port string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout) // 待機の締め切り時刻を決める
	for time.Now().Before(deadline) {   // 締め切りまで繰り返す
		if CheckPort(host, port) { // ポートが空いたら
			return true // 成功として返す
		}
		time.Sleep(1 * time.Second) // 1秒待って再確認
	}
	return false // 時間切れ
}

// GetEndpoint は与えられた path を、環境に応じたベースURLと結合して絶対URLを返す。
func GetEndpoint(path string) string {
	var baseURL string
	baseURL = "http://0.0.0.0:8080/"          // 既定（ローカル）のベースURL
	env := os.Getenv("APP_ENV")               // 実行環境を環境変数から取得
	if env == "stage" {                       // ステージング環境なら
		baseURL = "http://stage.localhost:8080/" // ベースURLを差し替える
	}
	p, _ := url.Parse(path)                    // path を URL として解析
	b, _ := url.Parse(baseURL)                 // baseURL を URL として解析
	return b.ResolveReference(p).String()      // baseURL を基準に path を解決した絶対URLを返す
}
