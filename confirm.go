// go run main.go [-flag など] で実行。
// 実験の手順はこの後に記載。

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func dump(label string) {
	fmt.Println("===", label, "===")
	fmt.Printf("app.port (GetInt) = %d\n", viper.GetInt("app.port"))
	all := viper.AllSettings()
	b, _ := json.MarshalIndent(all, "", "  ")
	fmt.Println(string(b))
	fmt.Println()
}

func main() {
	// ---- 0) Flag（標準flagで例示。pflag/Cobraでも同様の概念）
	var portFlag = flag.Int("port", 0, "port via flag (overrides env/config/default)")
	flag.Parse()

	// ---- 1) Default（最下位）
	viper.SetDefault("app.port", 80000)

	// ---- 2) Config（設定ファイル層：config.yaml を読む）
	// 例: config.yaml に `app: { port: 9000 }`
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		// 設定ファイルが無ければ続行（環境変数やデフォルトで動く）
		var notFound viper.ConfigFileNotFoundError
		if !os.IsNotExist(err) && !(errorAs(err, &notFound)) {
			log.Printf("read config warning: %v", err)
		}
	}

	// ---- 3) Env（環境変数層）
	//    AutomaticEnv: Get 時に環境変数を毎回照会
	//    KeyReplacer:  app.port -> APP_PORT のような変換
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// ---- 4) Flag（フラグ層）
	// viper.BindPFlag を使うのが定石だが、標準flag利用時は手動反映でも挙動確認できる
	if *portFlag != 0 {
		viper.Set("app.port", *portFlag) // ← Set は最上位だが、実験の便宜上ここで反映
		// ※BindPFlagを使う場合は Set ではなく Bind で “Flag 層” として載せる
	}

	// ---- 5) Set（最上位層：コードからの明示 Set）
	// 実験したい場合は以下をコメントイン
	// viper.Set("app.port", 12000)

	// ダンプ
	dump("effective")
}

// Go 1.20未満でも動くよう errorAs を薄く実装（簡易）
func errorAs(err error, target *viper.ConfigFileNotFoundError) bool {
	_, ok := err.(viper.ConfigFileNotFoundError)
	return ok
}
