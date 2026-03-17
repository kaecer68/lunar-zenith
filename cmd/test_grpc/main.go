package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	lunarv1 "github.com/kaecer68/lunar-zenith/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 連接 gRPC 服務器
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := lunarv1.NewLunarServiceClient(conn)

	// 設置超時上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 測試 GetCalendar
	req := &lunarv1.GetCalendarRequest{Date: "2024-02-10"}
	resp, err := client.GetCalendar(ctx, req)
	if err != nil {
		log.Fatalf("GetCalendar failed: %v", err)
	}

	// 輸出完整響應
	jsonResp, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println("GetCalendar Response:")
	fmt.Println(string(jsonResp))

	// 驗證 SolarTerm 字段
	fmt.Println("\n=== SolarTerm 驗證 ===")
	if resp.SolarTerm == nil {
		fmt.Println("❌ SolarTerm 為 null！")
	} else {
		fmt.Printf("✅ SolarTerm.Index: %d\n", resp.SolarTerm.Index)
		fmt.Printf("✅ SolarTerm.Name: %s\n", resp.SolarTerm.Name)
		fmt.Printf("✅ SolarTerm.Longitude: %.4f\n", resp.SolarTerm.Longitude)
	}
}
