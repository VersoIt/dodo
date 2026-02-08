package grpc

import (
	"context"

	analytics_pb "github.com/versoit/diploma/services/analytics/api/proto/pb"
	"github.com/versoit/diploma/services/analytics/usecase"
	"google.golang.org/grpc"
)

type AnalyticsHandler struct {
	analytics_pb.UnimplementedKpiServiceServer
	uc *usecase.AnalyticsUseCase
}

func NewAnalyticsHandler(uc *usecase.AnalyticsUseCase) *AnalyticsHandler {
	return &AnalyticsHandler{uc: uc}
}

func (h *AnalyticsHandler) Register(server *grpc.Server) {
	analytics_pb.RegisterKpiServiceServer(server, h)
}

func (h *AnalyticsHandler) GetManagerKPI(ctx context.Context, req *analytics_pb.KpiRequest) (*analytics_pb.KpiResponse, error) {
	kpi, err := h.uc.GetManagerPerformance(ctx, req.ManagerId)
	if err != nil {
		return nil, err
	}

	return &analytics_pb.KpiResponse{
		ManagerId:   kpi.ManagerID(),
		FactRevenue: kpi.Fact().InexactFloat64(),
		PlanRevenue: kpi.Plan().InexactFloat64(),
		KpiPercent:  kpi.CalculateKPIPercent().InexactFloat64(),
		HasBonus:    kpi.HasBonus(),
	}, nil
}