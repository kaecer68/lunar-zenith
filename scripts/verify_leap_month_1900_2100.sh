#!/bin/bash
# 閏月驗證腳本 (1900-2100)
# 驗證關鍵閏月年份

echo "=== 閏月驗證 (1900-2100) ==="
echo ""

# 關鍵測試案例
declare -A test_cases=(
    ["2023-03-22"]="2023年閏二月初一"
    ["2025-07-25"]="2025年閏六月初一"
    ["2033-08-25"]="2033年八月初一（無閏七月）"
    ["2033-12-23"]="2033年閏十一月初一"
)

echo "注意：此腳本需要 lunar-zenith REST 服務運行於 :8080"
echo ""

# 檢查服務
if ! curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo "✗ lunar-zenith 服務未運行"
    echo "請執行: make dev"
    exit 1
fi

echo "✓ 服務運行中"
echo ""

# 測試
echo "日期\t\t\t預期\t\t\t\t實際結果"
echo "--------------------------------------------------------------------------------"

for date in "${!test_cases[@]}"; do
    desc="${test_cases[$date]}"
    
    response=$(curl -s "http://localhost:8080/v1/calendar?date=$date" 2>/dev/null)
    if [ $? -ne 0 ]; then
        echo "$date\t查詢失敗"
        continue
    fi
    
    # 提取農曆資訊
    lunar_year=$(echo "$response" | grep -o '"year":[0-9]*' | head -1 | cut -d':' -f2)
    lunar_month=$(echo "$response" | grep -o '"month":[0-9]*' | head -1 | cut -d':' -f2)
    lunar_day=$(echo "$response" | grep -o '"day":[0-9]*' | head -1 | cut -d':' -f2)
    is_leap=$(echo "$response" | grep -o '"is_leap":[^,}]*' | head -1 | cut -d':' -f2)
    
    leap_str=""
    if [ "$is_leap" = "true" ]; then
        leap_str="閏"
    fi
    
    actual="${lunar_year}年${leap_str}${lunar_month}月${lunar_day}日"
    
    printf "%-12s\t%-30s\t%s\n" "$date" "$desc" "$actual"
done

echo ""
echo "說明："
echo "- 2023年閏二月：正常閏月"
echo "- 2025年閏六月：正常閏月"
echo "- 2033年無閏七月：經典2033年問題（冬至必須在十一月）"
echo "- 2033年閏十一月：2033-2034冬至年的閏月"
