package types

type Accounts struct {
	Username  string
	Password  string
	SudoPerms bool
}

type DiskConfig struct {
	DiskPath               string
	HasSwap                bool
	SwapSize               uint64 // Swap size in megabytes
	IsUEFI                 bool
	BootPartSize           uint64 // Boot Partition size in megabytes
	HomePartitionSeparated bool
	RootSize               uint64 // Root partition size in megabytes
	HomePartitionSize      uint64 // Home partition size in megabytes
}
