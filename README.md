# IPlayGames Go SDK

High-level Go SDK for the IPlayGames Game Aggregator API.

## Installation

```bash
go get github.com/iplaygamesai/sdk-wrapper-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    iplaygames "github.com/iplaygamesai/sdk-wrapper-go"
    "github.com/iplaygamesai/sdk-wrapper-go/flows"
)

func main() {
    client, err := iplaygames.NewClient(iplaygames.ClientOptions{
        APIKey:  "your-api-key",
        BaseURL: "https://api.iplaygames.ai",
    })
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Get games
    gamesResp := client.Games().List(ctx, flows.ListParams{
        Search: "bonanza",
    })
    if !gamesResp.Success {
        log.Fatal(gamesResp.Error)
    }
    fmt.Printf("Found %d games\n", len(gamesResp.Games))

    // Start a game session
    sessionResp := client.Sessions().Start(ctx, flows.StartSessionParams{
        GameID:      123,
        PlayerID:    "player_456",
        Currency:    "USD",
        CountryCode: "US",
        IPAddress:   "192.168.1.1",
    })
    if !sessionResp.Success {
        log.Fatal(sessionResp.Error)
    }

    // Redirect player to game
    fmt.Println("Game URL:", sessionResp.GameURL)
}
```

## Configuration

```go
client, err := iplaygames.NewClient(iplaygames.ClientOptions{
    APIKey:        "your-api-key",            // Required
    BaseURL:       "https://api.iplaygames.ai", // Optional
    WebhookSecret: "your-secret",             // Optional, for webhook verification
})
```

## Available Flows

### Games

```go
ctx := context.Background()

// List games with filters
gamesResp := client.Games().List(ctx, flows.ListParams{
    Search:     "bonanza",
    Provider:   "pragmatic",
    Type:       "slots",
    PerPage:    20,
})
if gamesResp.Success {
    for _, game := range gamesResp.Games {
        fmt.Printf("Game: %s by %s\n", game.Title, game.Producer)
    }
}

// Get single game
gameResp := client.Games().Get(ctx, 123)

// Convenience methods
pragmaticGames := client.Games().ByProducer(ctx, 42, flows.ListParams{})
liveGames := client.Games().ByCategory(ctx, "live", flows.ListParams{})
searchResults := client.Games().Search(ctx, "sweet bonanza", flows.ListParams{})
```

### Sessions

```go
ctx := context.Background()

// Start a game session
sessionResp := client.Sessions().Start(ctx, flows.StartSessionParams{
    GameID:      123,
    PlayerID:    "player_456",
    Currency:    "USD",
    CountryCode: "US",
    IPAddress:   "192.168.1.1",
    Locale:      "en",
    Device:      "mobile",
    ReturnURL:   "https://casino.com/lobby",
})
if sessionResp.Success {
    fmt.Println("Game URL:", sessionResp.GameURL)
    fmt.Println("Session ID:", sessionResp.SessionID)
}

// Get session status
statusResp := client.Sessions().Status(ctx, sessionResp.SessionID)

// End session
endResp := client.Sessions().End(ctx, sessionResp.SessionID)

// Start demo session
demoResp := client.Sessions().StartDemo(ctx, 123, flows.StartSessionParams{})
```

### Jackpot

```go
ctx := context.Background()

// Get configuration
configResp := client.Jackpot().GetConfiguration(ctx)

// Get all pools
poolsResp := client.Jackpot().GetPools(ctx)

// Get specific pool
dailyPoolResp := client.Jackpot().GetPool(ctx, "daily")

// Get winners
winnersResp := client.Jackpot().GetWinners(ctx, "pool_123")

// Manage games
addResp := client.Jackpot().AddGames(ctx, "daily", []int{1, 2, 3})
removeResp := client.Jackpot().RemoveGames(ctx, "daily", []int{1})

// Get contribution history
contribResp := client.Jackpot().GetContributions(ctx, flows.ContributionFilters{
    PlayerID: "player_456",
})
```

### Promotions

