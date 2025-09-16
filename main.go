package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

func mains() {
	// 1) 読み込みファイルの指定
	viper.SetConfigName("config") // 設定ファイル名（拡張子なし）
	viper.SetConfigType("json")   // ファイル形式を明示（yaml形式）
	viper.AddConfigPath(".")      // 現在のディレクトリを検索パスに追加

	// 2) 読み込み指示
	err := viper.ReadInConfig() // 設定ファイルを読み込む
	if err != nil {
		log.Fatal("環境変数の設定ファイルがありません")
	}

	// 環境変数の取得
	fmt.Println("---------- AllSettings() ---------")
	m1 := viper.AllSettings()
	fmt.Printf("%#v\n", m1)

	// これでも取れる
	fmt.Println("---------- AllKeys() + range ---------")
	for _, k := range viper.AllKeys() {
		fmt.Printf("%#v\n", k)
	}

	// 親情報の取得
	m2 := viper.GetStringMap("server") // map[string]any
	fmt.Printf("%#v\n", m2)

	// サブツリー
}
