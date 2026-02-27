package usecase

import (
	"fmt"
	"time"

	"github.com/versoit/diploma/be/orders/internal/domain"
	"github.com/xuri/excelize/v2"
)

const (
	sheetDashboard = "Dashboard"
	sheetOrders    = "Orders"
)

type ReportGenerator struct{}

func NewReportGenerator() *ReportGenerator {
	return &ReportGenerator{}
}

func (g *ReportGenerator) Generate(analytics *AnalyticsResult, orders []*domain.Order) ([]byte, error) {
	f := excelize.NewFile()
	defer func() {
		_ = f.Close()
	}()

	// 1. Dashboard Sheet (Standard approach for small data)
	if err := g.createDashboardSheet(f, analytics); err != nil {
		return nil, fmt.Errorf("dashboard sheet: %w", err)
	}

	// 2. Orders Sheet (Streamwriter approach for potentially large data)
	if err := g.createOrdersSheet(f, orders); err != nil {
		return nil, fmt.Errorf("orders sheet: %w", err)
	}

	// Clean up default sheet
	if f.GetSheetName(0) == "Sheet1" {
		if err := f.DeleteSheet("Sheet1"); err != nil {
			return nil, fmt.Errorf("delete sheet1: %w", err)
		}
	}

	// Set active sheet
	index, err := f.GetSheetIndex(sheetDashboard)
	if err != nil {
		return nil, fmt.Errorf("get sheet index: %w", err)
	}
	f.SetActiveSheet(index)

	// Buffer output
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to write excel buffer: %w", err)
	}

	return buf.Bytes(), nil
}

func (g *ReportGenerator) createDashboardSheet(f *excelize.File, analytics *AnalyticsResult) error {
	_, err := f.NewSheet(sheetDashboard)
	if err != nil {
		return err
	}

	// Define Styles
	styleTitle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 18, Color: "#1F2937"},
	})
	if err != nil {
		return err
	}

	styleHeader, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "#374151"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#E5E7EB"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "left", Indent: 1},
		Border:    []excelize.Border{{Type: "bottom", Color: "#D1D5DB", Style: 1}},
	})
	if err != nil {
		return err
	}

	styleCurrency, err := f.NewStyle(&excelize.Style{
		NumFmt:    4, // #,##0.00
		Font:      &excelize.Font{Bold: true, Color: "#059669"},
		Alignment: &excelize.Alignment{Horizontal: "right"},
	})
	if err != nil {
		return err
	}

	styleValue, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "right"},
	})
	if err != nil {
		return err
	}

	// Title
	if err := f.SetCellValue(sheetDashboard, "A1", "Sales Performance Report"); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetDashboard, "A1", "A1", styleTitle); err != nil {
		return err
	}
	if err := f.MergeCell(sheetDashboard, "A1", "D1"); err != nil {
		return err
	}

	// Metadata
	metaRow := 2
	if err := f.SetCellValue(sheetDashboard, fmt.Sprintf("A%d", metaRow), "Generated at:"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheetDashboard, fmt.Sprintf("B%d", metaRow), time.Now().Format(time.RFC1123)); err != nil {
		return err
	}

	// Summary
	summaryStartRow := 4
	if err := f.SetCellValue(sheetDashboard, fmt.Sprintf("A%d", summaryStartRow), "Kpi Metrics"); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetDashboard, fmt.Sprintf("A%d", summaryStartRow), fmt.Sprintf("B%d", summaryStartRow), styleHeader); err != nil {
		return err
	}

	metrics := []struct {
		Label string
		Value interface{}
		Style int
	}{
		{"Total Orders", analytics.OrdersCount, styleValue},
		{"Total Revenue", analytics.TotalRevenue, styleCurrency},
		{"Average Order Value", analytics.AvgCheck, styleCurrency},
		{"Current Cooking", analytics.CookingCount, styleValue},
		{"Current Delivering", analytics.DeliveringCount, styleValue},
	}

	for i, m := range metrics {
		row := summaryStartRow + 1 + i
		_ = f.SetCellValue(sheetDashboard, fmt.Sprintf("A%d", row), m.Label)
		_ = f.SetCellValue(sheetDashboard, fmt.Sprintf("B%d", row), m.Value)
		_ = f.SetCellStyle(sheetDashboard, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), m.Style)
	}

	// Top Products
	tpStartRow := summaryStartRow + len(metrics) + 2
	if err := f.SetCellValue(sheetDashboard, fmt.Sprintf("A%d", tpStartRow), "Top Selling Products"); err != nil {
		return err
	}
	if err := f.SetCellStyle(sheetDashboard, fmt.Sprintf("A%d", tpStartRow), fmt.Sprintf("C%d", tpStartRow), styleHeader); err != nil {
		return err
	}

	headers := []string{"Product", "Sold (Qty)", "Revenue"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, tpStartRow+1)
		_ = f.SetCellValue(sheetDashboard, cell, h)
		_ = f.SetCellStyle(sheetDashboard, cell, cell, styleHeader)
	}

	for i, p := range analytics.TopProducts {
		row := tpStartRow + 2 + i
		_ = f.SetCellValue(sheetDashboard, fmt.Sprintf("A%d", row), p.Name)
		_ = f.SetCellValue(sheetDashboard, fmt.Sprintf("B%d", row), p.Count)
		_ = f.SetCellValue(sheetDashboard, fmt.Sprintf("C%d", row), p.Revenue)
		_ = f.SetCellStyle(sheetDashboard, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), styleCurrency)
	}

	// Adjust column widths
	_ = f.SetColWidth(sheetDashboard, "A", "A", 30)
	_ = f.SetColWidth(sheetDashboard, "B", "B", 25)
	_ = f.SetColWidth(sheetDashboard, "C", "C", 20)

	return nil
}

