package binlookup

import (
	"os"
	"os/exec"
	"sync"

	"log"
)

// ----- from go-flutter hover----
type BinLookup struct {
	Name                string
	InstallInstructions string
	fullPath            string
	once                sync.Once
}

func (b *BinLookup) FullPath() string {
	b.once.Do(func() {
		var err error
		b.fullPath, err = exec.LookPath(b.Name)
		if err != nil {
			log.Printf("Failed to lookup `%s` executable: %s. %s", b.Name, err, b.InstallInstructions)
			os.Exit(1)
		}
	})
	return b.fullPath
}
