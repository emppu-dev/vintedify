# Vintedify
Vintedify is an application that searches for products on Vinted and sends notifications via Discord or Telegram.

## Installation
1. Clone the repository.
2. Install depencies:
```bash
go mod tidy
```
3. Configure the `.env` file
4. Run the application:
```
make build
make run
```

## Configuration
- `COUNTRY`: The country code for Vinted region (e.g, `fi` or `de`).
- `SEARCH`: The search term for products.
- `DISCORD_WEBHOOK`: The Discord webhook URL for the notifications.
- `TELEGRAM_BOT_TOKEN`: The Telegram bot token.
- `TELEGRAM_CHAT_ID`: The Telegram chat ID.

## Tested Countries
| Country Code | Country Name       | Flag   | Tested |
|--------------|--------------------|--------|--------|
| AT           | Austria            | 🇦🇹    |        |
| BE           | Belgium            | 🇧🇪    |        |
| CZ           | Czech Republic     | 🇨🇿    |        |
| DE           | Germany            | 🇩🇪    | ✅     |
| DK           | Denmark            | 🇩🇰    |        |
| ES           | Spain              | 🇪🇸    |        |
| FI           | Finland            | 🇫🇮    | ✅     |
| FR           | France             | 🇫🇷    |        |
| GR           | Greece             | 🇬🇷    |        |
| HR           | Croatia            | 🇭🇷    |        |
| HU           | Hungary            | 🇭🇺    |        |
| IE           | Ireland            | 🇮🇪    |        |
| IT           | Italy              | 🇮🇹    |        |
| LT           | Lithuania          | 🇱🇹    |        |
| LU           | Luxembourg         | 🇱🇺    |        |
| NL           | Netherlands        | 🇳🇱    |        |
| PL           | Poland             | 🇵🇱    |        |
| PT           | Portugal           | 🇵🇹    |        |
| RO           | Romania            | 🇷🇴    |        |
| SE           | Sweden             | 🇸🇪    | ✅     |
| SK           | Slovakia           | 🇸🇰    |        |
| GB           | United Kingdom     | 🇬🇧    |        |
| US           | United States      | 🇺🇸    |        |

## License
This project is licensed under the terms of the GNU General Public License v3.0.