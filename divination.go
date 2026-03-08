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

// IsYang returns true if the line is yang (odd number).
func (l Line) IsYang() bool {
	return l == YoungYang || l == OldYang
}

// IsMoving returns true if the line is a moving line (old yin or old yang).
func (l Line) IsMoving() bool {
	return l == OldYin || l == OldYang
}

// Symbol returns the visual representation of the line.
func (l Line) Symbol() string {
	if l.IsYang() {
		return "━━━━━"
	}
	return "━━　━━"
}

// Label returns description like "老陽 (9)".
func (l Line) Label() string {
	switch l {
	case OldYang:
		return "老陽 (9)"
	case YoungYang:
		return "少陽 (7)"
	case YoungYin:
		return "少陰 (8)"
	case OldYin:
		return "老陰 (6)"
	default:
		return fmt.Sprintf("未知(%d)", l)
	}
}

type Hexagram struct {
	Original    []Line
	Changed     []Line
	MovingLines []int
}

// doFourOps 在「一變（四營）」時，對當前 stalks 執行一次分二、掛一、揲四、歸奇
func doFourOps(stalks int, rnd *rand.Rand) int {
	if stalks <= 1 {
		return 0
	}
	a := rnd.Intn(stalks-1) + 1
	b := stalks - a

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

	remA := a % 4
	if remA == 0 {
		remA = 4
	}
	remB := b % 4
	if remB == 0 {
		remB = 4
	}

	totalTaken := 1 + remA + remB
	newA := a - remA
	newB := b - remB
	newStalks := newA + newB
	if newStalks != stalks-totalTaken {
		newStalks = stalks - totalTaken
	}
	return newStalks
}

// makeOneLine 以 49 策做三次四營，回傳該爻的數字（6|7|8|9）
func makeOneLine(initial int, rnd *rand.Rand) Line {
	stalks := initial
	for i := 0; i < 3; i++ {
		stalks = doFourOps(stalks, rnd)
	}
	if stalks%4 != 0 {
		stalks = (stalks / 4) * 4
	}
	val := stalks / 4
	return Line(val)
}

// generateHexagram 產生本卦與之卦（由下而上，第一爻最下面）
func generateHexagram(seed int64) *Hexagram {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	rnd := rand.New(rand.NewSource(seed))

	hex := &Hexagram{
		Original:    make([]Line, 0, 6),
		Changed:     make([]Line, 6),
		MovingLines: make([]int, 0, 6),
	}

	for i := 0; i < 6; i++ {
		line := makeOneLine(49, rnd)
		hex.Original = append(hex.Original, line)
	}

	for i, line := range hex.Original {
		switch line {
		case OldYang:
			hex.Changed[i] = YoungYin
			hex.MovingLines = append(hex.MovingLines, i+1)
		case OldYin:
			hex.Changed[i] = YoungYang
			hex.MovingLines = append(hex.MovingLines, i+1)
		default:
			hex.Changed[i] = line
		}
	}

	return hex
}

// InterpretZhuXi 朱熹法解卦
func (h *Hexagram) InterpretZhuXi() string {
	n := len(h.MovingLines)

	switch n {
	case 0:
		return "參考本卦卦辭"
	case 1:
		return fmt.Sprintf("參考本卦變爻爻辭：第 %d 爻", h.MovingLines[0])
	case 2:
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
		return fmt.Sprintf("參考本卦第 %d 爻（上爻為主，第 %d 爻為參照）", upper, lower)
	case 3:
		return "參考變卦卦辭（以本卦為輔助）"
	case 4:
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
		return fmt.Sprintf("參考變卦第 %d 爻（下不變爻，上不變爻可參考：第 %d 爻）", lower, upper)
	case 5:
		staticLines := h.findStaticLines()
		return fmt.Sprintf("變卦唯一不變爻爻辭：第 %d 爻", staticLines[0])
	case 6:
		if h.isQianOrKun() {
			if h.isAllYang() {
				return "參考乾卦：用九"
			}
			return "參考坤卦：用六"
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

// ToShareText 產生可分享的純文字結果
func (h *Hexagram) ToShareText(question string, createdAt time.Time) string {
	text := "【蓍草占卜】\n"
	text += fmt.Sprintf("時間：%s\n", createdAt.Format("2006-01-02 15:04"))
	if question != "" {
		text += fmt.Sprintf("事由：%s\n", question)
	}
	text += "\n本卦（由上到下）：\n"
	for i := 5; i >= 0; i-- {
		marker := "  "
		if h.Original[i].IsMoving() {
			marker = "◯"
		}
		text += fmt.Sprintf("  第%d爻  %s  %s  %s\n", i+1, h.Original[i].Symbol(), h.Original[i].Label(), marker)
	}

	if len(h.MovingLines) > 0 {
		text += "\n之卦（由上到下）：\n"
		for i := 5; i >= 0; i-- {
			text += fmt.Sprintf("  第%d爻  %s  %s\n", i+1, h.Changed[i].Symbol(), h.Changed[i].Label())
		}
	}

	text += "\n" + h.InterpretZhuXi() + "\n"
	return text
}

// hexagramUnicode 從六爻（由下到上）產生卦象 Unicode 符號
func hexagramUnicode(lines []Line) string {
	// 六爻轉為上下卦的三位二進位（陽=1, 陰=0）
	// lines[0] 是初爻（最下），lines[5] 是上爻
	upper := trigramIndex(lines[0], lines[1], lines[2])
	lower := trigramIndex(lines[3], lines[4], lines[5])

	// 八卦順序：坤(0) 剝(1)... 用 King Wen 序查表
	// 行=上卦，列=下卦，值=King Wen 序號（0-based）
	// 卦序：乾☰=111, 兌☱=110, 離☲=101, 震☳=100, 巽☴=011, 坎☵=010, 艮☶=001, 坤☷=000
	kingWen := [8][8]int{
		//  坤   艮   坎   巽   震   離   兌   乾    ← lower trigram
		{1, 23, 8, 20, 16, 35, 45, 12},   // 坤 upper
		{15, 52, 39, 53, 62, 56, 31, 33}, // 艮
		{7, 4, 29, 59, 40, 64, 47, 6},    // 坎
		{46, 18, 48, 57, 32, 50, 28, 44}, // 巽
		{24, 27, 3, 42, 51, 21, 17, 25},  // 震
		{36, 22, 63, 37, 55, 30, 49, 13}, // 離
		{19, 41, 60, 61, 54, 38, 58, 10}, // 兌
		{2, 26, 5, 9, 34, 14, 43, 11},    // 乾
	}

	num := kingWen[upper][lower] // 1-based King Wen 序
	return string(rune(0x4DC0 + num - 1))
}

// trigramIndex 將三爻轉為 0-7 的索引
// 坤=000(0), 艮=001(1), 坎=010(2), 巽=011(3),
// 震=100(4), 離=101(5), 兌=110(6), 乾=111(7)
func trigramIndex(bottom, middle, top Line) int {
	idx := 0
	if bottom.IsYang() {
		idx |= 4
	}
	if middle.IsYang() {
		idx |= 2
	}
	if top.IsYang() {
		idx |= 1
	}
	return idx
}
