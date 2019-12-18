package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

const matrixPathA string = "matrixA.txt"
const matrixPathB string = "matrixB.txt"
const matrixPathAns string = "matrixAns.txt"

func main() {
	A := readMatrix(matrixPathA)
	B := readMatrix(matrixPathB)
	//C := readMatrix(matrixPathAns)
	var matrix Matrix

	for i := 1; i <= 100; i++ {
		start := float64(time.Now().UnixNano())
		for j := 0; j < 10; j++ {
			matrix = multiplyMatrix(A, B, i)
		}
		fmt.Println((float64(time.Now().UnixNano()) - start) / float64(time.Millisecond) / 10)
	}
	//printMatrix(matrix)
	print(matrix.n)
	/*if checkMatrix(matrix, C) {
		println("Correct answer")
	}*/
}

func printMatrix(matrix Matrix) {
	for i := 0; i < matrix.n; i++ {
		for j := 0; j < matrix.m; j++ {
			fmt.Printf("%15.2f ", matrix.matrix[i][j])
		}
		fmt.Println()
	}
}

func checkMatrix(myAns, realAns Matrix) bool {
	for i := 0; i < myAns.n; i++ {
		for j := 0; j < myAns.m; j++ {
			if math.Abs(myAns.matrix[i][j]-realAns.matrix[i][j]) >= 0.00001 {
				return false
			}
		}
	}

	return true
}

func readMatrix(path string) Matrix {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	scanner.Scan()
	n64, _ := strconv.ParseInt(scanner.Text(), 10, 32)
	scanner.Scan()
	m64, _ := strconv.ParseInt(scanner.Text(), 10, 32)

	n := int(n64)
	m := int(m64)
	matrix := make([][]float64, n)

	for i := 0; i < n; i++ {
		matrix[i] = make([]float64, m)
		for j := 0; j < m; j++ {
			scanner.Scan()
			value, _ := strconv.ParseFloat(scanner.Text(), 64)
			matrix[i][j] = value
		}
	}

	return Matrix{matrix: matrix, n: n, m: m}
}
