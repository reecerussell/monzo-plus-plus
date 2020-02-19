package repository

type PermissionsRepository interface {
	LoadCollections() map[int][]string
}
