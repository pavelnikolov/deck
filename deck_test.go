package deck

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	d := New()

	if d == nil {
		t.Fatal("New() returned nil, want non-nil deck")
	}

	if got, want := d.Len(), 52; got != want {
		t.Errorf("New().Len() = %d, want %d", got, want)
	}

	if got, want := d.IsEmpty(), false; got != want {
		t.Errorf("New().IsEmpty() = %v, want %v", got, want)
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
					t.Errorf("NewMultiple(%d) got nil error, want error", tt.count)
				}
				return
			}

			if err != nil {
				t.Fatalf("NewMultiple(%d) got error: %v, want nil", tt.count, err)
			}

			if got, want := d.Len(), tt.wantCards; got != want {
				t.Errorf("NewMultiple(%d).Len() = %d, want %d", tt.count, got, want)
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

	if got, want := d.Len(), 52; got != want {
		t.Errorf("New().Len() = %d, want %d", got, want)
	}

	_, _ = d.Draw()
	if got, want := d.Len(), 51; got != want {
		t.Errorf("After Draw(), deck.Len() = %d, want %d", got, want)
	}
}

func TestDeckIsEmpty(t *testing.T) {
	d := &Deck{cards: []Card{}}

	if got, want := d.IsEmpty(), true; got != want {
		t.Errorf("Empty deck.IsEmpty() = %v, want %v", got, want)
	}

	d.Add(NewCard(Ace, Spades))
	if got, want := d.IsEmpty(), false; got != want {
		t.Errorf("After Add(), deck.IsEmpty() = %v, want %v", got, want)
	}
}

func TestDeckShuffle(t *testing.T) {
	d1 := New()
	d2 := New()

	original := d1.Cards()

	d1.Shuffle()
	shuffled := d1.Cards()

	if got, want := d1.Len(), 52; got != want {
		t.Errorf("After Shuffle(), deck.Len() = %d, want %d", got, want)
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
		t.Error("After Shuffle(), card order unchanged (very unlikely - shuffle may not be working)")
	}

	// Verify second deck unaffected
	d2Cards := d2.Cards()
	for i := range original {
		if got, want := d2Cards[i], original[i]; got != want {
			t.Errorf("After d1.Shuffle(), d2.Cards()[%d] = %v, want %v (d2 should be unaffected)", i, got, want)
			break
		}
	}
}

func TestDeckShuffleWithSeed(t *testing.T) {
	d1 := New()
	d2 := New()

	seed := int64(12345)
	d1.ShuffleWithSeed(seed)
	d2.ShuffleWithSeed(seed)

	cards1 := d1.Cards()
	cards2 := d2.Cards()

	// Verify same seed produces same order
	for i := range cards1 {
		if got, want := cards2[i], cards1[i]; got != want {
			t.Errorf("After ShuffleWithSeed(%d), decks differ at index %d: got %v, want %v (same seed should produce same order)", seed, i, got, want)
		}
	}

	// Verify different seed produces different order
	d2.ShuffleWithSeed(54321)
	cards2 = d2.Cards()

	sameOrder := true
	for i := range cards1 {
		if cards1[i] != cards2[i] {
			sameOrder = false
			break
		}
	}

	if sameOrder {
		t.Error("After ShuffleWithSeed(54321), deck order same as ShuffleWithSeed(12345) (different seeds should produce different orders)")
	}
}

func TestDeckDraw(t *testing.T) {
	d := New()
	topCard, _ := d.Peek()

	card, err := d.Draw()
	if err != nil {
		t.Fatalf("Draw() unexpected error: %v", err)
	}

	if got, want := card, topCard; got != want {
		t.Errorf("Draw() = %v, want %v (top card)", got, want)
	}

	if got, want := d.Len(), 51; got != want {
		t.Errorf("After Draw(), deck.Len() = %d, want %d", got, want)
	}
}

func TestDeckDrawEmpty(t *testing.T) {
	d := &Deck{cards: []Card{}}

	_, err := d.Draw()
	if err == nil {
		t.Error("Draw() from empty deck got nil error, want error")
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
					t.Errorf("DrawN(%d) got nil error, want error", tt.n)
				}
				return
			}

			if err != nil {
				t.Fatalf("DrawN(%d) unexpected error: %v", tt.n, err)
			}

			if got, want := len(cards), tt.n; got != want {
				t.Errorf("DrawN(%d) = %d cards, want %d", tt.n, got, want)
			}

			if got, want := d.Len(), tt.deckSize-tt.n; got != want {
				t.Errorf("After DrawN(%d), deck.Len() = %d, want %d", tt.n, got, want)
			}
		})
	}
}

func TestDeckPeek(t *testing.T) {
	d := New()
	initialLen := d.Len()

	card1, err := d.Peek()
	if err != nil {
		t.Fatalf("Peek() got error: %v, want nil", err)
	}

	if got, want := d.Len(), initialLen; got != want {
		t.Errorf("After Peek(), deck.Len() = %d, want %d (peek should not modify deck)", got, want)
	}

	card2, _ := d.Peek()
	if got, want := card2, card1; got != want {
		t.Errorf("Second Peek() = %v, want %v (consecutive peeks should return same card)", got, want)
	}
}

func TestDeckPeekEmpty(t *testing.T) {
	d := &Deck{cards: []Card{}}

	_, err := d.Peek()
	if err == nil {
		t.Error("Peek() on empty deck got nil error, want error")
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
					t.Errorf("PeekN(%d) got nil error, want error", tt.n)
				}
				return
			}

			if err != nil {
				t.Fatalf("PeekN(%d) got error: %v, want nil", tt.n, err)
			}

			if got, want := len(cards), tt.n; got != want {
				t.Errorf("PeekN(%d) = %d cards, want %d", tt.n, got, want)
			}

			if got, want := d.Len(), initialLen; got != want {
				t.Errorf("After PeekN(%d), deck.Len() = %d, want %d (peek should not modify deck)", tt.n, got, want)
			}
		})
	}
}

