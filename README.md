<img src="https://elytrium.net/src/img/elytrium.webp" alt="Elytrium" align="right">

## In development

# Elling - Elytrium Billing
[![Join our Discord](https://img.shields.io/discord/775778822334709780.svg?logo=discord&label=Discord)](https://ely.su/discord)

Module-based billing platform made with Go<br>
The main idea of this product - make a stable billing platform for high-loads<br>
This is only the back-end side of the API! Check out [elling-app](https://github.com/Elytrium/elling-app) for the front-end.

## Module system

Elling - module-based billing. You can create your own module, just export your module to the variable ``Module`` in your Go plugin

```go
type Module interface {
    OnInit()
    GetName() string
    OnRegisterMethods() map[string]routing.Method
    OnDBMigration() []interface{}
    OnSmallTick()
    OnBigTick()
}
```

### See more

- [elling-npd](https://github.com/Elytrium/elling-npd): Payments module for self-employed people


### Basic modules

Basic modules - really simple modules, you can configure them editing their .yml files

- ```basic/oauth```: Module for OAuth authorization support
  - Create the folder with name "oauth"
  - Create the .yml file there
  - Fill it (example):
  - ```yml
    display-name: Discord
    name: discord
    oauth-gen-request: "https://discord.com/api/oauth2/authorize?client_id=793481663077548032&redirect_uri=https%3A%2F%2Fsrv.cool%2Finternal%2Foauth&response_type=code&scope=identify"
    need-verify: true
    verify-request:
      url: https://discord.com/api/oauth2/token
      method: POST
      headers:
        Content-Type: application/x-www-form-urlencoded
      data: "client_id=793481663077548032&client_secret=whoopsy&grant_type=authorization_code&code={token}&redirect_uri=https%3A%2F%2Fsrv.cool%2Finternal%2Foauth"
      response-type: JSON
      response-value-path:
      - access_token
    get-data-request:
      url: https://discord.com/api/oauth2/@me
      method: GET
      headers:
        Authorization: Bearer {token}
      response-type: JSON
      response-value-path:
      - user.username
      - user.id
    ```
- ```basic/topup```: Simple top-up module
    - Create the folder with name "topup"
    - Create the .yml file there
    - Fill it (example):
    - ```yml
      name: hevav-pay
      display-name: hevav.pay 
      account-limit: 10
      ttl: 3600000
      pay-string: https://hevav.dev/pay/{topUpId}
      create-request:
        url: https://hevav.dev/pay
        method: PUT
        headers: 
          Authorization: Bearer 1234567890
        data: id={topUpID}&amount={amount}&user[name]={user_name}&user[id]={balance_id}&expiryDate={date}
        response-type: NONE
      check-request:
        url: https://hevav.dev/payverify
        method: POST
        headers: 
          Authorization: Bearer 1234567890
        data: id={topUpID}&amount={amount}&user[name]={user_name}&user[id]={balance_id}&expiryDate={date}
        response-type: PLAIN
      check-request-success-string: OK
      reject-request:
        url: https://hevav.dev/pay/{topUpId}
        method: DELETE
        headers: 
          Authorization: Bearer 1234567890
        data: id={topUpID}
        response-type: NONE
      ```

## Donation

Your donations are really appreciated. Donations wallets/links/cards:

- MasterCard Debit Card (Tinkoff Bank): ``5536 9140 0599 1975``
- Qiwi Wallet: ``PFORG`` or [this link](https://my.qiwi.com/form/Petr-YSpyiLt9c6)
- YooMoney Wallet: ``4100 1721 8467 044`` or [this link](https://yoomoney.ru/quickpay/shop-widget?writer=seller&targets=Donation&targets-hint=&default-sum=&button-text=11&payment-type-choice=on&mobile-payment-type-choice=on&hint=&successURL=&quickpay=shop&account=410017218467044)
- PayPal: ``ogurec332@mail.ru``
