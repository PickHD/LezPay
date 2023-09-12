# LezPay
E-Wallet using golang.

## Architect Overview :
![ARCH](https://raw.github.com/PickHD/LezPay/main/lezpay_arch.png)

## Main Features : 
1. **Register & Login, Reset Password Customer & Merchant** _Coming Soon_
2. **Topup Wallet Customer** (Only Support VA Bank For Now) _Coming Soon_
3. **Payout Wallet Customer** (Supported Merchant Only) _Coming Soon_
3. **Dashboard Customer** (Show Current Balances, Show List History Topup/Payout) _Coming Soon_
4. **Dashboard Merchant** (Show Total Transactions, Show List History Transaction) _Coming Soon_

## Tech Used :
1. Golang
2. PostgreSQL
3. Redis
4. Kafka
5. GRPC
6. Docker
7. Jaeger Tracer

## Prerequisites : 
1. Make sure Docker & Docker Compose already installed on your machine
2. Rename `example.env` to `.env` on folder `./cmd/v1` every services
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
