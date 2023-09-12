# net_watch

net_watchは、お気に入りのNetflixの番組や映画を管理・追跡するためのツールです。このプロジェクトには、APIサーバーとNetflixのタイトル検索のためのCLIツールの2つの部分が含まれています。


## 機能

- タイトル一覧の取得
- シーズン一覧の取得
- エピソード一覧の取得
- エピソードの更新
- CLIツールを使用したNetflixタイトルの検索と保存

## セットアップ & インストール

### 1. リポジトリをクローン:

```bash
git clone https://github.com/monjofight/net_watch
cd net_watch
```

### 2. 必要な依存関係をインストール:

```bash
go mod download
```

## アプリケーションの実行

### 環境変数の設定:

`.env.example` ファイルを `.env` にコピーして、必要な環境変数を設定してください。

### アプリケーションを実行:

```bash
go run cmd/server/main.go
```

### CLIツールを実行:

```bash
go run cmd/cli/main.go
```

### ローカルでUpdate関数を実行:

```bash
export FUNCTION_TARGET=Update
go run cmd/main.go
curl localhost:8080
```

## APIエンドポイント

- `GET /titles`: すべてのタイトルを取得します。
- `POST /titles/:titleId/watch`: タイトルのすべてのエピソードを視聴済みとしてマークします。
- `POST /titles/:titleId/unwatch`: タイトルのすべてのエピソードの視聴マークを解除します。
- `GET /titles/:titleId/seasons`: タイトルに関連するシーズンの一覧を取得します。
- `POST /titles/:titleId/seasons/:seasonId/watch`: 指定されたシーズンのすべてのエピソードを視聴済みとしてマークします。
- `POST /titles/:titleId/seasons/:seasonId/unwatch`: 指定されたシーズンのすべてのエピソードの視聴マークを解除します。
- `GET /titles/:titleId/seasons/:seasonId/episodes`: 指定されたシーズンのすべてのエピソードを取得します。
- `POST /titles/:titleId/seasons/:seasonId/episodes/:episodeId/watch`: 指定されたエピソードを視聴済みとしてマークします。
- `POST /titles/:titleId/seasons/:seasonId/episodes/:episodeId/unwatch`: 指定されたエピソードの視聴マークを解除します。
- `POST /update`: すべてのタイトルのエピソードを更新します。