func TestDeckAdd(t *testing.T) {
	d := New()
	card := NewCard(Ace, Hearts)
	initialLen := d.Len()

	d.Add(card)

	if got, want := d.Len(), initialLen+1; got != want {
		t.Errorf("After Add(), deck.Len() = %d, want %d", got, want)
	}

	// Draw all initial cards
	for i := 0; i < initialLen; i++ {
		_, _ = d.Draw()
	}

	// Check that added card is at the bottom
	got, _ := d.Draw()
	if got != card {
		t.Errorf("After Add(), bottom card = %v, want %v (added card should be at bottom)", got, card)
	}
}

func TestDeckAddToTop(t *testing.T) {
	d := New()
	card := NewCard(Ace, Hearts)
	initialLen := d.Len()

	d.AddToTop(card)

	if got, want := d.Len(), initialLen+1; got != want {
		t.Errorf("After AddToTop(), deck.Len() = %d, want %d", got, want)
	}

	got, _ := d.Draw()
	if got != card {
		t.Errorf("After AddToTop(), top card = %v, want %v (added card should be at top)", got, card)
	}
}

func TestDeckSort(t *testing.T) {
	d := New()
	d.Shuffle()

	d.Sort()

	cards := d.Cards()

	// Check sorting order
	for i := 1; i < len(cards); i++ {
		prev := cards[i-1]
		curr := cards[i]

		if got, want := prev.Suit() > curr.Suit(), false; got != want {
			t.Errorf("After Sort(), cards[%d].Suit() > cards[%d].Suit() = %v (%v > %v), want sorted by suit", i-1, i, got, prev.Suit(), curr.Suit())
		}

		if prev.Suit() == curr.Suit() {
			if got, want := prev.Rank() > curr.Rank(), false; got != want {
				t.Errorf("After Sort(), cards[%d].Rank() > cards[%d].Rank() = %v (%v > %v), want sorted by rank within suit", i-1, i, got, prev.Rank(), curr.Rank())
			}
		}
	}

	// Check first card
	if got, want := cards[0].Rank(), Ace; got != want {
		t.Errorf("After Sort(), first card rank = %v, want %v", got, want)
	}
	if got, want := cards[0].Suit(), Spades; got != want {
		t.Errorf("After Sort(), first card suit = %v, want %v", got, want)
	}

	// Check last card
	if got, want := cards[51].Rank(), King; got != want {
		t.Errorf("After Sort(), last card rank = %v, want %v", got, want)
	}
	if got, want := cards[51].Suit(), Clubs; got != want {
		t.Errorf("After Sort(), last card suit = %v, want %v", got, want)
	}
}

func TestDeckCards(t *testing.T) {
	d := New()
	cards := d.Cards()

	if got, want := len(cards), 52; got != want {
		t.Errorf("Cards() returned %d cards, want %d", got, want)
	}

	// Modify returned slice
	cards[0] = NewCard(Ace, Hearts)

	// Verify deck is not affected
	originalTop, _ := d.Peek()
	if got := originalTop.Rank() == Ace && originalTop.Suit() == Hearts; got {
		t.Errorf("After modifying Cards() slice, deck was affected (top card = %v), want original deck unchanged", originalTop)
	}
}

func TestDeckString(t *testing.T) {
	d := &Deck{cards: []Card{}}
	if got, want := d.String(), "Empty Deck"; got != want {
		t.Errorf("Empty deck String() = %q, want %q", got, want)
	}

	d.Add(NewCard(Ace, Spades))
	d.Add(NewCard(King, Hearts))

	str := d.String()
	if str == "Empty Deck" {
		t.Errorf("Non-empty deck String() = %q, want non-empty string", str)
	}

	if got, want := len(str), 0; got == want {
		t.Errorf("Non-empty deck String() length = %d, want > 0", got)
	}
}

func TestDeckFilter(t *testing.T) {
	d := New()

	aces := d.Filter(func(c Card) bool {
		return c.Rank() == Ace
	})

	if got, want := aces.Len(), 4; got != want {
		t.Errorf("Filter(Rank==Ace).Len() = %d, want %d", got, want)
	}

	for _, card := range aces.Cards() {
		if got, want := card.Rank(), Ace; got != want {
			t.Errorf("Filter(Rank==Ace) contains card with rank = %v, want %v (all cards should be Aces)", got, want)
		}
	}

	hearts := d.Filter(func(c Card) bool {
		return c.Suit() == Hearts
	})

	if got, want := hearts.Len(), 13; got != want {
		t.Errorf("Filter(Suit==Hearts).Len() = %d, want %d", got, want)
	}

	if got, want := d.Len(), 52; got != want {
		t.Errorf("After Filter(), original deck.Len() = %d, want %d (original should be unchanged)", got, want)
	}
}

func TestSecureShuffle(t *testing.T) {
	d1 := New()
	d2 := New()

	original := d1.Cards()

	d1.SecureShuffle()
	shuffled := d1.Cards()

	if got, want := d1.Len(), 52; got != want {
		t.Errorf("After SecureShuffle(), deck.Len() = %d, want %d", got, want)
	}

	sameOrder := true
	for i := range original {
		if original[i] != shuffled[i] {
			sameOrder = false
			break
		}
	}

	if sameOrder {
		t.Error("After SecureShuffle(), card order unchanged (very unlikely - secure shuffle may not be working)")
	}

	d2Cards := d2.Cards()
	for i := range original {
		if got, want := d2Cards[i], original[i]; got != want {
			t.Errorf("After d1.SecureShuffle(), d2.Cards()[%d] = %v, want %v (d2 should be unaffected)", i, got, want)
			break
		}
	}
}

