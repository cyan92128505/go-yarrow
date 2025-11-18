package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Line 值為 6,7,8,9 （古法：6=老陰、7=少陽、8=少陰、9=老陽）
type Line int

const (
	OldYin    Line = 6
	YoungYang Line = 7
	YoungYin  Line = 8
	OldYang   Line = 9
)

// doFourOps 在「一變（四營）」時，對當前 stalks 執行一次分二、掛一、揲四、歸奇
// 回傳新合併後剩餘的 stalks（已扣除取出的那些）
func doFourOps(stalks int, rnd *rand.Rand) int {
	if stalks <= 1 {
		return 0
	}
	// 分二：隨機分成兩堆，雙方至少各 1 支
	// 產生 x in [1, stalks-1]
	a := rnd.Intn(stalks-1) + 1
	b := stalks - a

	// 掛一：從其中一堆取出一策（按照你提供內容，不強求左右，隨機選一堆）
	// 若某堆只有 0 則不可能因為我們保證 >=1
	// 選擇 a 或 b 隨機
	var takeOneFromA bool
	if rnd.Intn(2) == 0 {
		takeOneFromA = true
	} else {
		takeOneFromA = false
	}

	if takeOneFromA {
		a -= 1
	} else {
		b -= 1
	}
	// 揲四 & 歸奇：將各堆以 4 為單位分組，取各堆的餘數（若整除則視為 4）
	remA := a % 4
	if remA == 0 {
		remA = 4
	}
	remB := b % 4
	if remB == 0 {
		remB = 4
	}

	// 三處（掛一 + 歸奇A + 歸奇B）共取走
	totalTaken := 1 + remA + remB

	// 剩下的合併回成下一輪的 stalks
	newA := a - remA
	newB := b - remB
	newStalks := newA + newB
	// 因為我們已從一堆「掛一」拿走 1，所以 totalTaken 已包含那 1
	// newStalks 應該等於 stalks - totalTaken
	if newStalks != stalks-totalTaken {
		// 保險起見（理論上不會觸發）
		newStalks = stalks - totalTaken
	}
	return newStalks
}

// makeOneLine 以當前 stalks（通常初始為 49）做三次四營，最後回傳該爻的數字（6|7|8|9）
// 以及最後剩餘的 stalks（通常會是 24/28/32/36 -> 除以4 得到 6/7/8/9）
func makeOneLine(initial int, rnd *rand.Rand) (Line, int) {
	stalks := initial
	// 做三次四營
	for i := 0; i < 3; i++ {
		stalks = doFourOps(stalks, rnd)
	}
	// 三次做完之後剩下的 stalks 應為可被4整除
	if stalks%4 != 0 {
		// 理論上不會發生；若發生，四捨五入到最近的可被4整數（保險處理）
		stalks = (stalks / 4) * 4
	}
	val := stalks / 4 // 會是 6,7,8,9
	return Line(val), stalks
}

// generateHexagram 產生一個本卦（由下而上，第一爻最下面）
// initialStalks 一般為 49
// seed 若為 0，會用 time.Now().UnixNano()
func generateHexagram(initialStalks int, seed int64) ([]Line, []Line, error) {
	if initialStalks <= 0 {
		return nil, nil, fmt.Errorf("initialStalks 必須 > 0")
	}
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	rnd := rand.New(rand.NewSource(seed))

	lines := make([]Line, 0, 6)
	// 每爻都從 49 開始（說明內提到：再將49策合成一堆，「太極」不動，再做三變）
	for i := 0; i < 6; i++ {
		line, _ := makeOneLine(initialStalks, rnd)
		lines = append(lines, line) // 由下而上順序 append
	}

	// 計算之卦（變卦）：九->陰、六->陽，其餘不變（7、8 不變）
	changeLines := make([]Line, 6)
	for i, l := range lines {
		if l == OldYang { // 9 -> 變為陰 (畫成兩橫)
			changeLines[i] = OldYin // 6
		} else if l == OldYin { // 6 -> 變為陽 (畫成一橫)
			changeLines[i] = OldYang // 9
		} else {
			changeLines[i] = l
		}
	}
	return lines, changeLines, nil
}

// prettyPrintLine 把 Line 轉成人易看得文字與符號（陽: 一、陰: - -）
func prettyPrintLine(l Line) (string, string) {
	switch l {
	case OldYang:
		return "老陽 (9)", "───"
	case YoungYang:
		return "少陽 (7)", "───"
	case YoungYin:
		return "少陰 (8)", "─ ─"
	case OldYin:
		return "老陰 (6)", "─ ─"
	default:
		return fmt.Sprintf("未知(%d)", l), ""
	}
}

func main() {
	seed := time.Now().UnixNano()
	fmt.Printf("模擬蓍草占卜 (每爻三變；每變都做四營)。seed=%d\n\n", seed)

	lines, changeLines, err := generateHexagram(49, seed)
	if err != nil {
		panic(err)
	}

	// 由上到下印出本卦（第6爻在上）
	fmt.Println("=== 本卦 (由上到下) ===")
	for i := 5; i >= 0; i-- {
		desc, bar := prettyPrintLine(lines[i])
		fmt.Printf("第%d爻: %s\t%s\n", i+1, desc, bar)
	}
	fmt.Println()

	fmt.Println("=== 之卦 / 變卦 (由上到下) ===")
	for i := 5; i >= 0; i-- {
		desc, bar := prettyPrintLine(changeLines[i])
		fmt.Printf("第%d爻: %s\t%s\n", i+1, desc, bar)
	}
	fmt.Println()

	// 同時輸出由下而上的數列（方便與傳統記錄對應）
	fmt.Println("=== 本卦（由下到上、數值） ===")
	for i := 0; i < 6; i++ {
		fmt.Printf("%d ", lines[i])
	}
	fmt.Println("\n=== 之卦（由下到上、數值） ===")
	for i := 0; i < 6; i++ {
		fmt.Printf("%d ", changeLines[i])
	}
	fmt.Println()
}
