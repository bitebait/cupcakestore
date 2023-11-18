package controllers

const baseLayout = "layouts/base"
const storeLayout = "layouts/store"

func selectLayout(isStaff bool, isUserProfile bool) string {
	switch {
	case isStaff:
		return baseLayout
	case isUserProfile:
		return storeLayout
	default:
		return ""
	}
}

func selectRedirectPath(isStaff bool) string {
	if isStaff {
		return "/users"
	}
	return "/"
}
