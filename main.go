package main

import (
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// TODO move this to a new package?
//go:embed templates
var templateFS embed.FS

//go:embed css
var cssFS embed.FS

type Results struct {
	AverageCallsUntilWin int
	WinsForEachShape     map[int]int
}

func main() {

	port := flag.String("port", "8888", "default http port")
	flag.Parse()

	// CSS is also static, but separated out so it works regardless of other static files
	http.Handle("/css/", http.FileServer(http.FS(cssFS)))
	// everything else is the main template
	http.HandleFunc("/", serveTemplate)
	http.HandleFunc("/board", parseBoard)
	log.Println("Ready at http://localhost:" + *port)
	log.Println(http.ListenAndServe(":"+*port, nil))

}
func parseBoard(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	b, err := io.ReadAll(req.Body)
	if err != nil {
		if err != nil {
			_, err := fmt.Fprintf(w, "{\"error\":\"%s\"}", err)
			if err != nil {
				log.Printf("error handling error: %s", err)
			}
			return
		}
	}
	//log.Println(string(b))

	// FIXME on error, send back a 4xx

	var data []map[string]string
	err = json.Unmarshal(b, &data)
	if err != nil {
		_, err = fmt.Fprintf(w, "{\"error\":\"%s\"}", err)
		if err != nil {
			log.Printf("error handling error: %s", err)
		}
		return
	}
	g, err := createGameFromDataMap(data)
	if err != nil {
		_, err = fmt.Fprintf(w, "{\"error\":\"%s\"}", err)
		if err != nil {
			log.Printf("error handling error: %s", err)
		}
		return
	}
	wins, shapes := computeAveragePlaysUntilWin(g, 100)

	results := Results{AverageCallsUntilWin: wins, WinsForEachShape: shapes}
	log.Println(results)

	// send some JSON back with the details
	jsonBytes, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		_, err = fmt.Fprintf(w, "{\"error\":\"%s\"}", err)
		if err != nil {
			log.Printf("error handling error: %s", err)
		}
		return
	}
	_, err = fmt.Fprintf(w, string(jsonBytes))
	if err != nil {
		log.Printf("error handling error: %s", err)
	}

}

// FIXME clearly I have a bug around the free square
// if I specify a square as "12" a shape that needs it never wins
// i think what i need to do is just ignore 12 on input- don't mark it required or called
// don't mark it anything special
// In the HTML on start, set that square to 0 and disable the input

func createGameFromDataMap(userBoard []map[string]string) (*Game, error) {
	g := newGame()
	g.BingoNumToIndex = make(map[int]int)
	// this is an array of maps (pairs of name=x, value=y)
	// each map has name=square_id_{squareid} or name= square_needed_{squareId}_{shapeId}
	// and value = Square.Number or "on" if that square is needed for the shapeId
	// "square_needed_1_#" value = "on"
	// "square_id_#"   value = "#"

	for _, m := range userBoard {
		n := m["name"]
		v := m["value"]
		if strings.HasPrefix(n, "square_id") {
			var squareIdx int
			_, err := fmt.Sscanf(n, "square_id_%d", &squareIdx)
			if err != nil {
				return nil, errors.New("cannot parse board: squarenum")
			}
			squareNum, err := strconv.Atoi(v)
			if err != nil {
				return nil, errors.New("cannot parse board: squarenum")
			}

			// Create a map from bingo num to square index
			g.BingoNumToIndex[squareNum] = squareIdx

			// can't modify it in place, I don't remember why, or is that just structs?
			sq := g.Squares[squareIdx]
			sq.Number = squareNum
			g.Squares[squareIdx] = sq

		} else if strings.HasPrefix(n, "square_needed") {
			var squareIdx int
			var shapeIdx int
			_, err := fmt.Sscanf(n, "square_needed_%d_%d", &squareIdx, &shapeIdx)
			if err != nil {
				return nil, errors.New("cannot parse board: square_needed")
			}
			// don't actually care about the value because if it was false it wouldn't be here
			sq := g.Squares[squareIdx]
			sq.Needed[shapeIdx] = true
			g.Squares[squareIdx] = sq
			g.Shapes[shapeIdx] = true // record that we need this shape
		} else if strings.HasPrefix(n, "called") {
			// the numbers we are getting in are *called* numbers, not board indices
			// so we need to find the square with that number and check it
			var calledNum int
			_, err := fmt.Sscanf(v, "called_%d", &calledNum)
			if err != nil {
				return nil, errors.New("cannot parse board: called: %" + n)
			}
			g.Called[calledNum] = true
			g.PreCalled[calledNum] = true
			for i, square := range g.Squares {
				if square.Number == calledNum {
					square.Called = true
					square.PreCalled = true
					g.Squares[i] = square
					break
				}
			}
		} else {
			log.Printf("don't know what to do with %s", n)
		}

	}
	// remove any mapping to the free square
	delete(g.BingoNumToIndex, 12)
	return g, nil
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("b").Funcs(template.FuncMap{"N": N}).ParseFS(templateFS, "templates/*.html")
	if err != nil {
		log.Println(err)
	}

	var rows [][]Square
	rows = make([][]Square, 5)
	for row := 0; row < 5; row++ {
		rows[row] = make([]Square, 5)
		for col := 0; col < 5; col++ {
			rows[row][col] = newSquare((row)*5 + col)
		}
	}
	rows[0][0].Number = 0
	rows[2][2].Number = 0

	err = tmpl.ExecuteTemplate(w, "board", struct {
		Rows [][]Square
	}{
		Rows: rows,
	})
	if err != nil {
		log.Println(err)
	}
}

func N(start, end int) (stream chan int) {
	stream = make(chan int)
	go func() {
		for i := start; i <= end; i++ {
			stream <- i
		}
		close(stream)
	}()
	return
}

// returns map of shapeIds to wins
func computeAveragePlaysUntilWin(g *Game, games int) (avgCalls int, winningShapes map[int]int) {
	var totalcalls int
	winningShapes = make(map[int]int)
	for i, _ := range g.Shapes {
		winningShapes[i] = 0
	}
	for i := 0; i < games; i++ {
		g.reset()
		calls, winners, err := playUntilWin(g)
		if err != nil {
			log.Printf("game failed: %v", err)
			continue
		}
		totalcalls += calls
		for _, winner := range winners {
			winningShapes[winner]++
		}
	}
	//avgCalls = float64(totalcalls / games)
	avgCalls = totalcalls / games
	return avgCalls, winningShapes
}

func playUntilWin(g *Game) (numCalls int, wonShapes []int, err error) {

	var won bool
	wonShapes = []int{}
	callednums := []int{}

	// check the degnerate case where we have already won
	for s, _ := range g.Shapes {
		if g.winner(s) {
			wonShapes = append(wonShapes, s)
			won = true
		}
	}

	for !won {
		sq, err := g.callRandomSquare()
		if err != nil {
			return 0, []int{}, errors.New("out of numbers without a winner")
		}
		callednums = append(callednums, sq)
		numCalls++
		won = g.playSquare(sq)
		// which shape won?
		for s, _ := range g.Shapes {
			if g.winner(s) {
				wonShapes = append(wonShapes, s)
			}
		}
	}

	//log.Printf("shapes %v won after %d calls: %v\n", wonShapes, numCalls, callednums)
	return
}
