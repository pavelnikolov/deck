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
	RedJoker
	BlackJoker
)

// String returns the string representation of a Rank.
func (r Rank) String() string {
	return [...]string{"", "Ace", "2", "3", "4", "5", "6", "7", "8", "9", "10", "Jack", "Queen", "King", "Joker", "Joker"}[r]
}

const (
	// suitShift is the number of bits to shift for suit encoding.
	suitShift = 6
	// rankMask is the bitmask for extracting rank from the card encoding.
	rankMask = 0x3F
)

const (
	// maxCardsPerPlayer is the maximum number of cards per player allowed in Deal.
	// Set to 52 to handle single-player scenarios with a full deck.
	maxCardsPerPlayer = 52
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

// NewRedJoker creates a red joker card (Hearts suit, Rank 14).
func NewRedJoker() Card {
	return NewCard(RedJoker, Hearts)
}

// NewBlackJoker creates a black joker card (Spades suit, Rank 15).
func NewBlackJoker() Card {
	return NewCard(BlackJoker, Spades)
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
	rank := c.Rank()

	// Handle jokers with color
	if rank == RedJoker {
		return "Joker (Red)"
	}
	if rank == BlackJoker {
		return "Joker (Black)"
	}

	return fmt.Sprintf("%s of %s", c.Rank(), c.Suit())
}

// ShortString returns a compact representation of a Card.
func (c Card) ShortString() string {
	rank := c.Rank()

	// Handle jokers specially
	if rank == RedJoker {
		return "JKR"
	}
	if rank == BlackJoker {
		return "JKB"
	}

	return fmt.Sprintf("%s%s", c.Rank(), c.Suit().Symbol())
}

// IsJoker returns true if the card is a joker (Rank >= 14).
func (c Card) IsJoker() bool {
	return c.Rank() >= RedJoker
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
		for s := Spades; s <= Clubs; s++ {
			for r := Ace; r <= King; r++ {
				cards = append(cards, NewCard(r, s))
			}
		}
	}
	return &Deck{cards: cards}, nil
}

// NewWithJokers creates a standard 54-card deck (52 regular cards + 2 jokers).
// The jokers are added at the end: one red joker (Hearts) and one black joker (Spades).
func NewWithJokers() *Deck {
	cards := make([]Card, 0, 54)
	// Add all 52 regular cards
	for suit := Spades; suit <= Clubs; suit++ {
		for rank := Ace; rank <= King; rank++ {
			cards = append(cards, NewCard(rank, suit))
		}
	}
	// Add jokers
	cards = append(cards, NewRedJoker(), NewBlackJoker())
	return &Deck{cards: cards}
}

