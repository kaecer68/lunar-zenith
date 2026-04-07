package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	lunarv1 "github.com/kaecer68/lunar-zenith/gen"
	runtimecfg "github.com/kaecer68/lunar-zenith/internal/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	grpcPort, err := runtimecfg.GetRequiredPort("LUNAR_GRPC_PORT", "GRPC_PORT")
	if err != nil {
		log.Fatal(fmt.Errorf("load gRPC port: %w", err))
	}
	address, err := grpcTarget(grpcPort)
	if err != nil {
		log.Fatal(err)
	}

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

func grpcTarget(port string) (string, error) {
	if value := strings.TrimSpace(os.Getenv("LUNAR_GRPC_TARGET")); value != "" {
		return value, nil
	}
	if value := strings.TrimSpace(os.Getenv("GRPC_TARGET")); value != "" {
		return value, nil
	}

	host := strings.TrimSpace(os.Getenv("LUNAR_GRPC_HOST"))
	if host == "" {
		host = strings.TrimSpace(os.Getenv("GRPC_HOST"))
	}
	if host == "" {
		return "", fmt.Errorf("missing gRPC target: set LUNAR_GRPC_TARGET or LUNAR_GRPC_HOST")
	}

	return net.JoinHostPort(host, port), nil
}