func TestCustomShuffler(t *testing.T) {
	d1 := New()
	d2 := New()

	shuffler1 := NewSeededShuffler(12345)
	shuffler2 := NewSeededShuffler(12345)

	d1.ShuffleWith(shuffler1)
	d2.ShuffleWith(shuffler2)

	cards1 := d1.Cards()
	cards2 := d2.Cards()

	for i := range cards1 {
		if got, want := cards2[i], cards1[i]; got != want {
			t.Errorf("After ShuffleWith(same seeded shuffler), decks differ at index %d: got %v, want %v (same seed should produce same order)", i, got, want)
		}
	}
}

func TestDeckMarshalBinary(t *testing.T) {
	d := New()
	d.Shuffle()

	data, err := d.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary() got error: %v, want nil", err)
	}

	expectedSize := 4 + 52 // 4 byte header + 52 cards
	if got, want := len(data), expectedSize; got != want {
		t.Errorf("MarshalBinary() returned %d bytes, want %d", got, want)
	}

	d2 := &Deck{}
	if err := d2.UnmarshalBinary(data); err != nil {
		t.Fatalf("UnmarshalBinary() got error: %v, want nil", err)
	}

	if got, want := d2.Len(), d.Len(); got != want {
		t.Errorf("After UnmarshalBinary(), deck.Len() = %d, want %d", got, want)
	}

	cards1 := d.Cards()
	cards2 := d2.Cards()

	for i := range cards1 {
		if got, want := cards2[i], cards1[i]; got != want {
			t.Errorf("After UnmarshalBinary(), card[%d] = %v, want %v", i, got, want)
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
				t.Errorf("UnmarshalBinary(%v) got nil error, want error", tt.data)
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
			if got, want := tt.deck.Size(), tt.want; got != want {
				t.Errorf("%s: Size() = %d bytes, want %d bytes", tt.name, got, want)
			}
		})
	}
}

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

func TestNewWithJokers(t *testing.T) {
	d := NewWithJokers()

	if d == nil {
		t.Fatal("NewWithJokers() returned nil, want non-nil deck")
	}

	if got, want := d.Len(), 54; got != want {
		t.Errorf("NewWithJokers().Len() = %d, want %d", got, want)
	}

	var regularCards int
	var redJokers int
	var blackJokers int

	allCards, _ := d.PeekN(54)
	for _, c := range allCards {
		rank := c.Rank()

		if rank == RedJoker {
			redJokers++
		} else if rank == BlackJoker {
			blackJokers++
		} else if rank >= Ace && rank <= King {
			regularCards++
		}
	}

	if got, want := regularCards, 52; got != want {
		t.Errorf("NewWithJokers() contains %d regular cards, want %d", got, want)
	}
	if got, want := redJokers, 1; got != want {
		t.Errorf("NewWithJokers() contains %d red jokers, want %d", got, want)
	}
	if got, want := blackJokers, 1; got != want {
		t.Errorf("NewWithJokers() contains %d black jokers, want %d", got, want)
	}
}

func TestNewMultipleWithJokers(t *testing.T) {
	tests := []struct {
		name      string
		count     int
		wantCards int
		wantErr   bool
	}{
		{"single deck", 1, 54, false},
		{"two decks", 2, 108, false},
		{"five decks", 5, 270, false},
		{"zero decks", 0, 0, true},
		{"negative decks", -1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewMultipleWithJokers(tt.count)
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewMultipleWithJokers(%d) got nil error, want error", tt.count)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewMultipleWithJokers(%d) got error: %v, want nil", tt.count, err)
			}

			if got, want := d.Len(), tt.wantCards; got != want {
				t.Errorf("NewMultipleWithJokers(%d).Len() = %d, want %d", tt.count, got, want)
			}
		})
	}
}

func TestAddJoker(t *testing.T) {
	d := New() // Start with 52 cards

	if got, want := d.Len(), 52; got != want {
		t.Fatalf("New().Len() = %d, want %d", got, want)
	}

	// Add red joker
	d.AddJoker(RedJoker)
	if got, want := d.Len(), 53; got != want {
		t.Errorf("After AddJoker(RedJoker), deck.Len() = %d, want %d", got, want)
	}

	d.AddJoker(BlackJoker)
	if got, want := d.Len(), 54; got != want {
		t.Errorf("After AddJoker(BlackJoker), deck.Len() = %d, want %d", got, want)
	}

	cards, _ := d.PeekN(54)
	lastTwo := cards[52:]

	if got, want := lastTwo[0].Rank(), RedJoker; got != want {
		t.Errorf("After AddJoker, card[52].Rank() = %v, want %v", got, want)
	}
	if got, want := lastTwo[1].Rank(), BlackJoker; got != want {
		t.Errorf("After AddJoker, card[53].Rank() = %v, want %v", got, want)
	}
}

func TestIsJoker(t *testing.T) {
	tests := []struct {
		name    string
		card    Card
		isJoker bool
	}{
		{"Ace of Spades", NewCard(Ace, Spades), false},
		{"King of Hearts", NewCard(King, Hearts), false},
		{"Queen of Clubs", NewCard(Queen, Clubs), false},
		{"Red Joker", NewRedJoker(), true},
		{"Black Joker", NewBlackJoker(), true},
		{"Red Joker via NewCard", NewCard(RedJoker, Hearts), true},
		{"Black Joker via NewCard", NewCard(BlackJoker, Spades), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, want := tt.card.IsJoker(), tt.isJoker; got != want {
				t.Errorf("%s: IsJoker() = %v, want %v", tt.name, got, want)
			}
		})
	}
}

