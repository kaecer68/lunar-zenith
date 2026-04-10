package service

import (
	"context"

	lunarv1 "github.com/kaecer68/lunar-zenith/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GrpcServer 實現 gRPC LunarService 服務器
type GrpcServer struct {
	lunarv1.UnimplementedLunarServiceServer
	Aggregator *Aggregator
}

// NewGrpcServer 創建 gRPC 服務器
func NewGrpcServer(agg *Aggregator) *GrpcServer {
	return &GrpcServer{
		Aggregator: agg,
	}
}

// GetCalendar 獲取完整曆法數據（gRPC 實現）
func (s *GrpcServer) GetCalendar(ctx context.Context, req *lunarv1.GetCalendarRequest) (*lunarv1.GetCalendarResponse, error) {
	t, err := resolveCalendarQueryTime(req.Date)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, invalidCalendarDateMessage)
	}

	return toCalendarGRPCResponse(s.Aggregator.GetCalendar(t)), nil
}
