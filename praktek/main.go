package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Player struct {
	Dices               []int
	Score               int
	CurrentDiceOneCount int
	CurrentDiceSixCount int
}

func main() {
	log.SetFlags(log.Llongfile)

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter number of players: ")
	playerCount, err := ReadNumberInput(reader)

	if err != nil {
		log.Fatalf("Error reading player count: %v", err)
	}

	if playerCount <= 1 {
		fmt.Printf("Number of players must be greater than 1\n")
	}

	fmt.Print("Enter number of dices: ")
	diceCount, err := ReadNumberInput(reader)

	if err != nil {
		log.Fatalf("Error reading dice count: %v", err)
	}

	if diceCount <= 0 {
		fmt.Printf("Number of dices must be greater than 0\n")
	}

	var players []Player

	for i := 0; i < playerCount; i++ {
		players = append(players, Player{
			Score:               0,
			CurrentDiceOneCount: 0,
		})

		for j := 0; j < diceCount; j++ {
			players[i].Dices = append(players[i].Dices, 0)
		}
	}

	fmt.Printf("\nPlayers: %v, Dices for each player: %v\n", playerCount, diceCount)

	fmt.Println("==== Game Start ====")
	for {
		participatingPlayers := 0
		for i := 0; i < len(players); i++ {
			if len(players[i].Dices) > 0 {
				participatingPlayers++
			}
		}

		if participatingPlayers <= 1 {
			break
		}

		fmt.Println("=====================")
		fmt.Printf("Press enter to roll dices...")
		_, _ = reader.ReadString('\n')
		fmt.Println("Roll Dices...")
		for i := 0; i < len(players); i++ {
			players[i].CurrentDiceOneCount = 0
			players[i].CurrentDiceSixCount = 0
			for j := 0; j < len(players[i].Dices); j++ {
				rand.Seed(time.Now().UnixNano())
				players[i].Dices[j] = rand.Intn(6) + 1

				if players[i].Dices[j] == 6 {
					players[i].CurrentDiceSixCount += 1
				} else if players[i].Dices[j] == 1 {
					players[i].CurrentDiceOneCount += 1
				}
			}

			fmt.Printf("Player %d (%d) : %v\n", i+1, players[i].Score, players[i].Dices)

			//dice with number 6 evaluation
			if players[i].CurrentDiceSixCount > 0 {
				players[i].Score += players[i].CurrentDiceSixCount
				players[i].Dices = players[i].Dices[:len(players[i].Dices)-players[i].CurrentDiceSixCount]
			}
		}

		//dice with number 1 evaluation with the consideration that players sit in a circle
		var diceMoveMap = map[int]int{}
		for i := 0; i < len(players); i++ {
			if players[i].CurrentDiceOneCount > 0 {
				var playerNumber, previousPlayerIndex, nextPlayerIndex int = 0, 0, 0

				//when the current player is the first player, previous player is the last player
				if i == 0 {
					previousPlayerIndex = len(players) - 1
					nextPlayerIndex = i + 1
				} else if i == len(players)-1 {
					//when the current player is the last player, next player is the first player
					previousPlayerIndex = i - 1
					nextPlayerIndex = 0
				} else {
					previousPlayerIndex = i - 1
					nextPlayerIndex = i + 1
				}

				//make player choose on which player to give their dice if both players beside them have dice(s)
				if len(players[previousPlayerIndex].Dices) > 0 && len(players[nextPlayerIndex].Dices) > 0 {
					for {
						fmt.Printf("Player %d have %d dice(s) with 1, who do you want to give it to (player %d or player %d)? ", i+1, players[i].CurrentDiceOneCount, previousPlayerIndex+1, nextPlayerIndex+1)
						playerNumber, err = ReadNumberInput(reader)
						if err != nil {
							log.Printf("\nError reading player count: %v\n", err)
							continue
						}

						if playerNumber != previousPlayerIndex+1 && playerNumber != nextPlayerIndex+1 {
							fmt.Printf("\nInvalid player number\n")
							continue
						}
						break
					}
				} else if len(players[nextPlayerIndex].Dices) > 0 {
					//automatically choose next player to give their dice with number 1 if only that player has dice(s)
					playerNumber = nextPlayerIndex + 1
				} else if len(players[previousPlayerIndex].Dices) > 0 {
					//automatically choose previous player to give their dice with number 1 if only that player has dice(s)
					playerNumber = previousPlayerIndex + 1
				}

				if playerNumber > 0 {
					diceMoveMap[i] = playerNumber - 1
					players[i].Dices = players[i].Dices[:len(players[i].Dices)-players[i].CurrentDiceOneCount]
				}
			}
		}

		for from, to := range diceMoveMap {
			fmt.Printf("Player %d give %d dices with number 1 to player %d\n", from+1, players[from].CurrentDiceOneCount, to+1)
			for j := 0; j < players[from].CurrentDiceOneCount; j++ {
				players[to].Dices = append(players[to].Dices, 0)
			}
		}
		fmt.Println("=====================")
	}

	fmt.Println("==== Game End ====")
	fmt.Printf("\n")

	highestScore := 0
	scoreMap := map[int][]int{}
	for i := 0; i < len(players); i++ {
		fmt.Printf("Player %d score : %d\n", i+1, players[i].Score)
		scoreMap[players[i].Score] = append(scoreMap[players[i].Score], i+1)
		if highestScore < players[i].Score {
			highestScore = players[i].Score
		}
	}

	fmt.Printf("\n==== Winner ====\n")
	if len(scoreMap[highestScore]) > 1 {
		fmt.Printf("Players %v win with score %d\n", scoreMap[highestScore], highestScore)
	} else {
		fmt.Printf("Player %d win with score %d\n", scoreMap[highestScore][0], highestScore)
	}
}

func ReadNumberInput(reader *bufio.Reader) (int, error) {
	input, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	number, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return 0, err
	}

	return number, nil
}
