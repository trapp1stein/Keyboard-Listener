# Keyboard-Listener

A cross-platform keyboard listener library written in Go.

This project provides low-level keyboard input handling with support for Windows and Linux systems.  
You can easily open the keyboard listener, read key events, and close it properly.

---

## Project Structure

📦keylistener
 ┣ 📂keyboard
 ┃ ┣ 📜keyboard.go
 ┃ ┣ 📜keyboard_common.go
 ┃ ┣ 📜keyboard_windows.go
 ┃ ┣ 📜syscalls.go
 ┃ ┣ 📜syscalls_linux.go
 ┃ ┗ 📜terminfo.go
 ┣ 📜go.mod
 ┣ 📜go.sum
 ┣ 📜README.md
 ┗ 📜test.go

 ---

## Prerequisites

- Go 1.24 or newer installed  
- Compatible OS: Windows, Linux

---

## Installation

Clone this repository and navigate into it:

```bash
git clone https://github.com/trapp1stein/Keyboard-Listener.git
cd Keyboard-Listener

```


## How to Use
The keyboard package provides functions to open the keyboard listener, get key events, and close the listener.

Running the Example (test.go)
Run the example program to test the keyboard listener:

```bash
go run test.go

```

What to expect:

The program will start and print Keyboard listener test started.

It will wait for you to press a key.

Once you press a key, it will display the rune and key code.

Then it will exit cleanly.

Example Code Snippet
Here is a minimal usage example from test.go:

```go

package main

import (
	"fmt"
	"Keyboard-Listener/keyboard"
)

func main() {
	err := keyboard.Open()
	if err != nil {
		fmt.Println("Error opening keyboard:", err)
		return
	}
	defer keyboard.Close()

	fmt.Println("Press any key...")

	ch, key, err := keyboard.GetKey()
	if err != nil {
		fmt.Println("Error reading key:", err)
		return
	}

	fmt.Printf("You pressed rune: %q, key code: %v\n", ch, key)
}
```

## Contributing
Contributions, issues, and feature requests are welcome!
Feel free to open a pull request or issue.

## Thank you for using Keyboard-Listener!
Happy coding! 🚀