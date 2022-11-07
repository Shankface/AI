package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Player Type
type Player struct {
	name      string
	computer  bool
	color     string
	opp_color string
	eval_type string
	heur      [][]int
}

// Game Type
type Game struct {
	board    *Board
	player1  *Player
	player2  *Player
	time_lim int
}

// Function to play game
func (g *Game) play() {
	var winner string
	var game_over bool
	// Keep switching turns until neither player has any legal moves left
	for {
		fmt.Println("\n----------------------------------------------------------------------")
		fmt.Println("Current Board")
		g.board.printBoard()
		winner, game_over = g.board.check_game_over() // Check if game over
		if game_over {
			fmt.Println("Neither player has any legal moves. Game Over.")
			break
		} else if g.board.player_turn == "B" { // Black's turn
			fmt.Println("---- Player 1 (Black's) Turn ----")
			g.board.choose_move(*g.player1, g.time_lim)
			g.board.player_turn = "W"
		} else if g.board.player_turn == "W" { // White's turn
			fmt.Println("---- Player 2 (White's) Turn ----")
			g.board.choose_move(*g.player2, g.time_lim)
			g.board.player_turn = "B"
		}
	}
	// Display final score and winner
	win := "Player 1 (Black)"
	if winner == "W" {
		win = "Player 2 (White)"
	}
	fmt.Printf("------ Final Score ------\nPlayer 1 (Black): %v\nPlayer 2 (White): %v\nWinner: %v\n", len(g.board.black_pos), len(g.board.white_pos), win)
}

// Board type
type Board struct {
	squares     [8][8]string
	directions  map[string][]int
	player_turn string
	black_pos   [][]int
	white_pos   [][]int
}

// Function for intializing board
func (b *Board) initBoard(path string) {
	DIRECTIONS := make(map[string][]int)
	DIRECTIONS["UP"] = []int{-1, 0}
	DIRECTIONS["DOWN"] = []int{1, 0}
	DIRECTIONS["LEFT"] = []int{0, -1}
	DIRECTIONS["RIGHT"] = []int{0, 1}
	DIRECTIONS["UP_RIGHT"] = []int{-1, 1}
	DIRECTIONS["UP_LEFT"] = []int{-1, -1}
	DIRECTIONS["DOWN_RIGHT"] = []int{1, 1}
	DIRECTIONS["DOWN_LEFT"] = []int{1, -1}
	b.directions = DIRECTIONS

	var content []string
	b.player_turn = "B"
	colors := []string{"E", "B", "W"}

	if path != "default" {
		cont, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		content = strings.Fields(string(cont))
		if content[64] == "2" {
			b.player_turn = "W"
		}
	}

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if path != "default" { // If user inputted preset board, fill in board
				n, _ := strconv.Atoi(content[(8*i)+j])
				b.squares[i][j] = colors[n]
			} else { // If not, initialize to beginning board
				b.squares[i][j] = "E"
				b.placePieces([][]int{{3, 3}, {4, 4}}, "W")
				b.placePieces([][]int{{4, 3}, {3, 4}}, "B")
			}
		}
	}
	b.update_board_data()
}

// Function to print board to terminal
func (b *Board) printBoard() {
	fmt.Print("\u001b[44m  ")
	for i := 0; i < 8; i++ {
		fmt.Printf("%v %v%v", "\u001b[44m", i, "\u001b[0m")
	}
	fmt.Println("")
	for i := 0; i < 8; i++ {
		fmt.Printf("%v %v", "\u001b[44m", i)
		for j := 0; j < 8; j++ {
			switch b.squares[i][j] {
			case "E":
				fmt.Print("\u001b[42;1m  ")
			case "W":
				fmt.Print("\u001b[47m  ")
			case "B":
				fmt.Print("\u001b[40m  ")
			default:
				fmt.Print("Unkown")
			}
		}
		fmt.Println("\u001b[0m")
	}
	fmt.Println("")
}

// Whenever new pieces are placed on board, update board data
func (b *Board) update_board_data() {
	b.white_pos = [][]int{}
	b.black_pos = [][]int{}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if b.squares[i][j] == "W" {
				b.white_pos = append(b.white_pos, []int{i, j})
			} else if b.squares[i][j] == "B" {
				b.black_pos = append(b.black_pos, []int{i, j})
			}
		}
	}
}

