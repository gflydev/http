package http

const (
	// ====================================================================
	// ===================== HTTP Context Constants =======================
	// ====================================================================

	// DataPathID key in Context's Data for ID extracted from path parameter
	DataPathID string = "__path_id__"
	// DataRequest key in Context's Data for raw request data
	DataRequest string = "__request__"
	// DataFilter key in Context's Data for filtering parameters
	DataFilter string = "__filter__"
)
