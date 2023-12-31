package eval

import (
	"fmt"

	"github.com/uchijo/plmfa-based-regex/model"
)

var logs = []Log{}

func isInLoop(log Log) bool {
	for _, v := range logs {
		if v.Alike(log) {
			return true
		}
	}
	return false
}

func Search(states model.StateList, input InputBuffer, start string, eSem bool, showLog bool) bool {
	// ログ初期化
	logs = []Log{}
	return search(states, input, start, PosMemoryList{}, CapMemoryList{}, 0, eSem, showLog)
}

// search returns true if plmfa accepts input
func search(
	st model.StateList,
	input InputBuffer,
	currentId string,
	posMem PosMemoryList,
	capMem CapMemoryList,
	depth int,
	epsilonSem bool,
	showLog bool,
) bool {
	// 無限ループ検知
	currentLog := Log{
		Input:          input.Input,
		CurrentId:      currentId,
		PositiveMemory: posMem,
		CaptureMemory:  capMem,
	}
	if isInLoop(currentLog) {
		return false
	}
	logs = append(logs, currentLog)

	if showLog {
		for i := 0; i < depth; i++ {
			fmt.Printf("  ")
		}
		fmt.Printf("buffer: %v, state: %v, pos_memory: %+v, cap_memory: %+v\n", input.Input, currentId, posMem, capMem)
	}
	// fixme: エラー処理直す
	curState, _ := st.StateById(currentId)

	// 終状態でない & 行先がない -> マッチ失敗
	if !curState.IsEnd && len(curState.Moves) == 0 {
		return false
	}

	// 終状態 & 入力文字がない -> マッチ成功
	if curState.IsEnd && input.Len() == 0 {
		return true
	}

	// 入力文字がない & イプシロンでたどり着けるゴールがある -> マッチ
	searchEpsLog = []string{}
	goalReachable := searchEps(st, currentId)
	if goalReachable && input.Len() == 0 {
		return true
	}

	for _, v := range curState.Moves {
		switch v.MType {
		case model.Epsilon:
			hasGoal := search(st, input, v.MoveTo, posMem, capMem, depth+1, epsilonSem, showLog)
			if hasGoal {
				return true
			}

		case model.Consumption:
			if ok, toConsume := input.CanConsume(v.Input); ok {
				consumed, _ := input.Consumed(toConsume)
				appendedPos := posMem.Appended(toConsume)
				appendedCap := capMem.Appended(toConsume)
				hasGoal := search(st, consumed, v.MoveTo, appendedPos, appendedCap, depth+1, epsilonSem, showLog)
				if hasGoal {
					return true
				}
			}

		case model.PosMem:
			if v.PLInst.Inst == model.Open {
				opened := posMem.OpenedMem(v.PLInst.MemIndex)
				hasGoal := search(st, input, v.MoveTo, opened, capMem, depth+1, epsilonSem, showLog)
				if hasGoal {
					return true
				}
			} else if v.PLInst.Inst == model.Close {
				closed, memContent := posMem.ClosedMem(v.PLInst.MemIndex)
				appended, _ := input.Appended(memContent)
				hasGoal := search(st, appended, v.MoveTo, closed, capMem, depth+1, epsilonSem, showLog)
				if hasGoal {
					return true
				}
			}

		case model.CapMem:
			if v.CInst.Inst == model.Open {
				opened := capMem.OpenedMem(v.CInst.MemIndex)
				hasGoal := search(st, input, v.MoveTo, posMem, opened, depth+1, epsilonSem, showLog)
				if hasGoal {
					return true
				}
			} else if v.CInst.Inst == model.Close {
				closed := capMem.ClosedMem(v.CInst.MemIndex)
				hasGoal := search(st, input, v.MoveTo, posMem, closed, depth+1, epsilonSem, showLog)
				if hasGoal {
					return true
				}
			}

		case model.Ref:
			mem, err := capMem.Content(v.RefIndex, epsilonSem)
			// エラーが返ってくるということはキャプチャを拾えなかったということ
			if err != nil {
				continue
			}
			memContainer := stringContainer(mem)
			if ok, toConsume := input.CanConsume(memContainer); ok {
				consumed, _ := input.Consumed(mem)
				// posMem, capMemはメモリの集合
				// なので、refでアクセス中のメモリとか気にせず消費したものは記録する必要あり
				appendedPos := posMem.Appended(toConsume)
				appendedCap := capMem.Appended(toConsume)
				hasGoal := search(st, consumed, v.MoveTo, appendedPos, appendedCap, depth+1, epsilonSem, showLog)
				if hasGoal {
					return true
				}
			}

		case model.ArbitraryConsumption:
			// 1文字は残ってないと消費できない
			if input.Len() >= 1 {
				toConsume := input.Input[:1]
				consumed, _ := input.Consumed(toConsume)
				appendedPos := posMem.Appended(toConsume)
				appendedCap := capMem.Appended(toConsume)
				hasGoal := search(st, consumed, v.MoveTo, appendedPos, appendedCap, depth+1, epsilonSem, showLog)
				if hasGoal {
					return true
				}
			}
		}
	}

	if showLog {
		for i := 0; i < depth; i++ {
			fmt.Printf("  ")
		}
		fmt.Println("backtrack!")
	}
	return false
}

var searchEpsLog []string
// searchEps returns true if goal state is reachable with only epsilon transitions
// initialize searchEpsLog before call.
func searchEps(st model.StateList, currentId string) bool {
	// 既に探索済みidに当たったら無限ループ -> この枝は失敗
	for _, v := range searchEpsLog {
		if v == currentId {
			return false
		}
	}
	searchEpsLog = append(searchEpsLog, currentId)

	curSt, _ := st.StateById(currentId)
	if curSt.IsEnd {
		return true
	}

	eps := curSt.ExtractMove(model.Epsilon)
	for _, v := range eps {
		hasGoal := searchEps(st, v.MoveTo)
		if hasGoal {
			return true
		}
	}

	return false
}
