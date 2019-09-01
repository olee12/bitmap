package roaring

import (
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"

	"git.garena.com/common/gocommon"
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
	gocommon.LoggerInit("log/error.log", 3600*24, 1024*1024*128, 10, 3)

	b.ReportAllocs()

	rand.Seed(1)

	bitmap := NewBitmap64()
	b.N = 1e9
	gocommon.LogDetailf("the value of b.N: %v", b.N)
	PrintMemUsage()
	time.Sleep(time.Millisecond * 300)
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

	gocommon.LogDetailf("bitmap size: %d", len(bitmap.keys))
	arrayCount, bitmapCount := 0, 0
	for _, container := range bitmap.containers {
		switch container.(type) {
		case *arrayContainer:
			arrayCount++
		case *bitmapContainer:
			bitmapCount++
		}
	}
	time.Sleep(time.Millisecond * 200)
	gocommon.LogDetailf("arrayCount: %d, bitmapCount: %d", arrayCount, bitmapCount)
	gocommon.LogDetailf("total set bits := %v", count)
	PrintMemUsage()
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
	gocommon.LogDetailf("Alloc = %v MiB", bToMb(m.Alloc))
	gocommon.LogDetailf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	gocommon.LogDetailf("\tSys = %v MiB", bToMb(m.Sys))
	gocommon.LogDetailf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
