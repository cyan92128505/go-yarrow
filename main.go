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

type Hexagram struct {
	Original    []Line
	Changed     []Line
	MovingLines []int
}

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
func generateHexagram(initialStalks int, seed int64) (*Hexagram, error) {
	if initialStalks <= 0 {
		return nil, fmt.Errorf("initialStalks 必須 > 0")
	}
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	rnd := rand.New(rand.NewSource(seed))

	hex := &Hexagram{
		Original:    make([]Line, 0, 6),
		Changed:     make([]Line, 6),
		MovingLines: make([]int, 0, 6),
	}

	// 每爻都從 49 開始（說明內提到：再將49策合成一堆，「太極」不動，再做三變）
	for i := 0; i < 6; i++ {
		line, _ := makeOneLine(initialStalks, rnd)
		hex.Original = append(hex.Original, line) // 由下而上順序 append
	}

	// 計算之卦（變卦）：九(老陽)->八(少陰)、六(老陰)->七(少陽)，其餘不變（7、8 不變）
	for i, line := range hex.Original {
		switch line {
		case OldYang: // 9 -> 變為陰 (畫成兩橫)
			hex.Changed[i] = YoungYin
			hex.MovingLines = append(hex.MovingLines, i+1)
		case OldYin: // // 6 -> 變為陽 (畫成一橫)
			hex.Changed[i] = YoungYang
			hex.MovingLines = append(hex.MovingLines, i+1)
		default:
			hex.Changed[i] = line
		}
	}

	return hex, nil
}

func (h *Hexagram) InterpretZhuXi() string {
	n := len(h.MovingLines)

	switch n {
	case 0:
		return "參考本卦卦辭"
	case 1:
		return fmt.Sprintf("參考本卦變爻爻辭: 第 %d 爻", h.MovingLines[0])
	case 2:
		// Take the upper moving line as primary
		upper := h.MovingLines[0]
		lower := h.MovingLines[0]
		for _, pos := range h.MovingLines {
			if pos > upper {
				upper = pos
			}
			if pos < lower {
				lower = pos
			}
		}
		return fmt.Sprintf("參考本卦的第 %d 爻（上爻為主爻，第 %d 爻為參照爻）", upper, lower)
	case 3:
		return "參考變卦卦辭（以本卦為參考）。"
	case 4:
		// Find two non-moving lines, take the lower one as primary
		staticLines := h.findStaticLines()
		lower := staticLines[0]
		upper := staticLines[0]
		for _, pos := range staticLines {
			if pos < lower {
				lower = pos
			}
			if pos > upper {
				upper = pos
			}
		}
		return fmt.Sprintf("參考變卦的第 %d 爻 (變卦下不變爻爻辭, 上不變爻可參考: 第 %d 爻)", lower, upper)
	case 5:
		// Only one static line remains
		staticLines := h.findStaticLines()
		return fmt.Sprintf("變卦唯一不變爻爻辭: 第 %d 爻", staticLines[0])
	case 6:
		// All lines moving
		if h.isQianOrKun() {
			if h.isAllYang() {
				return "參考乾卦:用九"
			}
			return "參考坤卦:用六"
		}
		return "參考變卦卦辭"
	}

	return ""
}

func (h *Hexagram) findStaticLines() []int {
	staticMap := make(map[int]bool)
	for i := 1; i <= 6; i++ {
		staticMap[i] = true
	}
	for _, pos := range h.MovingLines {
		delete(staticMap, pos)
	}

	static := make([]int, 0, 6-len(h.MovingLines))
	for i := 1; i <= 6; i++ {
		if staticMap[i] {
			static = append(static, i)
		}
	}
	return static
}

func (h *Hexagram) isQianOrKun() bool {
	return h.isAllYang() || h.isAllYin()
}

func (h *Hexagram) isAllYang() bool {
	for _, line := range h.Original {
		if line == OldYin || line == YoungYin {
			return false
		}
	}
	return true
}

func (h *Hexagram) isAllYin() bool {
	for _, line := range h.Original {
		if line == OldYang || line == YoungYang {
			return false
		}
	}
	return true
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

	hex, err := generateHexagram(49, seed)
	if err != nil {
		panic(err)
	}

	// 由上到下印出本卦（第6爻在上）
	fmt.Println("=== 本卦 (由上到下) ===")
	for i := 5; i >= 0; i-- {
		desc, bar := prettyPrintLine(hex.Original[i])
		fmt.Printf("第%d爻: %s\t%s\n", i+1, desc, bar)
	}
	fmt.Println()

	fmt.Println("=== 之卦 / 變卦 (由上到下) ===")
	for i := 5; i >= 0; i-- {
		desc, bar := prettyPrintLine(hex.Changed[i])
		fmt.Printf("第%d爻: %s\t%s\n", i+1, desc, bar)
	}
	fmt.Println()

	// 同時輸出由下而上的數列（方便與傳統記錄對應）
	fmt.Println("=== 本卦（由下到上、數值） ===")
	for i := 0; i < 6; i++ {
		fmt.Printf("%d ", hex.Original[i])
	}
	fmt.Println("\n=== 之卦（由下到上、數值） ===")
	for i := 0; i < 6; i++ {
		fmt.Printf("%d ", hex.Changed[i])
	}
	fmt.Println()

	fmt.Println("\n=== 爻變 ===")
	fmt.Println(hex.InterpretZhuXi())
	fmt.Println()
}
