package main

import bootstrap "golearn-structured/bootsrap"

func main() {
	// Titik masuk proses HTTP app: semua wiring dilakukan di bootstrap.Run().
	bootstrap.Run()
}
