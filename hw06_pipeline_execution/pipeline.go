package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func wrapChannel(in In, done In) Out {
	merged := make(Bi)

	go func() {
		defer close(merged)
		for {
			select {
			case <-done:
				return
			default:
			}

			select {
			case <-done:
				return
			case x, ok := <-in:
				if !ok {
					return
				}

				select {
				case <-done:
					return
				default:
				}

				select {
				case <-done:
					return
				case merged <- x:
				}
			}
		}
	}()

	return merged
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	dataCh := wrapChannel(in, done)

	for i := range stages {
		dataCh = stages[i](wrapChannel(dataCh, done))
	}

	return dataCh
}
