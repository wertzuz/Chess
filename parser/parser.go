package parser

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/ImVulcrum/Chess/pieces"
)

func clean_pgn(input_string string) string { //cleans a given string so that the compiler can read it

	//make spaces after the points
	parts := strings.Split(input_string, ".")
	input_string = strings.Join(parts, ". ")

	//remove tag section
	var string_without_tags = input_string
	index := strings.LastIndex(input_string, "]")                    //find the last instance of a closed bracket
	if index != -1 && !strings.Contains(input_string[index:], "[") { //Check if there are no more open brackets after the first closed bracket
		string_without_tags = input_string[index+1:] //Extract the substring starting from the character after the first closed bracket
	} //else: this string does not have a comment section

	//bring everything to one line and remove additional spaces between characters
	reg := regexp.MustCompile(`\s+`)
	cleaned := reg.ReplaceAllString(string_without_tags, " ")
	cleaned = strings.TrimSpace(cleaned)

	//remove exvery instance of x (indicates a capture) and plus (indicates a check) from the string cuz otherwise the parsing would be overly complicated
	var result strings.Builder
	for i := 0; i < len(cleaned); i++ {
		if !(string(cleaned[i]) == "x" || string(cleaned[i]) == "+") {
			result.WriteByte(cleaned[i])
		}
	}
	//remove comments
	re := regexp.MustCompile(`\{[^}]*\}`)
	cleaned_string := re.ReplaceAllString(result.String(), "")
	return cleaned_string
}

func Create_Array_Of_Moves(input_string string) []string {

	var cleaned_string = clean_pgn(input_string)

	//create a list with every move
	moves := strings.Split(cleaned_string, " ")
	var cleanedMoves []string
	for _, move := range moves {
		if move != "" && !unicode.IsDigit(rune(move[0])) { //wenn der string mit einer nummer startet, das heißt enweder wenn es sich um die zahl des moves handelt oder wenn das endergebnis aufgeschrieben wird
			cleanedMoves = append(cleanedMoves, move)
		}
	}
	return cleanedMoves
}

func Get_Correct_Move(move string, pieces_a [64]pieces.Piece, current_king_index int) (int, int, string) {
	var field [2]uint16
	var piece_executing_move = 64
	var index_of_correct_legal_move = 64
	var pawn_promotion_to_piece = "A" //A indicates that there is no pawn promotion

	if move[len(move)-1] != 'O' { //if the first element isn't a "O", this is not a castle, therefore it must be a normal move
		if move[len(move)-2:len(move)-1] == "=" { //if the secodnf last element of the string is a "=", this move contains a Pawn promotion
			pawn_promotion_to_piece = move[len(move)-1:]
			move = move[:len(move)-2]
		}

		field = Get_Field_From_Move(move)

		if firstChar := rune(move[0]); !unicode.IsUpper(firstChar) { //pawn move cuz the string does not start with an uppercase letter
			if len(move) == 2 { // Simple Pawn Move
				piece_executing_move, index_of_correct_legal_move = Get_Piece_Index_And_Move_Index(pieces_a, field, pieces_a[current_king_index].Is_White_Piece(), "", "0", 64)
			} else if len(move) == 3 { //pawn take move or pawn move that could be executed by more than one pieces
				piece_executing_move, index_of_correct_legal_move = Get_Piece_Index_And_Move_Index(pieces_a, field, pieces_a[current_king_index].Is_White_Piece(), "", string(move[0]), 64)
			}
		} else { //Piece move cuz the move string starts with an uppercase letter
			if len(move) == 3 { // simple piece move, the first character of the move string indicates the piece that is moving
				piece_executing_move, index_of_correct_legal_move = Get_Piece_Index_And_Move_Index(pieces_a, field, pieces_a[current_king_index].Is_White_Piece(), string(move[0]), "0", 64)
			} else if len(move) == 4 { //piece move with position, the first character of the move string indicates the piece that is moving, the second additional information to the starting position
				piece_executing_move, index_of_correct_legal_move = Get_Piece_Index_And_Move_Index(pieces_a, field, pieces_a[current_king_index].Is_White_Piece(), string(move[0]), string(move[1]), 64)
			}
		}

	} else {
		if move == "O-O" { //short castle
			piece_executing_move, index_of_correct_legal_move = Get_Piece_Index_And_Move_Index(pieces_a, [2]uint16{6, pieces_a[current_king_index].Give_Pos()[1]}, pieces_a[current_king_index].Is_White_Piece(), "K", "0", 64)
		} else if move == "O-O-O" { //long castle
			piece_executing_move, index_of_correct_legal_move = Get_Piece_Index_And_Move_Index(pieces_a, [2]uint16{2, pieces_a[current_king_index].Give_Pos()[1]}, pieces_a[current_king_index].Is_White_Piece(), "K", "0", 64)
		} else {
			fmt.Println("Error while Reading Premove File: Expected either (O-O) or (O-O-O), got", move, "instead")
		}
	}
	if piece_executing_move == 64 { //64 means there is no piece matching the given specifications
		fmt.Println("\n-------------------------------------------------------------------\n" + `there is no piece with the move: "` + move + `"` + "\n-------------------------------------------------------------------\n")
	}
	return piece_executing_move, index_of_correct_legal_move, pawn_promotion_to_piece
}

