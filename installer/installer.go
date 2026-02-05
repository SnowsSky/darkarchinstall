package installer

import (
	"darkarchinstall/fs"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh/spinner"
)

var RootPartition string = ""
var EFIPartition string = ""
var SwapPartition string = ""

func Setup(disk *string, bootloader string, de []string) {
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
	action := func() {
		// format disk
		err := fs.FormatDisk(RootPartition, EFIPartition, SwapPartition)
		if err != nil {
			fmt.Println("Failed", err)
			os.Exit(1)
		}
	}
	if err := spinner.New().Title("Formating Partitions...").Action(action).Run(); err != nil {
		fmt.Println("Failed:", err)
		os.Exit(1)
	}
	// mount partitions
	action = func() {
		// format disk
		err := fs.MountPartitions(RootPartition, EFIPartition, SwapPartition)
		if err != nil {
			fmt.Println("Failed", err)
			os.Exit(1)
		}
	}

	if err := spinner.New().Title("Mouting Partitions...").Action(action).Run(); err != nil {
		fmt.Println("Failed:", err)
		os.Exit(1)
	}

}

func Install(bootloader string, de []string) {
	to_install := strings.Builder{}
	to_install.WriteString("base linux linux-firmware efibootmgr lightdm")
	to_install.WriteString(" " + bootloader)

	if len(de) > 0 {
		to_install.WriteString(" ")
		to_install.WriteString(strings.Join(de, " "))
	}
	//to_install.String()
	cmd := exec.Command("pacstrap", "-K", "/mnt", to_install.String())

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
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
