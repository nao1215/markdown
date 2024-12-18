package arch

// Icon is the icon for a group.
// By default, architecture diagram supports the following icons: cloud, database, disk, internet, server.
// TODO: Add AWS icons.
type Icon string

const (
	// IconCloud is the cloud icon.
	IconCloud Icon = "cloud"
	// IconDatabase is the database icon.
	IconDatabase Icon = "database"
	// IconDisk is the disk icon.
	IconDisk Icon = "disk"
	// IconInternet is the internet icon.
	IconInternet Icon = "internet"
	// IconServer is the server icon.
	IconServer Icon = "server"
)
