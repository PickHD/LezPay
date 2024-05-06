# LezPay
E-Wallet using golang.

## Architect Overview :
![ARCH](https://raw.github.com/PickHD/LezPay/main/lezpay_arch.png)

## Main Features : 
- Authentication Customer & Merchant ✅
- Dashboard Customer & Merchant ✅
- Topup Wallet Customer
- Transfer Wallet Customer to Customer
- Payout Wallet Customer to Merchant app
- Redeem Profit Merchant

## Tech Stack :
1. Golang
2. PostgresDB
3. RedisDB
4. External Mail Service (Gmail)
5. Kafka
6. GRPC
7. Jaeger Tracer

## Prerequisites : 
1. Make sure Docker & Docker Compose already installed on your machine
2. Rename `example.env` to `.env` on folder `./cmd` every services
3. Make sure to uncheck comment & fill your **SMTP configuration** on auth env

## Setup :
1. To build & run all services in background using command : 
    ``` 
    make run
    ```
3. If you want to stop all services then run :
    ```
    make stop
    ```
4. Last if want to stop & remove entire services then run :
    ```
    make remove
    ```
