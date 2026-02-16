package handler

import (
	"context"
	"errors"
	"log/slog"

	"github.com/shopspring/decimal"
	analyticspb "github.com/versoit/diploma/services/analytics/api/proto/pb"
	"github.com/versoit/diploma/services/analytics/internal/domain"
	"github.com/versoit/diploma/services/analytics/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AnalyticsHandler struct {
	analyticspb.UnimplementedKpiServiceServer
	uc  *usecase.AnalyticsUseCase
	log *slog.Logger
}

func NewAnalyticsHandler(uc *usecase.AnalyticsUseCase, log *slog.Logger) *AnalyticsHandler {
	return &AnalyticsHandler{
		uc:  uc,
		log: log,
	}
}

func (h *AnalyticsHandler) Register(server *grpc.Server) {
	analyticspb.RegisterKpiServiceServer(server, h)
}

func (h *AnalyticsHandler) GetManagerKPI(ctx context.Context, req *analyticspb.KpiRequest) (*analyticspb.KpiResponse, error) {
	h.log.Info("Fetching manager KPI", slog.String("manager_id", req.ManagerId))

	kpi, err := h.uc.GetManagerPerformance(ctx, req.ManagerId)
	if err != nil {
		if errors.Is(err, domain.ErrKPINotFound) {
			return nil, status.Errorf(codes.NotFound, "KPI not found: %v", err)
		}
		h.log.Error("Failed to fetch KPI", slog.String("manager_id", req.ManagerId), slog.Any("error", err))
		return nil, status.Errorf(codes.Internal, "get performance: %v", err)
	}

	return &analyticspb.KpiResponse{
		ManagerId:   kpi.ManagerID(),
		FactRevenue: kpi.Fact().InexactFloat64(),
		PlanRevenue: kpi.Plan().InexactFloat64(),
		KpiPercent:  kpi.CalculateKPIPercent().InexactFloat64(),
		HasBonus:    kpi.HasBonus(),
	}, nil
}

func (h *AnalyticsHandler) RecordSale(ctx context.Context, req *analyticspb.RecordSaleRequest) (*analyticspb.RecordSaleResponse, error) {
	h.log.Info("Recording sale", slog.String("manager_id", req.ManagerId), slog.Float64("amount", req.Amount))

	err := h.uc.RecordSale(ctx, req.ManagerId, decimal.NewFromFloat(req.Amount))
	if err != nil {
		h.log.Error("Failed to record sale", slog.String("manager_id", req.ManagerId), slog.Any("error", err))
		return nil, status.Errorf(codes.Internal, "record sale: %v", err)
	}

	return &analyticspb.RecordSaleResponse{Success: true}, nil
}