func Write_PGN_File(pgn_moves_a []string, name_player_white string, name_player_black string) {
	var moves_counter = 1
	pgn_moves_a = pgn_moves_a[1:]
	year, month, day := time.Now().Date()

	//create the header vars
	var date = `[Date "` + strconv.Itoa(year) + "." + strconv.Itoa(int(month)) + "." + strconv.Itoa(day) + `"]`
	var identifier = `[Identifier "` + (strconv.Itoa(int(time.Now().UnixMilli()))) + `"]` //important for the file naming
	var white = `[White "` + name_player_white + `"]`
	var black = `[Black "` + name_player_black + `"]`
	var pgn_string string

	for i := 0; i < len(pgn_moves_a); i++ { //create the pgn string
		if i%2 == 0 { //every two moves a move counter must be included
			pgn_string = pgn_string + strconv.Itoa(moves_counter) + ". "
			moves_counter++
		}
		pgn_string = pgn_string + pgn_moves_a[i] + " "
	}

	//create the pgn file
	file, _ := os.Create("./" + name_player_white + " vs " + name_player_black + " " + strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day) + "-" + strconv.Itoa(int(time.Now().UnixMilli())) + ".pgn")

	file.WriteString(date + "\n" + identifier + "\n" + white + "\n" + black + "\n\n" + pgn_string)
}

func Get_Piece_Index_And_Move_Index(pieces_a [64]pieces.Piece, field [2]uint16, white_is_current_player bool, piece_type string, position string, exclude_piece_index int) (int, int) {
	var current_piece_type = "A"

	for i := 0; i < len(pieces_a); i++ { //loop trough the pieces array
		if pieces_a[i] != nil && pieces_a[i].Is_White_Piece() == white_is_current_player && i != exclude_piece_index { //if there is a piece that matches the color of the current player and it is not the piece that must be excluded
			for k := 0; k < len(pieces_a[i].Give_Legal_Moves()); k++ { //loop trough the moves of this piece
				if pieces_a[i].Give_Legal_Moves()[k][0] == field[0] && pieces_a[i].Give_Legal_Moves()[k][1] == field[1] { //there is a piece in the correct color with the given move
					cord, is_x_cord := Translate_PGN_Field_Notation(position)                                                               //get the cord
					if (is_x_cord && cord == pieces_a[i].Give_Pos()[0]) || (!is_x_cord && cord == pieces_a[i].Give_Pos()[1]) || cord == 8 { //check if the piece has the given x or y cord or has no cord specifictaion indicated by cord beeing 8
						switch pieces_a[i].(type) { //set the piece type
						case *pieces.Rook:
							current_piece_type = "R"
						case *pieces.King:
							current_piece_type = "K"
						case *pieces.Pawn:
							current_piece_type = ""
						case *pieces.Queen:
							current_piece_type = "Q"
						case *pieces.Bishop:
							current_piece_type = "B"
						case *pieces.Knight:
							current_piece_type = "N"
						default:
							fmt.Println("Error in Parser while iterating through pieces array: Unexpected piece type")
						}
						if current_piece_type == piece_type {
							return i, k
						}
					}
				}
			}
		}
	}
	if exclude_piece_index == 64 {
		fmt.Println("Error in Parser: there is no piece is the pieces array that matches the specifications given, which means that the given pgn file is corrupted")
	}
	return 64, 0 //if there is no piece mathcing the given specifications return 64
}

func Translate_PGN_Field_Notation(cord_string string) (uint16, bool) { //translates a pgn string of len 1 to uint16 values readable for the compiler
	var cord uint16
	var is_x_cord bool

	if len(cord_string) != 1 {
		cord = 8
		fmt.Println("Error: Unexpected lenght of string while trying to convert it from pgn field notation to a square notation")
	} else {
		if unicode.IsDigit(rune(cord_string[0])) { //if it is a y cord, it must be reversed, cuz the engine works from top to bottom when it comes to the y cords of the board, which is contrary to the way a normal chess board is named
			is_x_cord = false
			num, _ := strconv.Atoi(cord_string)
			num = num - 8
			num = -num
			cord = uint16(num)
		} else { //if it is a x cord, just convert to the integer
			is_x_cord = true
			cord = uint16(cord_string[0] - 'a')
		}
	}
	return cord, is_x_cord
}

func Get_Field_From_Move(move string) [2]uint16 { //get a field (readable for the compiler) from a pgn move string
	var field [2]uint16

	var x_cord = move[len(move)-2 : len(move)-1]
	var y_cord = move[len(move)-1:]

	x, _ := Translate_PGN_Field_Notation(x_cord)
	y, _ := Translate_PGN_Field_Notation(y_cord)

	field = [2]uint16{x, y}
	return field
}

func Translate_Field_Cord_To_PGN_String(field_cord uint16, is_x_cord bool) string { //translates a field cord (one digit) (chess engine language) to pgn string
	if is_x_cord {
		return string(rune(int("a"[0]) + int(field_cord)))
	} else {
		return strconv.Itoa(-1 * (int(field_cord) - 8))
	}
}

func Get_Move_As_String_From_Field(field [2]uint16) string {
	return Translate_Field_Cord_To_PGN_String(field[0], true) + Translate_Field_Cord_To_PGN_String(field[1], false)
}
