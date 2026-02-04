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
    "fmt"
    "log"

    iplaygames "github.com/iplaygamesai/sdk-wrapper-go"
)

func main() {
    client := iplaygames.NewClient(&iplaygames.ClientOptions{
        APIKey:  "your-api-key",
        BaseURL: "https://api.gamehub.com", // Configurable!
    })

    // Get games
    games, err := client.Games().List(&iplaygames.GameListParams{
        Currency: "USD",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found %d games\n", len(games))

    // Start a game session
    session, err := client.Sessions().Start(&iplaygames.SessionStartParams{
        GameID:      123,
        PlayerID:    "player_456",
        Currency:    "USD",
        CountryCode: "US",
        IPAddress:   "192.168.1.1",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Redirect player to game
    fmt.Println("Game URL:", session.GameURL)
}
```

## Configuration

```go
client := iplaygames.NewClient(&iplaygames.ClientOptions{
    APIKey:        "your-api-key",           // Required
    BaseURL:       "https://api.gamehub.com", // Optional, defaults to https://api.gamehub.com
    Timeout:       30 * time.Second,          // Optional, request timeout
    WebhookSecret: "your-secret",             // Optional, for webhook verification
})
```

## Available Flows

### Games

```go
// List games with filters
games, err := client.Games().List(&iplaygames.GameListParams{
    Currency: "USD",
    Country:  "US",
    Category: "slots",
    Search:   "bonanza",
})

// Get single game
game, err := client.Games().Get(123)

// Convenience methods
pragmaticGames, err := client.Games().ByProducer("Pragmatic Play")
liveGames, err := client.Games().ByCategory("live")
searchResults, err := client.Games().Search("sweet bonanza")
playerGames, err := client.Games().ForPlayer("USD", "US")
```

### Sessions

```go
// Start a game session
session, err := client.Sessions().Start(&iplaygames.SessionStartParams{
    GameID:      123,
    PlayerID:    "player_456",
    Currency:    "USD",
    CountryCode: "US",
    IPAddress:   r.RemoteAddr,
    Locale:      "en",
    Device:      "mobile",
    ReturnURL:   "https://casino.com/lobby",
})

// Get session status
status, err := client.Sessions().Status(session.SessionID)

// End session
err = client.Sessions().End(session.SessionID)

// Start demo session
demo, err := client.Sessions().StartDemo(123)
```

### Jackpot

```go
// Get configuration
config, err := client.Jackpot().GetConfiguration()

// Get all pools
pools, err := client.Jackpot().GetPools()

// Get specific pool
dailyPool, err := client.Jackpot().GetPool("daily")
weeklyPool, err := client.Jackpot().GetPool("weekly")

// Get winners
winners, err := client.Jackpot().GetWinners("daily")

// Manage games
err = client.Jackpot().AddGames("daily", []int{1, 2, 3})
err = client.Jackpot().RemoveGames("daily", []int{1})
```

### Promotions

```go
// List promotions
promotions, err := client.Promotions().List(&iplaygames.PromotionListParams{
    Status: "active",
})

// Get promotion details
promo, err := client.Promotions().Get(1)

// Get leaderboard
leaderboard, err := client.Promotions().GetLeaderboard(1)

// Opt-in player
err = client.Promotions().OptIn(1, "player_456", "USD")
```

### Jackpot Widgets

```go
// 1. Register your domain
domain, err := client.JackpotWidget().RegisterDomain("casino.example.com")
domainToken := domain.DomainToken

// 2. Create anonymous token (view-only)
token, err := client.JackpotWidget().CreateAnonymousToken(domainToken)

// 3. Create player token (can start game sessions)
playerToken, err := client.JackpotWidget().CreatePlayerToken(
    domainToken,
    "player_456",
    "USD",
)

// 4. Get embed code for your frontend
embedCode := client.JackpotWidget().GetEmbedCode(token.Token, &iplaygames.WidgetEmbedOptions{
    Theme:     "dark",
    Container: "jackpot-widget",
})
```

### Promotion Widgets

```go
// Same flow as jackpot widgets
domain, err := client.PromotionWidget().RegisterDomain("casino.example.com")
token, err := client.PromotionWidget().CreatePlayerToken(
    domain.DomainToken,
    "player_456",
    "USD",
)
embedCode := client.PromotionWidget().GetEmbedCode(token.Token, nil)
```

### Multi-Session (TikTok-style Game Swiping)

```go
// Start multi-session
multiSession, err := client.MultiSession().Start(&iplaygames.MultiSessionStartParams{
    PlayerID:    "player_456",
    Currency:    "USD",
    CountryCode: "US",
    IPAddress:   r.RemoteAddr,
    Device:      "mobile",
})

// Get iframe HTML to embed the swipe UI
iframe := client.MultiSession().GetIframe(multiSession.SwipeURL, &iplaygames.IframeOptions{
    Width:  "100%",
    Height: "100vh",
})

// Get status
status, err := client.MultiSession().Status(multiSession.MultiSessionID)

// End when player leaves
err = client.MultiSession().End(multiSession.MultiSessionID)
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
    client = iplaygames.NewClient(&iplaygames.ClientOptions{
        APIKey:        os.Getenv("GAMEHUB_API_KEY"),
        WebhookSecret: os.Getenv("GAMEHUB_WEBHOOK_SECRET"),
    })
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Failed to read body", http.StatusBadRequest)
        return
    }

    signature := r.Header.Get("X-Signature")
    handler := client.Webhooks()

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
        response = handleAuthenticate(webhook)
    case webhooks.TypeBalanceCheck:
        response = handleBalanceCheck(webhook)
    case webhooks.TypeBet:
        response = handleBet(webhook)
    case webhooks.TypeWin:
        response = handleWin(webhook)
    case webhooks.TypeRollback:
        response = handleRollback(webhook)
    case webhooks.TypeReward:
        response = handleReward(webhook)
    default:
        http.Error(w, "Unknown webhook type", http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func handleAuthenticate(webhook *webhooks.Payload) interface{} {
    player, err := findPlayer(webhook.PlayerID)
    if err != nil {
        return client.Webhooks().PlayerNotFoundResponse()
    }

    balance := player.GetBalance(webhook.Currency)

    return client.Webhooks().SuccessResponse(balance, map[string]interface{}{
        "player_name": player.Name,
    })
}

func handleBet(webhook *webhooks.Payload) interface{} {
    player, _ := findPlayer(webhook.PlayerID)
    balance := player.GetBalance(webhook.Currency)
    betAmount := webhook.GetAmountInDollars()

    // Check funds
    if balance < betAmount {
        return client.Webhooks().InsufficientFundsResponse(balance)
    }

    // Check idempotency
    if transactionExists(webhook.TransactionID) {
        return client.Webhooks().AlreadyProcessedResponse(balance)
    }

    // Process bet
    player.Debit(betAmount, webhook.Currency)
    createTransaction(&Transaction{
        ExternalID: webhook.TransactionID,
        PlayerID:   webhook.PlayerID,
        Type:       "bet",
        Amount:     betAmount,
        Currency:   webhook.Currency,
    })

    newBalance := player.GetBalance(webhook.Currency)
    return client.Webhooks().SuccessResponse(newBalance, nil)
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
webhook.TransactionID             // Unique transaction ID
webhook.Amount                    // Amount in cents
webhook.GetAmountInDollars()      // Amount in dollars
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

## Error Handling

```go
session, err := client.Sessions().Start(params)
if err != nil {
    if apiErr, ok := err.(*iplaygames.APIError); ok {
        fmt.Printf("API Error: %d - %s\n", apiErr.Status, apiErr.Message)
        fmt.Printf("Details: %v\n", apiErr.Data)
    } else {
        fmt.Printf("Unexpected error: %v\n", err)
    }
    return
}
```

## Context Support

All methods accept a context for cancellation and timeouts:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

games, err := client.Games().ListWithContext(ctx, params)
```

## Running Tests

```bash
go test ./...
```

## License

MIT
