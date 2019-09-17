package sparse

import (
	"log"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

func generateIndex() uint64 {
	e := rand.Int()
	if e&1 == 1 {
		i := rand.Intn(3e9)
		return uint64(1e9 + i)
	}
	return uint64(rand.Intn(1e9))
}

func BenchmarkBitmap64(b *testing.B) {
	rand.Seed(1)
	b.ReportAllocs()
	bitmap := NewBitmap64()
	b.N = 1e9
	log.Printf("the value of b.N: %v", b.N)
	PrintMemUsage()
	time.Sleep(300 * time.Millisecond)
	loadBitmap(bitmap, uint64(b.N), 12)
	PrintMemUsage()
	b.ResetTimer()
	b.StartTimer()
	rand.Seed(1)
	count := 0
	for i := 0; i < b.N; i++ {
		x := generateIndex()
		if bitmap.IsSet(x) == true {
			count++
		}
	}
	b.StopTimer()
	log.Printf("bitmap size: %d", bitmap.Size())
	log.Printf("bitmap set count: %d", count)

	PrintMemUsage()
	return
}

func loadBitmap(bitmap *Bitmap64, size uint64, goroutines int) {
	lastGoRoutine := goroutines - 1
	stride := size / uint64(goroutines)

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(g uint64) {
			defer wg.Done()
			start := g * stride
			end := start + stride
			if g == uint64(lastGoRoutine) {
				end = size
			}
			for i := start; i <= end; i++ {
				x := generateIndex()
				bitmap.Set(x)
			}
		}(uint64(i))
	}
	wg.Wait()
}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	log.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	log.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	log.Printf("\tSys = %v MiB", bToMb(m.Sys))
	log.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
