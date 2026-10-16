package routes

func BuildRoute(method string, baseUrl string, path string) string {
	return method + " " + baseUrl + path
}
