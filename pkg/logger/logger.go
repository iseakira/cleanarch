// package logger はアプリ共通のロガー。内部で Uber の zap をラップし、
// Info/Debug/Warn/... という薄い関数経由でどこからでも同じロガーを使えるようにする。
package logger

import (
	"os" // 環境変数の読み取りに使用

	"go.uber.org/zap" // 高速な構造化ロギングライブラリ
)

var (
	ZapLogger        *zap.Logger        // 生の zap ロガー（型付き・高速なAPI）。外部公開して直接使えるように
	zapSugaredLogger *zap.SugaredLogger // 使い勝手を優先した "Sugar" 版（可変長引数で気軽に書ける）。パッケージ内部用
)

// init はパッケージ読み込み時に自動で一度だけ実行され、ロガーを初期化する。
func init() {
	cfg := zap.NewProductionConfig()   // 本番向け設定（JSON出力・INFO以上）を基本にする
	logFile := os.Getenv("APP_LOG_FILE") // ログ出力先ファイルを環境変数から取得
	if logFile != "" {                   // 指定があれば
		cfg.OutputPaths = []string{"stderr", logFile} // 標準エラーとファイルの両方に出力
	}

	ZapLogger = zap.Must(cfg.Build())              // 設定からロガーを構築（失敗時は panic）
	if os.Getenv("APP_ENV") == "development" {     // 開発環境なら
		ZapLogger = zap.Must(zap.NewDevelopment()) // 人間が読みやすい開発用ロガーに差し替え
	}
	zapSugaredLogger = ZapLogger.Sugar()           // Sugar 版を生成し、以降の各関数で使う
}

// Sync はバッファに溜まったログを書き切る。プログラム終了前に呼ぶのが定石。
func Sync() {
	err := zapSugaredLogger.Sync()
	if err != nil {
		zap.Error(err) // ※ 実際にはログ出力していない（下部の注記参照）
	}
}

// Info は情報レベルのログを出力する。keysAndValues は key, value, key, value... の構造化フィールド。
func Info(msg string, keysAndValues ...interface{}) {
	zapSugaredLogger.Infow(msg, keysAndValues...)
}

// Debug はデバッグレベルのログを出力する（本番設定では通常出ない）。
func Debug(msg string, keysAndValues ...interface{}) {
	zapSugaredLogger.Debugw(msg, keysAndValues...)
}

// Warn は警告レベルのログを出力する。
func Warn(msg string, keysAndValues ...interface{}) {
	zapSugaredLogger.Warnw(msg, keysAndValues...)
}

// Error はエラーレベルのログを出力する。
func Error(msg string, keysAndValues ...interface{}) {
	zapSugaredLogger.Errorw(msg, keysAndValues...)
}

// Fatal はエラーを出力した後、os.Exit(1) でプロセスを終了させる。
func Fatal(msg string, keysAndValues ...interface{}) {
	zapSugaredLogger.Fatalw(msg, keysAndValues...)
}

// Panic はエラーを出力した後、panic を発生させる。
func Panic(msg string, keysAndValues ...interface{}) {
	zapSugaredLogger.Panicw(msg, keysAndValues...)
}
