package employee

import (
	"database/sql"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"skillpass-server-go/internal/rbac"
)

type Handler struct {
	svc   *Service
	rbacS *rbac.Service
}

func NewHandler(db *sql.DB, rbacS *rbac.Service) *Handler {
	return &Handler{svc: NewService(db), rbacS: rbacS}
}

func mustParseCompanyID(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.GetString("companyId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return uuid.Nil, false
	}
	return id, true
}

func (h *Handler) Create(c *gin.Context) {
	companyID, ok := mustParseCompanyID(c)
	if !ok {
		return
	}

	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	emp, err := h.svc.Create(c.Request.Context(), companyID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create employee"})
		return
	}

	c.JSON(http.StatusCreated, emp)
}

func (h *Handler) Get(c *gin.Context) {
	companyID, ok := mustParseCompanyID(c)
	if !ok {
		return
	}
	employeeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	requesterEmployeeID := c.GetString("employeeId")
	if requesterEmployeeID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Employee ID not found in context"})
		return
	}
	requesterUUID, err := uuid.Parse(requesterEmployeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid employee ID in context"})
		return
	}

	if employeeID != requesterUUID {
		// C6: requesting another employee's record. The requester must
		// have an elevated permission beyond `view_self` — otherwise any
		// employee could read PII (national ID, NPWP, bank account,
		// salary) of any other employee by passing the target ID.
		has, err := h.rbacS.HasAnyPermission(c.Request.Context(), requesterUUID, []string{
			"employee.view",
			"employee.view_team",
		})
		if err != nil {
			slog.Error("permission check failed", "requester", requesterUUID, "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Permission check failed"})
			return
		}
		if !has {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions to view other employees"})
			return
		}
	}

	emp, err := h.svc.Get(c.Request.Context(), companyID, employeeID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get employee"})
		return
	}

	c.JSON(http.StatusOK, emp)
}

// Invite creates a login account for an employee and returns a one-time
// temporary password for HR to share. Requires employee.update.
func (h *Handler) Invite(c *gin.Context) {
	companyID, ok := mustParseCompanyID(c)
	if !ok {
		return
	}
	employeeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	email, tempPassword, err := h.svc.InviteUser(c.Request.Context(), companyID, employeeID)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		case err == ErrAlreadyLinked:
			c.JSON(http.StatusConflict, gin.H{"error": "Employee already has a login account"})
		case err == ErrEmailTaken:
			c.JSON(http.StatusConflict, gin.H{"error": "A user with this email already exists"})
		default:
			slog.Error("invite failed", "employeeID", employeeID, "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create login"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"email": email, "tempPassword": tempPassword})
}

// GetMe returns the current user's own employee record (self-service).
// Available to any active company member; no extra permission needed.
func (h *Handler) GetMe(c *gin.Context) {
	companyID, ok := mustParseCompanyID(c)
	if !ok {
		return
	}
	employeeID, err := uuid.Parse(c.GetString("employeeId"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Employee ID not found in context"})
		return
	}
	emp, err := h.svc.Get(c.Request.Context(), companyID, employeeID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get employee"})
		return
	}
	c.JSON(http.StatusOK, emp)
}

// SelfUpdateRequest is the whitelist of fields an employee may edit on their
// own record. Employment terms (status, department, salary, etc.) are
// deliberately excluded — those remain HR-controlled.
type SelfUpdateRequest struct {
	Phone                    *string `json:"phone,omitempty"`
	DateOfBirth              *string `json:"dateOfBirth,omitempty"`
	Gender                   *string `json:"gender,omitempty"`
	MaritalStatus            *string `json:"maritalStatus,omitempty"`
	Address                  *string `json:"address,omitempty"`
	City                     *string `json:"city,omitempty"`
	Province                 *string `json:"province,omitempty"`
	PostalCode               *string `json:"postalCode,omitempty"`
	NationalID               *string `json:"nationalId,omitempty"`
	NPWP                     *string `json:"npwp,omitempty"`
	BankName                 *string `json:"bankName,omitempty"`
	BankAccountNumber        *string `json:"bankAccountNumber,omitempty"`
	BankAccountHolder        *string `json:"bankAccountHolder,omitempty"`
	EmergencyContactName     *string `json:"emergencyContactName,omitempty"`
	EmergencyContactPhone    *string `json:"emergencyContactPhone,omitempty"`
	EmergencyContactRelation *string `json:"emergencyContactRelation,omitempty"`
} //@name SelfUpdateRequest

// UpdateMe lets the current user edit their own personal fields (self-service).
func (h *Handler) UpdateMe(c *gin.Context) {
	companyID, ok := mustParseCompanyID(c)
	if !ok {
		return
	}
	employeeID, err := uuid.Parse(c.GetString("employeeId"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Employee ID not found in context"})
		return
	}

	var req SelfUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	upd := UpdateRequest{
		Phone:                    req.Phone,
		DateOfBirth:              req.DateOfBirth,
		Gender:                   req.Gender,
		MaritalStatus:            req.MaritalStatus,
		Address:                  req.Address,
		City:                     req.City,
		Province:                 req.Province,
		PostalCode:               req.PostalCode,
		NationalID:               req.NationalID,
		NPWP:                     req.NPWP,
		BankName:                 req.BankName,
		BankAccountNumber:        req.BankAccountNumber,
		BankAccountHolder:        req.BankAccountHolder,
		EmergencyContactName:     req.EmergencyContactName,
		EmergencyContactPhone:    req.EmergencyContactPhone,
		EmergencyContactRelation: req.EmergencyContactRelation,
	}
	emp, err := h.svc.Update(c.Request.Context(), companyID, employeeID, upd)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
		return
	}
	c.JSON(http.StatusOK, emp)
}

func (h *Handler) List(c *gin.Context) {
	companyID, ok := mustParseCompanyID(c)
	if !ok {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	status := c.Query("status")
	if status != "" {
		validStatuses := map[string]bool{
			"active":    true,
			"resigned":  true,
			"terminated": true,
			"on_leave":  true,
		}
		if !validStatuses[status] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
			return
		}
	}

	params := ListParams{
		CompanyID: companyID,
		Status:    status,
		Search:    c.Query("search"),
		Page:      page,
		PageSize:  pageSize,
	}

	if deptID := c.Query("departmentId"); deptID != "" {
		id, err := uuid.Parse(deptID)
		if err == nil {
			params.DepartmentID = &id
		}
	}
	if branchID := c.Query("branchId"); branchID != "" {
		id, err := uuid.Parse(branchID)
		if err == nil {
			params.BranchID = &id
		}
	}

	result, err := h.svc.List(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list employees"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) Update(c *gin.Context) {
	companyID, ok := mustParseCompanyID(c)
	if !ok {
		return
	}
	employeeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	emp, err := h.svc.Update(c.Request.Context(), companyID, employeeID, req)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update employee"})
		return
	}

	c.JSON(http.StatusOK, emp)
}
