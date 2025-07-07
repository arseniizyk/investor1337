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
    -   [Aim.Market](https://aim.market/)
    -   [CSFloat](https://csfloat.com/)
    -   [BUFF163](https://buff.163.com/)

---

## 🚀 Usage

### 1. Create a Telegram bot

-   Use [@BotFather](https://t.me/BotFather) to generate a bot token.

### 2. Get API tokens and Cookies

-   Register and obtain API tokens for:

    -   [CSGO Market](https://market.csgo.com/en/)
    -   [Lis-Skins](https://lis-skins.com/)

-   **Note**: This step may lead to a ban by CSFLOAT
-   Log in to [CSFloat](https://csfloat.com/) and retrieve your session cookie from the request.

### 3. Create `.env` file in the root of project

```ini
TELEGRAM_BOT=your_telegram_token
CSMARKET=your_csmarket_token
LISSKINS=your_lisskins_token
CSFLOAT_COOKIES=your_csfloat_session_cookie
```

### 4. Run from cmd/app

```go
go run main.go
```

## Thanks to

-   [somespecialone/steam-item-name-ids](https://github.com/somespecialone/steam-item-name-ids) for steam ids
-   [ModestSerhat/cs2-marketplace-ids](https://github.com/ModestSerhat/cs2-marketplace-ids) for buff163 ids
