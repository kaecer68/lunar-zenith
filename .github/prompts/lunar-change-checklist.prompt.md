---
description: "為曆法演算法修改建立執行前檢查清單（JDE/Delta-T、閏月邊界、契約同步、測試範圍）"
name: "Lunar Change Checklist"
argument-hint: "可選：填入本次變更目標，例如『調整 2033 邊界年閏月判定』"
agent: "agent"
---
針對目前工作區與你提供的目標，產生一份「曆法修改前檢查清單」。

若命令後方有補充內容，將其視為本次變更目標；若沒有，從目前變更與檔案上下文推定目標。

請輸出四個區塊：

1. 變更風險分類
- 判斷屬於哪一類：天文計算、農曆轉換、閏月/邊界、契約/傳輸層、其他。
- 用一句話說明主要風險來源。

2. 必做前置檢查
- 是否涉及 `pkg/zodiac/**/*.go`，若是，確認是否需要套用 zodiac-boundary 規範。
- 是否涉及 API 欄位/輸出，若是，列出 contract-first 先後順序。
- 是否有 UT/JDE 混用風險，若有，標註要檢查的函式與資料流。

3. 最小驗證命令
- 先列「最小必要」命令，再列「完整驗證」命令。
- 至少涵蓋：
  - `make verify-contracts`（若有契約或埠號影響）
  - `go test ./pkg/zodiac -run "TestLunarEngineLeap|Test.*Leap" -v`（若涉閏月）
  - `go test ./pkg/zodiac -run "TestLunarEngineEdge|Test.*Edge" -v`（若涉邊界）
  - `make test`（整體回歸）

4. 完成定義（DoD）
- 用 5 條以內列出可合併條件。
- 包含：測試通過、無 generated file 手改、契約與實作一致、無越界宣稱（未驗證年份不宣稱正確）。

輸出要求：
- 精簡、可執行、避免泛泛建議。
- 若資訊不足，明確列出缺少哪些上下文，並提供保守預設清單。
