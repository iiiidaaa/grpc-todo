# What is this ?
gRPC通信とfirebaseのauthentication,firestoreのテスト環境を用意します

# Usage
## firebaseのセットアップ
firebaseプロジェクトを追加して、以下のことを行なってください
Authentication -> メールパスワード認証を有効化  
Database -> Firestoreを開始  
プロジェクト設定 -> 全般 から`アプリを追加`ボタンを押して、登録後に表示されるスクリプトを`js/test.html`の`{your firebase config}`と置換してください  
プロジェクト設定 -> サービスアカウント からFirebase Admin SDK(go)の新しい設定鍵を生成でjsonファイルをDL  
jsonファイルをDL後、`service-acount-file.json`という名前でtodoディレクトリに配置してください

## gRPCサーバーの立ち上げ
以下のコマンドでgRPCサーバーが開始されます
```bash
docker-compose up -d
docker exec -it app bash
docker-compose exec app sh -c "go run server/*.go"
```
## firebaseユーザーの作成・ログイン
`js/test.html`をブラウザから開き、任意のメアド/パスワードでsign upを実行してください  
firebaseのauthenticationにユーザーが作成されます  
firestoreを操作する用のclientスクリプトを動作させるためには、ログイン後の画面で表示される`accessToken`のところからJWTを引っこ抜いてください  

## firestoreの利用
`client/main.go`の以下を上で取得したIdTokenに置き換えてください
```accessToken = "{youraccesstoken}"```

その後、以下のコマンドを実行することでfirestoreへデータの追加・削除が実施可能です  
(その一環でユーザーを作ったりもします)  
ユーザーを残したい場合などは、`os.Exit(0)`とかを挟んで削除前で止めるなど適宜やってください