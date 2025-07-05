# [Investor1337_bot](https://t.me/investor1337_bot)

A Telegram bot for monitoring CS2 item prices on multiple markets.  
Originally built for cases.

---

## 📦 Markets

-   Monitors prices from:
    -   [Steam](https://steamcommunity.com/market/)
    -   [CSGO Market](https://market.csgo.com/en/)
    -   [Lis-Skins](https://lis-skins.com/)
    -   [CS.Money](https://cs.money/)

---

## 🚀 Usage

### 1. Create a Telegram bot

-   Use [@BotFather](https://t.me/BotFather) to generate a bot token.

### 2. Get API tokens

-   Register and obtain API tokens for:
    -   [CSGO Market](https://market.csgo.com/en/)
    -   [Lis-Skins](https://lis-skins.com/)

### 3. Create `.env` file in the root of project

```ini
TELEGRAM_BOT=your_telegram_token
CSMARKET=your_csmarket_token
LISSKINS=your_lisskins_token
```

### 4. Run from cmd/app

```go
go run main.go
```