func TestJokerString(t *testing.T) {
	redJoker := NewRedJoker()
	blackJoker := NewBlackJoker()

	if got, want := redJoker.String(), "Joker (Red)"; got != want {
		t.Errorf("NewRedJoker().String() = %q, want %q", got, want)
	}
	if got, want := blackJoker.String(), "Joker (Black)"; got != want {
		t.Errorf("NewBlackJoker().String() = %q, want %q", got, want)
	}

	if got, want := redJoker.ShortString(), "JKR"; got != want {
		t.Errorf("NewRedJoker().ShortString() = %q, want %q", got, want)
	}
	if got, want := blackJoker.ShortString(), "JKB"; got != want {
		t.Errorf("NewBlackJoker().ShortString() = %q, want %q", got, want)
	}
}

func TestSortWithJokers(t *testing.T) {
	d := NewWithJokers()
	d.ShuffleWithSeed(42) // Shuffle first
	d.Sort()

	cards, _ := d.PeekN(54)

	if got, want := cards[52].Rank(), RedJoker; got != want {
		t.Errorf("After Sort(), card[52].Rank() = %v, want %v (Red Joker should be second to last)", got, want)
	}
	if got, want := cards[53].Rank(), BlackJoker; got != want {
		t.Errorf("After Sort(), card[53].Rank() = %v, want %v (Black Joker should be last)", got, want)
	}

	kingOfClubs := NewCard(King, Clubs)
	if got, want := cards[51], kingOfClubs; got != want {
		t.Errorf("After Sort(), card[51] = %s, want %s (King of Clubs should be before jokers)", got.String(), want.String())
	}
}

func TestShuffleWithJokers(t *testing.T) {
	d1 := NewWithJokers()
	d2 := NewWithJokers()

	cards1Before, _ := d1.PeekN(54)
	cards2Before, _ := d2.PeekN(54)

	for i := 0; i < 54; i++ {
		if got, want := cards2Before[i], cards1Before[i]; got != want {
			t.Fatalf("Before shuffle, d2.card[%d] = %v, want %v (decks should start with same order)", i, got, want)
		}
	}

	d1.ShuffleWithSeed(42)
	d2.ShuffleWithSeed(99)

	cards1After, _ := d1.PeekN(54)
	cards2After, _ := d2.PeekN(54)

	var differences int
	for i := 0; i < 54; i++ {
		if cards1After[i] != cards2After[i] {
			differences++
		}
	}

	if got, want := differences < 10, false; got != want {
		t.Errorf("After ShuffleWithSeed(different seeds), found %d differences, want >= 10 (different seeds should produce different orders)", differences)
	}

	if got, want := d1.Len(), 54; got != want {
		t.Errorf("After ShuffleWithSeed, d1.Len() = %d, want %d", got, want)
	}
	if got, want := d2.Len(), 54; got != want {
		t.Errorf("After ShuffleWithSeed, d2.Len() = %d, want %d", got, want)
	}
}

func TestDrawWithJokers(t *testing.T) {
	d := NewWithJokers()

	drawnCards := make([]Card, 0, 54)
	for i := 0; i < 54; i++ {
		card, err := d.Draw()
		if err != nil {
			t.Fatalf("Draw() on card %d got error: %v, want nil", i, err)
		}
		drawnCards = append(drawnCards, card)
	}

	if got, want := d.Len(), 0; got != want {
		t.Errorf("After drawing all cards, deck.Len() = %d, want %d (deck should be empty)", got, want)
	}

	var jc int
	for _, card := range drawnCards {
		if card.IsJoker() {
			jc++
		}
	}

	if got, want := jc, 2; got != want {
		t.Errorf("Drew %d jokers from NewWithJokers(), want %d", got, want)
	}
}

func TestFilterJokers(t *testing.T) {
	d := NewWithJokers()

	jokers := d.Filter(func(c Card) bool {
		return c.IsJoker()
	})

	if got, want := jokers.Len(), 2; got != want {
		t.Errorf("Filter(IsJoker).Len() = %d, want %d", got, want)
	}

	regularCards := d.Filter(func(c Card) bool {
		return !c.IsJoker()
	})

	if got, want := regularCards.Len(), 52; got != want {
		t.Errorf("Filter(!IsJoker).Len() = %d, want %d", got, want)
	}
}

func TestMarshalJokers(t *testing.T) {
	redJoker := NewRedJoker()
	blackJoker := NewBlackJoker()

	// Red Joker: Hearts (0x01) suit, Rank 14 (0x0E) = 0x4E
	redJokerByte := byte(redJoker)
	if got, want := redJokerByte, byte(0x4E); got != want {
		t.Errorf("byte(NewRedJoker()) = 0x%02X, want 0x%02X", got, want)
	}

	// Black Joker: Spades (0x00) suit, Rank 15 (0x0F) = 0x0F
	blackJokerByte := byte(blackJoker)
	if got, want := blackJokerByte, byte(0x0F); got != want {
		t.Errorf("byte(NewBlackJoker()) = 0x%02X, want 0x%02X", got, want)
	}

	d := NewWithJokers()
	data, err := d.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary() got error: %v, want nil", err)
	}

	// Check that we have 4 bytes (length) + 54 bytes (cards) = 58 bytes
	if got, want := len(data), 58; got != want {
		t.Fatalf("MarshalBinary() returned %d bytes, want %d (4 length + 54 cards)", got, want)
	}

	// In sorted order, jokers come last
	d.Sort()
	sortedData, _ := d.MarshalBinary()

	// Skip the 4-byte length prefix
	if got, want := sortedData[56], byte(0x4E); got != want {
		t.Errorf("After Sort(), MarshalBinary()[56] = 0x%02X, want 0x%02X (red joker)", got, want)
	}
	if got, want := sortedData[57], byte(0x0F); got != want {
		t.Errorf("After Sort(), MarshalBinary()[57] = 0x%02X, want 0x%02X (black joker)", got, want)
	}
}

