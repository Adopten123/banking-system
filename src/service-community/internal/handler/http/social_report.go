package http

import (
	"encoding/json"
	"net/http"

	"github.com/Adopten123/banking-system/service-community/internal/domain"
)

// @Summary      File a report
// @Description  Creates a report for content (post, comment, message, user)
// @Tags         social
// @Accept       json
// @Produce      json
// @Param        X-User-ID header string true "UUID of the current user (mock authorization)"
// @Param        input body ReportRequest true "Report data"
// @Success      201 {object} domain.Report
// @Failure      400 "Bad request"
// @Failure      401 "Unauthorized"
// @Failure      500 "Internal server error"
// @Router       /api/v1/reports [post]
func (h *SocialHandler) report(w http.ResponseWriter, r *http.Request) {
	reporterID, err := getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req ReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	report, err := h.service.FileReport(r.Context(), domain.CreateReportInput{
		ReporterID: reporterID,
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
		Reason:     req.Reason,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(report)
}
