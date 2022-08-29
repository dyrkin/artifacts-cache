package file

import "testing"

func Test_RemoveCwd(t *testing.T) {
	cwd := "/home/user/pkg/repository"
	path := "/home/user/pkg/repository/hello/world.txt"
	if RemoveCwd(cwd, path) != "hello/world.txt" {
		t.Errorf("RemoveCwd(%s, %s) = %s, want %s", cwd, path, RemoveCwd(cwd, path), "hello/world.txt")
	}

	cwd = "/home/user/pkg/repository/"
	path = "/home/user/pkg/repository/hello/world.txt"
	if RemoveCwd(cwd, path) != "hello/world.txt" {
		t.Errorf("RemoveCwd(%s, %s) = %s, want %s", cwd, path, RemoveCwd(cwd, path), "hello/world.txt")
	}
}
