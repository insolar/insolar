package index

// RecentObjects contains data about last manipulations aroun objects
type RecentObjects struct {
	Fetched map[string]*RecentObjectsMeta
	Updated map[string]*RecentObjectsMeta
}

// RecentObjectsMeta contains meta-info about index
type RecentObjectsMeta struct {
	TTL int
}
