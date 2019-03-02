# sample-cosmos-app
tendermint上でcosmos-sdkを使って従来のDNSシステムをモデル化したsample applicationを作りました。
ユーザーは未使用の名前を購入したり、自分の名前を売買することができます。

参考:https://cosmos.network/docs/tutorial/

Tendermintの概要は、[こちら](./docs/tendermint-overview.md)を参照してください。

**ユーザーができること**
- 好きな文字列(ドメイン)の購入(buy-name)
- 他のユーザーが所有する文字列の購入(buy-name)
- ドメインに対して、好きな文字列(IP)の紐づけ(set-name)

## Getting started
```
$ git clone https://github.com/EG-easy/sample-cosmos-app.git
$ cd sample-cosmos-app
```

### Local Environment on Mac OS
**requirement**
Golang 

```
$ make get_tools && make get_vendor_deps
$ make install
```

### Using Docker
```
$ make start
```

### Demo
ユーザーを2人作って、nameの売買およびそれに紐づくIPの設定を行う。

参考:https://cosmos.network/docs/tutorial/build-run.html#running-the-live-network-and-using-the-commands


設定をする
```
$ sh scripts/start.sh #シェルスクリプトを実行すると、2人分のpasswordを設定することが求められるので、適宜設定する。
$ nsd start #applicationのdeamonが起動する
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

ドメインの所有権を確認する
```
nscli query nameservice whois jack.id
```

別のユーザーであるaliceがそのドメインを10mycoinsで購入する
```
nscli tx nameservice buy-name jack.id 10mycoin --from alice 
```


## 今後の展望
このapplicationはあくまで、cosmos-sdkの動作確認用であり、コンセンサスはtendermintに部分を預けているだけなので、validatorの数を増やしたり、手数料の設定などを独自の設定を追加していきたい。
