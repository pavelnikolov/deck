package deck

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	d := New()

	if d == nil {
		t.Fatal("New() returned nil")
	}

	if d.Len() != 52 {
		t.Errorf("New deck should have 52 cards, got %d", d.Len())
	}

	if d.IsEmpty() {
		t.Error("New deck should not be empty")
	}
}

func TestNewMultiple(t *testing.T) {
	tests := []struct {
		name      string
		count     int
		wantCards int
		wantErr   bool
	}{
		{"single deck", 1, 52, false},
		{"two decks", 2, 104, false},
		{"five decks", 5, 260, false},
		{"zero decks", 0, 0, true},
		{"negative decks", -1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewMultiple(tt.count)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if d.Len() != tt.wantCards {
				t.Errorf("expected %d cards, got %d", tt.wantCards, d.Len())
			}
		})
	}
}

func TestSuitString(t *testing.T) {
	tests := []struct {
		suit Suit
		want string
	}{
		{Spades, "Spades"},
		{Hearts, "Hearts"},
		{Diamonds, "Diamonds"},
		{Clubs, "Clubs"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.suit.String(); got != tt.want {
				t.Errorf("Suit.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSuitSymbol(t *testing.T) {
	tests := []struct {
		suit Suit
		want string
	}{
		{Spades, "♠"},
		{Hearts, "♥"},
		{Diamonds, "♦"},
		{Clubs, "♣"},
	}

	for _, tt := range tests {
		t.Run(tt.suit.String(), func(t *testing.T) {
			if got := tt.suit.Symbol(); got != tt.want {
				t.Errorf("Suit.Symbol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRankString(t *testing.T) {
	tests := []struct {
		rank Rank
		want string
	}{
		{Ace, "Ace"},
		{Two, "2"},
		{Three, "3"},
		{Ten, "10"},
		{Jack, "Jack"},
		{Queen, "Queen"},
		{King, "King"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.rank.String(); got != tt.want {
				t.Errorf("Rank.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCardString(t *testing.T) {
	card := NewCard(Ace, Spades)
	want := "Ace of Spades"

	if got := card.String(); got != want {
		t.Errorf("Card.String() = %v, want %v", got, want)
	}
}

func TestCardShortString(t *testing.T) {
	tests := []struct {
		card Card
		want string
	}{
		{NewCard(Ace, Spades), "Ace♠"},
		{NewCard(King, Hearts), "King♥"},
		{NewCard(Ten, Diamonds), "10♦"},
		{NewCard(Two, Clubs), "2♣"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.card.ShortString(); got != tt.want {
				t.Errorf("Card.ShortString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCardEncoding(t *testing.T) {
	tests := []struct {
		rank Rank
		suit Suit
	}{
		{Ace, Spades},
		{King, Hearts},
		{Ten, Diamonds},
		{Two, Clubs},
		{Queen, Spades},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s of %s", tt.rank, tt.suit), func(t *testing.T) {
			card := NewCard(tt.rank, tt.suit)

			if card.Rank() != tt.rank {
				t.Errorf("Card.Rank() = %v, want %v", card.Rank(), tt.rank)
			}

			if card.Suit() != tt.suit {
				t.Errorf("Card.Suit() = %v, want %v", card.Suit(), tt.suit)
			}
		})
	}
}

func TestDeckLen(t *testing.T) {
	d := New()

	if d.Len() != 52 {
		t.Errorf("expected length 52, got %d", d.Len())
	}

	_, _ = d.Draw()
	if d.Len() != 51 {
		t.Errorf("after drawing, expected length 51, got %d", d.Len())
	}
}

func TestDeckIsEmpty(t *testing.T) {
	d := &Deck{cards: []Card{}}

	if !d.IsEmpty() {
		t.Error("empty deck should return true for IsEmpty()")
	}

	d.Add(NewCard(Ace, Spades))
	if d.IsEmpty() {
		t.Error("deck with cards should return false for IsEmpty()")
	}
}

func TestDeckShuffle(t *testing.T) {
	d1 := New()
	d2 := New()

	// Get original order
	original := d1.Cards()

	// Shuffle first deck
	d1.Shuffle()
	shuffled := d1.Cards()

	// Check all cards are still present
	if d1.Len() != 52 {
		t.Errorf("deck should still have 52 cards after shuffle, got %d", d1.Len())
	}

	// Check that order changed (with very high probability)
	sameOrder := true
	for i := range original {
		if original[i] != shuffled[i] {
			sameOrder = false
			break
		}
	}

	if sameOrder {
		t.Error("shuffle did not change card order (very unlikely)")
	}

	// Verify d2 is still in original order
	d2Cards := d2.Cards()
	for i := range original {
		if original[i] != d2Cards[i] {
			t.Error("second deck should not be affected by first deck shuffle")
			break
		}
	}
}

func TestDeckShuffleWithSeed(t *testing.T) {
	d1 := New()
	d2 := New()

	// Shuffle both with same seed
	seed := int64(12345)
	d1.ShuffleWithSeed(seed)
	d2.ShuffleWithSeed(seed)

	// They should have identical order
	cards1 := d1.Cards()
	cards2 := d2.Cards()

	for i := range cards1 {
		if cards1[i] != cards2[i] {
			t.Errorf("decks shuffled with same seed should have identical order at index %d", i)
		}
	}

	// Shuffle d2 with different seed
	d2.ShuffleWithSeed(54321)
	cards2 = d2.Cards()

	// They should now differ
	sameOrder := true
	for i := range cards1 {
		if cards1[i] != cards2[i] {
			sameOrder = false
			break
		}
	}

	if sameOrder {
		t.Error("decks shuffled with different seeds should have different orders")
	}
}

func TestDeckDraw(t *testing.T) {
	d := New()
	topCard, _ := d.Peek()

	card, err := d.Draw()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if card != topCard {
		t.Error("drawn card should be the top card")
	}

	if d.Len() != 51 {
		t.Errorf("deck should have 51 cards after draw, got %d", d.Len())
	}
}

func TestDeckDrawEmpty(t *testing.T) {
	d := &Deck{cards: []Card{}}

	_, err := d.Draw()
	if err == nil {
		t.Error("expected error when drawing from empty deck")
	}
}

func TestDeckDrawN(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		deckSize int
		wantErr  bool
	}{
		{"draw 5 from full deck", 5, 52, false},
		{"draw all cards", 52, 52, false},
		{"draw 0 cards", 0, 52, false},
		{"draw more than available", 53, 52, true},
		{"draw negative", -1, 52, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New()

			cards, err := d.DrawN(tt.n)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(cards) != tt.n {
				t.Errorf("expected %d cards, got %d", tt.n, len(cards))
			}

			if d.Len() != tt.deckSize-tt.n {
				t.Errorf("expected deck size %d, got %d", tt.deckSize-tt.n, d.Len())
			}
		})
	}
}

func TestDeckPeek(t *testing.T) {
	d := New()
	initialLen := d.Len()

	card1, err := d.Peek()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if d.Len() != initialLen {
		t.Error("peek should not modify deck size")
	}

	card2, _ := d.Peek()
	if card1 != card2 {
		t.Error("consecutive peeks should return same card")
	}
}

func TestDeckPeekEmpty(t *testing.T) {
	d := &Deck{cards: []Card{}}

	_, err := d.Peek()
	if err == nil {
		t.Error("expected error when peeking at empty deck")
	}
}

func TestDeckPeekN(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		deckSize int
		wantErr  bool
	}{
		{"peek 5 cards", 5, 52, false},
		{"peek all cards", 52, 52, false},
		{"peek 0 cards", 0, 52, false},
		{"peek more than available", 53, 52, true},
		{"peek negative", -1, 52, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New()
			initialLen := d.Len()

			cards, err := d.PeekN(tt.n)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(cards) != tt.n {
				t.Errorf("expected %d cards, got %d", tt.n, len(cards))
			}

			if d.Len() != initialLen {
				t.Error("peek should not modify deck size")
			}
		})
	}
}

func TestDeckAdd(t *testing.T) {
	d := New()
	card := NewCard(Ace, Hearts)
	initialLen := d.Len()

	d.Add(card)

	if d.Len() != initialLen+1 {
		t.Errorf("expected deck size %d, got %d", initialLen+1, d.Len())
	}

	// Draw all cards except the last one
	for i := 0; i < initialLen; i++ {
		_, _ = d.Draw()
	}

	// The added card should be at the bottom
	lastCard, _ := d.Draw()
	if lastCard != card {
		t.Error("added card should be at the bottom of the deck")
	}
}

func TestDeckAddToTop(t *testing.T) {
	d := New()
	card := NewCard(Ace, Hearts)
	initialLen := d.Len()

	d.AddToTop(card)

	if d.Len() != initialLen+1 {
		t.Errorf("expected deck size %d, got %d", initialLen+1, d.Len())
	}

	// The added card should be at the top
	topCard, _ := d.Draw()
	if topCard != card {
		t.Error("added card should be at the top of the deck")
	}
}

func TestDeckSort(t *testing.T) {
	d := New()
	d.Shuffle()

	d.Sort()

	cards := d.Cards()

	// Verify sorting by suit then rank
	for i := 1; i < len(cards); i++ {
		prev := cards[i-1]
		curr := cards[i]

		if prev.Suit() > curr.Suit() {
			t.Errorf("cards not sorted by suit at index %d", i)
		}

		if prev.Suit() == curr.Suit() && prev.Rank() > curr.Rank() {
			t.Errorf("cards not sorted by rank within suit at index %d", i)
		}
	}

	// Verify first card is Ace of Spades
	if cards[0].Rank() != Ace || cards[0].Suit() != Spades {
		t.Errorf("first card should be Ace of Spades, got %v", cards[0])
	}

	// Verify last card is King of Clubs
	if cards[51].Rank() != King || cards[51].Suit() != Clubs {
		t.Errorf("last card should be King of Clubs, got %v", cards[51])
	}
}

func TestDeckCards(t *testing.T) {
	d := New()
	cards := d.Cards()

	if len(cards) != 52 {
		t.Errorf("expected 52 cards, got %d", len(cards))
	}

	// Modify returned slice
	cards[0] = NewCard(Ace, Hearts)

	// Original deck should be unchanged
	originalTop, _ := d.Peek()
	if originalTop.Rank() == Ace && originalTop.Suit() == Hearts {
		t.Error("modifying returned slice should not affect original deck")
	}
}

func TestDeckString(t *testing.T) {
	d := &Deck{cards: []Card{}}
	if got := d.String(); got != "Empty Deck" {
		t.Errorf("empty deck string = %v, want 'Empty Deck'", got)
	}

	d.Add(NewCard(Ace, Spades))
	d.Add(NewCard(King, Hearts))

	str := d.String()
	if str == "Empty Deck" {
		t.Error("non-empty deck should not return 'Empty Deck'")
	}

	// Should contain count
	if len(str) == 0 {
		t.Error("deck string should not be empty")
	}
}

func TestDeckFilter(t *testing.T) {
	d := New()

	// Filter for only Aces
	aces := d.Filter(func(c Card) bool {
		return c.Rank() == Ace
	})

	if aces.Len() != 4 {
		t.Errorf("expected 4 aces, got %d", aces.Len())
	}

	// Verify all cards are Aces
	for _, card := range aces.Cards() {
		if card.Rank() != Ace {
			t.Errorf("filtered deck should only contain Aces, got %v", card)
		}
	}

	// Filter for Hearts
	hearts := d.Filter(func(c Card) bool {
		return c.Suit() == Hearts
	})

	if hearts.Len() != 13 {
		t.Errorf("expected 13 hearts, got %d", hearts.Len())
	} // Original deck should be unchanged
	if d.Len() != 52 {
		t.Errorf("original deck should still have 52 cards, got %d", d.Len())
	}
}

func TestSecureShuffle(t *testing.T) {
	d1 := New()
	d2 := New()

	// Get original order
	original := d1.Cards()

	// Secure shuffle first deck
	d1.SecureShuffle()
	shuffled := d1.Cards()

	// Check all cards are still present
	if d1.Len() != 52 {
		t.Errorf("deck should still have 52 cards after shuffle, got %d", d1.Len())
	}

	// Check that order changed (with very high probability)
	sameOrder := true
	for i := range original {
		if original[i] != shuffled[i] {
			sameOrder = false
			break
		}
	}

	if sameOrder {
		t.Error("secure shuffle did not change card order (very unlikely)")
	}

	// Verify d2 is still in original order
	d2Cards := d2.Cards()
	for i := range original {
		if original[i] != d2Cards[i] {
			t.Error("second deck should not be affected by first deck shuffle")
			break
		}
	}
}

func TestCustomShuffler(t *testing.T) {
	d1 := New()
	d2 := New()

	// Use seeded shuffler for deterministic results
	shuffler1 := NewSeededShuffler(12345)
	shuffler2 := NewSeededShuffler(12345)

	d1.ShuffleWith(shuffler1)
	d2.ShuffleWith(shuffler2)

	// Both decks should have identical order
	cards1 := d1.Cards()
	cards2 := d2.Cards()

	for i := range cards1 {
		if cards1[i] != cards2[i] {
			t.Errorf("decks shuffled with same seeded shuffler should have identical order at index %d", i)
		}
	}
}

func TestDeckMarshalBinary(t *testing.T) {
	d := New()
	d.Shuffle()

	// Marshal
	data, err := d.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary() error = %v", err)
	}

	// Check size
	expectedSize := 4 + 52 // 4 byte header + 52 cards
	if len(data) != expectedSize {
		t.Errorf("expected binary size %d, got %d", expectedSize, len(data))
	}

	// Unmarshal into new deck
	d2 := &Deck{}
	if err := d2.UnmarshalBinary(data); err != nil {
		t.Fatalf("UnmarshalBinary() error = %v", err)
	}

	// Verify decks are identical
	if d.Len() != d2.Len() {
		t.Errorf("unmarshaled deck length %d, want %d", d2.Len(), d.Len())
	}

	cards1 := d.Cards()
	cards2 := d2.Cards()

	for i := range cards1 {
		if cards1[i] != cards2[i] {
			t.Errorf("card at index %d differs after unmarshal", i)
		}
	}
}

func TestDeckUnmarshalBinaryErrors(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"too short", []byte{0x01}},
		{"mismatched length", []byte{0x05, 0x00, 0x00, 0x00, 0x01}}, // says 5 cards, provides 1
		{"empty", []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Deck{}
			if err := d.UnmarshalBinary(tt.data); err == nil {
				t.Error("expected error but got nil")
			}
		})
	}
}

func TestDeckSize(t *testing.T) {
	tests := []struct {
		name string
		deck *Deck
		want int
	}{
		{"empty deck", &Deck{cards: []Card{}}, 4},
		{"full deck", New(), 56},
		{"double deck", func() *Deck { d, _ := NewMultiple(2); return d }(), 108},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.deck.Size(); got != tt.want {
				t.Errorf("deck.Size() = %v bytes, want %v bytes", got, tt.want)
			}
		})
	}
}

// Benchmark tests
func BenchmarkNew(b *testing.B) {
	for b.Loop() {
		_ = New()
	}
}

func BenchmarkShuffle(b *testing.B) {
	d := New()
	for b.Loop() {
		d.Shuffle()
	}
}

func BenchmarkDraw(b *testing.B) {
	for b.Loop() {
		d := New()
		for !d.IsEmpty() {
			_, _ = d.Draw()
		}
	}
}

func BenchmarkSort(b *testing.B) {
	d := New()
	d.Shuffle()
	for b.Loop() {
		d.Sort()
	}
}

func BenchmarkFilter(b *testing.B) {
	d := New()
	for b.Loop() {
		_ = d.Filter(func(c Card) bool {
			return c.Rank() == Ace
		})
	}
}

func BenchmarkSecureShuffle(b *testing.B) {
	d := New()
	for b.Loop() {
		d.SecureShuffle()
	}
}

func BenchmarkMarshalBinary(b *testing.B) {
	d := New()
	for b.Loop() {
		_, _ = d.MarshalBinary()
	}
}

func BenchmarkUnmarshalBinary(b *testing.B) {
	d := New()
	data, _ := d.MarshalBinary()
	for b.Loop() {
		d2 := &Deck{}
		_ = d2.UnmarshalBinary(data)
	}
}
