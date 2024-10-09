package	main

import (
	"fmt"
	"time"
	"os"
	"os/exec"
	"github.com/eiannone/keyboard"
)

type Player struct{
	Position int 
	Speed int
}

type Ball struct{
	Position [2]int
	Velocity [2]int
}

type Game struct{
	Player1 Player
	Player2 Player
	Ball Ball
	Width int
	Height int
}

func main(){
	game := Game{
		Player1 : Player{Position : game.Height / 2, Speed: 5},
		Player2: Player{Position : game.Height / 2, Speed: 5},
		Ball: Ball{Position: [2]int{game.Width / 2, game.Height / 2}, Velocity: [2]int{-2, 2}}
		Width: 80,
		Height: 80,
	}
	for {
		UpdatePlayerPosition(&game.Player1, 'w', 's')
		UpdatePlayerPosition(&game.Player2, 'i', 'k')
		updateBallposition(&game.Ball)
		diplayGame(game)
		time.Sleep(50 * time.Millisecond)
	}
}

func UpdatePlayerPosition(player *Player, upKey, downKey rune){
	if isKeypressed(upKey){
		player.Position = max(player.Position - player.Speed, 0)
	}
	else if isKeypressed(downKey){
		player.Position = min(player.Position + player.Speed, game.Height - 1)
	}

}

func isKeypressed(key rune) bool{
	return false
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}