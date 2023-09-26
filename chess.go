package main

import (
	"fmt"
	"image"
	"os"
	"time"

	"./buttons"
	"./clipboard"
	gfx "./gfxw"
	"./imaging"
	"./parser"
	"./pieces"
	"./time_counter"
)

func main() {
	fmt.Println("start game")
	var restart_window bool = false

restart_marker:

	var use_clipboard_as_premoves = false
	var a uint16 = 100
	var w_x, w_y uint16 = 10 * a, 8 * a
	var duration_of_premove_animation int = 1
	var deselect_piece_after_clicking = false
	var game_history_can_be_changed = true
	var game_timer int64 = 500000

	//~ var restart bool
	var white_is_current_player bool = false
	var player_change bool = true
	var current_king_index int
	var checkmate bool
	var check bool
	var current_piece pieces.Piece
	var temp_current_piece pieces.Piece
	var piece_index int
	var current_legal_moves [][3]uint16
	var moves_counter int16
	var current_field [2]uint16
	var dragging bool
	var ending_premoves bool = true
	var piece_is_selected uint16 = 64
	var premoves string = get_clipboard_if_asked(use_clipboard_as_premoves)
	var promotion uint16 = 0

	pieces_a, white_king_index, black_king_index, one_move_back, one_move_forward, restart_button, pause_button, moves_a, white_time_counter, black_time_counter := initialize(w_x, w_y, a, restart_window, game_timer)
	premoves_array := parser.Create_Array_Of_Moves(premoves)

	draw_pieces(pieces_a, w_x, w_y, a)

	display_message(2, a, false)

	m_channel := make(chan [4]int16, 1)
	gor_status := make(chan bool)

	go mouse_handler(m_channel, gor_status)

	for { //gameloop

		if player_change {
			if int(moves_counter) == len(moves_a) { //apppend the moves_array when a new move was made
				moves_a = append_moves_array(moves_a, pieces_a)
			} else if !array_one_is_equal_to_array_two(pieces_a, moves_a[moves_counter]) { //if there was a move changed
				if game_history_can_be_changed { //if it's allowed to change a move in the history, this block replaces the move in the moves array with the current one and cuts of the rest of the array
					fmt.Println("there was a move changed in the game history")
					moves_a[moves_counter] = pieces.Copy_Array(pieces_a)
					moves_a = moves_a[:moves_counter+1]
				} else { //if it's not allowed it just deletes the move made and resets
					fmt.Println("changed game history although this is not allowed")
					pieces_a = pieces.Copy_Array(moves_a[moves_counter])
				}
			}

			//restart = false
			white_is_current_player, current_king_index, player_change, moves_counter = change_player(white_is_current_player, white_king_index, black_king_index, moves_counter) // this function sets the current_king_index
			pieces_a, checkmate = pieces.Calc_Moves_With_Check(pieces_a, moves_counter, current_king_index)                                                                       //calce the move
			check = pieces_a[current_king_index].(*pieces.King).Is_In_Check(pieces_a, moves_counter)                                                                              //check if the current king is in check

			if !use_clipboard_as_premoves || duration_of_premove_animation != 0 { //wenn keine premoves mehr da sind oder premoves da sind und diese gezeichnet werden sollen dann wird das board gezeichnet
				draw_board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, false, current_king_index, check)
			}

			if checkmate && check { //game end / restart
				if check {
					display_message(0, a, white_is_current_player)
				} else {
					display_message(1, a, white_is_current_player)
				}
				restart_window = true
				gor_status <- true
				goto restart_marker
			}

			//timer changing --> can be moved into to change player function
			if pause_button.Give_State() {
				pause_button.Switch(0, 0, 0)
			}
			if white_is_current_player {
				black_time_counter.Stop_Counting()
				white_time_counter.Init_Counting()

			} else if !white_is_current_player {
				white_time_counter.Stop_Counting()
				black_time_counter.Init_Counting()
			}
		}

		if len(premoves_array) > 0 { //if there are premoves the program premoves
			var promotion uint16 = 0
			fmt.Println("premove:", premoves_array[0])
			piece_executing_move, index_of_move, piece_promoting_to := parser.Get_Correct_Move(premoves_array[0], pieces_a, current_king_index)
			pieces_a, promotion = pieces.Move_Piece_To(pieces_a[piece_executing_move], pieces_a[piece_executing_move].Give_Legal_Moves()[index_of_move], moves_counter, pieces_a)
			if promotion != 64 {
				pieces_a, _ = pawn_promotion(w_x, w_y, a, piece_executing_move, pieces_a, piece_promoting_to, m_channel)
			}
			premoves_array = premoves_array[1:]

			time.Sleep(time.Duration(duration_of_premove_animation) * time.Millisecond)
			player_change = true
		} else if ending_premoves { //soll im optimalfall nur einmal ausgeführt werden
			draw_board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, false, current_king_index, check)
			use_clipboard_as_premoves = false
			ending_premoves = false
			fmt.Println("----- End of Premoves -----")
		}

		if !use_clipboard_as_premoves { //wenn die andere restart methode ausgeführt wird, dann sollt hier noch noch überprüft werden ob restart == false ist, da es sonst zu fehlern kommt
			select {
			case mouse_input := <-m_channel:

				if dragging && !(mouse_input[0] == 1 && mouse_input[1] == 0) { //wenn nichts gehalten wird dann löschen
					draw_board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, true, current_king_index, check)
					dragging = false
				}

				if mouse_input[0] == 1 && mouse_input[1] == 1 {
					if one_move_back.Is_Clicked(uint16(mouse_input[2]), uint16(mouse_input[3])) {
						if moves_counter >= 2 {
							moves_counter = moves_counter - 2
							pieces_a = pieces.Copy_Array(moves_a[moves_counter])
							player_change = true
						} else {
							fmt.Println("first move has been reached")
						}
					} else if one_move_forward.Is_Clicked(uint16(mouse_input[2]), uint16(mouse_input[3])) {
						if int(moves_counter) < len(moves_a) {
							pieces_a = pieces.Copy_Array(moves_a[moves_counter])
							player_change = true
						} else {
							fmt.Println("last move has been reached")
						}
					} else if restart_button.Is_Clicked(uint16(mouse_input[2]), uint16(mouse_input[3])) {
						restart_window = true
						gor_status <- true
						goto restart_marker
					} else if pause_button.Is_Clicked(uint16(mouse_input[2]), uint16(mouse_input[3])) {
						var pausing bool = pause_button.Switch(38, 37, 34)
						if pausing {
							white_time_counter.Stop_Counting()
							black_time_counter.Stop_Counting()
						} else {
							if white_is_current_player {
								white_time_counter.Init_Counting()
							} else {
								black_time_counter.Init_Counting()
							}
						}
					} else {

						current_field = calc_field(a, uint16(mouse_input[2]), uint16(mouse_input[3]), 0)
						temp_current_piece, piece_index = get_current_piece(pieces_a, current_field)

						//auswählen eines pieces
						if temp_current_piece != nil && temp_current_piece.Is_White_Piece() == white_is_current_player && (!deselect_piece_after_clicking || (current_piece == nil || temp_current_piece.Give_Pos() != current_piece.Give_Pos())) {
							current_piece = temp_current_piece
							current_legal_moves = current_piece.Give_Legal_Moves()
							promotion = 0
							draw_board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, true, current_king_index, check)
							piece_is_selected = uint16(piece_index)

						} else if piece_is_selected != 64 {
							//überprüfen ob auf ein Feld in Legal Moves geklickt wurde
							pieces_a, piece_is_selected, player_change, promotion = move_if_current_field_is_in_legal_moves(current_field, pieces_a, promotion, piece_is_selected, a, w_x, w_y, current_king_index, check, moves_counter, m_channel)

							//Deselect sobald ein Piece ausgewählt ist: wenn auf kein Piece geklickt wurde, oder auf ein gegnerisches Piece geklickt wurde, wird das ausgewählte Piece deselected
							//mit der Option deselect_piece_after_clicking kann eingestellt werde, ob auch deselected werden soll wenn auf dasslebe piece nocheinmal geklickt wurde
							if (temp_current_piece == nil) || (temp_current_piece != nil && (temp_current_piece.Give_Pos() == current_piece.Give_Pos() || temp_current_piece.Is_White_Piece() != white_is_current_player)) {
								current_piece = nil
								draw_board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, false, current_king_index, check)
								piece_is_selected = 64
							}
						}
					}

				} else if mouse_input[1] == -1 && mouse_input[0] == 1 {
					current_field = calc_field(a, uint16(mouse_input[2]), uint16(mouse_input[3]), 0)
					temp_current_piece, _ = get_current_piece(pieces_a, current_field)

					if piece_is_selected != 64 {
						//überprüfen ob in current legal moves
						pieces_a, piece_is_selected, player_change, promotion = move_if_current_field_is_in_legal_moves(current_field, pieces_a, promotion, piece_is_selected, a, w_x, w_y, current_king_index, check, moves_counter, m_channel)

						//Deselect sobald ein Piece ausgewählt ist: wenn auf keinem Piece losgelassen wurde oder auf ein generisches Piece losgelassen wurde oder wenn auf einem eigenen piece losgelassen wurde
						if (temp_current_piece == nil) || (temp_current_piece != nil && ((temp_current_piece.Is_White_Piece() != white_is_current_player) || (temp_current_piece.Is_White_Piece() == white_is_current_player && temp_current_piece.Give_Pos() != current_piece.Give_Pos()))) {
							current_piece = nil
							draw_board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, false, current_king_index, check)
							piece_is_selected = 64
						}
					}
				} else if mouse_input[1] == 0 && mouse_input[0] == 1 && current_piece != nil && current_piece.Is_White_Piece() == white_is_current_player {
					//wenn die maustaste gehalten wird, wird ein ghosttpiece gemalt, welches der maus folgt
					dragging = true
					pieces.Draw_To_Point(current_piece, w_x, w_y, a, uint16(mouse_input[2]), uint16(mouse_input[3]), -int16(a/2), -int16(a/2), 50, uint16(mouse_input[2]))
				}

			default:
				time.Sleep(5 * time.Millisecond)
			}
			if draw_timers(white_time_counter, black_time_counter, a) {
				display_message(0, a, white_is_current_player)
				restart_window = true
				gor_status <- true
				goto restart_marker
			}
		}
	}
}

