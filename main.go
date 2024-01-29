package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var tpl = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Route Planner</title>
    <style>
        /* Add some basic styling for visualization */
        .grid {
            display: grid;
            grid-template-columns: repeat(10, 50px);
            grid-template-rows: repeat(10, 50px);
            gap: 1px;
        }
        .cell {
            width: 50px;
            height: 50px;
            border: 1px solid #ccc;
            display: flex;
            justify-content: center;
            align-items: center;
            font-size: 0.8em;
        }
        .obstacle {
            background-color: #333;
            color: #fff;
        }
    </style>
</head>
<body>
    <h1>Route Planner</h1>
    <form method="post">
        <label for="startX">Start X:</label>
        <input type="number" name="startX" required>
        <label for="startY">Start Y:</label>
        <input type="number" name="startY" required>
        <br>
        <label for="endX">End X:</label>
        <input type="number" name="endX" required>
        <label for="endY">End Y:</label>
        <input type="number" name="endY" required>
        <br>
        <button type="submit">Find Path</button>
    </form>
    <div class="grid">
        <!-- Visualize the grid here -->
        {{range .Grid}}
            {{range .}}
                <div class="cell {{if .Obstacle}}obstacle{{end}}">{{.Text}}</div>
            {{end}}
        {{end}}
    </div>
</body>
</html>
`))

type Cell struct {
	Text     string
	Obstacle bool
}

type PageVariables struct {
	Grid [][]Cell
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		startX, _ := strconv.Atoi(r.FormValue("startX"))
		startY, _ := strconv.Atoi(r.FormValue("startY"))
		endX, _ := strconv.Atoi(r.FormValue("endX"))
		endY, _ := strconv.Atoi(r.FormValue("endY"))

		grid := aStarAlgorithm(startX, startY, endX, endY)

		pageVariables := PageVariables{Grid: grid}
		err := tpl.Execute(w, pageVariables)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	err := tpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func aStarAlgorithm(startX, startY, endX, endY int) [][]Cell {
	grid := make([][]Cell, 10)
	for i := range grid {
		grid[i] = make([]Cell, 10)
		for j := range grid[i] {
			grid[i][j] = Cell{Text: "", Obstacle: false}
		}
	}

	grid[startX][startY].Text = "S"
	grid[endX][endY].Text = "E"

	path := findPath(startX, startY, endX, endY)
	for _, point := range path {
		grid[point[0]][point[1]].Text = "X"
	}

	return grid
}

func findPath(startX, startY, endX, endY int) [][]int {
	return [][]int{
		{startX, startY},
		{startX + 1, startY},
		{startX + 2, startY},
		{endX - 2, endY},
		{endX - 1, endY},
		{endX, endY},
	}
}
