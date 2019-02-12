# kranx_test_1

add account with id=98, user=user3, balance=6545:
```
curl -H "Content-Type: application/json" \
--data '{"id": "98", "user": "user3", "balance": 6545}' \
http://localhost:8080/account/  
```
add account with id=105, user=user1, balance=7846:
```
curl -H "Content-Type: application/json" \
--data '{"id": "105", "user": "user1", "balance": 7846}' \
http://localhost:8080/account/  
```
get all accounts:
```
curl http://localhost:8080/account/  
```
-> {"id":{"105":{"id":"105","user":"user1","balance":7846},"98":{"id":"98","user":"user3","balance":6545}}}

get account with id=98:
```
curl http://localhost:8080/account/98  
```
-> {"id":"98","user":"user3","balance":817}

move 1432 coins from id=98 to id=105:
```
curl -H "Content-Type: application/json" \
--data '{"from_id": "98", "to_id": "105", "amount": 1432}' \
http://localhost:8080/transaction/  
```