func BenchmarkNewWithJokers(b *testing.B) {
	for b.Loop() {
		_ = NewWithJokers()
	}
}

func BenchmarkSortWithJokers(b *testing.B) {
	d := NewWithJokers()
	d.Shuffle()
	for b.Loop() {
		d.Sort()
	}
}

func TestDealValidation(t *testing.T) {
	tests := []struct {
		name           string
		numPlayers     int
		cardsPerPlayer int
		deckSize       int
		wantErr        string
	}{
		{
			name:           "zero players",
			numPlayers:     0,
			cardsPerPlayer: 5,
			deckSize:       52,
			wantErr:        "number of players must be at least 1",
		},
		{
			name:           "negative players",
			numPlayers:     -1,
			cardsPerPlayer: 5,
			deckSize:       52,
			wantErr:        "number of players must be at least 1",
		},
		{
			name:           "zero cards per player",
			numPlayers:     4,
			cardsPerPlayer: 0,
			deckSize:       52,
			wantErr:        "cards per player must be at least 1",
		},
		{
			name:           "negative cards per player",
			numPlayers:     4,
			cardsPerPlayer: -1,
			deckSize:       52,
			wantErr:        "cards per player must be at least 1",
		},
		{
			name:           "too many cards per player",
			numPlayers:     1,
			cardsPerPlayer: 53,
			deckSize:       52,
			wantErr:        "cards per player exceeds maximum of 52",
		},
		{
			name:           "insufficient cards in deck",
			numPlayers:     4,
			cardsPerPlayer: 5,
			deckSize:       19,
			wantErr:        "insufficient cards: need 20, have 19",
		},
		{
			name:           "empty deck",
			numPlayers:     1,
			cardsPerPlayer: 1,
			deckSize:       0,
			wantErr:        "insufficient cards: need 1, have 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New()
			for d.Len() > tt.deckSize {
				_, _ = d.Draw()
			}
			for d.Len() < tt.deckSize {
				d.Add(NewCard(Ace, Spades))
			}

			originalLen := d.Len()
			hands, err := d.Deal(tt.numPlayers, tt.cardsPerPlayer)
			if err == nil {
				t.Fatalf("Deal(%d, %d) got nil error, want %q", tt.numPlayers, tt.cardsPerPlayer, tt.wantErr)
			}

			if got, want := err.Error(), tt.wantErr; got != want {
				t.Errorf("Deal(%d, %d) error = %q, want %q", tt.numPlayers, tt.cardsPerPlayer, got, want)
			}

			if hands != nil {
				t.Errorf("Deal(%d, %d) returned hands = %v, want nil when error occurs", tt.numPlayers, tt.cardsPerPlayer, hands)
			}

			if got, want := d.Len(), originalLen; got != want {
				t.Errorf("After Deal(%d, %d) error, deck.Len() = %d, want %d (deck should be unchanged)", tt.numPlayers, tt.cardsPerPlayer, got, want)
			}
		})
	}
}

func TestDealHappyPath(t *testing.T) {
	tests := []struct {
		name           string
		numPlayers     int
		cardsPerPlayer int
		deckSize       int
		wantHands      int
		wantCardsEach  int
		wantRemaining  int
	}{
		{
			name:           "poker: 4 players, 5 cards each",
			numPlayers:     4,
			cardsPerPlayer: 5,
			deckSize:       52,
			wantHands:      4,
			wantCardsEach:  5,
			wantRemaining:  32,
		},
		{
			name:           "blackjack: 5 players, 2 cards each",
			numPlayers:     5,
			cardsPerPlayer: 2,
			deckSize:       52,
			wantHands:      5,
			wantCardsEach:  2,
			wantRemaining:  42,
		},
		{
			name:           "single player: 7 cards",
			numPlayers:     1,
			cardsPerPlayer: 7,
			deckSize:       52,
			wantHands:      1,
			wantCardsEach:  7,
			wantRemaining:  45,
		},
		{
			name:           "exact deck size: 4 players, 13 cards each",
			numPlayers:     4,
			cardsPerPlayer: 13,
			deckSize:       52,
			wantHands:      4,
			wantCardsEach:  13,
			wantRemaining:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New()
			hands, err := d.Deal(tt.numPlayers, tt.cardsPerPlayer)

			if err != nil {
				t.Fatalf("Deal(%d, %d) got error: %v, want nil", tt.numPlayers, tt.cardsPerPlayer, err)
			}

			if got, want := len(hands), tt.wantHands; got != want {
				t.Errorf("Deal(%d, %d) = %d hands, want %d", tt.numPlayers, tt.cardsPerPlayer, got, want)
			}

			for i, hand := range hands {
				if got, want := len(hand), tt.wantCardsEach; got != want {
					t.Errorf("Deal(%d, %d)[%d] = %d cards, want %d", tt.numPlayers, tt.cardsPerPlayer, i, got, want)
				}
			}

			if got, want := d.Len(), tt.wantRemaining; got != want {
				t.Errorf("After Deal(%d, %d), deck.Len() = %d, want %d", tt.numPlayers, tt.cardsPerPlayer, got, want)
			}
		})
	}
}

