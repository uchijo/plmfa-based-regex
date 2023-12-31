package model

import "github.com/google/uuid"

type RegCapture struct {
	// captureする正規表現
	Content RegExp
	// メモリのindex
	MemoryIndex int
}

func (rc RegCapture) States(startId string) ([]State, string, error) {
	// fixme: エラー処理直す
	contentSId, _ := uuid.NewRandom()
	entryState := State{
		Id:    startId,
		IsEnd: false,
		Moves: []Move{
			{
				MType:  CapMem,
				Input:  nil,
				MoveTo: contentSId.String(),
				CInst: CaptureInstr{
					Inst:     Open,
					MemIndex: rc.MemoryIndex,
				},
			},
		},
	}

	cs, endId, _ := rc.Content.States(contentSId.String())

	exitId, _ := uuid.NewRandom()
	exitState := State{
		Id:    endId,
		IsEnd: false,
		Moves: []Move{
			{
				MType:  CapMem,
				Input:  nil,
				MoveTo: exitId.String(),
				CInst: CaptureInstr{
					Inst:     Close,
					MemIndex: rc.MemoryIndex,
				},
			},
		},
	}

	cs = append(cs, entryState, exitState)

	return cs, exitId.String(), nil
}
