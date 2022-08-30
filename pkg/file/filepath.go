package file

func RemoveCwd(cwd string, path string) string {
	if cwd[len(cwd)-1] == '/' {
		return path[len(cwd):]
	}
	return path[len(cwd)+1:]
}
