PAX BRITANNICA
==============

Notes : all balances are in lowest unit e.g. wei
but VALUES are in Units (e.g. BAT)

Errors
---

```
Status : "Error"
Error  : "Description of error"
Result : only if it makes sense
```

1. Create New Address
---

`/admin/createUserAddress`

parameter : user (string)

`http://127.0.0.1:9000/admin/createUserAddress?user=dave`

```
{
    "Status": "OK",
    "Error": "",
    "User": "dave",
    "Result": "0x4933f5A58e1Dd13F39Cb662DFf4566eC1259C187"
}
```
or
```
{
    "Status": "ERROR",
    "Error": "User already exists",
    "User": "dave",
    "Result": "0x4933f5A58e1Dd13F39Cb662DFf4566eC1259C187"
}
```

Result returns the user address if it already exists or can be created


2. Get users address (public key)
---

`/admin/getUserAddress`

paramater : user (string)

`http://127.0.0.1:9000/admin/getUserAddress?user=dave`

```
{
    "Status": "OK",
    "Error": "",
    "User": "dave",
    "Result": "0x4933f5A58e1Dd13F39Cb662DFf4566eC1259C187"
}
```
or
```
{
    "Status": "ERROR",
    "Error": "User does not exist",
    "User": "millie",
    "Result": ""
}
```

3. Get user's private key (warning: allows account to be compromised)
---

`/admin/getUserPrivateKey`

parameter : user (string)

`http://127.0.0.1:9000/admin/getUserPrivateKey?user=dave`

```
{
    "Status": "OK",
    "Error": "",
    "User": "dave",
    "Result": "0x3605fdbb8053dfceefa743f95024e3a11d86a32aca1580eb216ba4893b40380d"
}
```

4. Get user's ether balance
---

`/admin/getUserEtherBalance`

parameter : user (string)

`http://127.0.0.1:9000/admin/getUserEtherBalance?user=dave`


```
{
    "Status": "OK",
    "Error": "",
    "User": "dave",
    "Result": "0"
}
```

5. Get user's ether AND token balance
---

`/admin/getUserBalance`

parameter : user (string)

`http://127.0.0.1:9000/admin/getUserBalance?user=bloosie`

returns

```
{
    "Status": "OK",
    "Result": [
        {
            "Token": "StorjToken",
            "Balance": 100000
        },
        {
            "Token": "ETH",
            "Balance": 7366325000000000
        }
    ]
}
```

6. Send ether from user account
---

`/admin/sendEtherFromUser`

parameters :
	user - user ID to send from (string)
	destination - address to send to
	value - value to send (string) in ether (i.e. 2.5 is OK)
	gasLimit - gasLimit to supply
	gasPrice - gasPrice to offer
	data - data to add to the txn...

```
http://127.0.0.1:9000/admin/sendEtherFromUser?user=bloosie&destination=0x31efd75bc0b5fbafc6015bd50590f4fdab6a3f22&value=0.000345326000000000&gasLimit=21000&gasPrice=100000000
```

```
{
    "Hash": "0xa6fa1cba62c9d321901e0419935f3612370ed4b40ccf8faa161bc48a3b95c675"
}
```

7. Add a token to the known token database
---

`/admin/addTokenToSystem`

parameter : address

`http://127.0.0.1:8009/admin/addTokenToSystem?address=0xba2184520A1cC49a6159c57e61E1844E085615B6`

```
{
    "Name": "HelloGold Token",
    "Symbol": "HGT",
    "Address": "0xba2184520A1cC49a6159c57e61E1844E085615B6",
    "Decimals": 8
}
```

8. List tokens in database
---

`/admin/listTokens`

parameters : none

`http://127.0.0.1:8009/admin/listTokens`

```
[
    {
        "Name": "HelloGold Token",
        "Symbol": "HGT",
        "Address": "0xba2184520A1cC49a6159c57e61E1844E085615B6",
        "Decimals": 8
    }
]
```

9. Get user's token AND ether balances
---

`/admin/getTotalUserBalances`

Get total token holding across ALL users for ALL tokens AND ether
---

```
{
    "Status": "OK",
    "Result": [
        {
            "Token": "HelloGold Gold Backed Token",
            "Balance": 0
        },
        {
            "Token": "HelloGold Token",
            "Balance": 0
        },
        {
            "Token": "Jetcoin",
            "Balance": 0
        },
        {
            "Token": "StorjToken",
            "Balance": 200000
        },
        {
            "Token": "ETH",
            "Balance": 7366325000000000
        }
    ]
}
```

10 Get Token Balances 
---

`/admin/getTotalTokenBalances`


requires one of

`symbol` - TLA e.g. HGT

or

`address` - contract address e.g. 

`http://127.0.0.1:9000/admin/getTotalTokenBalances?symbol=STORJ`

```
{
    "Status": "OK",
    "Result": [
        {
            "Token": "StorjToken",
            "Balance": 100000
        }
    ]
}
```

11. Get 
`/admin/getTokenBalance`

requires one of

`symbol` - TLA e.g. HGT

or

`address` - contract address e.g. 

`http://127.0.0.1:9000/admin/getTokenBalance?user=bloosie&symbol=STORJ`

```
{
    "Status": "OK",
    "Result": [
        {
        "Token": "StorjToken",
        "Balance": 100000
        }
    ]
}
```


12. Get the result of a transaction
---

`/admin/getTxStatus`

gets pass/fail status of a transaction

parameter - `hash` - hash of transactions 

```
http://127.0.0.1:8009/admin/getTxStatus?hash=0xfaa97aa80db0f9326647cd611666f7bc853ae6a133fb17218183e19d314d2f31
```

```
{
    "Status": "OK",
    "Result": false
}
```

13. Send tokens from a user account
---

`/admin/sendTokenFromUser`

parameters : 
user
symbol / address
destination
value
gasPrice
gasLimit

`http://127.0.0.1:9000/admin/sendTokenFromUser?user=bloosie&symbol=STORJ&destination=0x31EFd75bc0b5fbafc6015Bd50590f4fDab6a3F22&value=10&gasPrice=10000000&gasLimit=200000`

```
{
    "Hash": "0x50e88585c87c6d706878eb0ec11b09bd22570280dd17b69f27f1181aebff75dd"
}
```

14. Get Info about the token
---

`/admin/getTokenInfo`

`http://127.0.0.1:9000/admin/getTokenInfo?symbol=HGT`

```
{
    "Status": "",
    "Result": [
    {
        "Name": "HelloGold Token",
        "Symbol": "HGT",
        "Address": "0xba2184520A1cC49a6159c57e61E1844E085615B6",
        "Decimals": 8
    }
    ]
}
```