func mouse_handler(m_channel chan [4]int16, gor_status chan bool) {
	for {
		select {
		case quit := <-gor_status:
			if quit {
				return
			}
		default:
			button, status, m_x, m_y := gfx.MausLesen1()
			if !(button == 0 && status == 0) {
				select {
				case temp := <-m_channel: //stellt sicher dass leer ist
					if temp[0] == 1 && temp[1] == 1 { //wenn es ein klicken befehl ist, dann soll dieser nicht überschrieben werden, da es sonst zu bugs kommen kann und "verschluckt" wird, dass ein piece ausgewählt wurde
						m_channel <- temp
					} else {
						m_channel <- [4]int16{int16(button), int16(status), int16(m_x), int16(m_y)} //schreibt nur wenn leer ist
					}
				default:
					m_channel <- [4]int16{int16(button), int16(status), int16(m_x), int16(m_y)}
				}
			}
		}
	}
}

func get_clipboard_if_asked(use_clipboard bool) string {
	var premoves string

	if use_clipboard {
		err := clipboard.Init()
		if err != nil {
			panic(err)
		}
		premoves = string(clipboard.Read(clipboard.FmtText))
	} else {
		premoves = ""
	}
	return premoves
}

func draw_timers(white_time_counter, black_time_counter time_counter.Counter, a uint16) bool {
	gfx.SetzeFont("./resources/fonts/firamono.ttf", int((a*10)/38))

	gfx.UpdateAus()
	gfx.Stiftfarbe(86, 82, 77)
	gfx.Vollrechteck(81*a/10, 60*a/10, 2*a-2*a/10, a-65*a/100)
	gfx.Stiftfarbe(48, 46, 43)
	gfx.Vollrechteck(89*a/10, 60*a/10, a-8*a/10, a-65*a/100)
	gfx.Stiftfarbe(255, 255, 255)
	white_string, white_has_no_time := white_time_counter.Return_Current_Counter()
	black_string, black_has_no_time := black_time_counter.Return_Current_Counter()
	black_time_counter.Return_Current_Counter()
	gfx.SchreibeFont(81*a/10, 60*a/10, white_string)
	gfx.SchreibeFont(91*a/10, 60*a/10, black_string)
	gfx.UpdateAn()

	if white_has_no_time || black_has_no_time {
		return true
	} else {
		return false
	}
}

