package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	lunarv1 "github.com/kaecer68/lunar-zenith/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getEnvWithFallback(primary, fallback string) string {
	if v := os.Getenv(primary); v != "" {
		return v
	}
	if v := os.Getenv(fallback); v != "" {
		return v
	}
	return loadPortFromEnvFile(primary)
}

func loadPortFromEnvFile(key string) string {
	data, err := os.ReadFile(".env.ports")
	if err != nil {
		log.Fatalf("錯誤：找不到 .env.ports 檔案。請先執行 'make sync-contracts' 同步契約。%v", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 && strings.TrimSpace(parts[0]) == key {
			return strings.TrimSpace(parts[1])
		}
	}

	log.Fatalf("錯誤：在 .env.ports 中找不到 %s。請執行 'make sync-contracts' 重新生成。", key)
	return ""
}

func main() {
	grpcPort := getEnvWithFallback("LUNAR_GRPC_PORT", "GRPC_PORT")
	address := fmt.Sprintf("localhost:%s", grpcPort)

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := lunarv1.NewLunarServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &lunarv1.GetCalendarRequest{Date: "2024-02-10"}
	resp, err := client.GetCalendar(ctx, req)
	if err != nil {
		log.Fatalf("GetCalendar failed: %v", err)
	}

	jsonResp, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println("GetCalendar Response:")
	fmt.Println(string(jsonResp))

	fmt.Println("\n=== SolarTerm 驗證 ===")
	if resp.SolarTerm == nil {
		fmt.Println("❌ SolarTerm 為 null！")
	} else {
		fmt.Printf("✅ SolarTerm.Index: %d\n", resp.SolarTerm.Index)
		fmt.Printf("✅ SolarTerm.Name: %s\n", resp.SolarTerm.Name)
		fmt.Printf("✅ SolarTerm.Longitude: %.4f\n", resp.SolarTerm.Longitude)
	}
}
