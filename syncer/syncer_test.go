package syncer

import (
	"fmt"
	"os"
	"testing"
)

type Test struct {
	Pc   uintptr
	File string
	Line int
	Ok   bool
}

func TestSync(t *testing.T) {

	//	s := New()
	//	s.NewPair(SyncConfig{true, true, false}, "D:\\pictures", "E:\\pictures")
	//	s.NewPair(SyncConfig{true, true, false}, "D:\\video", "E:\\pictures")
	//	b, _ := json.Marshal(s.SyncPairs)
	//	ioutil.WriteFile("config.json", b, 777)
	_, err := os.Stat(`E:\\Picturesl`)
	fmt.Println(!os.IsNotExist(err))
}
