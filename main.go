package main

import (
	"fmt"
	"math/rand"
	//"os"
	//"bufio"
	"github.com/google/uuid"
  "time"
)

type GameState struct {
	Id        string
	PlayerOne Player
	PlayerTwo Player
	CurrRoll  int
	GameOver  bool
}

type Player struct {
	Id      string
	Score   int
	CurTurn bool
	Board   [9]int
	Winner  bool
}

func (GameState) StartTurn(game *GameState) {
	if game.PlayerOne.CurTurn {
		fmt.Printf("It is Player 1's turn. They have rolled a %d. Where would you like to place this dice?\n", game.CurrRoll)
	} else {
		fmt.Printf("It is Player 2's turn. They have rolled a %d. Where would you like to place this dice?\n", game.CurrRoll)
	}
}

func (GameState) EndTurn(position int, game *GameState) {
	if position < 0 || position > 8 {
		fmt.Println("Invalid position. You must choose an unoccupied index between 0 and 8. Try again.")
		return
	}
	column := position % 3

	if game.PlayerOne.CurTurn {
		game.PlayerOne.Board[position] = game.CurrRoll
	} else {
		game.PlayerTwo.Board[position] = game.CurrRoll
	}
	game.ScoreGame(game, column)
}

func (GameState) ScoreGame(game *GameState, column int) {
	for i := column; i < 9; i += 3 {
		if !game.PlayerOne.CurTurn && game.PlayerOne.Board[i] == game.CurrRoll {
			game.PlayerOne.Board[i] = 0
		} else if !game.PlayerTwo.CurTurn && game.PlayerTwo.Board[i] == game.CurrRoll {
			game.PlayerTwo.Board[i] = 0
		}
	}

	game.ScoreBoard(&game.PlayerOne)
	game.ScoreBoard(&game.PlayerTwo)

	game.GameOverCheck(game)

	if !game.GameOver {
		game.RollDice(game)
		game.PlayerOne.CurTurn = !game.PlayerOne.CurTurn
		game.PlayerTwo.CurTurn = !game.PlayerTwo.CurTurn
	}
}

func (GameState) ScoreBoard(player *Player) {
	board := player.Board
	newScore := 0
	for i := 0; i < 3; i++ {
		if one, two, three := board[i], board[i+3], board[i+6]; one == two && two == three {
			newScore = newScore + (one+two+three)*3
		} else if one == two || two == three || one == three {
			if one == two {
				newScore = newScore + three + (one*2)*2
			} else if two == three {
				newScore = newScore + one + (two*2)*2
			} else if one == three {
				newScore = newScore + two + (one*2)*2
			}
		} else {
			newScore = newScore + one + two + three
		}
	}
	player.Score = newScore
}

func (GameState) GameOverCheck(game *GameState) {
	playerOneDiceCount := 0
	playerTwoDiceCount := 0

	for i, di := range game.PlayerOne.Board {
		if di > 0 {
			playerOneDiceCount++
		}

		if game.PlayerTwo.Board[i] > 0 {
			playerTwoDiceCount++
		}
	}

	if playerOneDiceCount == 9 || playerTwoDiceCount == 9 {
		game.GameOver = true
	}
}

func (GameState) RollDice(game *GameState) {
	game.CurrRoll = rand.Intn(6) + 1
}

func main() {
	state := NewGame()
	RandomAutoPlay(state)
}

func NewGame() GameState {
	playerOne := NewPlayer(true)
	playerTwo := NewPlayer(false)

	state := GameState{
		Id:        uuid.New().String(),
		PlayerOne: playerOne,
		PlayerTwo: playerTwo,
		CurrRoll:  rand.Intn(6) + 1,
		GameOver:  false}

	return state
}

func NewPlayer(starter bool) Player {
  return Player{
		Id:      uuid.New().String(),
		Score:   0,
		CurTurn: starter,
		Board:   [9]int{},
		Winner:  false}
}

func RenderGame(state GameState) {
	playerOneBoard := RenderBoard(state.PlayerOne.Board)
	playerTwoBoard := RenderBoard(state.PlayerTwo.Board)

	fmt.Println("Player One Score: ", state.PlayerOne.Score)
	fmt.Println(playerOneBoard)
	fmt.Printf(playerTwoBoard)
	fmt.Println("Player Two Score: ", state.PlayerTwo.Score)
}

func RenderBoard(board [9]int) string {
	var boardstr string
	for i, v := range board {
		switch v {
		case 0:
			boardstr = boardstr + "[ ]"
		case 1:
			boardstr = boardstr + "[⚀]"
		case 2:
			boardstr = boardstr + "[⚁]"
		case 3:
			boardstr = boardstr + "[⚂]"
		case 4:
			boardstr = boardstr + "[⚃]"
		case 5:
			boardstr = boardstr + "[⚄]"
		case 6:
			boardstr = boardstr + "[⚅]"
		default:
			fmt.Println("Dice value out of range of D6")
		}

		if i == 2 || i == 5 || i == 8 {
			boardstr = boardstr + "\n"
		}
	}
	return boardstr
}

func RandomAutoPlay(state GameState) {
	for !state.GameOver {
		RenderGame(state)
		state.StartTurn(&state)
    time.Sleep(1 * time.Second)
		var move int
		var availableMoves []int
		if state.PlayerOne.CurTurn {
			for i, v := range state.PlayerOne.Board {
				if v == 0 {
					availableMoves = append(availableMoves, i)
				}
			}
		} else {
			for i, v := range state.PlayerTwo.Board {
				if v == 0 {
					availableMoves = append(availableMoves, i)
				}
			}
		}
		move = availableMoves[rand.Intn(len(availableMoves))]
		fmt.Println(len(availableMoves))
		state.EndTurn(move, &state)
	}
	
  RenderGame(state)

	if state.GameOver && state.PlayerOne.Score > state.PlayerTwo.Score {
		fmt.Printf("Game over, Player One wins %v to %v", state.PlayerOne.Score, state.PlayerTwo.Score)
	} else {
		fmt.Printf("Game over, Player Two wins %v to %v", state.PlayerTwo.Score, state.PlayerOne.Score)
	}
}
