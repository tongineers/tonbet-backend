package workerpool_test

import (
	"log"
	"sync"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tongineers/tonbet-backend/pkg/workerpool"
)

func TestWorkerPool(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test WorkerPool Suite")
}

type (
	Task struct {
		ID int
		wg *sync.WaitGroup
	}
)

func (t *Task) Run() {
	defer t.wg.Done()
	log.Printf("Running task with ID %d", t.ID)
	time.Sleep(3 * time.Second)
	log.Printf("Task with ID %d stopped.", t.ID)
}

var _ = Describe("Test WorkerPool", func() {
	var wp *workerpool.WorkerPool

	wg := &sync.WaitGroup{}
	wg.Add(3)

	var (
		tasks = []workerpool.Task{
			&Task{ID: 1, wg: wg},
			&Task{ID: 2, wg: wg},
			&Task{ID: 3, wg: wg},
		}
	)

	BeforeEach(func() {
		wp = workerpool.NewWorkerPool(len(tasks))
		wp.Start()
	})

	Describe("", func() {
		JustBeforeEach(func() {
			wp.Submit(tasks[0])
			wp.Submit(tasks[1])
			wp.Submit(tasks[2])
		})

		When("have new active bets", func() {
			BeforeEach(func() {

			})

			It("should not have error occured", func() {
				wg.Wait()
			})
		})
	})
})