// Function to place list of pieces of certain color on board
func (b *Board) placePieces(pos [][]int, color string) {
	for _, coord := range pos {
		b.squares[coord[0]][coord[1]] = color
	}
	b.update_board_data()
}

// Function to check if game is over and return the winner or draw if it is over
func (b *Board) check_game_over() (string, bool) {
	legal_B_moves, _ := b.find_legal_moves("B")
	legal_W_moves, _ := b.find_legal_moves("W")
	if len(legal_B_moves) == 0 && len(legal_W_moves) == 0 {
		if len(b.black_pos) > len(b.white_pos) {
			return "B", true
		} else if len(b.black_pos) < len(b.white_pos) {
			return "W", true
		} else if len(b.black_pos) == len(b.white_pos) {
			return "Draw", true
		}
	}
	return "None", false

}

// Function to check if the baord is in early game, middle game, or late game.
func (b *Board) check_game_state() string {
	tot_pieces := len(b.black_pos) + len(b.white_pos)
	if tot_pieces < 20 {
		return "EARLY_GAME"
	} else if tot_pieces < 50 {
		return "MID_GAME"
	} else {
		return "LATE_GAME"
	}
}

// Function to check if inputted coordinate is a valid board space
func (b *Board) valid_square(pos []int) bool {
	if pos[0] > 7 || pos[0] < 0 || pos[1] > 7 || pos[1] < 0 {
		return false
	}
	return true
}

// Function to return all flipped tiles when a piece of certain color is placed in certain space on board
func (b *Board) check_bracket(pos []int, color string, opp_color string) [][]int {
	total_flipped_spaces := [][]int{pos}
	flipped_spaces := [][]int{}
	curr_row, curr_col := 0, 0
	row_inc, col_inc := 0, 0

	for _, incs := range b.directions {
		flipped_spaces = [][]int{}
		curr_row, curr_col = pos[0], pos[1]
		row_inc, col_inc = incs[0], incs[1]

		for i := 0; i < 8; i++ {
			curr_row += row_inc
			curr_col += col_inc
			if b.valid_square([]int{curr_row, curr_col}) == false || b.squares[curr_row][curr_col] == "E" {
				break
			} else if b.squares[curr_row][curr_col] == opp_color {
				flipped_spaces = append(flipped_spaces, []int{curr_row, curr_col})
			} else if b.squares[curr_row][curr_col] == color {
				total_flipped_spaces = append(total_flipped_spaces, flipped_spaces...)
				break
			}
		}
	}
	return total_flipped_spaces
}

// Function to check if certain possible spot has been checked already during 'find_legal_move' function
func already_checked(checked [][]int, pos []int) bool {
	for _, p := range checked {
		if reflect.DeepEqual(pos, p) {
			return true
		}
	}
	return false
}

// Function to find all legal moves for certain color and returns possible moves along with the tiles they would flip
func (b *Board) find_legal_moves(color string) ([][]int, [][][]int) {
	var opponent_pos [][]int
	var opp_color string
	var temp_pos []int
	var checked [][]int
	var legal_moves [][]int
	var flipped_tiles [][][]int
	var flipped [][]int

	if color == "W" {
		opponent_pos = b.black_pos
		opp_color = "B"
	} else if color == "B" {
		opponent_pos = b.white_pos
		opp_color = "W"
	}

	for _, pos := range opponent_pos {
		for j := -1; j < 2; j++ {
			for k := -1; k < 2; k++ {

				temp_pos = []int{pos[0] + j, pos[1] + k}
				if b.valid_square(temp_pos) && b.squares[temp_pos[0]][temp_pos[1]] == "E" && already_checked(checked, temp_pos) == false {
					flipped = b.check_bracket(temp_pos, color, opp_color)

					if len(flipped) > 1 {
						legal_moves = append(legal_moves, temp_pos)
						flipped_tiles = append(flipped_tiles, flipped)
					}
					checked = append(checked, temp_pos)
				}
			}
		}
	}
	return legal_moves, flipped_tiles
}

