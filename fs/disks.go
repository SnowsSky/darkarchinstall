package fs

import (
	"os"
	"os/exec"
	"path"

	"github.com/diskfs/go-diskfs"
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
	cmd := exec.Command("cfdisk", disk)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}
	cmd = exec.Command("clear")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	return nil
}

func GetDiskLabelType(disktoopen string) (string, error) {
	disk, err := diskfs.Open(disktoopen)
	if err != nil {
		return "", err
	}
	p, _ := disk.GetPartitionTable()
	return p.Type(), nil
}

func GetDiskFilesystems(disktocheck string) {

}
