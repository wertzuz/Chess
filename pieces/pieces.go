package pieces

import (
	"fmt"
	. "gfxw"
)

type Piece interface {
	Calc_Moves(pieces_a [64]Piece)
	Piece_Is_White() bool
	Give_Legal_Moves() [][2]uint16
	Give_Pos() [2]uint16
	Move_To(new_position [2]uint16)
	Is_White_Piece() bool
	Append_Legal_Moves(new_legal_move [2]uint16)
	Clear_Legal_Moves() //kann wahrscheinlich weg
}

func (c *ChessObject) Give_Legal_Moves() [][2]uint16 {
	return c.Legal_Moves
}

func (c *ChessObject) Clear_Legal_Moves() {
	c.Legal_Moves = nil
}

func (c *ChessObject) Append_Legal_Moves(new_legal_move [2]uint16) {
	c.Legal_Moves = append(c.Legal_Moves, new_legal_move)
}

func (c *ChessObject) Is_White_Piece() bool {
	return c.White
}

func Draw(piece Piece, w_x, w_y, a uint16) {
	UpdateAus()
	Archivieren()
	LadeBild(0, 0, "C:\\Users\\liamw\\Documents\\_Privat\\_Go\\Chess\\Pieces2.bmp")

	switch piece.(type) {
	case *Pawn:
		if piece.Is_White_Piece() {
			Clipboard_kopieren(0, a, a, a)
		} else {
			Clipboard_kopieren(0, 0, a, a)
		}
	case *Knight:
		if piece.Is_White_Piece() {
			Clipboard_kopieren(a, a, a, a)
		} else {
			Clipboard_kopieren(a, 0, a, a)
		}
	case *Bishop:
		if piece.Is_White_Piece() {
			Clipboard_kopieren(2*a, a, a, a)
		} else {
			Clipboard_kopieren(2*a, 0, a, a)
		}
	case *Rook:
		if piece.Is_White_Piece() {
			Clipboard_kopieren(3*a, a, a, a)
		} else {
			Clipboard_kopieren(3*a, 0, a, a)
		}
	case *Queen:
		if piece.Is_White_Piece() {
			Clipboard_kopieren(4*a, a, a, a)
		} else {
			Clipboard_kopieren(4*a, 0, a, a)
		}
	case *King:
		if piece.Is_White_Piece() {
			Clipboard_kopieren(5*a, a, a, a)
		} else {
			Clipboard_kopieren(5*a, 0, a, a)
		}
	default:
		fmt.Println("Unknown Piece type")
	}

	Restaurieren(0, 0, w_x, w_y)
	UpdateAn()
	Clipboard_einfuegenMitColorKey(piece.Give_Pos()[0]*a, piece.Give_Pos()[1]*a, 5, 5, 5)

}

func (c *ChessObject) Move_To(new_position [2]uint16) {
	c.Position = new_position
}

func (c *ChessObject) Give_Pos() [2]uint16 {
	return c.Position
}

func (c *ChessObject) Piece_Is_White() bool {
	return c.White
}

type Positioning struct { //datentyp Positioning
	Position [2]uint16
}

type ChessObject struct { //datentyp ChessObject erbt vom datentyp Positioning
	Positioning
	White       bool
	Legal_Moves [][2]uint16
}

type Pawn struct { //alle Schachobjekte erben wiederum vom datentyp ChessObject
	Has_moved bool
	ChessObject
}

func NewPawn(x, y uint16, is_white bool) *Pawn {
	return &Pawn{
		ChessObject: ChessObject{Positioning: Positioning{Position: [2]uint16{x, y}}, White: is_white},
	}
}

type Knight struct {
	ChessObject
}

func NewKnight(x, y uint16, is_white bool) *Knight {
	return &Knight{
		ChessObject: ChessObject{Positioning: Positioning{Position: [2]uint16{x, y}}, White: is_white},
	}
}

type Bishop struct {
	ChessObject
}