// Function that prompts the computer or user to choose a move
func (b *Board) choose_move(player Player, time_lim int) {

	legal_moves, flipped_tiles := b.find_legal_moves(player.color)

	if len(legal_moves) > 0 {

		fmt.Printf("Legal Moves for %v: [row column]\n", player.name)
		for i, move := range legal_moves {
			fmt.Printf("Move %v: %v\n", i, move)
		}

		// If player is computer
		if player.computer == true {
			index := 0
			if len(legal_moves) > 1 {
				s1 := rand.NewSource(time.Now().UnixNano())
				r1 := rand.New(s1)
				if player.eval_type != "random" {
					best_moves, depth, time := iterative_search(b, player, time_lim, flipped_tiles) // Do iterative deepening of minmax search
					r := r1.Intn(len(best_moves))
					index = best_moves[r] // if multiple moves are tied for best, choose randomlly from them
					fmt.Printf("Computer searched depth %v to completion. Total time searching: %v sec\n", depth, time)
				} else { // if the computer's evaluation type is random, choose from possible moves
					index = r1.Intn(len(legal_moves))
				}
			}
			fmt.Printf("Computer Chose Move Number: %v %v\n", index, legal_moves[index])
			b.placePieces(flipped_tiles[index], player.color)

		} else { // If player is human
			move_num, _ := strconv.Atoi(get_user_input("Input Move Number: ", "int", 0, len(legal_moves)-1))
			b.placePieces(flipped_tiles[move_num], player.color)
		}
	} else {
		fmt.Printf("No Legal Moves for %v, skipping to other player.\n", player.name)
	}
}

// Function that dynamically evaluates board score depending on stage of the game
func (b *Board) eval_score(player Player) int {
	// return static heuristic score if eval type of player is "static"
	if player.eval_type == "static" {
		return b.heur_score(player)
	} else if player.eval_type == "dynamic" {
		switch b.check_game_state() {
		case "EARLY_GAME":
			// Corners and mobility are weighted higher and personal piece amount is weighted negatively in early game
			return ((20 * b.corner_score(player)) + (5 * b.around_empty_corner_score(player)) + (6 * b.mobility_score(player)) + (-1 * b.piece_diff_score(player)))

		case "MID_GAME":
			// Personal piece count is no longer weighted negatively and corners and mobility are weighted high in middle game
			return ((20 * b.corner_score(player)) + (5 * b.around_empty_corner_score(player)) + (6 * b.mobility_score(player))) // + (1 * b.piece_diff_score(player)))

		case "LATE_GAME":
			// Corner and Piece count are wieghted heavily in late game
			return ((20 * b.corner_score(player)) + (5 * b.around_empty_corner_score(player)) + (4 * b.mobility_score(player)) + (10 * b.piece_diff_score(player)))

		default:
		}
	}
	return 0
}

// Function to calculate static heuristic score of player based on pieces on the board
func (b *Board) heur_score(player Player) int {
	score := 0
	spaces := b.black_pos
	if player.color == "W" {
		spaces = b.white_pos
	}
	for _, pos := range spaces {
		score += player.heur[pos[0]][pos[1]]
	}
	return score
}

// Function to calculate corner score based on player's corner peices vs opponent's corner pieces
func (b *Board) corner_score(player Player) int {
	score, opp_score := 0, 0
	corners := [][]int{{0, 0}, {0, 7}, {7, 0}, {7, 7}}
	for i := range corners {
		if b.squares[corners[i][0]][corners[i][1]] == player.color {
			score++
		} else if b.squares[corners[i][0]][corners[i][1]] == player.opp_color {
			opp_score++
		}
	}
	return 100 * (score - opp_score) / (opp_score + score + 1)
}

// Function to calculate player's pieces around empty corners compared to opponents pieces around empty corners
func (b *Board) around_empty_corner_score(player Player) int {
	score, opp_score := 0, 0
	corners := [][]int{{0, 0}, {0, 7}, {7, 0}, {7, 7}}
	edges := [][]int{{0, 1}, {1, 0}, {1, 1}, {0, 6}, {1, 6}, {1, 7}, {6, 0}, {6, 1}, {7, 1}, {6, 6}, {6, 7}, {7, 6}}
	for i, corner := range corners {
		if b.squares[corner[0]][corner[1]] == "E" {
			for j := i * 3; j < ((i * 3) + 3); j++ {
				if b.squares[edges[j][0]][edges[j][1]] == player.color {
					score++
				} else if b.squares[edges[j][0]][edges[j][1]] == player.opp_color {
					opp_score++
				}
			}
		}
	}

	return -100 * (score - opp_score) / (opp_score + score + 1)
}

