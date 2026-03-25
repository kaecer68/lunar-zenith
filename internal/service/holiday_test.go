package service

import (
	"strings"
	"testing"
)

func TestHolidayService(t *testing.T) {
	s := NewHolidayService()
	s.holidays["20240101"] = TaiwanHoliday{
		Date:      "20240101",
		Name:      "元旦",
		IsHoliday: true,
		Category:  TypePublicHoliday,
	}
	s.holidays["20240217"] = TaiwanHoliday{
		Date:      "20240217",
		Name:      "補班",
		IsHoliday: false,
		Category:  TypeWorkday,
	}

	// 測試放假
	isHol, name := s.IsHoliday("20240101")
	if !isHol || name != "元旦" {
		t.Errorf("20240101 應該是放假的元旦, got %v, %s", isHol, name)
	}

	// 測試補班
	isHol2, name2 := s.IsHoliday("20240217")
	if isHol2 || name2 != "" {
		t.Errorf("20240217 應該是補班不放假且不顯示名稱, got %v, %s", isHol2, name2)
	}

	// 測試非記錄平日
	isHol3, _ := s.IsHoliday("20240520")
	if isHol3 {
		t.Errorf("未記錄日期不應預設為假期")
	}

	// 測試非記錄週末（預設放假）
	isHol4, name4 := s.IsHoliday("20240519")
	if !isHol4 || name4 != "" {
		t.Errorf("未記錄週末應預設放假但不顯示名稱, got %v, %s", isHol4, name4)
	}
}

func TestHolidayServiceWithFallback(t *testing.T) {
	tw := NewHolidayService()
	tw.holidays["20241010"] = TaiwanHoliday{
		Date:      "20241010",
		Name:      "國慶日",
		IsHoliday: true,
		Category:  TypePublicHoliday,
	}

	cn := NewHolidayServiceWithFallback(tw)
	cn.holidays["20241010"] = TaiwanHoliday{
		Date:      "20241010",
		Name:      "一般工作日",
		IsHoliday: false,
		Category:  TypeWorkday,
	}
	cn.holidays["20241001"] = TaiwanHoliday{
		Date:      "20241001",
		Name:      "国庆节",
		IsHoliday: true,
		Category:  TypePublicHoliday,
	}

	isHol1, name1 := cn.IsHoliday("20241001")
	if !isHol1 || name1 != "国庆节" {
		t.Errorf("20241001 應該使用大陸覆寫假期, got %v, %s", isHol1, name1)
	}

	isHol2, name2 := cn.IsHoliday("20241010")
	if isHol2 || name2 != "" {
		t.Errorf("20241010 應該覆寫為工作日且不顯示名稱, got %v, %s", isHol2, name2)
	}

	tw.holidays["20240101"] = TaiwanHoliday{
		Date:      "20240101",
		Name:      "元旦",
		IsHoliday: true,
		Category:  TypePublicHoliday,
	}
	isHol3, name3 := cn.IsHoliday("20240101")
	if !isHol3 || name3 != "元旦" {
		t.Errorf("20240101 應該回退到台灣假期, got %v, %s", isHol3, name3)
	}
}

func TestHolidayJSONCoverage(t *testing.T) {
	tw := NewHolidayService()
	if err := tw.LoadFromJSON("../../configs/holidays_tw_2024_2026.json"); err != nil {
		t.Errorf("載入台灣假期 JSON 失敗: %v", err)
	}

	cn := NewHolidayServiceWithFallback(tw)
	if err := cn.LoadFromJSON("../../configs/holidays_cn_2024_2026.json"); err != nil {
		t.Errorf("載入大陸假期 JSON 失敗: %v", err)
	}

	if isHol, name := tw.IsHoliday("20240210"); !isHol || name == "" {
		t.Errorf("台灣春節應為假日, got %v, %s", isHol, name)
	}

	if isHol, _ := tw.IsHoliday("20240217"); isHol {
		t.Errorf("台灣補班日不應為假日")
	}

	if isHol, name := cn.IsHoliday("20241001"); !isHol || name != "国庆节" {
		t.Errorf("大陸國慶日應為假日, got %v, %s", isHol, name)
	}

	if isHol, _ := cn.IsHoliday("20241010"); isHol {
		t.Errorf("大陸 10/10 應覆寫為工作日")
	}

	if isHol, name := cn.IsHoliday("20240101"); !isHol || name != "元旦" {
		t.Errorf("大陸應回退使用台灣共有假日元旦, got %v, %s", isHol, name)
	}
}

func TestTaiwanObservanceNonHolidayVisible(t *testing.T) {
	s := NewHolidayService()

	isHol, name := s.IsHoliday("20240314")
	if isHol {
		t.Errorf("反侵略日不應為放假日")
	}
	if name != "反侵略日" {
		t.Errorf("反侵略日應顯示節日名稱, got %s", name)
	}
}

func TestTaiwanObservanceMultipleNames(t *testing.T) {
	s := NewHolidayService()

	isHol, name := s.IsHoliday("20240312")
	if isHol {
		t.Errorf("20240312 不應為放假日")
	}
	if name != "國父逝世紀念日、植樹節" && name != "植樹節、國父逝世紀念日" {
		t.Errorf("20240312 應同時顯示紀念日與節日, got %s", name)
	}
}

func TestChinaObservanceOverrides(t *testing.T) {
	tw := NewHolidayService()
	cn := NewHolidayServiceWithFallback(tw)

	if isHol, name := tw.IsHoliday("20240801"); isHol || name != "原住民族日" {
		t.Errorf("台灣 8/1 應為原住民族日, got %v, %s", isHol, name)
	}

	if isHol, name := cn.IsHoliday("20240801"); isHol || name != "建軍節" {
		t.Errorf("大陸 8/1 應為建軍節且非放假, got %v, %s", isHol, name)
	}

	if isHol, name := tw.IsHoliday("20240808"); isHol || name != "父親節" {
		t.Errorf("台灣 8/8 應為父親節, got %v, %s", isHol, name)
	}

	if isHol, name := cn.IsHoliday("20240808"); isHol || name != "" {
		t.Errorf("大陸 8/8 不應顯示台灣父親節, got %v, %s", isHol, name)
	}

	if isHol, name := cn.IsHoliday("20240601"); isHol || name != "兒童節" {
		t.Errorf("大陸兒童節應為 6/1, got %v, %s", isHol, name)
	}

	if isHol, name := cn.IsHoliday("20240616"); isHol || name != "父親節" {
		t.Errorf("大陸父親節應為六月第三個週日, got %v, %s", isHol, name)
	}

	if isHol, name := cn.IsHoliday("20240910"); isHol || name != "教師節" {
		t.Errorf("大陸教師節應為 9/10, got %v, %s", isHol, name)
	}

	if isHol, name := cn.IsHoliday("20240928"); strings.Contains(name, "教師節") {
		t.Errorf("大陸不應在 9/28 顯示教師節, got %v, %s", isHol, name)
	}

	if isHol, name := cn.IsHoliday("20240404"); !isHol || name != "清明節" {
		t.Errorf("大陸 4/4 應顯示清明節且放假, got %v, %s", isHol, name)
	}
}
