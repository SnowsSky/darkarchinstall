package installer

import (
	"darkarchinstall/fs"
	"fmt"
)

var RootPartition string = ""
var EFIPartition string = ""
var SwapPartition string = ""

func Install(disk *string) {
	partitions := fs.GetPartitionOfDisk(*disk)
	for _, partition := range partitions {
		_, err := fs.GetPartitionType(partition)
		if err != nil {
			fmt.Println(err)
		}
		if fs.PartitionTypeLinuxFileSystem == 1 {
			RootPartition = partition
		}
		if fs.PartitionTypeEFI == 1 {
			EFIPartition = partition

		}
		if fs.PartitionTypeLinuxSwap == 1 {
			SwapPartition = partition
		}

	}

	// format disk
	err := fs.FormatDisk(RootPartition, EFIPartition, SwapPartition)
	if err != nil {
		fmt.Println(err)
		return
	}
}

/*
1 -> boot partition (if gpt)? ✅
2 -> root part ? ✅
3 -> swap ?✅
4 -> Homepart ? ❌


*/
