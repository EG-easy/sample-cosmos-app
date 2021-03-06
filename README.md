# sample-cosmos-app
tendermint上でcosmos-sdkを使って従来のDNSシステムをモデル化したsample applicationです。
ユーザーは未使用の名前を購入したり、自分の名前を売買することができます。

参考:https://cosmos.network/docs/tutorial/

Tendermintの概要は、[こちら](./docs/tendermint-overview.md)を参照してください。

**ユーザーができること**
- 好きな文字列(ドメイン)の購入(buy-name)
- 他のユーザーが所有する文字列の購入(buy-name)
- ドメインに対して、好きな文字列(IP)の紐づけ(set-name)
- ドメインの名前解決(resolve)
- ドメインの所有権の確認(whois)

## Getting started

### Local Environment on Mac OS
**requirement**
Golang 

```
$ mkdir -p $GOPATH/src/github.com/EG-easy
$ cd $GOPATH/src/github.com/EG-easy
$ git clone https://github.com/EG-easy/sample-cosmos-app.git
$ cd sample-cosmos-app
$ make get_tools && make get_vendor_deps
$ make install
```

### Using Docker
```
$ git clone https://github.com/EG-easy/sample-cosmos-app.git
$ cd sample-cosmos-app
$ make build
$ make start
# コンテナのなかで
$ make install 
```

ちなみに起動しているコンテナの中には、`docker exec -it sample-cosmos-app /bin/sh`で入れる。

## Demo
ユーザーを2人作って、nameの売買およびそれに紐づくIPの設定を行う。

参考:https://cosmos.network/docs/tutorial/build-run.html#running-the-live-network-and-using-the-commands

まずは、ユーザーの設定をして、アプリケーションを開始する。
```
# シェルスクリプトを実行すると、2人分のpasswordを設定することが求められるので、適宜設定する。
$ sh scripts/start.sh
# applicationのdeamonが起動する
$ nsd start 
```

ここで、別のwindowを開いて、次のコマンドで残高を確認する
```
$ nscli query account $(nscli keys show jack -a) 
$ nscli query account $(nscli keys show alice -a)
```
次に、jackとして、jack.idという文字列(ドメイン)を購入する
```
$ nscli tx nameservice buy-name jack.id 5mycoin --from jack
```

そのドメインに対して、文字列(IP)を紐づける
```
$ nscli tx nameservice set-name jack.id 8.8.8.8 --from jack 
```

ドメインの名前解決をする
```
$ nscli query nameservice resolve jack.id
```

ドメインの所有権を確認する
```
$nscli query nameservice whois jack.id
```

別のユーザーであるaliceがそのドメインをjackより高値で購入する
```
$ nscli tx nameservice buy-name jack.id 10mycoin --from alice 
```

## 今後の展望
このapplicationはあくまで、cosmos-sdkの動作確認用であり、コンセンサスはtendermintに部分を預けているだけなので、validatorの数を増やしたり、手数料の設定などを独自の設定を追加していきたい。
