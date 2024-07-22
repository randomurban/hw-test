package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var (
		inStage  In
		outStage Out
	)
	inStage = in
	for _, stage := range stages {
		outStage = stage(doneFn(done, inStage))
		inStage = outStage
	}
	return outStage
}

func doneFn(done In, in In) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		for v := range in {
			select {
			case <-done:
				return
			default:
			}
			out <- v
		}
	}()
	return out
}