func array_one_is_equal_to_array_two(array_a [64]pieces.Piece, array_b [64]pieces.Piece) bool {
	for i := 0; i < len(array_a); i++ {
		if fmt.Sprint(array_a[i]) != fmt.Sprint(array_b[i]) || fmt.Sprintf("%T", array_a[i]) != fmt.Sprintf("%T", array_b[i]) {
			return false
		}
	}
	return true
}

func get_current_piece(pieces_a [64]pieces.Piece, current_field [2]uint16) (pieces.Piece, int) {
	//gibt das
	var temp_current_piece pieces.Piece = nil
	var piece_index int
	for piece_index = 0; piece_index < len(pieces_a); piece_index++ {
		if pieces_a[piece_index] != nil {
			if current_field == pieces_a[piece_index].Give_Pos() {
				temp_current_piece = pieces_a[piece_index]
				break
			}
		}
	}
	return temp_current_piece, piece_index
}

func move_if_current_field_is_in_legal_moves(current_field [2]uint16, pieces_a [64]pieces.Piece, promotion uint16, piece_is_selected uint16, a, w_x, w_y uint16, current_king_index int, check bool, moves_counter int16, m_channel chan [4]int16) ([64]pieces.Piece, uint16, bool, uint16) {
	current_piece := pieces_a[piece_is_selected]
	current_legal_moves := current_piece.Give_Legal_Moves()
	for k := 0; k < len(current_legal_moves); k++ {
		if current_field == [2]uint16{current_legal_moves[k][0], current_legal_moves[k][1]} { //wenn das der Fall ist, wird das Piece bewegt
			//moved something
			pieces_a, promotion = pieces.Move_Piece_To(current_piece, current_legal_moves[k], moves_counter, pieces_a)
			if promotion != 64 {
				draw_board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, false, current_king_index, check)
				pieces_a, _ = pawn_promotion(w_x, w_y, a, int(piece_is_selected), pieces_a, "A", m_channel)
			}
			piece_is_selected = 64
			promotion = 0
			return pieces_a, 64, true, 0
		}
	}
	return pieces_a, piece_is_selected, false, promotion
}