// Function to calculate players mobility compared to opponent's mobility
func (b *Board) mobility_score(player Player) int {
	moves, _ := b.find_legal_moves(player.color)
	opp_moves, _ := b.find_legal_moves(player.opp_color)
	return int(100 * (len(moves) - len(opp_moves)) / (len(opp_moves) + len(moves) + 1))
}

// Function to calculate players number of pieces compared to opponent's pieces
func (b *Board) piece_diff_score(player Player) int {
	score := len(b.black_pos)
	opp_score := len(b.white_pos)
	if player.color == "W" {
		score = len(b.white_pos)
		opp_score = len(b.black_pos)
	}
	return 100 * (score - opp_score) / (opp_score + score)
}

// Function that uses iterative deepening and alpha-beta search and returns best moves given current board and player
func iterative_search(board *Board, player Player, max_sec int, flipped_tiles [][][]int) ([]int, int, float64) {
	var new_board Board
	start := time.Now()
	last_completed_depth, curr_max_depth := 0, 0
	var best_moves []int
	last_time := 0.0

	for {
		temp_best_moves := []int{}
		max_score := float64(math.MinInt64)
		complete := true

		// If program has searched entire remaining game, break
		if curr_max_depth == 64-(len(board.white_pos)+len(board.black_pos)) {
			break
		}

		// When starting search of new depth, it takes time it took to search a move up to previous depth, multiplies it, and checks to see if it could finish in time remaining
		if (last_time*1.4)*float64(len(flipped_tiles))+time.Since(start).Seconds() > (float64(max_sec) * .8) {
			complete = false
			break
		}

		curr_max_depth++
		for i := range flipped_tiles {
			//fmt.Printf("Move: %v - ", i)

			s := time.Now()
			new_board = *board
			new_board.placePieces(flipped_tiles[i], player.color)
			score := alpha_beta_minimax(new_board, player, new_board.eval_score(player), player.opp_color, math.MinInt64, math.MaxInt64, 0, curr_max_depth, false)
			if score > max_score {
				temp_best_moves = []int{i}
				max_score = score
			} else if score == max_score {
				temp_best_moves = append(temp_best_moves, i)
			}
			last_time = float64(time.Since(s).Seconds())
			//fmt.Printf("%v | %v\n", last_time, time.Since(start).Seconds())

			// After finishing search of a move up to current depth, it checks to see if it could finish searching rest of moves up to current depth in the remaining time
			if last_time*float64(len(flipped_tiles)-(i+1))+time.Since(start).Seconds() > (float64(max_sec) * .8) {
				complete = false
				break
			}
		}
		// If it completed search of current depth, plan to return the best moves it calculated unless the next depth was searched to completion
		if complete {
			last_completed_depth = curr_max_depth
			best_moves = temp_best_moves
		} else { //otherwise print that search was cutoff during current depth
			fmt.Printf("Timelimit caused cuttoff during search of depth %v\n", curr_max_depth)
			break
		}
	}
	return best_moves, last_completed_depth, time.Since(start).Seconds()
}

