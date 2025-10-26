package http

const (
	// ====================================================================
	// ======================== Common Constants ==========================
	// ====================================================================

	// UserKey key in Context's Data for ID extracted from path parameter
	UserKey string = "__user__"

	// DataKey key in Context's Data for processed/transformed request data
	DataKey string = "__data__"

	// ====================================================================
	// ===================== HTTP Context Constants =======================
	// ====================================================================

	// PathIDKey key in Context's Data for ID extracted from path parameter
	PathIDKey string = "__path_id__"
	// RequestKey key in Context's Data for raw request data
	RequestKey string = "__request__"
	// FilterKey key in Context's Data for filtering parameters
	FilterKey string = "__filter__"
)