func TestDealSequentialDistribution(t *testing.T) {
	d := New()
	originalCards := d.Cards()

	hands, err := d.Deal(3, 2)
	if err != nil {
		t.Fatalf("Deal(3, 2) got error: %v, want nil", err)
	}

	// Check player 1 got cards[0:2]
	if got, want := hands[0][0], originalCards[0]; got != want {
		t.Errorf("Deal(3, 2)[0][0] = %v, want %v", got, want)
	}
	if got, want := hands[0][1], originalCards[1]; got != want {
		t.Errorf("Deal(3, 2)[0][1] = %v, want %v", got, want)
	}

	// Check player 2 got cards[2:4]
	if got, want := hands[1][0], originalCards[2]; got != want {
		t.Errorf("Deal(3, 2)[1][0] = %v, want %v", got, want)
	}
	if got, want := hands[1][1], originalCards[3]; got != want {
		t.Errorf("Deal(3, 2)[1][1] = %v, want %v", got, want)
	}

	// Check player 3 got cards[4:6]
	if got, want := hands[2][0], originalCards[4]; got != want {
		t.Errorf("Deal(3, 2)[2][0] = %v, want %v", got, want)
	}
	if got, want := hands[2][1], originalCards[5]; got != want {
		t.Errorf("Deal(3, 2)[2][1] = %v, want %v", got, want)
	}
}

func TestDealAtomicOperation(t *testing.T) {
	tests := []struct {
		name           string
		numPlayers     int
		cardsPerPlayer int
		setupDeck      func() *Deck
	}{
		{
			name:           "insufficient cards",
			numPlayers:     4,
			cardsPerPlayer: 5,
			setupDeck: func() *Deck {
				d := New()
				for i := 0; i < 33; i++ {
					_, _ = d.Draw()
				}
				return d
			},
		},
		{
			name:           "zero players",
			numPlayers:     0,
			cardsPerPlayer: 5,
			setupDeck:      func() *Deck { return New() },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.setupDeck()
			originalCards := d.Cards()
			originalLen := d.Len()

			_, err := d.Deal(tt.numPlayers, tt.cardsPerPlayer)
			if err == nil {
				t.Fatalf("Deal(%d, %d) got nil error, want error", tt.numPlayers, tt.cardsPerPlayer)
			}

			if got, want := d.Len(), originalLen; got != want {
				t.Errorf("After Deal error, deck.Len() = %d, want %d (deck should be unchanged)", got, want)
			}

			currentCards := d.Cards()
			if got, want := len(currentCards), len(originalCards); got != want {
				t.Errorf("After Deal error, deck has %d cards, want %d", got, want)
			}

			for i := range originalCards {
				if got, want := currentCards[i], originalCards[i]; got != want {
					t.Errorf("After Deal error, card[%d] = %v, want %v (unchanged)", i, got, want)
				}
			}
		})
	}
}

func TestDealIndependentHands(t *testing.T) {
	d := New()
	hands, err := d.Deal(3, 5)
	if err != nil {
		t.Fatalf("Deal(3, 5) unexpected error: %v", err)
	}

	originalCard := hands[0][0]
	hands[0][0] = NewCard(Ace, Spades)

	// Verify other hands are not affected
	for i := 1; i < len(hands); i++ {
		for j := range hands[i] {
			if got := hands[i][j]; got == NewCard(Ace, Spades) && got != originalCard {
				t.Errorf("After modifying hands[0], hands[%d][%d] = %v (affected by change)", i, j, got)
			}
		}
	}

	remainingCards := d.Cards()
	for i, card := range remainingCards {
		if got := card; got == NewCard(Ace, Spades) && i < 15 {
			t.Errorf("After modifying hands[0], deck card[%d] = %v (deck was affected)", i, got)
		}
	}
}

func BenchmarkDeal(b *testing.B) {
	tests := []struct {
		name           string
		numPlayers     int
		cardsPerPlayer int
	}{
		{"poker_4p_5c", 4, 5},
		{"blackjack_5p_2c", 5, 2},
		{"bridge_4p_13c", 4, 13},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for b.Loop() {
				d := New()
				_, err := d.Deal(tt.numPlayers, tt.cardsPerPlayer)
				if err != nil {
					b.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestMustDraw_Success(t *testing.T) {
	d := New()
	card := d.MustDraw()

	if got, want := card.Rank(), Ace; got != want {
		t.Errorf("MustDraw().Rank() = %v, want %v", got, want)
	}
	if got, want := card.Suit(), Spades; got != want {
		t.Errorf("MustDraw().Suit() = %v, want %v", got, want)
	}
	if got, want := d.Len(), 51; got != want {
		t.Errorf("After MustDraw(), deck.Len() = %d, want %d", got, want)
	}
}

func TestMustDraw_EmptyDeckPanic(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("MustDraw() on empty deck did not panic, want panic")
		}
		expectedMsg := "cannot draw from empty deck"
		if got, want := r.(string), expectedMsg; got != want {
			t.Errorf("MustDraw() panic = %q, want %q", got, want)
		}
	}()

	d := &Deck{}
	d.MustDraw()
}

func TestMustDrawN_Success(t *testing.T) {
	d := New()
	cards := d.MustDrawN(5)

	if got, want := len(cards), 5; got != want {
		t.Errorf("MustDrawN(5) = %d cards, want %d", got, want)
	}
	if got, want := d.Len(), 47; got != want {
		t.Errorf("After MustDrawN(5), deck.Len() = %d, want %d", got, want)
	}
}

func TestMustDrawN_Panics(t *testing.T) {
	tests := []struct {
		name        string
		deckSize    int
		drawCount   int
		expectPanic string
	}{
		{"negative count", 52, -1, "cannot draw negative number of cards: -1"},
		{"insufficient cards", 52, 53, "not enough cards in deck: have 52, need 53"},
		{"empty deck", 0, 1, "not enough cards in deck: have 0, need 1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if r == nil {
					t.Fatalf("MustDrawN(%d) did not panic, want panic", tt.drawCount)
				}
				if got, want := r.(string), tt.expectPanic; got != want {
					t.Errorf("MustDrawN(%d) panic = %q, want %q", tt.drawCount, got, want)
				}
			}()

			d := New()
			if tt.deckSize == 0 {
				d = &Deck{}
			}
			d.MustDrawN(tt.drawCount)
		})
	}
}

