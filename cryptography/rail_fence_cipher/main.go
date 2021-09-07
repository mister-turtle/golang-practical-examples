package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {

	argPlainText := flag.String("p", "", "Plaintext message to encode")
	argCipherText := flag.String("c", "", "Cipher text to decode")
	argRails := flag.Int("r", 2, "Number of rails to use")
	flag.Parse()

	if (*argPlainText == "" && *argCipherText == "") || (*argPlainText != "" && *argCipherText != "") {
		fmt.Printf("Please provide either a cipher, or a plaintext\n")
		flag.Usage()
		os.Exit(1)
	}

	if *argPlainText != "" {
		fmt.Printf("\"%s\"\n", Encrypt(*argPlainText, *argRails))
		os.Exit(0)
	}

	fmt.Printf("\"%s\"\n", Decrypt(*argCipherText, *argRails))
}

// Encrypt implements a basic rail fence cipher encryption
func Encrypt(plaintext string, rails int) string {

	reverse := false
	railNumber := 0

	if len(plaintext) == 0 {
		return ""
	}

	rail := make(map[int][]rune)
	for _, char := range plaintext {
		rail[railNumber] = append(rail[railNumber], char)

		if railNumber == rails-1 {
			reverse = true
		} else if railNumber == 0 {
			reverse = false
		}

		if reverse {
			railNumber--
			continue
		}
		railNumber++
	}

	var output string
	for i := 0; i < rails; i++ {
		output += string(rail[i])
	}
	return output
}

// NaiveDecrypt implements a basic decryption by ranging over the number of rails and calculating the offset
// of the decrypted position by stepping up and down the rails and incrementing the offset counter accordingly.
// When the target rail is reached again, the offset is used, and position indicators reset.
func NaiveDecrypt(ciphertext string, rails int) string {

	var output = make([]rune, len(ciphertext))

	var position = 0

	for targetRail := 0; targetRail < rails; targetRail++ {

		railNumber := targetRail
		reverse := false
		targetPosition := targetRail

		for targetPosition != len(output) {

			if railNumber == targetRail {
				output[targetPosition] = rune(ciphertext[position])
				position++
			}

			if railNumber == rails-1 {
				reverse = true
			} else if railNumber == 0 {
				reverse = false
			}

			if reverse {
				railNumber--
				targetPosition++
				continue
			}
			railNumber++
			targetPosition++
		}
	}
	return string(output)
}

// Decrypt implements the decryption process in a slightly smarter way. Each position in the ciphertext is checked and an offset calculated
// based on the number of rails, the current direction of travel, and the last offset used for that rail. If the calculated offset is beyond the length of the origina ciphertext
// it is assumed the character was on the next rail up. This allows us to avoid manually walking the "zig-zag for each position"
func Decrypt(ciphertext string, rails int) string {

	var output = make([]rune, len(ciphertext))

	directionDown := true
	currentRail := 0
	position := 0
	railOffsetPosition := 0
	previousTargetOffset := 0
	for position != len(ciphertext) {

		var newOffset int
		if directionDown {

			// calculate the next offset in the output slice for the current rail character when moving down
			newOffset = previousTargetOffset + (rails-(currentRail+1))*2

			// we should calculate the next offset by travelling upwards over the rails IF:
			// if we're not on the top rail AND
			// we're on an odd offset AND
			// we're past the point where rails should zig-zag OR
			// we're on the last rail
			if (currentRail != 0 && (railOffsetPosition&1) == 1) && position >= rails-1 || currentRail == rails-1 {
				directionDown = false
			}
		} else {

			// calculate the next offset in the output slice for the current rail character when moving up
			newOffset = previousTargetOffset + (currentRail * 2)

			// we should switch to searching down the rails, as long as we're not on the bottom rail
			if currentRail != rails-1 {
				directionDown = true
			}
		}

		// if this is the first item in the current rail, the rail number will dictate its offset from the start of the slice
		if railOffsetPosition == 0 {
			newOffset = currentRail
		}

		// if the calculated position is beyond the length of the input text, the character must have been for a different rail
		// reset the direction, and offset into the rail, and increment the current rail
		if newOffset > len(ciphertext)-1 {
			directionDown = true
			currentRail++
			railOffsetPosition = 0
			previousTargetOffset = currentRail
			continue
		}

		// add the character in the ciphertext to the calculated output offset, increment the position and railOffset
		previousTargetOffset = newOffset
		output[newOffset] = rune(ciphertext[position])
		position++
		railOffsetPosition++
	}

	return string(output)
}