func NewBishop(x, y uint16, is_white bool) *Bishop {
	return &Bishop{
		ChessObject: ChessObject{Positioning: Positioning{Position: [2]uint16{x, y}}, White: is_white},
	}
}

type Rook struct {
	Has_moved bool
	ChessObject
}

func NewRook(x, y uint16, is_white bool) *Rook {
	return &Rook{
		ChessObject: ChessObject{Positioning: Positioning{Position: [2]uint16{x, y}}, White: is_white},
	}
}

type Queen struct {
	ChessObject
}

func NewQueen(x, y uint16, is_white bool) *Queen {
	return &Queen{
		ChessObject: ChessObject{Positioning: Positioning{Position: [2]uint16{x, y}}, White: is_white},
	}
}

type King struct {
	Has_moved bool
	ChessObject
}

func NewKing(x, y uint16, is_white bool) *King {
	return &King{
		ChessObject: ChessObject{Positioning: Positioning{Position: [2]uint16{x, y}}, White: is_white},
	}
}

func (p *Pawn) Calc_Moves(pieces_a [64]Piece) { //en passant --> nur unmittelbar nach dem bauern zweier zug, es darf kein anderer zug dazwischen liegen
	fmt.Printf("Moves of Pawn")

	var blocking_piece_1 bool
	var blocking_piece_2 bool
	new_legal_move_1 := [2]uint16{10, 10}
	new_legal_move_2 := [2]uint16{10, 10}
	new_legal_move_3 := [2]uint16{10, 10}
	new_legal_move_4 := [2]uint16{10, 10}

	if p.Is_White_Piece() && p.Position[1] != 0 {
		new_legal_move_1 = [2]uint16{p.Position[0], p.Position[1] - 1}
		if p.Position[1] > 1 && !p.Has_moved {
			new_legal_move_2 = [2]uint16{p.Position[0], p.Position[1] - 2}
		}
		new_legal_move_3 = [2]uint16{p.Position[0] + 1, p.Position[1] - 1}
		new_legal_move_4 = [2]uint16{p.Position[0] - 1, p.Position[1] - 1}
	} else if p.Position[1] != 7 {
		new_legal_move_1 = [2]uint16{p.Position[0], p.Position[1] + 1}
		if p.Position[1] < 6 && !p.Has_moved {
			new_legal_move_2 = [2]uint16{p.Position[0], p.Position[1] + 2}
		}
		new_legal_move_3 = [2]uint16{p.Position[0] + 1, p.Position[1] + 1}
		new_legal_move_4 = [2]uint16{p.Position[0] - 1, p.Position[1] + 1}
	}
	for i := 0; i < len(pieces_a) && (!blocking_piece_1 || !blocking_piece_2); i++ {
		if pieces_a[i] != nil {
			if pieces_a[i].Give_Pos() == new_legal_move_1 {
				blocking_piece_1 = true
			} else if pieces_a[i].Give_Pos() == new_legal_move_2 {
				blocking_piece_2 = true
			} else if pieces_a[i].Give_Pos() == new_legal_move_3 && pieces_a[i].Is_White_Piece() != p.Is_White_Piece() {
				p.Append_Legal_Moves(new_legal_move_3)
			} else if pieces_a[i].Give_Pos() == new_legal_move_4 && pieces_a[i].Is_White_Piece() != p.Is_White_Piece() {
				p.Append_Legal_Moves(new_legal_move_4)
			}
		}
	}

	if !blocking_piece_1 && new_legal_move_1 != [2]uint16{10, 10} { //es steht nichts im weg direkt davor
		p.Append_Legal_Moves(new_legal_move_1)
	}
	if !blocking_piece_2 && new_legal_move_2 != [2]uint16{10, 10} { //es steht nichts im weg direkt davor
		p.Append_Legal_Moves(new_legal_move_2)
	}

}

func (p *Knight) Calc_Moves(pieces_a [64]Piece) {
	fmt.Printf("Moves of Knight")
}