func TestMustDeal_Success(t *testing.T) {
	tests := []struct {
		name           string
		numPlayers     int
		cardsPerPlayer int
		wantRemaining  int
	}{
		{"bridge", 4, 13, 0},
		{"poker", 4, 5, 32},
		{"blackjack", 5, 2, 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New()
			hands := d.MustDeal(tt.numPlayers, tt.cardsPerPlayer)

			if got, want := len(hands), tt.numPlayers; got != want {
				t.Errorf("MustDeal(%d, %d) = %d hands, want %d", tt.numPlayers, tt.cardsPerPlayer, got, want)
			}

			for i, hand := range hands {
				if got, want := len(hand), tt.cardsPerPlayer; got != want {
					t.Errorf("MustDeal(%d, %d)[%d] = %d cards, want %d", tt.numPlayers, tt.cardsPerPlayer, i, got, want)
				}
			}

			if got, want := d.Len(), tt.wantRemaining; got != want {
				t.Errorf("After MustDeal(%d, %d), deck.Len() = %d, want %d", tt.numPlayers, tt.cardsPerPlayer, got, want)
			}
		})
	}
}

func TestMustDeal_Panics(t *testing.T) {
	tests := []struct {
		name           string
		numPlayers     int
		cardsPerPlayer int
		deckSize       int
		expectPanic    string
	}{
		{"zero players", 0, 5, 52, "number of players must be at least 1"},
		{"negative players", -1, 5, 52, "number of players must be at least 1"},
		{"zero cards", 4, 0, 52, "cards per player must be at least 1"},
		{"negative cards", 4, -1, 52, "cards per player must be at least 1"},
		{"too many cards per player", 1, 53, 52, "cards per player exceeds maximum of 52"},
		{"insufficient cards", 4, 5, 15, "insufficient cards: need 20, have 15"},
		{"insufficient cards (too many players)", 53, 1, 52, "insufficient cards: need 53, have 52"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if r == nil {
					t.Errorf("MustDeal(%d, %d) did not panic, want panic", tt.numPlayers, tt.cardsPerPlayer)
					return
				}
				if got, want := r.(string), tt.expectPanic; got != want {
					t.Errorf("MustDeal(%d, %d) panic message = %q, want %q", tt.numPlayers, tt.cardsPerPlayer, got, want)
				}
			}()

			d := New()
			for d.Len() > tt.deckSize {
				_, _ = d.Draw()
			}
			d.MustDeal(tt.numPlayers, tt.cardsPerPlayer)
		})
	}
}

func TestDealHands_HappyPath(t *testing.T) {
	tests := []struct {
		name          string
		handSizes     []int
		wantRemaining int
	}{
		{
			name:          "casino style",
			handSizes:     []int{2, 2, 2, 1},
			wantRemaining: 45,
		},
		{
			name:          "progressive",
			handSizes:     []int{2, 3, 5},
			wantRemaining: 42,
		},
		{
			name:          "single hand",
			handSizes:     []int{7},
			wantRemaining: 45,
		},
		{
			name:          "uniform",
			handSizes:     []int{5, 5, 5, 5},
			wantRemaining: 32,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New()
			hands, err := d.DealHands(tt.handSizes)

			if err != nil {
				t.Fatalf("DealHands(%v) unexpected error: %v", tt.handSizes, err)
			}

			if got, want := len(hands), len(tt.handSizes); got != want {
				t.Errorf("DealHands(%v) = %d hands, want %d", tt.handSizes, got, want)
			}

			for i, hand := range hands {
				if got, want := len(hand), tt.handSizes[i]; got != want {
					t.Errorf("DealHands(%v)[%d] = %d cards, want %d", tt.handSizes, i, got, want)
				}
			}

			if got, want := d.Len(), tt.wantRemaining; got != want {
				t.Errorf("After DealHands(%v), deck.Len() = %d, want %d", tt.handSizes, got, want)
			}
		})
	}
}

func TestDealHands_Validation(t *testing.T) {
	tests := []struct {
		name        string
		handSizes   []int
		deckSize    int
		expectedErr string
	}{
		{
			name:        "empty slice",
			handSizes:   []int{},
			deckSize:    52,
			expectedErr: "handSizes must contain at least one hand",
		},
		{
			name:        "zero value",
			handSizes:   []int{5, 0, 3},
			deckSize:    52,
			expectedErr: "hand size must be positive: got 0 at index 1",
		},
		{
			name:        "negative value",
			handSizes:   []int{2, -3, 4},
			deckSize:    52,
			expectedErr: "hand size must be positive: got -3 at index 1",
		},
		{
			name:        "hand too large",
			handSizes:   []int{53},
			deckSize:    52,
			expectedErr: "hand size (53) at index 0 exceeds maximum of 52",
		},
		{
			name:        "insufficient cards",
			handSizes:   []int{30, 30},
			deckSize:    52,
			expectedErr: "insufficient cards: need 60, have 52",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New()

			for d.Len() > tt.deckSize {
				_, _ = d.Draw()
			}

			initialSize := d.Len()

			hands, err := d.DealHands(tt.handSizes)
			if err == nil {
				t.Fatalf("DealHands(%v) got nil error, want error", tt.handSizes)
			}

			if got, want := err.Error(), tt.expectedErr; got != want {
				t.Errorf("DealHands(%v) error = %q, want %q", tt.handSizes, got, want)
			}

			if hands != nil {
				t.Errorf("DealHands(%v) with error returned non-nil hands, want nil", tt.handSizes)
			}

			// Verify deck unchanged (atomic operation)
			if got, want := d.Len(), initialSize; got != want {
				t.Errorf("After DealHands(%v) error, deck.Len() = %d, want %d (deck should be unchanged)", tt.handSizes, got, want)
			}
		})
	}
}