// NewMultipleWithJokers creates a deck with multiple 54-card decks (including jokers).
// Returns an error if count is less than 1.
func NewMultipleWithJokers(count int) (*Deck, error) {
	if count < 1 {
		return nil, fmt.Errorf("count must be at least 1, got %d", count)
	}

	cards := make([]Card, 0, 54*count)
	for i := 0; i < count; i++ {
		// Add all 52 regular cards
		for s := Spades; s <= Clubs; s++ {
			for r := Ace; r <= King; r++ {
				cards = append(cards, NewCard(r, s))
			}
		}
		// Add jokers for this deck
		cards = append(cards, NewRedJoker(), NewBlackJoker())
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

// MustDraw removes and returns the top card from the deck.
// It panics if the deck is empty.
//
// Use MustDraw when you are certain the deck is not empty, such as after
// creating a new deck or when you've verified the deck has cards.
// For runtime error handling, use Draw instead.
//
// Example:
//
//	d := deck.New()
//	card := d.MustDraw() // Safe - deck has 52 cards
//
// Panics with: "cannot draw from empty deck"
func (d *Deck) MustDraw() Card {
	card, err := d.Draw()
	if err != nil {
		panic(err.Error())
	}
	return card
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

// MustDrawN removes and returns n cards from the top of the deck.
// It panics if there are not enough cards or if n is negative.
//
// Use MustDrawN when you are certain the deck has at least n cards.
// For runtime error handling, use DrawN instead.
//
// Example:
//
//	d := deck.New()
//	hand := d.MustDrawN(5) // Safe - deck has 52 cards
//
// Panics with:
//   - "cannot draw negative number of cards: X" (if n < 0)
//   - "not enough cards in deck: have X, need Y" (if insufficient cards)
func (d *Deck) MustDrawN(n int) []Card {
	cards, err := d.DrawN(n)
	if err != nil {
		panic(err.Error())
	}
	return cards
}

// Deal distributes cards from the deck to multiple players.
// It removes n * cards from the top of the deck and returns
// them as a slice of hands (each hand is a slice of Cards).
// Cards are dealt in sequential blocks: player 1 gets the first cards each,
// player 2 gets the next cards each, and so on.
// If validation fails, the deck remains unchanged and an error is returned.
//
// Parameters:
//   - n: number of players to deal to
//   - cards: number of cards each player receives
//
// Returns:
//   - [][]Card: slice of hands, where each hand is an independent slice of Cards
//   - error: validation error if parameters are invalid or insufficient cards
//
// Example:
//
//	d := deck.New()
//	hands, err := d.Deal(4, 5) // Deal 4 hands of 5 cards each (poker)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// hands[0] contains player 1's 5 cards
//	// hands[1] contains player 2's 5 cards
//	// deck now has 32 cards remaining
func (d *Deck) Deal(n, cards int) ([][]Card, error) {
	if n < 1 {
		return nil, fmt.Errorf("number of players must be at least 1")
	}

	if cards < 1 {
		return nil, fmt.Errorf("cards per player must be at least 1")
	}

	if cards > maxCardsPerPlayer {
		return nil, fmt.Errorf("cards per player exceeds maximum of %d", maxCardsPerPlayer)
	}

	totalCards := n * cards
	if totalCards > len(d.cards) {
		return nil, fmt.Errorf("insufficient cards: need %d, have %d", totalCards, len(d.cards))
	}

	hands := make([][]Card, n)
	for i := 0; i < n; i++ {
		start := i * cards
		end := start + cards

		hand := make([]Card, cards)
		copy(hand, d.cards[start:end])
		hands[i] = hand
	}

	d.cards = d.cards[totalCards:]

	return hands, nil
}

// MustDeal distributes cards from the deck to multiple players.
// It panics if parameters are invalid or if there are insufficient cards.
//
// Cards are dealt in sequential blocks: player 1 gets the first cardsPerPlayer
// cards, player 2 gets the next cardsPerPlayer cards, and so on.
//
// Use MustDeal when you are certain the deck has enough cards, such as
// dealing a bridge game from a fresh 52-card deck. For runtime error
// handling, use Deal instead.
//
// Example:
//
//	d := deck.New()
//	hands := d.MustDeal(4, 13) // Bridge - safe with 52 cards
//
// Panics with:
//   - "number of players must be at least 1" (if numPlayers < 1)
//   - "cards per player must be at least 1" (if cardsPerPlayer < 1)
//   - "number of players exceeds maximum of 26" (if numPlayers > 26)
//   - "cards per player exceeds maximum of 52" (if cardsPerPlayer > 52)
//   - "insufficient cards: need X, have Y" (if insufficient cards)
func (d *Deck) MustDeal(numPlayers, cardsPerPlayer int) [][]Card {
	hands, err := d.Deal(numPlayers, cardsPerPlayer)
	if err != nil {
		panic(err.Error())
	}
	return hands
}

// DealHands distributes cards from the deck to multiple players with variable hand sizes.
// It removes sum(sizes) cards from the top of the deck and returns them as a slice
// of hands where each hand has a different number of cards as specified in sizes.
//
// Cards are dealt in sequential blocks: the first player gets the first sizes[0] cards,
// the second player gets the next sizes[1] cards, and so on.
//
// If validation fails, the deck remains unchanged and an error is returned.
//
// Parameters:
//   - sizes: slice of integers where sizes[i] specifies cards for player i
//     must be non-empty, all values must be positive and ≤ 52
//
// Returns:
//   - [][]Card: slice of hands where hands[i] contains sizes[i] cards
//   - error: validation error if parameters are invalid or insufficient cards
//
// Example:
//
//	d := deck.New()
//	// Deal casino-style: 3 players get 2 cards, dealer gets 1
//	hands, err := d.DealHands([]int{2, 2, 2, 1})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// hands[0-2] contain 2 cards each (players)
//	// hands[3] contains 1 card (dealer)
//	// deck now has 45 cards remaining (52 - 7 = 45)
func (d *Deck) DealHands(handSizes []int) ([][]Card, error) {
	// Validation: non-empty slice
	if len(handSizes) < 1 {
		return nil, fmt.Errorf("handSizes must contain at least one hand")
	}

	// Calculate total cards needed and validate each hand size
	totalCards := 0
	for i, handSize := range handSizes {
		if handSize <= 0 {
			return nil, fmt.Errorf("hand size must be positive: got %d at index %d", handSize, i)
		}
		if handSize > maxCardsPerPlayer {
			return nil, fmt.Errorf("hand size (%d) at index %d exceeds maximum of %d", handSize, i, maxCardsPerPlayer)
		}
		totalCards += handSize
	}

	// Validation: sufficient cards
	if totalCards > len(d.cards) {
		return nil, fmt.Errorf("insufficient cards: need %d, have %d", totalCards, len(d.cards))
	}

	// Allocate result slice
	hands := make([][]Card, len(handSizes))

	// Distribute cards in sequential blocks
	offset := 0
	for i, handSize := range handSizes {
		hand := make([]Card, handSize)
		copy(hand, d.cards[offset:offset+handSize])
		hands[i] = hand
		offset += handSize
	}

	// Remove dealt cards from deck
	d.cards = d.cards[offset:]

	return hands, nil
}

// MustDealHands distributes cards from the deck to multiple players with variable hand sizes.
// It panics if parameters are invalid or if there are insufficient cards.
//
// Cards are dealt in sequential blocks: the first player gets the first handSizes[0] cards,
// the second player gets the next handSizes[1] cards, and so on.
//
// Use MustDealHands when you are certain the deck has enough cards and the hand sizes
// are valid, such as dealing predetermined hands in a game setup. For runtime error
// handling, use DealHands instead.
//
// Example:
//
//	d := deck.New()
//	hands := d.MustDealHands([]int{2, 2, 2, 1}) // Casino-style - safe with 52 cards
//
// Panics with:
//   - "handSizes must contain at least one hand" (if slice is empty)
//   - "hand size must be positive: got X at index Y" (if handSizes[i] <= 0)
//   - "hand size (X) at index Y exceeds maximum of 52" (if handSizes[i] > 52)
//   - "insufficient cards: need X, have Y" (if insufficient cards)
func (d *Deck) MustDealHands(handSizes []int) [][]Card {
	hands, err := d.DealHands(handSizes)
	if err != nil {
		panic(err.Error())
	}
	return hands
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

// AddJoker adds a joker card to the bottom of the deck.
// The rank parameter should be RedJoker or BlackJoker.
func (d *Deck) AddJoker(rank Rank) {
	var joker Card
	if rank == RedJoker {
		joker = NewRedJoker()
	} else {
		joker = NewBlackJoker()
	}
	d.Add(joker)
}

// AddToTop adds a card to the top of the deck.
func (d *Deck) AddToTop(card Card) {
	d.cards = append([]Card{card}, d.cards...)
}

// Sort sorts the deck by suit (Spades, Hearts, Diamonds, Clubs) and then by rank.
// Jokers are sorted to the end of the deck (Red Joker before Black Joker).
func (d *Deck) Sort() {
	sort.Slice(d.cards, func(i, j int) bool {
		iRank, jRank := d.cards[i].Rank(), d.cards[j].Rank()
		iJoker, jJoker := iRank >= RedJoker, jRank >= RedJoker

		// If one is a joker and the other isn't, non-joker comes first
		if iJoker != jJoker {
			return !iJoker // non-joker (false) comes before joker (true)
		}

		// If both are jokers, sort by rank (Red Joker=14 < Black Joker=15)
		if iJoker && jJoker {
			return iRank < jRank
		}

		// Regular cards: sort by suit, then rank
		if d.cards[i].Suit() != d.cards[j].Suit() {
			return d.cards[i].Suit() < d.cards[j].Suit()
		}
		return iRank < jRank
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
