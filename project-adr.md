# YMB Cloz - Telegram Bot

Below is a proposed architecture for **YMB Cloz**, a Telegram bot to track local Dota games played among friends. The goal is to have an MVP (minimum viable product) that provides:

- An **admin panel** (local usage) to record match data  
- A **Telegram bot** to retrieve and display stats

## 1. Overview

1. **Frontend (TypeScript)**  
   - An admin panel (also TypeScript and React) for local or internal usage, providing forms to input game details.

2. **Backend (Go)**
   - A REST or GraphQL API server that handles requests for:
     - Storing new game records.
     - Retrieving player statistics.
   - Implements the logic of updating statistics (wins, losses, total games) after each game is recorded.
   - Hosts the **Telegram Bot** functionality (using a Go library such [telebot](https://github.com/tucnak/telebot)).

3. **Database (PostgreSQL)**  
   - Stores information about players, their stats, and game records.

4. **Telegram Bot**
   - Users interact with the bot via **commands** in a Telegram chat (e.g., `/top_winrate`, `/top_games`, etc.).
   - The bot queries the backend to gather the requested stats and responds in Telegram with formatted text or messages.

---

## 2. Data Model

The following entities should be enough for the MVP. Feel free to extend them if needed.

1. **Player**  
   - **id** (UUID or serial)  
   - **nickname** (string, unique)  
   - **games** (number of games played, can be computed or stored directly)  

2. **Game**  
   - **id** (UUID or serial)  
   - **timestamp** (date/time of the game)  
   - **radiantPlayers** (array of players or a join table)  
   - **direPlayers** (array of players or a join table)  
   - **winningTeam** (enum: `DIRE`, `RADIANT`)  
   - **roles** for each player (carry, mid, offlane, pos4, pos5)  
   - **captain** for each team (store as a special flag or a reference to a player)

> **Additional Table**: You could create a `GamePlayers` table to keep track of each player’s role, team, and captain flag for a given game. For instance:
> - `game_id` (FK to `Game`)
> - `player_id` (FK to `Player`)
> - `team` (`DIRE` / `RADIANT`)
> - `role` (`carry`, `mid`, etc.)
> - `is_captain` (boolean)

## 3. Backend (Go)

### Structure
- Use a framework like [Gin](https://github.com/gin-gonic/gin) or [Fiber](https://github.com/gofiber/fiber) for RESTful routes.  
- Implement the standard CRUD endpoints:
  - `POST /api/games` to add a new game
  - `GET /api/players/top-winrate` to get top players by win rate  
  - `GET /api/players/top-games-played` to get top players by total games  
  - `GET /api/players/top-role/{role}` to get top players on a specific role  
- Connect to PostgreSQL (via `database/sql`).

### Business Logic
- When a new game is recorded:
  1. Insert a new `Game` record.  
  2. For each player, create a record in `GamePlayers` with `team`, `role`, and `is_captain`. 

### Telegram Bot
- Integrate a bot in Go (using something like [telebot](https://github.com/tucnak/telebot) or [go-telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api)).  
- Configure a **webhook** or **long polling** to receive commands.  
- Commands can include:
  - `/top_winrate`  
  - `/top_games`  
  - `/top_role [role]`  
  - `/help`  

## 4. Frontend (TypeScript)

### Admin Panel
- Can be a separate local web page or part of the main front-end but protected (e.g., behind a login or run locally).  
- Provide a form to create a new game:
  - Select the **two teams** (Radiant, Dire).  
  - For each team, choose up to 5 players from a player dropdown.  
  - Mark which team won.  
  - Mark one captain per team.  
- Send the data to your Go API (`POST /api/games`) to insert a new game.

### Implementation Details
- Use a bundler (Vite) for the TS code.  
- You can serve the built static files via the Go server or deploy them to a static host.

## 5. Telegram Bot

- **Implementation** in Go:
  - Use a library like [telebot](https://github.com/tucnak/telebot).
  - Configure either a **webhook** or **long polling** to receive messages from Telegram.
- **Bot Commands** (examples):
  - `/start` – greets the user or provides help info.
  - `/top_winrate` – returns a list of players sorted by highest win rate.
  - `/top_games` – returns a list of players sorted by total games played.
  - `/top_role <role>` – returns a list of top players in the specified role (carry, mid, etc.).
  - `/help` – lists available commands and usage.
- **Response Format**:
  - The bot queries the backend for the required data and returns formatted text messages. 
  - Example response to `/top_winrate`:
    ```
    Top 5 players by win rate:
    1. PlayerA - 75% (9/12)
    2. PlayerB - 70% (7/10)
    ...
    ```


## 6. Deployment Suggestions

- **Server**  
  - For a beginner-friendly approach, you can deploy the Go server to [Railway](https://railway.app/), [Heroku](https://www.heroku.com/) (if you find a free tier alternative), or [Fly.io](https://fly.io/).  
  - Ensure you configure environment variables for database connection strings, bot tokens, etc.

- **Database**  
  - Use a managed Postgres solution (e.g., Railway, ElephantSQL, etc.).  
  - Create migrations (with a library like [golang-migrate](https://github.com/golang-migrate/migrate)) to set up the schema.

## 7. Additional Ideas / Extensions

- **Ranked Roles**: Show who’s the best *carry*, *mid*, *offlane*, *pos4*, *pos5*, based on either a rating or some other performance metric (like KDA, if you want to record it).  
- **Detailed Stats**: Track kills, deaths, assists, or other Dota stats for each game to offer richer insights.  
- **MVP vs. Extended**: For the MVP, keep it simple with just wins, losses, total games. Later, expand.  
- **Automated Stats**: If you have a local replay parsing mechanism or an API, you could automate retrieving game data. But for MVP, manual input is simpler.

## 8. Summary & Next Steps

1. **Set up the Go backend** with a framework and migrations for Postgres.  
2. **Create the minimal Admin Panel** to input game data (teams, roles, captains, winner).  
3. **Build the Telegram Bot** in Go, register commands, and set up either a webhook or polling.  
4. **Develop the Telegram WebApp** in TypeScript for reading and displaying top stats.  
5. **Deploy** the Go server and Postgres.  
6. **Configure** the Telegram Bot settings in the [BotFather](https://t.me/botfather). 
