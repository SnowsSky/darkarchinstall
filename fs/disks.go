package fs

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path"
	"strings"
)

type PartitionType int

const (
	PartitionTypeEFI PartitionType = iota
	PartitionTypeLinuxFileSystem
	PartitionTypeLinuxSwap
)

type ListBlockPayload struct {
	BlockDevices []struct{
		PartitionTypeName string `json:"parttypename"`
	} `json:"blockdevices"`
}

func GetDisks() ([]string, error) {
	var disks []string

	disksAndPartitions, err := os.ReadDir("/sys/block")
	if err != nil {
		return disks, err
	}
	for _, diskOrPartition := range disksAndPartitions {
		_, err := os.Stat(path.Join("/sys/block", diskOrPartition.Name(), "device"))
		if err == nil {
			disks = append(disks, path.Join("/dev", diskOrPartition.Name()))
		}
	}

	return disks, nil
}

func EditDisk(disk string) error {
	cmd := exec.Command("cfdisk", disk)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}
	cmd = exec.Command("clear")

	cmd.Stdout = os.Stdout
	cmd.Run()
	return nil
}

func DetectDiskLayout() string {
	return ""
}

func GetPartitionType(part string) (PartitionType, error) {
	ret, err := exec.Command("lsblk", "-JO", part).Output()
	if err != nil {
		return -1, err
	}

	var payload ListBlockPayload

	err = json.Unmarshal(ret, &payload); if err != nil {
		return -1, err
	}

	switch(strings.ToLower(payload.BlockDevices[0].PartitionTypeName)) {
	case "efi system":
		return PartitionTypeEFI, nil
	case "linux filesystem":
		return PartitionTypeLinuxFileSystem, nil
	case "linux swap":
		return PartitionTypeLinuxSwap, nil
	default:
		return -2, errors.New("unsupported partition type")
	}
}