// Function that implements minimax search with alpha-beta pruning
func alpha_beta_minimax(board Board, player Player, score int, color string, alpha float64, beta float64, depth int, max_depth int, max_player bool) float64 {
	var new_board Board

	// If max depth of current iterative deepening limit is reached, return evaluation score at that depth
	if depth > max_depth {
		return float64(score)
	}

	opp_color := "W"
	if color == "W" {
		opp_color = "B"
	}

	legal_moves, flipped_tiles := board.find_legal_moves(color)

	// If player has no legal moves
	if len(legal_moves) == 0 {
		winner, game_over := board.check_game_over()
		if game_over { // return -infinity or +infinity or 0 depending on if player lost or won or tied at this possible board state
			if winner == player.color {
				return math.MaxInt64
			} else if winner == player.opp_color {
				return math.MinInt64
			} else if winner == "Draw" {
				return 0
			}
		} else { //continue minmax search and just skip player if game is not over
			legal_moves, flipped_tiles = board.find_legal_moves(opp_color)
			max_player = !max_player
			temp := color
			color = opp_color
			opp_color = temp
		}
	}

	// Max player
	if max_player == true {
		max_eval := float64(math.MinInt64) // set max_eval to -inf
		for i, _ := range legal_moves {    // loop through possible moves
			new_board = board
			new_board.placePieces(flipped_tiles[i], color)
			eval := alpha_beta_minimax(new_board, player, new_board.eval_score(player), opp_color, alpha, beta, depth+1, max_depth, false)
			max_eval = math.Max(max_eval, eval)
			alpha = math.Max(alpha, max_eval)
			if beta <= alpha {
				break
			}
		}
		return max_eval

		// Min player
	} else if max_player == false {
		min_eval := float64(math.MaxInt64) // set min_eval to +inf
		for i, _ := range legal_moves {    // loop through possible moves
			new_board = board
			new_board.placePieces(flipped_tiles[i], color)
			eval := alpha_beta_minimax(new_board, player, new_board.eval_score(player), opp_color, alpha, beta, depth+1, max_depth, true)
			min_eval = math.Min(min_eval, eval)
			beta = math.Min(beta, min_eval)
			if beta <= alpha {
				break
			}
		}
		return min_eval
	}
	return 0
}

// Function that takes in input and checks if it's valid depending on certain conditions
func get_user_input(prompt string, typ string, lower_bound int, upper_bound int) string {
	var inp string
	for {
		fmt.Printf("%v", prompt)
		fmt.Scanln(&inp)
		switch typ {

		case "comp": // computer
			if inp == "y" {
				return "true"
			} else {
				return "false"
			}

		case "file": // checks if file is valid or if user inputted nothing on purpose
			if _, err := os.Stat(inp); err == nil {
				return inp
			} else if inp != "" {
				fmt.Printf("ERROR: File '%v' does not exist\n", inp)
			} else if inp == "" {
				return "default"
			}

		case "int": // checks if inputted string is valid integer within specified bounds
			if intgr, err := strconv.Atoi(inp); err == nil {
				if intgr >= lower_bound && intgr <= upper_bound {
					return inp
				} else {
					fmt.Printf("ERROR: '%v' is out of bounds\n", inp)
				}
			} else {
				fmt.Printf("ERROR: '%v' is not a valid integer\n", inp)
			}
		default:
		}
	}
}

// Main Function
func main() {

	// Static heuristic of board spaces
	// heur := [][]int{{200, -30, 30, 10, 10, 30, -30, 200},
	// 	{-30, -50, -10, -10, -10, -10, -50, -30},
	// 	{30, -10, 20, 5, 5, 20, -10, 30},
	// 	{10, -10, 5, 5, 5, 5, -10, 10},
	// 	{10, -10, 5, 5, 5, 5, -10, 10},
	// 	{30, -10, 20, 5, 5, 20, -10, 30},
	// 	{-30, -50, -10, -10, -10, -10, -50, -30},
	// 	{200, -30, 30, 10, 10, 30, -30, 200}}

	// Get user inputs for board preset file, which players are computers, and time limit for computer
	path := get_user_input("Input file name of premade game state or leave blank for new game: ", "file", 0, 0)
	p1_comp, _ := strconv.ParseBool(get_user_input("Is player 1 (Black) a computer? (y/n): ", "comp", 0, 0))
	p2_comp, _ := strconv.ParseBool(get_user_input("Is player 2 (White) a computer? (y/n): ", "comp", 0, 0))
	time_lim, _ := strconv.Atoi(get_user_input("Input time limit for turn: ", "int", 0, int(math.MaxInt64)))

	// Create game board and initialize board to default or to preset depending on user input
	board := &Board{}
	board.initBoard(path)

	// Create players with user specifications
	player1 := &Player{name: "Player 1 (Black)", computer: p1_comp, color: "B", opp_color: "W", eval_type: "dynamic", heur: [][]int{}}
	player2 := &Player{name: "Player 2 (White)", computer: p2_comp, color: "W", opp_color: "B", eval_type: "dynamic", heur: [][]int{}}

	// Create and start game
	game := &Game{
		board:    board,
		player1:  player1,
		player2:  player2,
		time_lim: time_lim,
	}

	game.play()

}
