package main

import (
	"fmt"
	"Keyboard-Listener/keyboard"
)

func main() {
	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Println("Lütfen bir tuşa basınız; (Çıkmak için ESC)")
	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
		}
		fmt.Printf("Key: rune %q, anahtar: %X\r\n", event.Rune, event.Key)
		if event.Key == keyboard.KeyEsc {
			break
		}
	}
}