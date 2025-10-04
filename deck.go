// Package deck provides functionality for working with standard playing cards.
// It implements a deck of 52 cards with support for common operations like
// shuffling, drawing, and sorting.
//
// The package is designed for building card games and provides:
//   - Secure shuffling using crypto/rand
//   - Custom RNG interface for deterministic shuffles
//   - Efficient binary encoding for network transfer
//   - Space-optimized Card representation (1 byte per card)
//
// Example usage:
//
//	d := deck.New()
//	d.SecureShuffle() // Cryptographically secure shuffle
//	card, _ := d.Draw()
package deck

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	mathrand "math/rand"
	"sort"
	"strings"
	"time"
)

// Suit represents the suit of a playing card.
type Suit uint8

const (
	Spades Suit = iota
	Hearts
	Diamonds
	Clubs
)

// String returns the string representation of a Suit.
func (s Suit) String() string {
	return [...]string{"Spades", "Hearts", "Diamonds", "Clubs"}[s]
}

// Symbol returns the Unicode symbol for a Suit.
func (s Suit) Symbol() string {
	return [...]string{"♠", "♥", "♦", "♣"}[s]
}

// Rank represents the rank of a playing card.
type Rank uint8

const (
	Ace Rank = iota + 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

// String returns the string representation of a Rank.
func (r Rank) String() string {
	return [...]string{"", "Ace", "2", "3", "4", "5", "6", "7", "8", "9", "10", "Jack", "Queen", "King"}[r]
}

const (
	// suitShift is the number of bits to shift for suit encoding.
	suitShift = 6
	// rankMask is the bitmask for extracting rank from the card encoding.
	rankMask = 0x3F
)

// Card represents a single playing card using an efficient 1-byte representation.
// This compact format is ideal for memory efficiency and network transfer.
// The upper 2 bits represent the suit (0-3), and the lower 6 bits represent the rank (1-13).
// Visual representation of the bit layout:
//
//	Bit position:  7 6 | 5 4 3 2 1 0
//	              └─┬─┘ └────┬─────┘
//	              Suit      Rank
//
// Examples:
//
//	Ace of Spades   (Rank=1,  Suit=0): 0b00_000001 = 0x01
//	King of Spades  (Rank=13, Suit=0): 0b00_001101 = 0x0D
//	Ace of Hearts   (Rank=1,  Suit=1): 0b01_000001 = 0x41
//	Queen of Clubs  (Rank=12, Suit=3): 0b11_001100 = 0xCC
type Card uint8

// NewCard creates a new Card from a Rank and Suit.
func NewCard(rank Rank, suit Suit) Card {
	return Card((suit << suitShift) | (Suit(rank) & rankMask))
}

// Rank returns the rank of the card.
func (c Card) Rank() Rank {
	return Rank(c & rankMask)
}

// Suit returns the suit of the card.
func (c Card) Suit() Suit {
	return Suit(c >> suitShift)
}

// String returns the string representation of a Card.
func (c Card) String() string {
	return fmt.Sprintf("%s of %s", c.Rank(), c.Suit())
}

// ShortString returns a compact representation of a Card.
func (c Card) ShortString() string {
	return fmt.Sprintf("%s%s", c.Rank(), c.Suit().Symbol())
}

// Shuffler is an interface for custom random number generators.
// Implement this interface to provide deterministic or custom shuffling behavior.
type Shuffler interface {
	// Shuffle randomizes the order of n elements in a slice.
	// The parameter n is the number of elements to shuffle.
	// The swap function should be called to exchange elements at positions i and j
	// according to the implementation's random number generation strategy.
	Shuffle(n int, swap func(i, j int))
}

// SecureShuffler uses crypto/rand for cryptographically secure shuffling.
// This is suitable for applications requiring unpredictable randomness,
// such as online card games or gambling applications.
type SecureShuffler struct{}

// Shuffle implements the Shuffler interface using crypto/rand.
func (s SecureShuffler) Shuffle(n int, swap func(i, j int)) {
	for i := n - 1; i > 0; i-- {
		// Choose a random number j in the range [0, i] with 0 allocations
		// which is more efficient than the [rand.Int] implementation.
		var b [8]byte
		_, _ = rand.Read(b[:]) // [rand.Read] never returns an error https://pkg.go.dev/crypto/rand#Read
		// modulo bias percentage is negligible: 1 / (2^64 / 52) ≈ 0.00000000028%
		j := int(binary.LittleEndian.Uint64(b[:]) % uint64(i+1))
		swap(i, j)
	}
}

// DefaultShuffler uses math/rand with time-based seeding.
// This is not secure enough and is only suitable for trivial applications.
type DefaultShuffler struct {
	rng *mathrand.Rand
}

// NewDefaultShuffler creates a new DefaultShuffler with time-based seed.
func NewDefaultShuffler() *DefaultShuffler {
	return &DefaultShuffler{
		rng: mathrand.New(mathrand.NewSource(time.Now().UnixNano())),
	}
}

// NewSeededShuffler creates a new DefaultShuffler with a specific seed.
// This is useful for reproducible shuffles in testing or replays.
func NewSeededShuffler(seed int64) *DefaultShuffler {
	return &DefaultShuffler{
		rng: mathrand.New(mathrand.NewSource(seed)),
	}
}

// Shuffle implements the Shuffler interface using math/rand.
func (s *DefaultShuffler) Shuffle(n int, swap func(i, j int)) {
	s.rng.Shuffle(n, swap)
}

// Deck represents a deck of playing cards.
// It uses a slice for efficient operations like shuffling and drawing.
type Deck struct {
	cards []Card
}

// New creates and returns a new standard 52-card deck.
// The deck is created in sorted order (Spades, Hearts, Diamonds, Clubs,
// each with Ace through King).
func New() *Deck {
	cards := make([]Card, 0, 52)
	for suit := Spades; suit <= Clubs; suit++ {
		for rank := Ace; rank <= King; rank++ {
			cards = append(cards, NewCard(rank, suit))
		}
	}
	return &Deck{cards: cards}
}

// NewMultiple creates a deck with multiple standard 52-card decks.
// Returns an error if count is less than 1.
func NewMultiple(count int) (*Deck, error) {
	if count < 1 {
		return nil, fmt.Errorf("count must be at least 1, got %d", count)
	}

	cards := make([]Card, 0, 52*count)
	for i := 0; i < count; i++ {
		for suit := Spades; suit <= Clubs; suit++ {
			for rank := Ace; rank <= King; rank++ {
				cards = append(cards, NewCard(rank, suit))
			}
		}
	}
	return &Deck{cards: cards}, nil
}

// Len returns the number of cards currently in the deck.
func (d *Deck) Len() int {
	return len(d.cards)
}

// IsEmpty returns true if the deck has no cards.
func (d *Deck) IsEmpty() bool {
	return len(d.cards) == 0
}

// Shuffle randomizes the order of cards in the deck using math/rand.
// It uses the current time as a seed for the random number generator.
// For cryptographically secure shuffling, use SecureShuffle instead.
func (d *Deck) Shuffle() {
	shuffler := NewDefaultShuffler()
	d.ShuffleWith(shuffler)
}

// SecureShuffle randomizes the order of cards using crypto/rand.
// This provides cryptographically secure randomness suitable for
// security-sensitive applications like online gambling.
func (d *Deck) SecureShuffle() {
	d.ShuffleWith(SecureShuffler{})
}

// ShuffleWithSeed randomizes the order of cards in the deck using the provided seed.
// This allows for reproducible shuffles when needed (e.g., for testing).
func (d *Deck) ShuffleWithSeed(seed int64) {
	shuffler := NewSeededShuffler(seed)
	d.ShuffleWith(shuffler)
}

// ShuffleWith randomizes the order of cards using a custom Shuffler.
// This allows clients to provide their own random number generation strategy.
func (d *Deck) ShuffleWith(shuffler Shuffler) {
	shuffler.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

// Draw removes and returns the top card from the deck.
// Returns an error if the deck is empty.
func (d *Deck) Draw() (Card, error) {
	if d.IsEmpty() {
		return Card(0), fmt.Errorf("cannot draw from empty deck")
	}

	card := d.cards[0]
	d.cards = d.cards[1:]
	return card, nil
}

// DrawN removes and returns n cards from the top of the deck.
// Returns an error if there are fewer than n cards in the deck.
func (d *Deck) DrawN(n int) ([]Card, error) {
	if n < 0 {
		return nil, fmt.Errorf("cannot draw negative number of cards: %d", n)
	}
	if n > len(d.cards) {
		return nil, fmt.Errorf("not enough cards in deck: have %d, need %d", len(d.cards), n)
	}

	cards := make([]Card, n)
	copy(cards, d.cards[:n])
	d.cards = d.cards[n:]
	return cards, nil
}

// Peek returns the top card without removing it from the deck.
// Returns an error if the deck is empty.
func (d *Deck) Peek() (Card, error) {
	if d.IsEmpty() {
		return Card(0), fmt.Errorf("cannot peek at empty deck")
	}
	return d.cards[0], nil
}

// PeekN returns the top n cards without removing them from the deck.
// Returns an error if there are fewer than n cards in the deck.
func (d *Deck) PeekN(n int) ([]Card, error) {
	if n < 0 {
		return nil, fmt.Errorf("cannot peek negative number of cards: %d", n)
	}
	if n > len(d.cards) {
		return nil, fmt.Errorf("not enough cards in deck: have %d, need %d", len(d.cards), n)
	}

	cards := make([]Card, n)
	copy(cards, d.cards[:n])
	return cards, nil
}

// Add adds a card to the bottom of the deck.
func (d *Deck) Add(card Card) {
	d.cards = append(d.cards, card)
}

// AddToTop adds a card to the top of the deck.
func (d *Deck) AddToTop(card Card) {
	d.cards = append([]Card{card}, d.cards...)
}

// Sort sorts the deck by suit (Spades, Hearts, Diamonds, Clubs) and then by rank.
func (d *Deck) Sort() {
	sort.Slice(d.cards, func(i, j int) bool {
		if d.cards[i].Suit() != d.cards[j].Suit() {
			return d.cards[i].Suit() < d.cards[j].Suit()
		}
		return d.cards[i].Rank() < d.cards[j].Rank()
	})
}

// Cards returns a copy of all cards in the deck.
// The returned slice is a copy to prevent external modification.
func (d *Deck) Cards() []Card {
	cards := make([]Card, len(d.cards))
	copy(cards, d.cards)
	return cards
}

// String returns a string representation of the deck.
func (d *Deck) String() string {
	if d.IsEmpty() {
		return "Empty Deck"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Deck (%d cards): [", len(d.cards)))
	for i, card := range d.cards {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(card.ShortString())
	}
	sb.WriteString("]")
	return sb.String()
}

// Filter applies a filter function to create a new deck containing only cards
// that satisfy the predicate.
func (d *Deck) Filter(predicate func(Card) bool) *Deck {
	filtered := make([]Card, 0, len(d.cards))
	for _, card := range d.cards {
		if predicate(card) {
			filtered = append(filtered, card)
		}
	}
	return &Deck{cards: filtered}
}

// MarshalBinary implements encoding.BinaryMarshaler.
// This provides efficient binary encoding for network transfer.
// Format: 4 bytes for length (uint32) + 1 byte per card.
func (d *Deck) MarshalBinary() ([]byte, error) {
	// 4 bytes for length + 1 byte per card
	data := make([]byte, 4+len(d.cards))
	binary.LittleEndian.PutUint32(data[0:4], uint32(len(d.cards)))
	for i, card := range d.cards {
		data[4+i] = byte(card)
	}
	return data, nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
// This decodes the binary format produced by MarshalBinary.
func (d *Deck) UnmarshalBinary(data []byte) error {
	if len(data) < 4 {
		return fmt.Errorf("invalid data: too short")
	}

	count := binary.LittleEndian.Uint32(data[0:4])
	if len(data) != int(4+count) {
		return fmt.Errorf("invalid data: expected %d bytes, got %d", 4+count, len(data))
	}

	d.cards = make([]Card, count)
	for i := uint32(0); i < count; i++ {
		d.cards[i] = Card(data[4+i])
	}
	return nil
}

// Size returns the byte size of the deck when marshaled.
// This is useful for network transfer size estimation.
func (d *Deck) Size() int {
	return 4 + len(d.cards) // 4 bytes header + 1 byte per card
}
