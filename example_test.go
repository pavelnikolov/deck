package deck_test

import (
	"fmt"

	"github.com/pavelnikolov/deck"
)

func ExampleNew() {
	d := deck.New()
	fmt.Printf("Created a new deck with %d cards\n", d.Len())
	// Output: Created a new deck with 52 cards
}

func ExampleDeck_Shuffle() {
	d := deck.New()
	d.Shuffle()
	fmt.Println("Deck shuffled")
	// Output: Deck shuffled
}

func ExampleDeck_Draw() {
	d := deck.New()
	card, err := d.Draw()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Drew: %s\n", card)
	fmt.Printf("Cards remaining: %d\n", d.Len())
	// Output: Drew: Ace of Spades
	// Cards remaining: 51
}

func ExampleDeck_DrawN() {
	d := deck.New()
	cards, err := d.DrawN(5)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Drew %d cards\n", len(cards))
	fmt.Printf("Cards remaining: %d\n", d.Len())
	// Output: Drew 5 cards
	// Cards remaining: 47
}

func ExampleDeck_Peek() {
	d := deck.New()
	card, err := d.Peek()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Top card: %s\n", card)
	fmt.Printf("Cards in deck: %d\n", d.Len())
	// Output: Top card: Ace of Spades
	// Cards in deck: 52
}

func ExampleDeck_Sort() {
	d := deck.New()
	d.Shuffle()
	d.Sort()

	// Peek at first and last cards
	first, _ := d.Peek()
	cards := d.Cards()
	last := cards[len(cards)-1]

	fmt.Printf("First card: %s\n", first)
	fmt.Printf("Last card: %s\n", last)
	// Output: First card: Ace of Spades
	// Last card: King of Clubs
}

func ExampleDeck_Filter() {
	d := deck.New()

	// Filter for only Aces
	aces := d.Filter(func(c deck.Card) bool {
		return c.Rank() == deck.Ace
	})

	fmt.Printf("Number of Aces: %d\n", aces.Len())
	// Output: Number of Aces: 4
}

func ExampleDeck_Add() {
	d := &deck.Deck{}
	card := deck.NewCard(deck.Ace, deck.Spades)
	d.Add(card)
	fmt.Printf("Deck has %d card(s)\n", d.Len())
	// Output: Deck has 1 card(s)
}

func ExampleNewMultiple() {
	d, err := deck.NewMultiple(2)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Created a deck with %d cards\n", d.Len())
	// Output: Created a deck with 104 cards
}

func ExampleCard_String() {
	card := deck.NewCard(deck.King, deck.Hearts)
	fmt.Println(card.String())
	// Output: King of Hearts
}

func ExampleCard_ShortString() {
	card := deck.NewCard(deck.Ace, deck.Spades)
	fmt.Println(card.ShortString())
	// Output: Ace♠
}

func ExampleDeck_ShuffleWithSeed() {
	// Create two decks with same seed for reproducible shuffle
	d1 := deck.New()
	d2 := deck.New()

	d1.ShuffleWithSeed(12345)
	d2.ShuffleWithSeed(12345)

	// Both decks will have identical order
	card1, _ := d1.Draw()
	card2, _ := d2.Draw()

	fmt.Printf("Same top card: %v\n", card1 == card2)
	// Output: Same top card: true
}

// Example demonstrating secure shuffle for security-sensitive applications
func ExampleDeck_SecureShuffle() {
	d := deck.New()

	// Use cryptographically secure shuffle
	d.SecureShuffle()

	card, _ := d.Draw()
	fmt.Printf("Drew a card: %s\n", card)
	// Output will vary due to secure randomness
}

// Example demonstrating custom RNG implementation
func ExampleDeck_ShuffleWith() {
	d := deck.New()

	// Use a seeded shuffler for reproducible results
	shuffler := deck.NewSeededShuffler(42)
	d.ShuffleWith(shuffler)

	// Always draws the same first card with this seed
	card, _ := d.Draw()
	fmt.Printf("First card with seed 42: %s\n", card)
	// Output will be consistent across runs
}

// Example showing binary marshaling for network transfer
func ExampleDeck_MarshalBinary() {
	d := deck.New()
	d.Shuffle()

	// Marshal to binary format (efficient for network transfer)
	data, err := d.MarshalBinary()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Deck serialized to %d bytes\n", len(data))
	fmt.Printf("Size calculated as: %d bytes\n", d.Size())
	// Output: Deck serialized to 56 bytes
	// Size calculated as: 56 bytes
}

// Example showing binary unmarshaling
func ExampleDeck_UnmarshalBinary() {
	// Create and marshal a deck
	d1 := deck.New()
	d1.Shuffle()
	data, _ := d1.MarshalBinary()

	// Unmarshal into a new deck
	d2 := &deck.Deck{}
	err := d2.UnmarshalBinary(data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Restored deck with %d cards\n", d2.Len())
	// Output: Restored deck with 52 cards
}

// Example: Texas Hold'em poker game setup
func ExampleDeck_poker() {
	d := deck.New()
	d.SecureShuffle() // Use secure shuffle for fair play

	// Deal 2 cards to 4 players
	players := 4
	hands := make([][]deck.Card, players)

	for i := 0; i < players; i++ {
		hand, _ := d.DrawN(2)
		hands[i] = hand
	}

	// Burn a card
	_, _ = d.Draw()

	// Deal the flop (3 community cards)
	flop, _ := d.DrawN(3)

	fmt.Printf("Players: %d\n", players)
	fmt.Printf("Flop: %d cards\n", len(flop))
	fmt.Printf("Cards remaining: %d\n", d.Len())
	// Output: Players: 4
	// Flop: 3 cards
	// Cards remaining: 40
}

// Example: Blackjack game setup
func ExampleDeck_blackjack() {
	// Blackjack typically uses 6 decks
	d, _ := deck.NewMultiple(6)
	d.SecureShuffle()

	// Deal initial hands (2 cards each for dealer and player)
	playerHand, _ := d.DrawN(2)
	dealerHand, _ := d.DrawN(2)

	fmt.Printf("Player: %s, %s\n", playerHand[0].ShortString(), playerHand[1].ShortString())
	fmt.Printf("Dealer shows: %s\n", dealerHand[0].ShortString())
	fmt.Printf("Cards remaining: %d\n", d.Len())
	// Output will vary based on shuffle
}

// Example: Filtering cards for a specific game
func ExampleDeck_Filter_suit() {
	d := deck.New()

	// Get all hearts for a hearts-only game variant
	hearts := d.Filter(func(c deck.Card) bool {
		return c.Suit() == deck.Hearts
	})

	fmt.Printf("Hearts in deck: %d\n", hearts.Len())
	// Output: Hearts in deck: 13
}

// Example: Network transfer of deck state
func ExampleDeck_network() {
	// Server side: create and shuffle deck
	serverDeck := deck.New()
	serverDeck.SecureShuffle()

	// Serialize for network transfer
	data, _ := serverDeck.MarshalBinary()

	fmt.Printf("Sending %d bytes over network\n", len(data))

	// Client side: receive and restore deck
	clientDeck := &deck.Deck{}
	_ = clientDeck.UnmarshalBinary(data)

	fmt.Printf("Client received deck with %d cards\n", clientDeck.Len())
	// Output: Sending 56 bytes over network
	// Client received deck with 52 cards
}
