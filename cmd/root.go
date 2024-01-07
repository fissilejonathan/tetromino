package cmd

import (
	"fmt"
	"os"

	"github.com/eiannone/keyboard"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var rootCmd = &cobra.Command{
	Use:   "tetro",
	Short: "cli tetromino",
	Run: func(cmd *cobra.Command, args []string) {
		input := make(chan keyboard.Key)
		done := make(chan struct{})

		err := keyboard.Open()
		if err != nil {
			panic(err)
		}
		defer func() {
			keyboard.Close()
		}()

		// handle user input
		go func() {
			defer close(input)
			for {
				_, key, err := keyboard.GetSingleKey()
				if err != nil {
					fmt.Println("Error reading input:", err)
					return
				}

				if key == keyboard.KeyEsc {
					close(done)
				}

				if key == keyboard.KeyArrowDown || key == keyboard.KeyArrowUp || key == keyboard.KeyArrowLeft || key == keyboard.KeyArrowRight {
					input <- key
				}
			}
		}()

		for {
			select {
			case userInput := <-input:
				fmt.Printf("Received: %c\n", userInput)
			case <-done:
				return
			default:
				width, height, err := terminalSize()
				if err != nil {
					fmt.Println("Error:", err)
					return
				}

				fmt.Printf("Terminal size: %d columns x %d rows\n", width, height)
			}
		}
	},
}

func terminalSize() (int, int, error) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 0, 0, err
	}
	return width, height, nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
