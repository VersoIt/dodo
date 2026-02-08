package grpc

import (
	"context"

	"github.com/versoit/diploma/services/analytics/api/proto/pb"
	"github.com/versoit/diploma/services/analytics/usecase"
	"google.golang.org/grpc"
)

type AnalyticsHandler struct {
	pb.UnimplementedKpiServiceServer
	uc *usecase.AnalyticsUseCase
}

func NewAnalyticsHandler(uc *usecase.AnalyticsUseCase) *AnalyticsHandler {
	return &AnalyticsHandler{uc: uc}
}

func (h *AnalyticsHandler) Register(server *grpc.Server) {
	pb.RegisterKpiServiceServer(server, h)
}

func (h *AnalyticsHandler) GetManagerKPI(ctx context.Context, req *pb.KpiRequest) (*pb.KpiResponse, error) {
	kpi, err := h.uc.GetManagerPerformance(ctx, req.ManagerId)
	if err != nil {
		return nil, err
	}

	return &pb.KpiResponse{
		ManagerId:    kpi.ManagerID(),
		FactRevenue:  kpi.Fact(),
		PlanRevenue:  kpi.Plan(),
		KpiPercent:   kpi.CalculateKPIPercent(),
		HasBonus:     kpi.HasBonus(),
	}, nil
}