func display_message(message_type uint8, a uint16, white_is_current_player bool) {

	gfx.Stiftfarbe(0, 0, 0)
	gfx.Transparenz(70)
	gfx.Vollrechteck(0, 0, a*8, a*8)
	gfx.Transparenz(50)
	gfx.Stiftfarbe(8, 8, 8)
	gfx.Vollrechteck((a), 3*a, 6*a, 2*a)

	gfx.Stiftfarbe(220, 220, 220)
	gfx.SetzeFont("./resources/fonts/junegull.ttf", int(5*a/10))

	if message_type == 0 {
		if !white_is_current_player {
			gfx.SchreibeFont(18*a/10, 31*a/10, "player black has")
		} else {
			gfx.SchreibeFont(17*a/10, 31*a/10, "player white has")
		}
		gfx.Stiftfarbe(136, 8, 8)
		gfx.SetzeFont("./resources/fonts/punk.ttf", int(a))
		gfx.SchreibeFont(29*a/10, 37*a/10, "lost")
	} else if message_type == 1 {
		gfx.SchreibeFont(295*a/100, 31*a/10, "That's a ")

		gfx.Stiftfarbe(0, 143, 230)
		gfx.SetzeFont("./resources/fonts/punk.ttf", int(a))
		gfx.SchreibeFont(144*a/100, 37*a/10, "stalemate")

	} else if message_type == 2 {
		gfx.SchreibeFont(180*a/100, 31*a/10, "Press any key to")

		gfx.Stiftfarbe(28, 205, 60)
		gfx.SetzeFont("./resources/fonts/punk.ttf", int(a))
		gfx.SchreibeFont(260*a/100, 37*a/10, "star t")
	}

	gfx.Transparenz(0)
	gfx.UpdateAn()
	gfx.TastaturLesen1()

}