```go
ctx := context.Background()

// List promotions
promoListResp := client.Promotions().List(ctx, "active", "")

// Get promotion details
promoResp := client.Promotions().Get(ctx, 1)

// Create a promotion
createResp := client.Promotions().Create(ctx, flows.PromotionData{
    Name:          "Summer Tournament",
    PromotionType: "tournament",
    CycleType:     "daily",
    StartsAt:      "2024-06-01T00:00:00Z",
    EndsAt:        "2024-06-30T23:59:59Z",
})

// Get leaderboard
leaderboardResp := client.Promotions().GetLeaderboard(ctx, 1, 10, 0)

// Opt-in player
optInResp := client.Promotions().OptIn(ctx, 1, "player_456", "USD")

// Manage games for promotion
manageResp := client.Promotions().ManageGames(ctx, 1, []int{1, 2, 3})
```

### Jackpot Widgets

```go
ctx := context.Background()

// 1. Register your domain
domainResp := client.JackpotWidget().RegisterDomain(ctx, "casino.example.com", "My Casino")
// Get domain token from response

// 2. List registered domains
domainsResp := client.JackpotWidget().ListDomains(ctx)

// 3. Create anonymous token (view-only)
anonTokenResp := client.JackpotWidget().CreateAnonymousToken(ctx, "domain_token_here")

// 4. Create player token (can interact)
playerTokenResp := client.JackpotWidget().CreatePlayerToken(ctx, "domain_token_here", "player_456", "USD")

// 5. Get embed code for your frontend
embedCode := client.JackpotWidget().GetEmbedCode("widget_token_here", flows.EmbedOptions{
    Theme:     "dark",
    Container: "jackpot-widget",
    PoolTypes: []string{"daily", "weekly"},
})
fmt.Println(embedCode)
// Output:
// <div id="jackpot-widget"></div>
// <script src="https://api.iplaygames.ai/widgets/jackpot.js"></script>
// <script>
//     IPlayGamesJackpotWidget.init({"container":"jackpot-widget","theme":"dark","pool_types":["daily","weekly"],"token":"widget_token_here"});
// </script>
```

### Promotion Widgets

```go
ctx := context.Background()

// Register domain
domainResp := client.PromotionWidget().RegisterDomain(ctx, "casino.example.com")

// Create player token
tokenResp := client.PromotionWidget().CreatePlayerToken(ctx, "domain_token", "player_456", "USD")

// Get embed code
embedCode := client.PromotionWidget().GetEmbedCode("widget_token", flows.PromotionEmbedOptions{
    Theme:        "dark",
    Container:    "promo-widget",
    PromotionIDs: []int{1, 2, 3},
})
```

### Multi-Session (TikTok-style Game Swiping)

```go
ctx := context.Background()

// Start multi-session
multiResp := client.MultiSession().Start(ctx, flows.StartMultiSessionParams{
    PlayerID:    "player_456",
    Currency:    "USD",
    CountryCode: "US",
    IPAddress:   "192.168.1.1",
    Device:      "mobile",
    GameIDs:     []string{"123", "456", "789"}, // Optional: specific games
})
if multiResp.Success {
    fmt.Println("Swipe URL:", multiResp.SwipeURL)
    fmt.Println("Total Games:", multiResp.TotalGames)
}

// Get iframe HTML to embed the swipe UI
iframe := client.MultiSession().GetIframe(multiResp.SwipeURL, flows.IframeOptions{
    Width:  "100%",
    Height: "100vh",
    ID:     "game-swiper",
})

// Get status
statusResp := client.MultiSession().Status(ctx, multiResp.MultiSessionID)

// End when player leaves
endResp := client.MultiSession().End(ctx, multiResp.MultiSessionID)
```

## Handling Webhooks

GameHub sends webhooks for transactions. Your casino must implement a webhook endpoint.

### Webhook Types

| Type | Description |
|------|-------------|
| `authenticate` | Verify player exists and get initial data |
| `balance_check` | Get current player balance |
| `bet` | Player placed a bet |
| `win` | Player won money |
| `rollback` | Undo a transaction |
| `reward` | Award from promotions/tournaments |

### Implementing Your Webhook Handler

