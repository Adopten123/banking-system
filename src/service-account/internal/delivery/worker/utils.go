package worker

func resolveEntityType(typeID int32) string {
	switch typeID {
	case 1:
		return "account"
	case 2:
		return "card"
	default:
		return "account"
	}
}