func pawn_promotion(w_x, w_y, a uint16, pawn_index int, pieces_a [64]pieces.Piece, premoved string, m_channel chan [4]int16) ([64]pieces.Piece, string) {
	var piece_promoted_to string
	var queen pieces.Piece = pieces.NewQueen(pieces_a[pawn_index].Give_Pos()[0], pieces_a[pawn_index].Give_Pos()[1], pieces_a[pawn_index].Is_White_Piece())
	var knight pieces.Piece = pieces.NewKnight(pieces_a[pawn_index].Give_Pos()[0], pieces_a[pawn_index].Give_Pos()[1], pieces_a[pawn_index].Is_White_Piece())
	var rook pieces.Piece = pieces.NewRook(pieces_a[pawn_index].Give_Pos()[0], pieces_a[pawn_index].Give_Pos()[1], pieces_a[pawn_index].Is_White_Piece())
	var bishop pieces.Piece = pieces.NewBishop(pieces_a[pawn_index].Give_Pos()[0], pieces_a[pawn_index].Give_Pos()[1], pieces_a[pawn_index].Is_White_Piece())

	if premoved == "A" {
		gfx.Stiftfarbe(0, 0, 0)
		gfx.Transparenz(70)
		gfx.Vollrechteck(0, 0, a*8, a*8)
		gfx.Transparenz(0)
		gfx.Stiftfarbe(221, 221, 221)
		gfx.Vollrechteck((2 * a), (4*a)-(5*a)/10, 4*a, a)

		pieces.Draw_To_Point(queen, w_x, w_y, a, (2 * a), (4*a)-(5*a)/10, 0, 0, 0, 0)
		pieces.Draw_To_Point(knight, w_x, w_y, a, (3 * a), (4*a)-(5*a)/10, 0, 0, 0, 0)
		pieces.Draw_To_Point(rook, w_x, w_y, a, (4 * a), (4*a)-(5*a)/10, 0, 0, 0, 0)
		pieces.Draw_To_Point(bishop, w_x, w_y, a, (5 * a), (4*a)-(5*a)/10, 0, 0, 0, 0)
		for {
			mouse_input := <-m_channel
			if mouse_input[0] == 1 && mouse_input[1] == 1 && uint16(mouse_input[2]) >= 2*a && uint16(mouse_input[2]) <= 6*a && uint16(mouse_input[3]) >= (4*a)-(5*a)/10 && uint16(mouse_input[3]) <= (4*a)+(5*a)/10 {
				x_field_pos := (calc_field(a, uint16(mouse_input[2]), uint16(mouse_input[3]), (5*a)/10))[0]
				if x_field_pos == 2 {
					piece_promoted_to = "Q"
					pieces_a[pawn_index] = queen
				} else if x_field_pos == 3 {
					piece_promoted_to = "K"
					pieces_a[pawn_index] = knight
				} else if x_field_pos == 4 {
					piece_promoted_to = "R"
					pieces_a[pawn_index] = rook
				} else if x_field_pos == 5 {
					piece_promoted_to = "B"
					pieces_a[pawn_index] = bishop
				}
				break
			}
		}

	} else {
		piece_promoted_to = premoved

		if premoved == "Q" {
			pieces_a[pawn_index] = queen
		} else if premoved == "N" {
			pieces_a[pawn_index] = knight
		} else if premoved == "R" {
			pieces_a[pawn_index] = rook
		} else if premoved == "B" {
			pieces_a[pawn_index] = bishop
		}
	}
	return pieces_a, piece_promoted_to
}

func draw_board(a, w_x, w_y uint16, current_piece pieces.Piece, current_legal_moves [][3]uint16, pieces_a [64]pieces.Piece, highlighting_is_activated bool, current_king_index int, check bool) {
	gfx.UpdateAus()
	draw_background(a)
	if check {
		highlight(a, pieces_a[current_king_index].Give_Pos(), 255, 0, 0)
	}
	if highlighting_is_activated {
		highlight(a, current_piece.Give_Pos(), 0, 50, 255)
		for k := 0; k < len(current_legal_moves); k++ {
			highlight(a, [2]uint16{current_legal_moves[k][0], current_legal_moves[k][1]}, 0, 255, 0)
		}
	}
	draw_pieces(pieces_a, w_x, w_y, a)
	gfx.UpdateAn()
}

func change_player(white_is_current_player bool, white_king_index, black_king_index int, moves_counter int16) (bool, int, bool, int16) {
	var current_king_index int
	moves_counter = moves_counter + 1
	if white_is_current_player {
		current_king_index = black_king_index
		white_is_current_player = false
	} else {
		current_king_index = white_king_index
		white_is_current_player = true
	}
	return white_is_current_player, current_king_index, false, moves_counter
}

