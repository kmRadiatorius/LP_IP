package main

// Matrix for storing matrix
type Matrix struct {
	matrix [][]float64
	n      int
	m      int
}

type resultSet struct {
	values []float64
	row    int
}

type repository struct {
	A          Matrix
	B          Matrix
	rowChan    chan int
	resultChan chan resultSet
	doneChan   chan bool
}

// Inicializuoja reikalingus kanalus ir matricas
func initRepos(A, B Matrix) repository {
	var repos repository
	repos.A = A
	repos.B = B
	repos.rowChan = make(chan int)
	repos.resultChan = make(chan resultSet)
	repos.doneChan = make(chan bool)

	return repos
}

// Paima matricas ir gijų skaičių kaip parametrus
// Paleidžia duotą gijų skaičių matricų daugybai
// Grąžina duotų matricų sandaugą
func multiplyMatrix(A, B Matrix, threads int) Matrix {
	// init
	repos := initRepos(A, B)
	matrix := Matrix{matrix: make([][]float64, repos.A.n, repos.B.m), n: repos.A.n, m: repos.B.m}
	for i := 0; i < threads; i++ {
		go manageMatrixMultiplication(repos)
	}

	// pass rows to calculate
	var row int = 0
	for i := 0; i < repos.A.n; {
		select {
		case result := <-repos.resultChan:
			matrix.matrix[result.row] = result.values
			row++
			break
		case repos.rowChan <- i:
			i++
			break
		}
	}

	for i := row; i < repos.A.n; i++ {
		result := <-repos.resultChan
		matrix.matrix[result.row] = result.values
	}

	for i := 0; i < threads; i++ {
		repos.doneChan <- true
	}

	return matrix
}

// Paleidžiama multiplyMatrix funkcijos
// Po vieną ima ir daugina matricos eilutes
// Baigia darbą, multiplyMatrix funkcijai davus baigimo signalą
func manageMatrixMultiplication(repos repository) {
	done := false
	var row int
	for !done {
		select {
		case done = <-repos.doneChan:
			break
		case row = <-repos.rowChan:
			var result resultSet
			result.values = multiplyRow(repos.A, repos.B, row)
			result.row = row
			// printVector(result.values, repos.B.m)
			repos.resultChan <- result
			break
		}
	}
}

// Paima dvi matricas ir eilutę, kurią reikia gauti sudauginus
// Grąžina daugybos rezultatą
func multiplyRow(A Matrix, B Matrix, row int) []float64 {
	matrixRow := make([]float64, B.m)
	for i := 0; i < B.m; i++ {
		matrixRow[i] = multiplyVectors(A, B, row, i)
	}

	// printVector(matrixRow, B.m)
	return matrixRow
}

// Paima dvi matricas, pirmos matricos eilutę bei antros matricos stulpelį
// Grąžina gautų dviejų vektorių skaliarinę sandaugą
func multiplyVectors(A, B Matrix, row, col int) float64 {
	var result float64
	for i := 0; i < A.m; i++ {
		result += A.matrix[row][i] * B.matrix[i][col]
	}
	return result
}

func printVector(v []float64, n int) {
	println("Print Vector:")
	for i := 0; i < n; i++ {
		print(v[i], " ")
	}
	println()
}
