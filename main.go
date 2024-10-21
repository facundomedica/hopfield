package main

import (
	"fmt"
	"math/rand"
)

// Size de 100 (bitmaps de 10x20)
const size = 100

type HopfieldNetwork struct {
	weights [size][size]int
	state   [size]int
}

// Train the network using Hebbian learning rule
func (h *HopfieldNetwork) Train(patterns [][size]int) {
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if i != j {
				sum := 0
				for _, pattern := range patterns {
					sum += pattern[i] * pattern[j]
				}
				h.weights[i][j] = sum
			} else {
				h.weights[i][j] = 0
			}
		}
	}
}

func (h *HopfieldNetwork) SetState(input [size]int) {
	h.state = input
}

func (h *HopfieldNetwork) Update() {
	for i := 0; i < size; i++ {
		sum := 0
		for j := 0; j < size; j++ {
			sum += h.weights[i][j] * h.state[j]
		}
		if sum >= 0 {
			h.state[i] = 1
		} else {
			h.state[i] = -1
		}
	}
}

// Helper function to display 10x10 bitmap
func displayState(state [size]int) {
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if state[i*10+j] == 1 {
				fmt.Print("█")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
	// Define some training patterns for 'A', 'V', and 'Z'
	patternA := [size]int{
		-1, -1, 1, 1, 1, 1, 1, -1, -1, -1,
		-1, 1, -1, -1, -1, -1, -1, 1, -1, -1,
		1, -1, -1, -1, -1, -1, -1, -1, 1, -1,
		1, 1, 1, 1, 1, 1, 1, -1, 1, -1,
		1, -1, -1, -1, -1, -1, -1, -1, 1, -1,
		1, -1, -1, -1, -1, -1, -1, -1, 1, -1,
		1, -1, -1, -1, -1, -1, -1, -1, 1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	}

	patternV := [size]int{
		1, -1, -1, -1, -1, -1, -1, -1, -1, 1,
		1, -1, -1, -1, -1, -1, -1, -1, -1, 1,
		-1, 1, -1, -1, -1, -1, -1, -1, 1, -1,
		-1, 1, -1, -1, -1, -1, -1, -1, 1, -1,
		-1, -1, 1, -1, -1, -1, -1, 1, -1, -1,
		-1, -1, 1, -1, -1, -1, -1, 1, -1, -1,
		-1, -1, -1, 1, -1, -1, 1, -1, -1, -1,
		-1, -1, -1, 1, -1, -1, 1, -1, -1, -1,
		-1, -1, -1, -1, 1, 1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	}

	patternZ := [size]int{
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, 1,
		-1, -1, -1, -1, -1, -1, -1, -1, 1, -1,
		-1, -1, -1, -1, -1, -1, -1, 1, -1, -1,
		-1, -1, -1, -1, -1, -1, 1, -1, -1, -1,
		-1, -1, -1, -1, -1, 1, -1, -1, -1, -1,
		-1, -1, -1, -1, 1, -1, -1, -1, -1, -1,
		-1, -1, -1, 1, -1, -1, -1, -1, -1, -1,
		-1, -1, 1, -1, -1, -1, -1, -1, -1, -1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	}

	network := HopfieldNetwork{}

	network.Train([][size]int{patternA, patternV, patternZ})

	noisyInput := patternV
	for i := 0; i < 30; i++ { // agregamos ruido
		index := rand.Intn(size)
		noisyInput[index] = -noisyInput[index]
	}

	fmt.Println("Input con ruido (10 bits cambiados al azar):")
	displayState(noisyInput)

	network.SetState(noisyInput)

	for i := 0; i < 100; i++ {
		network.Update()
	}

	// Display the recovered pattern
	fmt.Println("Patrón recuperado:")
	displayState(network.state)

	// Creo 100 inputs aleatorios, y chequeo que el output sea el esperado

	type input struct {
		input         [size]int
		expected      [size]int
		amountOfNoise int
	}

	inputs := []input{}
	for i := 0; i < 100; i++ {
		// selecciona A o V de manera aleatoria
		var pattern [size]int
		var expected [size]int
		if rand.Intn(2) == 0 {
			pattern = patternA
			expected = patternA
		} else {
			pattern = patternV
			expected = patternV
		}

		// agrega ruido, con hasta 30 bits cambiados
		noisyInput := pattern
		amountOfNoise := rand.Intn(40)
		for i := 0; i < amountOfNoise; i++ {
			index := rand.Intn(size)
			noisyInput[index] = -noisyInput[index]
		}

		inputs = append(inputs, input{noisyInput, expected, amountOfNoise})
	}

	failures := 0
	for i, input := range inputs {
		network.SetState(input.input)
		for i := 0; i < 100; i++ {
			network.Update()
		}
		if input.expected != network.state {
			failures++
			fmt.Printf("Test %d failed. Expected:\n", i)
			displayState(input.expected)
			fmt.Println("Got:")
			displayState(network.state)
			fmt.Println("Amount of noise:", input.amountOfNoise)
			fmt.Println("with input:")
			displayState(input.input)
			fmt.Println()
		}
	}

	fmt.Println("Test aleatorizado, cantidad de fallas:", failures, "de", len(inputs))

}