func initialize(w_x, w_y, a uint16, restart bool, game_timer int64) ([64]pieces.Piece, int, int, buttons.Button, buttons.Button, buttons.Button, buttons.Button, [][64]pieces.Piece, time_counter.Counter, time_counter.Counter) {
	var one_move_back buttons.Button
	var one_move_forward buttons.Button
	var restart_button buttons.Button
	var pause_button buttons.Button
	var moves_a [][64]pieces.Piece

	if !restart {
		gfx.Fenster(w_x, w_y)
		gfx.Fenstertitel("Chess")
		rescale_image(a)
	}

	one_move_back = *buttons.New(8*a+a/10, 7*a+a/10, a-a/5, a-a/5, "<", 38, 37, 34, 200, 200, 200, (a / 4), int(a/2))
	one_move_forward = *buttons.New(9*a+a/10, 7*a+a/10, a-a/5, a-a/5, ">", 38, 37, 34, 200, 200, 200, (a / 4), int(a/2))
	restart_button = *buttons.New(8*a+a/10, 65*a/10, a-a/15, 3*a/10, "restart", 86, 82, 77, 200, 200, 200, (a / 15), int(a/5))
	pause_button = *buttons.New(9*a+a/6, 65*a/10, a-a/4, 3*a/10, "pause", 86, 82, 77, 200, 200, 200, (a / 15), int(a/5))

	gfx.Stiftfarbe(86, 82, 77)
	gfx.Vollrechteck(8*a, 7*a, 2*a, a)
	gfx.Stiftfarbe(48, 46, 43)
	gfx.Vollrechteck(8*a, 6*a, 2*a, a)
	one_move_back.Draw()
	one_move_forward.Draw()
	restart_button.Draw()
	pause_button.Draw()

	gfx.Stiftfarbe(48, 46, 43)
	gfx.Vollrechteck(8*a, 0, 2*a, 6*a)
	draw_background(a)

	gfx.Stiftfarbe(124, 119, 111)
	gfx.Vollrechteck(8*a+a/10, a/10, a, 6*a-a/5)
	gfx.Stiftfarbe(22, 21, 20)
	gfx.Vollrechteck(9*a, a/10, a-a/10, 6*a-a/5)

	var pieces_a [64]pieces.Piece
	var white_king_index int = -1
	var black_king_index int = -1

	pieces_a[0] = pieces.NewRook(0, 0, false)
	pieces_a[1] = pieces.NewKnight(1, 0, false)
	pieces_a[2] = pieces.NewBishop(2, 0, false)
	pieces_a[3] = pieces.NewQueen(3, 0, false)
	pieces_a[4] = pieces.NewKing(4, 0, false)
	pieces_a[5] = pieces.NewBishop(5, 0, false)
	pieces_a[6] = pieces.NewKnight(6, 0, false)
	pieces_a[7] = pieces.NewRook(7, 0, false)

	var i uint16
	for i = 0; i < 8; i++ {
		pieces_a[i+8] = pieces.NewPawn(i, 1, false)
	}
	for i = 0; i < 8; i++ {
		pieces_a[i+16] = pieces.NewPawn(i, 6, true)
	}

	pieces_a[24] = pieces.NewRook(0, 7, true)
	pieces_a[25] = pieces.NewKnight(1, 7, true)
	pieces_a[26] = pieces.NewBishop(2, 7, true)
	pieces_a[27] = pieces.NewQueen(3, 7, true)
	pieces_a[28] = pieces.NewKing(4, 7, true)
	pieces_a[29] = pieces.NewBishop(5, 7, true)
	pieces_a[30] = pieces.NewKnight(6, 7, true)
	pieces_a[31] = pieces.NewRook(7, 7, true)

	//moves_a = append_moves_array(moves_a, pieces_a)

	for i := 0; i < len(pieces_a); i++ {
		if pieces_a[i] != nil {
			if king, ok := pieces_a[i].(*pieces.King); ok {
				if king.Is_White_Piece() {
					white_king_index = i
				} else if !king.Is_White_Piece() {
					black_king_index = i
				}
			}
		}
	}

	white_time_counter := time_counter.New(game_timer)
	black_time_counter := time_counter.New(game_timer)

	return pieces_a, white_king_index, black_king_index, one_move_back, one_move_forward, restart_button, pause_button, moves_a, white_time_counter, black_time_counter
}

