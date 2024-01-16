# find-shunkei

`find-shunkei` コマンドを利用することで、LAN内のShunkeiデバイスを検索することができます。

検索には DNS-SD (mDNS) を利用しているため、マルチキャストパケットが到達する環境(同一L2内など)である必要があります。

## Installation
### バイナリをダウンロード

[GitHub Releases](https://github.com/shunkei-jp/find-shunkei/releases) からバイナリをダウンロードし、
パスの通った場所に配置してください。

### go install

Goがインストールされている環境であれば `go install` を用いてインストールすることが可能です。

```sh
go install github.com/shunkei-jp/find-shunkei@latest
```

## Usage

`find-shunkei -h` でヘルプを参照できます。

```sh
find-shunkei # LAN内を検索
find-shunkei -rx # Shunkei VTX受信機のみを検索
find-shunkei -tx # Shunkei VTX送信機のみを検索
find-shunkei -t 5 # タイムアウトを5秒に設定して検索(デフォルトは2秒)
find-shunkei -1 # 最初の一台が見つかったら終了

find-shunkei -ip-only # IPアドレスのみを表示(スクリプトなどでの利用を想定)
find-shunkei -ip-only -rx -1 # 受信機のみを検索し、最初の一台が見つかったら、IPアドレスのみを表示して終了
```

### 終了コード

デバイスの検出に成功した場合、終了コード0にて正常終了します。
デバイスが一台も見つからなかった場合、終了コード9を返します。

エラーが発生した場合は終了コード1を返します。

## Development

開発時には `main.go` を直接実行してください。

```sh
go run main.go
```
