package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	tr "./ytapi"
)

func main() {
	k := os.Getenv("Y_TRANSLATE_KEY")

	client := tr.NewClient(k)

	sentence := "我爱你"

	ui, _ := client.Detect(sentence)
	log.Println("Detect language code is ", ui)
	allPossibleWay, _ := client.GetList(ui)
	log.Println("there are", len(allPossibleWay), "kind of ways to say", sentence)

	var wg sync.WaitGroup
	wg.Add(len(allPossibleWay))
	for k, v := range allPossibleWay {
		go func(k, v string) {
			translated, _ := client.Trans(sentence, k)
			chinese, _ := client.Trans(v, "zh")
			fmt.Println(chinese, "(", v, "):", translated)
			wg.Done()
		}(k, v)
	}
	wg.Wait()
}