func TestDealHands_BoundaryConditions(t *testing.T) {
	tests := []struct {
		name          string
		handSizes     []int
		wantRemaining int
	}{
		{
			name:          "exact deck size",
			handSizes:     []int{26, 26},
			wantRemaining: 0,
		},
		{
			name:          "maximum hand size",
			handSizes:     []int{52},
			wantRemaining: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New()
			hands, err := d.DealHands(tt.handSizes)

			if err != nil {
				t.Fatalf("DealHands(%v) unexpected error: %v", tt.handSizes, err)
			}

			if got, want := len(hands), len(tt.handSizes); got != want {
				t.Errorf("DealHands(%v) = %d hands, want %d", tt.handSizes, got, want)
			}

			for i, hand := range hands {
				if got, want := len(hand), tt.handSizes[i]; got != want {
					t.Errorf("DealHands(%v)[%d] = %d cards, want %d", tt.handSizes, i, got, want)
				}
			}

			if got, want := d.Len(), tt.wantRemaining; got != want {
				t.Errorf("After DealHands(%v), deck.Len() = %d, want %d", tt.handSizes, got, want)
			}
		})
	}
}

func TestDealHands_IndependentSlices(t *testing.T) {
	d := New()
	hands, err := d.DealHands([]int{5, 5})
	if err != nil {
		t.Fatalf("DealHands([5, 5]) unexpected error: %v", err)
	}

	originalHand2 := make([]Card, len(hands[1]))
	copy(originalHand2, hands[1])

	originalDeckSize := d.Len()
	if len(hands[0]) > 0 {
		hands[0][0] = NewCard(Ace, Hearts)
	}

	for i, card := range hands[1] {
		if got, want := card, originalHand2[i]; got != want {
			t.Errorf("After modifying hands[0], hands[1][%d] = %v, want %v (hands should be independent)", i, got, want)
		}
	}

	if got, want := d.Len(), originalDeckSize; got != want {
		t.Errorf("After modifying hands[0], deck.Len() = %d, want %d (deck should be unaffected)", got, want)
	}
}

func TestDealHands_EquivalentToDeal(t *testing.T) {
	d1 := New()
	d2 := New()

	hands1, err1 := d1.Deal(3, 5)
	if err1 != nil {
		t.Fatalf("Deal(3, 5) unexpected error: %v", err1)
	}

	hands2, err2 := d2.DealHands([]int{5, 5, 5})
	if err2 != nil {
		t.Fatalf("DealHands([5, 5, 5]) unexpected error: %v", err2)
	}

	if got, want := len(hands1), len(hands2); got != want {
		t.Fatalf("Deal(3, 5) returned %d hands, DealHands([5, 5, 5]) returned %d hands, want same", got, want)
	}

	for i := range hands1 {
		if got, want := len(hands1[i]), len(hands2[i]); got != want {
			t.Errorf("hand %d: Deal = %d cards, DealHands = %d cards, want same", i, got, want)
			continue
		}

		for j := range hands1[i] {
			if got, want := hands1[i][j], hands2[i][j]; got != want {
				t.Errorf("hand %d, card %d: Deal = %v, DealHands = %v, want same", i, j, got, want)
			}
		}
	}

	if got, want := d1.Len(), d2.Len(); got != want {
		t.Errorf("After dealing, Deal deck = %d cards, DealHands deck = %d cards, want same", got, want)
	}
}

func TestMustDealHands_Success(t *testing.T) {
	tests := []struct {
		name      string
		handSizes []int
		deckSize  int
	}{
		{"casino style", []int{2, 2, 2, 2}, 52},
		{"progressive", []int{3, 5, 7, 9}, 52},
		{"uniform", []int{5, 5, 5, 5}, 52},
		{"single hand", []int{10}, 52},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New()
			hands := d.MustDealHands(tt.handSizes)

			if got, want := len(hands), len(tt.handSizes); got != want {
				t.Errorf("MustDealHands(%v) = %d hands, want %d", tt.handSizes, got, want)
			}

			for i, hand := range hands {
				if got, want := len(hand), tt.handSizes[i]; got != want {
					t.Errorf("MustDealHands(%v)[%d] = %d cards, want %d", tt.handSizes, i, got, want)
				}
			}

			totalDealt := 0
			for _, size := range tt.handSizes {
				totalDealt += size
			}
			expectedRemaining := tt.deckSize - totalDealt
			if got, want := d.Len(), expectedRemaining; got != want {
				t.Errorf("After MustDealHands(%v), deck.Len() = %d, want %d", tt.handSizes, got, want)
			}
		})
	}
}

func TestMustDealHands_Panics(t *testing.T) {
	tests := []struct {
		name        string
		handSizes   []int
		deckSize    int
		expectPanic string
	}{
		{"empty slice", []int{}, 52, "handSizes must contain at least one hand"},
		{"zero value", []int{2, 0, 3}, 52, "hand size must be positive: got 0 at index 1"},
		{"negative value", []int{2, -1, 3}, 52, "hand size must be positive: got -1 at index 1"},
		{"hand too large", []int{53}, 52, "hand size (53) at index 0 exceeds maximum of 52"},
		{"insufficient cards", []int{5, 5, 5, 5}, 15, "insufficient cards: need 20, have 15"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if r == nil {
					t.Fatalf("MustDealHands(%v) did not panic, want panic", tt.handSizes)
				}
				if got, want := r.(string), tt.expectPanic; got != want {
					t.Errorf("MustDealHands(%v) panic = %q, want %q", tt.handSizes, got, want)
				}
			}()

			d := New()
			for d.Len() > tt.deckSize {
				_, _ = d.Draw()
			}
			d.MustDealHands(tt.handSizes)
		})
	}
}