func append_moves_array(moves_a [][64]pieces.Piece, pieces_a [64]pieces.Piece) [][64]pieces.Piece {
	moves_a = append(moves_a, pieces.Copy_Array(pieces_a))
	return moves_a
}

func draw_pieces(pieces_a [64]pieces.Piece, w_x, w_y, a uint16) {
	gfx.UpdateAus()
	for i := 0; i < len(pieces_a); i++ {
		if pieces_a[i] != nil {
			gfx.Archivieren()
			pieces.Draw(pieces_a[i], w_x, w_y, a)
			gfx.Restaurieren(0, 0, w_x, w_y)

			gfx.Clipboard_einfuegenMitColorKey(pieces_a[i].Give_Pos()[0]*a, pieces_a[i].Give_Pos()[1]*a, 5, 5, 5)
		}
	}
	gfx.Archivieren()
	gfx.UpdateAn()
}

func highlight(a uint16, pos [2]uint16, r, g, b uint8) {

	var cord_x uint16 = a * pos[0]
	var cord_y uint16 = a * pos[1]

	gfx.Transparenz(170)
	gfx.Stiftfarbe(r, g, b)
	gfx.Vollrechteck(cord_x, cord_y, a, a)
	gfx.Transparenz(0)
}

func draw_background(a uint16) {
	var f_x uint16 = 0
	var f_y uint16 = 0
	for i := 0; i <= 7; i++ {
		for k := 0; k <= 7; k++ {
			if k%2 == 0 {
				if i%2 == 0 {
					gfx.Stiftfarbe(240, 217, 181)
				} else {
					gfx.Stiftfarbe(181, 136, 99)
				}

			} else {
				if i%2 == 1 {
					gfx.Stiftfarbe(240, 217, 181)
				} else {
					gfx.Stiftfarbe(181, 136, 99)
				}
			}

			gfx.Vollrechteck(f_x, f_y, a, a)
			f_x = f_x + a
		}
		f_x = 0
		f_y = f_y + a
	}
}

func rescale_image(a uint16) {

	// Open the BMP file
	file, err := os.Open("./resources/images/Pieces_Source_Original.bmp")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Decode the BMP file
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decoding file:", err)
		return
	}

	// Specify the desired width and height
	var width int = 6 * int(a)
	var height int = 2 * int(a)

	// Resize the image to the specified dimensions
	resizedImg := imaging.Resize(img, width, height, imaging.NearestNeighbor)

	// Save the resized image to a new BMP file
	err = imaging.Save(resizedImg, "./resources/images/Pieces.bmp")
	if err != nil {
		fmt.Println("Error saving file:", err)
		return
	}
}

func calc_field(a, m_x, m_y, y_offset uint16) [2]uint16 {
	var current_field [2]uint16

	current_field[0] = m_x / a
	current_field[1] = (m_y + y_offset) / a
	return current_field
}

//the old restart works like this:
//~ pieces_a, white_king_index, black_king_index, moves_counter, check, white_is_current_player, restart, player_change, moves_a, white_time_counter, black_time_counter = restart_game(w_x, w_y, a, one_move_back, one_move_forward, game_timer)
//~ dragging = false
//~ empty_channel(m_channel)

// func restart_game(w_x, w_y, a uint16, one_move_back, one_move_forward buttons.Button, game_timer int64) ([64]pieces.Piece, int, int, int16, bool, bool, bool, bool, [][64]pieces.Piece, time_counter.Counter, time_counter.Counter) {
// 	//this restart function works in parts, the problem is that it does not clear the current piece, which means a piece is technicly already selected after restarting, it seems like clearing the channel is not needed
// 	pieces_a, white_king_index, black_king_index, _, _, moves_a, white_time_counter, black_time_counter := initialize(w_x, w_y, a, true, game_timer)

// 	return pieces_a, white_king_index, black_king_index, 0, false, false, true, true, moves_a, white_time_counter, black_time_counter
// }

// func empty_channel(m_channel chan [4]int16) {
// outer:
// 	for {
// 		select {
// 		case <-m_channel:
// 		default:
// 			break outer
// 		}
// 	}
// }

// func calc_a(w_x, w_y uint16) uint16 {
// 	var a uint16
// 	if w_x < w_y {
// 		a = w_x / 8
// 	} else {
// 		a = w_y / 8
// 	}
// 	return a
// }
