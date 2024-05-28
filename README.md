# 概要
Google OAuth2.0 の調査結果。

# シーケンス
OAuth2.0を使用してGoogleの認証を行う際のシーケンスは、以下のような手順で行われる。
1. クライアントの登録
  <br>最初に、Google Developers Consoleでアプリケーションを登録し、クライアントIDとクライアントシークレットを取得する。

2. ユーザーの認証の開始（リクエスト）
  <br>クライアントアプリケーションは、Googleの認証サーバーに対して認証リクエストを送信する。リクエストには以下の情報が含まれる。
    1. 必要なスコープ（アクセス権限）の指定
    2. クライアントID
    3. リダイレクトURI（認証後にユーザーをリダイレクトするURI）（※ユーザーのログインの数だけURLが違う）
3. ユーザーの認証ページへのリダイレクト（GoogleのUI）
  <br>Googleの認証サーバーは、認証ページへのリダイレクトを返す。ユーザーはGoogleの認証ページでログインし、アプリケーションにアクセスを許可する。
4. 認証コードの取得（レスポンス）
  <br>ユーザーがアクセスを許可すると、Googleの認証サーバーは指定されたリダイレクトURIに認証コードを返す。この認証コードは、アクセストークンを取得するために使用される。
5. アクセストークンの取得（リクエスト、レスポンス）
  <br>クライアントアプリケーションは、取得した認証コードとクライアントシークレットを使用して、アクセストークンをGoogleのトークンエンドポイントから要求する。
6. アクセストークンの利用
  <br>クライアントアプリケーションは、取得したアクセストークンを使用して、GoogleAPIにリクエストを送信する。アクセストークンは、APIへのアクセス権限を示す。

OAuth2.0プロトコルにはさまざまなフローがあり、利用ケースに応じて適切なフローを選択する必要がある。

``` mermaid 
sequenceDiagram
    participant USER as ユーザー
    participant UAGE as ユーザーエージェント<br>(ブラウザ)
    participant APSV as アプリケーション<br>サーバー
    participant OASV as OAuth認証<br>サーバー
    participant RSSV as リソース<br>サーバー

    OASV ->> OASV: register client ※1
    USER ->> UAGE: request resource
    UAGE ->> APSV: request resource
    APSV ->> APSV: build authorization url
    APSV -->> UAGE: authorization url
    UAGE ->> OASV: redirect to authorization url ※2
    OASV -->> UAGE: present approval form ※3
    UAGE -->> USER: present approval form
    USER ->> UAGE: approve
    UAGE ->> OASV: post approval ※4
    OASV -->> UAGE: autorization code
    UAGE ->> APSV: redirect to callback url<br>autorization code
    APSV ->> OASV: request access token<br>(exchange) ※5
    OASV -->> APSV: access token
    APSV ->> RSSV: request resource ※6
    RSSV -->> APSV: resource
    APSV -->> UAGE: resource
    UAGE -->> USER: resource
```

# auth urlについて
```
# play groundより取得
https://accounts.google.com/o/oauth2/v2/auth/oauthchooseaccount
  ?redirect_uri=https%3A%2F%2Fdevelopers.google.com%2Foauthplayground
  &prompt=consent
  &response_type=code
  &client_id=407408718192.apps.googleusercontent.com
  &scope=openid
  &access_type=offline
  &service=lso
  &o2v=2
  &theme=mn
  &ddm=0
  &flowName=GeneralOAuthFlow
```

## scopeについて
Google OAuth2 API v2だと以下が存在する。
- https://www.googleapis.com/auth/userinfo.email
- https://www.googleapis.com/auth/userinfo.profile
- openid(認証したいだけならこれだけで良いのかな...？)

%20(スペースのurlencode)で区切って複数指定可能。

## access_typeについて
access_type=offline パラメータは、OAuth 2.0の認証フローにおいて、refresh tokenを取得するための指定。<br>
長期間にわたってユーザーのデータにアクセスし続けるアプリケーションや、バッチ処理などで利用される。

# access token/refresh tokenについて

## request
```
# request header
Authorization: Bearer {access token}
```

## acccess tokenの有効性確認
凡そ以下の通り。
1. token検証endpoint(※googleにはない)で検証
1. resource serverに投げてresponse statusが200であること
1. expire dateを過ぎていないこと

