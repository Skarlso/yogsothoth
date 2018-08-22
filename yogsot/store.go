package yogsot

// Storer defines a storage medium.
//
// DO lacks the ability to tag all of there resources, this is why
// a storage facility is needed.
type Storer interface {
	saveAllResourcesForStack(stackname string, resources []Resource) error
	loadAllResourcesForStack(stackname string) ([]Resource, error)
	deleteAllResourcesForStack(stackname string) error
	stackExits(stackname string) (ok bool, err error)
}

// SQLiteStore is an implementation of Storer with a SQLite backed
// storage facility.
type SQLiteStore struct{}

// SaveAllResourcesForStack will save all resources that belong to a
// stackname. Thereby making possible to handle them together as one.
func (s *SQLiteStore) saveAllResourcesForStack(stackname string, resources []Resource) error {
	return nil
}

// LoadAllResourcesForStack will gather all the resources belonging to
// a stack in order to display them in a status table for example.
func (s *SQLiteStore) loadAllResourcesForStack(stackname string) ([]Resource, error) {
	resources := make([]Resource, 0)
	return resources, nil
}

// DeleteAllResourcesForStack will remove all resources for a stackname.
// This should only be called after the stack has been deleted.
func (s *SQLiteStore) deleteAllResourcesForStack(stackname string) error {
	return nil
}

// StackExists will return true if a given stackname has existing
// resources. This should be used as a quick check as it only returns
// a boolean.
func (s *SQLiteStore) stackExists(stackname string) (bool, error) {
	return false, nil
}
