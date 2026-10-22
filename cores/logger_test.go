package cores_test

import (
	"sync"
	"testing"

	"github.com/rachmanzz/fiber-starter/cores"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewLogger_Concurrent(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cores.NewLogger()
		}()
	}
	wg.Wait()

	assert.NotNil(t, cores.Logger)
	assert.NotNil(t, zap.L())
}
