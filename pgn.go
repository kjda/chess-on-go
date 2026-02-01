package chessongo

import (
	"fmt"
	"strings"
	"unicode"
)

// LoadPGN loads a PGN string into the board, playing all main-line moves and
// recording position history via Zobrist hashing. Variations and comments are
// ignored; only the main line is applied.
func (g *Game) LoadPGN(pgn string) error {
	fastPath := !strings.ContainsAny(pgn, "[{(")
	startFEN := STARTING_POSITION_FEN
	if !fastPath {
		if fen := extractFENFromPGN(pgn); fen != "" {
			startFEN = fen
		}
	}
	if err := g.LoadFen(startFEN); err != nil {
		return err
	}
	// Initial generation of legal moves so the loop can start matching immediately.
	// Subsequent generations are handled by g.MakeMove().
	g.GenerateLegalMoves()

	var tokens []string
	if fastPath {
		tokens = fastTokenizeMoves(pgn)
	} else {
		tokens = tokenizePGNMoves(pgn)
	}

	for _, tok := range tokens {
		tok = strings.TrimSpace(tok)
		if tok == "" {
			continue
		}
		// Skip results, move numbers, and NAGs
		if isPGNResult(tok) || strings.Contains(tok, "..") || strings.Contains(tok, ".") || strings.HasPrefix(tok, "$") {
			continue
		}

		tok = trimSANAnnotations(tok)
		if tok == "" {
			continue
		}

		target := getTargetSquare(tok)

		// b.GenerateLegalMoves() is already done by LoadFen (initially) and MakeMove (subsequently).
		matched := false
		for _, mv := range g.LegalMoves {
			if target != -1 && int(mv.To()) != target {
				continue
			}

			// Optimization: GetMoveSanWithoutSuffix avoids cloning the board (to check for check/mate)
			// which is very expensive. We strip annotations from the token anyway.
			san := trimSANAnnotations(g.GetMoveSanWithoutSuffix(mv))
			if san == tok {
				g.MakeMove(mv)
				matched = true
				break
			}
		}
		if !matched {
			return fmt.Errorf("pgn move not found: %s", tok)
		}
	}

	return nil
}

// LoadPGNGame is a helper that constructs a fresh board, loads the PGN, and
// returns the populated board.
func LoadPGNGame(pgn string) (*Game, error) {
	g := &Game{}
	if err := g.LoadPGN(pgn); err != nil {
		return nil, err
	}
	return g, nil
}

func extractFENFromPGN(pgn string) string {
	isSpace := func(b byte) bool { return b == ' ' || b == '\t' || b == '\r' }
	for i := 0; i < len(pgn); {
		lineStart := i
		// find line end
		lineEnd := strings.IndexByte(pgn[i:], '\n')
		if lineEnd == -1 {
			lineEnd = len(pgn)
		} else {
			lineEnd = i + lineEnd
		}

		// trim leading spaces manually
		for lineStart < lineEnd && isSpace(pgn[lineStart]) {
			lineStart++
		}

		// quick prefix check for [FEN
		if lineEnd-lineStart >= 5 && pgn[lineStart] == '[' && pgn[lineStart+1] == 'F' && pgn[lineStart+2] == 'E' && pgn[lineStart+3] == 'N' && pgn[lineStart+4] == ' ' {
			firstQuote := -1
			for j := lineStart + 5; j < lineEnd; j++ {
				if pgn[j] == '"' {
					firstQuote = j
					break
				}
			}
			if firstQuote != -1 {
				lastQuote := -1
				for j := lineEnd - 1; j > firstQuote; j-- {
					if pgn[j] == '"' {
						lastQuote = j
						break
					}
				}
				if lastQuote != -1 {
					fenStart := firstQuote + 1
					fenEnd := lastQuote
					for fenStart < fenEnd && isSpace(pgn[fenStart]) {
						fenStart++
					}
					for fenEnd > fenStart && isSpace(pgn[fenEnd-1]) {
						fenEnd--
					}
					return pgn[fenStart:fenEnd]
				}
			}
		}

		if lineEnd == len(pgn) {
			break
		}
		i = lineEnd + 1
	}
	return ""
}

func tokenizePGNMoves(pgn string) []string {
	var tokens []string
	var bld strings.Builder
	braceDepth, parenDepth := 0, 0
	inTag := false
	lineStart := true

	flush := func() {
		if bld.Len() > 0 {
			tokens = append(tokens, bld.String())
			bld.Reset()
		}
	}

	for i := 0; i < len(pgn); i++ {
		c := pgn[i]
		if lineStart {
			inTag = c == '['
		}
		if c == '\n' {
			lineStart = true
			inTag = false
			flush()
			continue
		}
		lineStart = false

		if inTag {
			continue
		}
		switch c {
		case '{':
			braceDepth++
			flush()
			continue
		case '}':
			if braceDepth > 0 {
				braceDepth--
			}
			flush()
			continue
		case '(':
			parenDepth++
			flush()
			continue
		case ')':
			if parenDepth > 0 {
				parenDepth--
			}
			flush()
			continue
		}

		if braceDepth > 0 || parenDepth > 0 {
			continue
		}

		if c == ' ' || c == '\t' || c == '\r' {
			flush()
			continue
		}

		bld.WriteByte(c)
	}
	flush()
	return tokens
}

// fastTokenizeMoves is a lightweight splitter for simple PGNs that only contain
// moves (e.g. "1. e4 e5 2. Nf3 Nc6") without tags, comments, or variations.
func fastTokenizeMoves(pgn string) []string {
	var tokens []string
	start := -1
	for i := 0; i <= len(pgn); i++ {
		var c byte
		if i < len(pgn) {
			c = pgn[i]
		}
		if i == len(pgn) || c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			if start >= 0 {
				tokens = append(tokens, pgn[start:i])
				start = -1
			}
			continue
		}
		if start == -1 {
			start = i
		}
	}
	return tokens
}

func isPGNResult(tok string) bool {
	switch tok {
	case "1-0", "0-1", "1/2-1/2", "*":
		return true
	default:
		return false
	}
}

func trimSANAnnotations(san string) string {
	// Drop trailing check/mate and annotation glyphs.
	san = strings.TrimRightFunc(san, func(r rune) bool {
		return strings.ContainsRune("+#!?", r)
	})
	// Remove leading move numbers if any slipped through (e.g., "12.Nf3").
	san = strings.TrimLeftFunc(san, func(r rune) bool {
		return unicode.IsDigit(r) || r == '.'
	})
	return san
}

func getTargetSquare(san string) int {
	for i := len(san) - 1; i >= 0; i-- {
		c := san[i]
		if c >= '1' && c <= '8' {
			if i > 0 {
				f := san[i-1]
				if f >= 'a' && f <= 'h' {
					col := int(f - 'a')
					row := int('8' - c)
					return row*8 + col
				}
			}
		}
	}
	return -1
}