## refresh tokenからaccess tokenの取得
googleの場合、[token refresh endpoint](https://oauth2.googleapis.com/token)に以下をpostする。
1. client_id: クライアントID
1. client_secret: クライアントシークレット
1. refresh_token: リフレッシュトークン
1. grant_type: refresh_token（リフレッシュトークンを使用して新しいアクセストークンを取得することを示す）    

## refresh tokenの保持
refresh tokenはサーバで保持する。

# Google OAuth2 API v2について
Google OAuth2 API v2のscope(前述)を全指定した上で、取得できる項目は以下の通り。

## request url
```
https://www.googleapis.com/userinfo/v2/me
```
## response
```
{
  "family_name": "xxx", 
  "name": "xxx", 
  "picture": "https://lh3.googleusercontent.com/a/xxx", 
  "locale": "ja", 
  "email": "xxx@gmail.com", 
  "given_name": "xxx", 
  "id": "xxx", 
  "verified_email": true
}
```

# リンク
- [Google Developers OAuth 2.0 Playground](https://developers.google.com/oauthplayground/) ← 便利
- [OAUTH2 – GET A TOKEN VIA REST GOOGLE SIGN IN](https://csdcorp.com/blog/coding/oauth2-get-a-token-via-rest-google-sign-in/)


# data-platform-api-google-account-user-info-requests-rmq-kube
data-platform-api-google-account-user-info-requests-rmq-kube は、周辺システム　を データ連携基盤 と統合することを目的に、API でSMS認証トークンデータを生成するマイクロサービスです。

* https://xxx.xxx.io/api/API_SMS_AUTHENTICATION_TOKEN_SRV/generates/

## 動作環境
data-platform-api-google-account-user-info-requests-rmq-kube の動作環境は、次の通りです。  
・ OS: LinuxOS （必須）  
・ CPU: ARM/AMD/Intel（いずれか必須）  

## 本レポジトリ が 対応する API サービス
data-platform-api-google-account-user-info-requests-rmq-kube が対応する APIサービス は、次のものです。

* APIサービス URL: https://xxx.xxx.io/api/API_SMS_AUTHENTICATION_TOKEN_SRV/generates/

## 本レポジトリ に 含まれる API名
data-platform-api-google-account-user-info-requests-rmq-kube には、次の API をコールするためのリソースが含まれています。  

* A_SMSAuthToken（SMS認証トークン - 入力データ）

## API への 値入力条件 の 初期値
data-platform-api-google-account-user-info-requests-rmq-kube において、API への値入力条件の初期値は、入力ファイルレイアウトの種別毎に、次の通りとなっています。  

## データ連携基盤のAPIの選択的コール
Latona および AION の データ連携基盤 関連リソースでは、Inputs フォルダ下の sample.json の accepter に取得したいデータの種別（＝APIの種別）を入力し、指定することができます。  
なお、同 accepter にAll(もしくは空白)の値を入力することで、全データ（＝全APIの種別）をまとめて取得することができます。  

* sample.jsonの記載例(1)  

accepter において 下記の例のように、データの種別（＝APIの種別）を指定します。  
ここでは、"SMSAuthToken" が指定されています。    
  
```
	"api_schema": "DPFMSMSAuthenticatinTokenGenerates",
	"accepter": ["SMSAuthToken"],
```
  
* 全データを取得する際のsample.jsonの記載例(2)  

全データを取得する場合、sample.json は以下のように記載します。  

```
	"api_schema": "DPFMSMSAuthenticatinTokenGenerates",
	"accepter": ["All"],
```

## 指定されたデータ種別のコール
accepter における データ種別 の指定に基づいて DPFM_API_Caller 内の caller.go で API がコールされます。  
caller.go の func() 毎 の 以下の箇所が、指定された API をコールするソースコードです。  

```
func (c *DPFMAPICaller) AsyncCreates(
	accepter []string,
	input *dpfm_api_input_reader.SDC,

	log *logger.Logger,
) []error {
	wg := sync.WaitGroup{}
	mtx := sync.Mutex{}
	errs := make([]error, 0, 5)
	exconfAllExist := false

	subFuncFin := make(chan error)
	exconfFin := make(chan error)

	wg.Add(1)
	go func() {
		defer wg.Done()
		var e []error
		exconfAllExist, e = c.confirmor.Conf(input, log)
		if len(e) != 0 {
			mtx.Lock()
			errs = append(errs, e...)
			mtx.Unlock()
			exconfFin <- xerrors.Errorf("exconf error")
			return
		}
		exconfFin <- nil
	}()

	for _, fn := range accepter {
		wg.Add(1)
		switch fn {
		case "Header":
			go c.headerCreate(&wg, &mtx, subFuncFin, log, errs, input)
		case "Item":
			errs = append(errs, xerrors.Errorf("accepter Item is not implement yet"))
		default:
			wg.Done()
		}
	}
```

## Output  
本マイクロサービスでは、[golang-logging-library-for-data-platform](https://github.com/latonaio/golang-logging-library-for-data-platform) により、以下のようなデータがJSON形式で出力されます。  
以下の sample.json の例は SMS認証トークン の SMS認証トークンデータ が生成された結果の JSON の例です。  
以下の項目のうち、"UserID" ～ "AuthenticationCode" は、/DPFM_API_Output_Formatter/type.go 内 の Type SMSAuthToken {} による出力結果です。"cursor" ～ "time"は、golang-logging-library による 定型フォーマットの出力結果です。  

```
XXX
```