func (g *ReportGenerator) createOrdersSheet(f *excelize.File, orders []*domain.Order) error {
	_, err := f.NewSheet(sheetOrders)
	if err != nil {
		return err
	}

	sw, err := f.NewStreamWriter(sheetOrders)
	if err != nil {
		return err
	}

	// Styles in StreamWriter are applied via cell values
	styleHeader, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#4F46E5"}, Pattern: 1},
	})
	if err != nil {
		return err
	}

	// Header row
	headers := []interface{}{
		excelize.Cell{StyleID: styleHeader, Value: "Order ID"},
		excelize.Cell{StyleID: styleHeader, Value: "Created At"},
		excelize.Cell{StyleID: styleHeader, Value: "Status"},
		excelize.Cell{StyleID: styleHeader, Value: "Customer"},
		excelize.Cell{StyleID: styleHeader, Value: "Address"},
		excelize.Cell{StyleID: styleHeader, Value: "Amount"},
		excelize.Cell{StyleID: styleHeader, Value: "Promo Code"},
		excelize.Cell{StyleID: styleHeader, Value: "Items Summary"},
	}

	if err := sw.SetRow("A1", headers); err != nil {
		return err
	}

	// Data rows
	for i, o := range orders {
		rowNum := i + 2
		cell, _ := excelize.CoordinatesToCellName(1, rowNum)

		var items []string
		for _, it := range o.Items() {
			items = append(items, fmt.Sprintf("%dx %s", it.Quantity(), it.ProductName()))
		}
		
		itemsStr := ""
		for j, s := range items {
			if j > 0 {
				itemsStr += ", "
			}
			itemsStr += s
		}

		addr := fmt.Sprintf("%s, %s %s", o.Address().City, o.Address().Street, o.Address().House)

		rowData := []interface{}{
			o.OrderNumber(),
			o.CreatedAt().Format("2006-01-02 15:04:05"),
			o.Status().String(),
			o.CustomerID(),
			addr,
			o.FinalPrice().InexactFloat64(),
			o.PromoCode(),
			itemsStr,
		}

		if err := sw.SetRow(cell, rowData); err != nil {
			return err
		}
	}

	if err := sw.Flush(); err != nil {
		return err
	}

	// Set column widths (must be done after StreamWriter for the sheet is finished or before)
	_ = f.SetColWidth(sheetOrders, "A", "A", 20)
	_ = f.SetColWidth(sheetOrders, "B", "B", 25)
	_ = f.SetColWidth(sheetOrders, "C", "C", 15)
	_ = f.SetColWidth(sheetOrders, "D", "D", 20)
	_ = f.SetColWidth(sheetOrders, "E", "E", 40)
	_ = f.SetColWidth(sheetOrders, "F", "F", 15)
	_ = f.SetColWidth(sheetOrders, "G", "G", 15)
	_ = f.SetColWidth(sheetOrders, "H", "H", 60)

	return nil
}

