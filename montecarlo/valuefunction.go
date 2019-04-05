package montecarlo

import (
	"fmt"
	"love-letter-ai/gamemaster"
	"love-letter-ai/players"
	"love-letter-ai/state"
)

type Value struct {
	// The maximum visits to a state isn't likely to be much above 65k, so let's save some RAM.
	// It's necessary to deal with overflow, but we don't care about loss of precision from occasionally dividing by two.
	sum   uint16
	count uint16
}

type ValueFunction [state.SpaceMagnitude]Value
type Action [state.SpaceMagnitude]uint8

func (vf *ValueFunction) Train(pl players.Player, episodes int) {
	for i := 0; i < episodes; i++ {
		if (i % 100000) == 0 {
			fmt.Printf("% 2.2f%% complete\r", float32(i)/float32(episodes)*100)
		}
		vf.Update(pl)
	}
	fmt.Println("100.0% complete")
}

func (vf *ValueFunction) Update(pl players.Player) {
	tr, err := gamemaster.TraceOneGame(pl)
	if err != nil {
		panic(err.Error())
	}

	for _, si := range tr.StateInfos {
		vf.SaveState(si)
	}
}

const numChans = 8
const numGenThreads = 8

// Update runs 8*episodes
func (vf *ValueFunction) ThreadedTrain(pl players.Player, episodes int) {
	chs := make([]chan gamemaster.StateInfo, numChans)
	doneStoring := make(chan bool, numChans)

	for i := range chs {
		ch := make(chan gamemaster.StateInfo, 128)
		chs[i] = ch

		go func() {
			for true {
				si, more := <-ch

				if !more {
					doneStoring <- true
					return
				}

				vf.SaveState(si)
			}
		}()
	}

	// Launch some threads to generate data
	doneGenerating := make(chan bool, numGenThreads)
	for i := 0; i < numGenThreads; i++ {
		go func() {
			for i := 0; i < episodes; i++ {
				err := gamemaster.TraceIntoChannels(pl, chs)
				if err != nil {
					panic(err.Error())
				}
			}
			doneGenerating <- true
		}()
	}

	// Wait for the data to all be generated
	for i := 0; i < numGenThreads; i++ {
		<-doneGenerating
	}
	close(doneGenerating)

	// Now close all of the generating channels because generation is done
	for _, ch := range chs {
		close(ch)
	}

	// Wait for the data to all be saved
	for i := 0; i < numChans; i++ {
		<-doneStoring
	}
	close(doneStoring)
}

func (vf *ValueFunction) Value(state int) float32 {
	return float32(vf[state].sum) / float32(vf[state].count)
}

func (vf *ValueFunction) SaveState(si gamemaster.StateInfo) {
	s := si.State
	if si.Won {
		vf[s].sum++
	}
	vf[s].count++

	// Check for overflow
	if vf[s].count == 0xFFFF {
		vf[s].sum /= 2
		vf[s].count /= 2
	}
}
