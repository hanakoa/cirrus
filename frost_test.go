package frost

import (
	"testing"
	"log"
)

func TestFrost(t *testing.T) {
	log.Println("Seeding available node IDs...")
	for i := 0; i < NumNodes; i++ {
		AvailableNodeIds <- i
	}
	log.Println("Finished seeding...")
	f := &FrostServer{}
	f.Run()
}