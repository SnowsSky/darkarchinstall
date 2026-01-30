package fs

import (
	"os"
	"os/exec"
	"path"
)

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
	err := exec.Command("cfdisk", disk).Run()
	if err != nil {
		return err
	}
	return nil
}