func (p *Rook) Calc_Moves(pieces_a [64]Piece) {
	p.Legal_Moves = nil

	for new_x := p.Position[0]; new_x < 7; {
		new_x++
		var current_pos [2]uint16 = [2]uint16{new_x, p.Position[1]}

		if check_if_piece_is_blocking(p, pieces_a, current_pos) {
			break
		}
	}

	for new_x := p.Position[0]; new_x != 0; {
		new_x--
		var current_pos [2]uint16 = [2]uint16{new_x, p.Position[1]}
		if check_if_piece_is_blocking(p, pieces_a, current_pos) {
			break
		}
		if new_x == 0 {
			break
		}

	}

	for new_y := p.Position[1]; new_y < 7; {
		new_y++
		var current_pos [2]uint16 = [2]uint16{p.Position[0], new_y}
		if check_if_piece_is_blocking(p, pieces_a, current_pos) {
			break
		}
	}

	for new_y := p.Position[1]; new_y != 0; {
		new_y--
		var current_pos [2]uint16 = [2]uint16{p.Position[0], new_y}
		if check_if_piece_is_blocking(p, pieces_a, current_pos) {
			break
		}
		if new_y == 0 {
			break
		}

	}
}

func (p *Bishop) Calc_Moves(pieces_a [64]Piece) {
	p.Legal_Moves = nil

	for new_x, new_y := p.Position[0], p.Position[1]; new_x < 7 && new_y < 7; {
		new_x++
		new_y++
		var current_pos [2]uint16 = [2]uint16{new_x, new_y}

		if check_if_piece_is_blocking(p, pieces_a, current_pos) {
			break
		}
	}

	for new_x, new_y := p.Position[0], p.Position[1]; new_x < 7 && new_y != 0; {
		new_x++
		new_y--
		var current_pos [2]uint16 = [2]uint16{new_x, new_y}

		if check_if_piece_is_blocking(p, pieces_a, current_pos) {
			break
		}
	}

	for new_x, new_y := p.Position[0], p.Position[1]; new_x != 0 && new_y < 7; {
		new_x--
		new_y++
		var current_pos [2]uint16 = [2]uint16{new_x, new_y}
		fmt.Println(current_pos)

		if check_if_piece_is_blocking(p, pieces_a, current_pos) {
			break
		}
	}

	for new_x, new_y := p.Position[0], p.Position[1]; new_x != 0 && new_y != 0; {
		new_x--
		new_y--
		var current_pos [2]uint16 = [2]uint16{new_x, new_y}
		fmt.Println(current_pos)

		if check_if_piece_is_blocking(p, pieces_a, current_pos) {
			break
		}
	}
}

func (p *Queen) Calc_Moves(pieces_a [64]Piece) {
	fmt.Printf("Moves of Queen")
}

func (p *King) Calc_Moves(pieces_a [64]Piece) {
	fmt.Printf("Moves of King")
}

func check_if_piece_is_blocking(p Piece, pieces_a [64]Piece, current_pos [2]uint16) bool {
	var blocking_piece Piece
	var var_break bool = false

	for i := 0; i < len(pieces_a) && blocking_piece == nil; i++ {
		if pieces_a[i] != nil {
			if pieces_a[i].Give_Pos() == current_pos {
				blocking_piece = pieces_a[i]
			}
		}
	}

	if blocking_piece == nil { //es steht nichts im weg
		p.Append_Legal_Moves(current_pos)
	} else if blocking_piece.Is_White_Piece() != p.Is_White_Piece() { //es steht etwas im weg, was aber geschlagen werden kann, daher wird danach jedoch gebreaked
		p.Append_Legal_Moves(current_pos)
		var_break = true
	} else if blocking_piece.Is_White_Piece() == p.Is_White_Piece() { //es steht etwas im weg, was aber nicht geschlagen werden kann, daher wird sofort gebreaked
		var_break = true
	} else {
		fmt.Println("fatal: Error in Calculating Moves Method")
	}
	return var_break
}