```go
package main

import (
    "encoding/json"
    "io"
    "net/http"
    "os"

    iplaygames "github.com/iplaygamesai/sdk-wrapper-go"
    "github.com/iplaygamesai/sdk-wrapper-go/webhooks"
)

var client *iplaygames.Client

func init() {
    var err error
    client, err = iplaygames.NewClient(iplaygames.ClientOptions{
        APIKey:        os.Getenv("GAMEHUB_API_KEY"),
        WebhookSecret: os.Getenv("GAMEHUB_WEBHOOK_SECRET"),
    })
    if err != nil {
        panic(err)
    }
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Failed to read body", http.StatusBadRequest)
        return
    }

    signature := r.Header.Get("X-Signature")
    handler, err := client.Webhooks()
    if err != nil {
        http.Error(w, "Webhook handler not configured", http.StatusInternalServerError)
        return
    }

    // Verify signature
    if !handler.Verify(string(body), signature) {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{"error": "Invalid signature"})
        return
    }

    // Parse webhook
    webhook, err := handler.Parse(string(body))
    if err != nil {
        http.Error(w, "Failed to parse webhook", http.StatusBadRequest)
        return
    }

    // Handle by type
    var response interface{}

    switch webhook.Type {
    case webhooks.TypeAuthenticate:
        response = handleAuthenticate(handler, webhook)
    case webhooks.TypeBalanceCheck:
        response = handleBalanceCheck(handler, webhook)
    case webhooks.TypeBet:
        response = handleBet(handler, webhook)
    case webhooks.TypeWin:
        response = handleWin(handler, webhook)
    case webhooks.TypeRollback:
        response = handleRollback(handler, webhook)
    case webhooks.TypeReward:
        response = handleReward(handler, webhook)
    default:
        http.Error(w, "Unknown webhook type", http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func handleAuthenticate(handler *webhooks.Handler, webhook *webhooks.Payload) interface{} {
    player, err := findPlayer(webhook.PlayerID)
    if err != nil {
        return handler.PlayerNotFoundResponse()
    }

    balance := player.GetBalance(webhook.Currency)
    return handler.SuccessResponse(balance, map[string]interface{}{
        "player_name": player.Name,
    })
}

func handleBet(handler *webhooks.Handler, webhook *webhooks.Payload) interface{} {
    player, _ := findPlayer(webhook.PlayerID)
    balance := player.GetBalance(webhook.Currency)
    betAmount := webhook.GetAmountInDollars()

    // Check funds
    if betAmount != nil && balance < *betAmount {
        return handler.InsufficientFundsResponse(balance)
    }

    // Check idempotency
    if webhook.TransactionID != nil && transactionExists(*webhook.TransactionID) {
        return handler.AlreadyProcessedResponse(balance)
    }

    // Process bet
    if betAmount != nil {
        player.Debit(*betAmount, webhook.Currency)
    }

    newBalance := player.GetBalance(webhook.Currency)
    return handler.SuccessResponse(newBalance, nil)
}

// ... implement other handlers similarly

func main() {
    http.HandleFunc("/webhooks/gamehub", webhookHandler)
    http.ListenAndServe(":8080", nil)
}
```

## Webhook Payload Fields

### Common Fields (all webhook types)

```go
webhook.Type        // "bet", "win", "rollback", "reward", "authenticate", "balance_check"
webhook.PlayerID    // Player's ID in your system
webhook.Currency    // "USD", "EUR", etc.
webhook.GameID      // Game ID (nullable)
webhook.GameType    // "slot", "live", "table", etc.
webhook.Timestamp   // ISO 8601 timestamp
```

### Transaction Fields (bet, win, rollback, reward)

```go
webhook.TransactionID             // Unique transaction ID (nullable)
webhook.Amount                    // Amount in cents (nullable)
webhook.GetAmountInDollars()      // Amount in dollars (nullable)
webhook.SessionID                 // Game session ID
webhook.RoundID                   // Game round ID
```

### Freespin Fields

```go
webhook.IsFreespin              // Is this a freespin round?
webhook.FreespinID              // Freespin campaign ID
webhook.FreespinTotal           // Total freespins awarded
webhook.FreespinsRemaining      // Remaining freespins
webhook.FreespinRoundNumber     // Current spin number
webhook.FreespinTotalWinnings   // Cumulative winnings
```

## Response Pattern

All flow methods return response structs with a consistent pattern:

```go
// Check success
resp := client.Games().List(ctx, params)
if !resp.Success {
    fmt.Println("Error:", resp.Error)
    return
}

// Use the data
for _, game := range resp.Games {
    fmt.Println(game.Title)
}
```

For generic responses using `ApiResponse`:

```go
resp := client.Jackpot().GetPools(ctx)
if resp.Success {
    pools := resp.Data["pools"]
    fmt.Println(pools)
}
```

## Running Tests

```bash
go test ./...
```

## License

MIT
