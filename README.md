# deck

A high-performance, idiomatic Go package for working with standard playing cards, designed for building card games.

[![Go Reference](https://pkg.go.dev/badge/github.com/pavelnikolov/deck.svg)](https://pkg.go.dev/github.com/pavelnikolov/deck)
[![Go Report Card](https://goreportcard.com/badge/github.com/pavelnikolov/deck)](https://goreportcard.com/report/github.com/pavelnikolov/deck)
[![Coverage](https://img.shields.io/badge/coverage-100%25-brightgreen)](https://github.com/pavelnikolov/deck)

## Features

- 🎯 **Idiomatic Go**: Follows best practices from Effective Go
- 🔒 **Secure Shuffling**: Cryptographically secure shuffling using `crypto/rand`
- 🎲 **Custom RNG Interface**: Bring-your-own random number generator
- 📦 **Space Optimized**: 1-byte card representation (vs 16 bytes for struct)
- 🌐 **Network Efficient**: Binary marshaling for efficient transfer (56 bytes for 52 cards)
- ✅ **Fully Tested**: 100% test coverage with comprehensive examples
- ⚡ **High Performance**: Optimized data structures and algorithms
- 📚 **Well Documented**: Complete godoc documentation with examples

## Installation

```bash
go get github.com/pavelnikolov/deck
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/pavelnikolov/deck"
)

func main() {
    // Create a new deck
    d := deck.New()
    
    // Secure shuffle for fair play
    d.SecureShuffle()
    
    // Draw a card
    card, _ := d.Draw()
    fmt.Printf("Drew: %s\n", card)
    
    // Draw multiple cards
    hand, _ := d.DrawN(5)
    fmt.Printf("Hand: %v\n", hand)
}
```

## Core Concepts

### Card Representation

Cards use an efficient 1-byte encoding:
- **Space**: 1 byte per card (vs 16 bytes for struct-based representation)
- **Network**: 52-card deck = 56 bytes (4-byte header + 52 bytes)
- **Performance**: Direct value comparison, no pointer indirection

```go
card := deck.NewCard(deck.Ace, deck.Spades)
fmt.Println(card.String())      // "Ace of Spades"
fmt.Println(card.ShortString()) // "Ace♠"
```

### Shuffling Options

#### 1. Default Shuffle (math/rand)

```go
d := deck.New()
d.Shuffle() // Uses time-based seed
```

#### 2. Secure Shuffle (crypto/rand)

```go
d := deck.New()
d.SecureShuffle() // Cryptographically secure - use for online games
```

#### 3. Seeded Shuffle (Reproducible)

```go
d := deck.New()
d.ShuffleWithSeed(12345) // Deterministic - use for testing/replays
```

#### 4. Custom Shuffler (BYO RNG)

```go
type MyShuffler struct{}

func (s MyShuffler) Shuffle(n int, swap func(i, j int)) {
    // Your custom shuffle logic
}

d := deck.New()
d.ShuffleWith(MyShuffler{})
```

## Game Examples

### Texas Hold'em Poker

```go
d := deck.New()
d.SecureShuffle()

// Deal to 4 players
for i := 0; i < 4; i++ {
    hand, _ := d.DrawN(2)
    fmt.Printf("Player %d: %v\n", i+1, hand)
}

// Burn and flop
d.Draw()
flop, _ := d.DrawN(3)
fmt.Printf("Flop: %v\n", flop)
```

### Blackjack (Multiple Decks)

```go
// Blackjack typically uses 6 decks
d, _ := deck.NewMultiple(6)
d.SecureShuffle()

playerHand, _ := d.DrawN(2)
dealerHand, _ := d.DrawN(2)

fmt.Printf("Player: %s, %s\n", playerHand[0], playerHand[1])
fmt.Printf("Dealer shows: %s\n", dealerHand[0])
```

## Network Transfer

### Efficient Binary Serialization

```go
// Server side
serverDeck := deck.New()
serverDeck.SecureShuffle()
data, _ := serverDeck.MarshalBinary()

// Send `data` over network (56 bytes for 52 cards)

// Client side
clientDeck := &deck.Deck{}
_ = clientDeck.UnmarshalBinary(data)
```

### Size Comparison

| Cards | This Package | JSON | Struct-based |
|-------|--------------|------|--------------|
| 52    | 56 bytes     | ~1KB | ~832 bytes   |
| 104   | 108 bytes    | ~2KB | ~1664 bytes  |
| 312   | 316 bytes    | ~6KB | ~5KB         |

## API Reference

### Creating Decks

```go
d := deck.New()                    // Standard 52-card deck
d, _ := deck.NewMultiple(6)        // Multiple decks (e.g., blackjack)
```

### Drawing/dealing Cards

```go
card, err := d.Draw()              // Draw one card
cards, err := d.DrawN(5)           // Draw multiple cards
card, err := d.Peek()              // Peek without removing
cards, err := d.PeekN(5)           // Peek multiple cards
hands, err := d.Deal(4, 5)         // Deal 4 hands of 5 cards each
```

### Must* Methods (Panic on Error)

For scenarios where you're certain the deck has sufficient cards, use the Must* variants that panic instead of returning errors:

```go
d := deck.New()

// No error checking needed for known-safe operations
card := d.MustDraw()               // Panics if deck is empty
hand := d.MustDrawN(5)             // Panics if insufficient cards
hands := d.MustDeal(4, 13)         // Panics if invalid params or insufficient cards
hands := d.MustDealHands(3, 3, 1)  // Panics if invalid params or insufficient cards
```

**When to use Must* methods:**

- ✅ Fresh deck with known card count
- ✅ After verifying deck has sufficient cards  
- ✅ Deterministic game scenarios (e.g., bridge deal: `MustDeal(4, 13)` with 52 cards)
- ❌ Unknown deck state at runtime
- ❌ When graceful error handling is needed

These methods follow Go conventions (like `regexp.MustCompile`) where failure indicates a programming error rather than a runtime condition.

### Deck Operations

```go
d.Shuffle()                        // Standard shuffle
d.SecureShuffle()                  // Cryptographically secure
d.ShuffleWithSeed(seed)            // Reproducible shuffle
d.ShuffleWith(shuffler)            // Custom shuffler

d.Sort()                           // Sort by suit then rank
d.Add(card)                        // Add to bottom
d.AddToTop(card)                   // Add to top
```

### Filtering

```go
// Get all aces
aces := d.Filter(func(c deck.Card) bool {
    return c.Rank() == deck.Ace
})

// Get all hearts
hearts := d.Filter(func(c deck.Card) bool {
    return c.Suit() == deck.Hearts
})
```

### Information

```go
count := d.Len()                   // Number of cards
empty := d.IsEmpty()               // Check if empty
cards := d.Cards()                 // Get copy of all cards
size := d.Size()                   // Binary size in bytes
str := d.String()                  // String representation
```

### Binary Marshaling

```go
data, err := d.MarshalBinary()     // Encode to bytes
err = d.UnmarshalBinary(data)      // Decode from bytes
```

## Performance

Benchmarks on Apple M1 Pro:

```sh
BenchmarkNew-10                 60.80 ns/op
BenchmarkShuffle-10             11793 ns/op
BenchmarkSecureShuffle-10        4959 ns/op
BenchmarkDraw-10                 268.7 ns/op
BenchmarkSort-10                 480.1 ns/op
BenchmarkFilter-10                77.15 ns/op
BenchmarkMarshalBinary-10         40.67 ns/op
BenchmarkUnmarshalBinary-10       53.08 ns/op
```

## Design Decisions

### Why 1-byte Card encoding?

- **Memory**: 94% reduction (1 byte vs 16 bytes for struct)
- **Network**: Direct binary transfer with minimal overhead
- **Performance**: Value-based comparison, CPU cache friendly
- **Simplicity**: No pointer indirection

### Why Shuffler interface?

- **Security**: Different games have different requirements
- **Testing**: Deterministic shuffles for test reproducibility
- **Flexibility**: Custom RNG for specific needs (e.g., regulatory compliance)
- **Best Practice**: Dependency injection for better testing

### Why only standard library?

- **Portability**: Works everywhere Go works
- **Stability**: No external dependency breakage
- **Size**: Minimal binary size
- **Trust**: Audited and maintained by Go team
