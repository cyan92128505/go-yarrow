# 蓍草占卜 App (Yarrow Stalk Divination)

使用 [go-drift](https://github.com/go-drift/drift) 框架，以純 Go 開發的易經蓍草占卜 App。

## 功能

- **起卦**：模擬傳統蓍草占卜法（大衍之數，三變十八營）
- **結果顯示**：本卦、之卦六爻圖形，變爻標記，朱熹法解卦
- **歷史記錄**：自動儲存所有占卜記錄，以列表顯示
- **刪除記錄**：長按或進入詳情頁刪除
- **分享結果**：將卦象結果以純文字分享

## 專案結構

```
yijing-app/
├── main.go          # App 入口、路由設定、全域初始化
├── divination.go    # 蓍草占卜核心演算法（你的原始邏輯）
├── store.go         # 資料持久化（JSON 檔案）
├── theme.go         # 主題配色（墨色水墨風格）
├── page_home.go     # 首頁：歷史記錄列表
├── page_divine.go   # 起卦頁：輸入事由、執行占卜
├── page_result.go   # 結果頁：卦象詳情、分享、刪除
├── go.mod           # Go module
└── README.md        # 本文件
```

## 環境需求

- Go 1.22+
- Android SDK（用於 Android 建置）
- drift CLI

## 安裝與執行

```bash
# 1. 安裝 drift CLI
go install github.com/go-drift/drift/cmd/drift@latest

# 2. 初始化一個新 drift 專案（取得腳手架）
drift init app
cd app

# 3. 用本專案的檔案替換 scaffold 的 Go 原始碼
#    把以下檔案複製到 app/ 根目錄：
#    - main.go, divination.go, store.go, theme.go
#    - page_home.go, page_divine.go, page_result.go
#    - go.mod (注意 module name 需要與 scaffold 一致，可能需要調整)

# 4. 同步依賴
go mod tidy

# 5. 在 Android 設備上運行
drift run android
```

## 注意事項

### drift init 之後的整合
`drift init` 會產生一個基礎 scaffold，其中包含 `main.go` 和 `go.mod`。
你需要：
1. 把 scaffold 的 `go.mod` 中的 module name 記下來
2. 用本專案的 Go 檔案替換 scaffold 的檔案
3. 確保 `go.mod` 的 module name 一致（或全域 replace）
4. `go mod tidy` 解決依賴

### 卦名對照表
目前程式只輸出數值（6/7/8/9）和朱熹法解卦提示。
你可以在 `divination.go` 中加入 64 卦對照表，將六爻數值映射為卦名。

### 資料儲存
資料以 JSON 檔案儲存在 App 的 Documents 目錄中。
檔案名稱：`divination_records.json`

## 技術細節

- 每爻從 49 策開始，經三變（四營）得出 6/7/8/9
- 6 (老陰) 和 9 (老陽) 為變爻
- 朱熹法根據變爻數量決定查閱方式
- 使用 `time.Now().UnixNano()` 作為隨機種子
