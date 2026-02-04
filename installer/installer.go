package installer

import (
	"darkarchinstall/fs"
	"fmt"
)

func Install(disk *string) {
	disktype, err := fs.GetDiskLabelType(*disk)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch disktype {
	case "gpt":

	}
}

/*
if gpt {
1 -> boot partition?
2 -> root part ?
3 -> swap ?
4 -> Homepart ?
} else if dos/mbr {
 root part ?
 swap ?
 homepart ?
}

*/
