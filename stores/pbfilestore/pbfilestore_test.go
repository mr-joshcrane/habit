package pbfilestore_test

import (
	"habit/stores/pbfilestore"
	"os"
	"testing"
)

func TestOpenReturnsEmptyFileStoreIfFileNotExists(t *testing.T) {
	t.Parallel()
	path := t.TempDir() + "/path_does_not_exist"
	_, err := pbfilestore.Open(path)
	if err != nil {
		t.Fatalf("Open incorrectly errored: %t", err)
	}
}

func TestOpenReturnsErrorIfInsufficientPermissions(t *testing.T) {
	t.Parallel()
	path := t.TempDir() + "/insufficient_perms"
	_, err := os.Create(path)
	if err != nil {
		t.Fatalf("Error creating test file")
	}
	err = os.Chmod(path, 0200)
	if err != nil {
		t.Fatalf("Unable to set perms on file")
	}
	_, err = pbfilestore.Open(path)
	if err == nil {
		t.Fatalf("Open incorrectly errored: %t", err)
	}
}