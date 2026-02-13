package installer

import (
	"darkarchinstall/fs"
	"darkarchinstall/types"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh/spinner"
)

var RootPartition string = ""
var EFIPartition string = ""
var SwapPartition string = ""

func Setup(disk *string, bootloader string, de []string, timezone *string, locale string, keymap string, hostname string, rootpasswd string, accounts []types.Account) {
	partitions := fs.GetPartitionOfDisk(*disk)
	for _, partition := range partitions {
		Parttype, err := fs.GetPartitionType(partition)
		if err != nil {
			fmt.Println(err)
		}
		switch Parttype {
		case fs.PartitionTypeLinuxFileSystem:
			RootPartition = partition
		case fs.PartitionTypeEFI:
			EFIPartition = partition
		case fs.PartitionTypeLinuxSwap:
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
	action = func() {
		err := Install(bootloader, de)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if err := spinner.New().Title("Installing DarkArch...").Action(action).Run(); err != nil {
		fmt.Println("Failed:", err)
		os.Exit(1)
	}
	cmd := exec.Command("genfstab", "-U", "/mnt", ">>", "/mnt/etc/fstab")

	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//after chroot
	action = func() {
		err := AddDarkArchRepos()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		EditOSRelease()
		SetTime(timezone)
		SetLocalisation(locale)
		SetKeymap(keymap)
		SetHostname(hostname)
		// set root password
		cmd = exec.Command("chpasswd")
		cmd.Stdin = strings.NewReader("root:" + rootpasswd + "\n")
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = SetupAccounts(accounts)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	if err := spinner.New().Title("System Configuration...").Action(action).Run(); err != nil {
		fmt.Println("Failed:", err)
		os.Exit(1)
	}
	action = func() {
		err := SetupBootloader(bootloader)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if err := spinner.New().Title("Installing bootloader...").Action(action).Run(); err != nil {
		fmt.Println("Failed:", err)
		os.Exit(1)
	}
	action = func() {
		err := InstallBlackArchRepos()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		InstallExtraPackages()
	}
	if err := spinner.New().Title("Installing extras features...").Action(action).Run(); err != nil {
		fmt.Println("Failed:", err)
		os.Exit(1)
	}
	action = func() {
		err := EnableServices()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	if err := spinner.New().Title("Enabling services...").Action(action).Run(); err != nil {
		fmt.Println("Failed:", err)
		os.Exit(1)
	}

	ExitInstall()
}

func AddDarkArchRepos() error {
	f, err := os.OpenFile("/mnt/etc/pacman.conf", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString("\n[darkarch]\nSigLevel = Optional TrustAll\nServer = https://raw.githubusercontent.com/darkarchlinux/DarkArchPackages/main/main/binaries/$arch\n#[darkarch-unstable]\n#SigLevel = Optional TrustAll\n#Server = https://raw.githubusercontent.com/darkarchlinux/DarkArchPackages/main/unstable/binaries/$arch"); err != nil {
		return err
	}
	cmd := exec.Command("arch-chroot", "/mnt", "pacman", "-Sy")

	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil

}

func EnableServices() error {
	cmd := exec.Command("arch-chroot", "/mnt", "systemctl", "enable", "NetworkManager")
	err := cmd.Run()
	if err != nil {
		return err
	}
	cmd = exec.Command("arch-chroot", "/mnt", "systemctl", "enable", "ligthdm")
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func ExitInstall() {
	cmd := exec.Command("umount", "-R", "/mnt")
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	cmd = exec.Command("reboot")
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func InstallBlackArchRepos() error {
	cmd := exec.Command("curl", "-O", "https://blackarch.org/strap.sh")
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("cp", "strap.sh", "/mnt/root/strap.sh")
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	cmd = exec.Command("arch-chroot", "/mnt", "chmod", "+x", "/root/strap.sh")
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("arch-chroot", "/mnt", "/root/strap.sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func SetupBootloader(bootloader string) error {
	if bootloader == "grub" {
		exec.Command("mkdir", "-p", "/mnt/boot/efi").Run()
		cmd := exec.Command("arch-chroot", "/mnt", "grub-install", "--target=x86_64-efi", "--efi-directory=/boot/efi", "--bootloader-id=GRUB")

		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			return err
		}
		cmd = exec.Command("arch-chroot", "/mnt", "grub-mkconfig", "-o", "/boot/grub/grub.cfg")

		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			return err
		}

	}
	return nil
}

func SetupAccounts(accounts []types.Account) error {
	for _, account := range accounts {
		if !account.SudoPerms {
			cmd := exec.Command("arch-chroot", "/mnt", "useradd", "-m", account.Username)

			cmd.Stderr = os.Stderr

			err := cmd.Run()
			if err != nil {
				return err
			}
		} else {
			cmd := exec.Command("arch-chroot", "/mnt", "useradd", "-mG", "wheel", account.Username)

			cmd.Stderr = os.Stderr

			err := cmd.Run()
			if err != nil {
				return err
			}
		}
		cmd := exec.Command("arch-chroot", "/mnt", "chpasswd")
		cmd.Stdin = strings.NewReader(account.Username + ":" + account.Password + "\n")
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil

}

func EditOSRelease() {
	data := []byte("NAME=\"DarkArch Linux\"\nPRETTY_NAME=\"DarkArch Linux\"\nID=darkarch\nID_LIKE=arch\nAINSI_COLOR=\"0;31\"\nHOME_URL=\"https://github.com/kam/darkarch\"")
	err := os.WriteFile("/mnt/etc/os-release", data, 0644)
	if err != nil {
		fmt.Println("Error writing os-release:", err)
		return
	}
}

func SetTime(timezone *string) {
	cmd := exec.Command("arch-chroot", "/mnt", "ln", "-sf", fmt.Sprintf("/usr/share/zoneinfo/%s", *timezone), "/etc/localtime")

	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	cmd = exec.Command("arch-chroot", "/mnt", "hwclock", "--systohc")

	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func SetLocalisation(locale string) {
	f, err := os.OpenFile("/mnt/etc/locale.gen", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString("\n" + locale); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	cmd := exec.Command("arch-chroot", "/mnt", "locale-gen")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	data := []byte("LANG=" + locale)
	err = os.WriteFile("/mnt/etc/locale.conf", data, 0644)
	if err != nil {
		fmt.Println("Error writing locales:", err)
		return
	}

}

func SetKeymap(keymap string) {
	data := []byte("KEYMAP=" + keymap)
	err := os.WriteFile("/mnt/etc/vconsole.conf", data, 0644)
	if err != nil {
		fmt.Println("Error writing locales:", err)
		return
	}
}

func SetHostname(hostname string) {
	data := []byte(hostname)
	err := os.WriteFile("/mnt/etc/hostname", data, 0644)
	if err != nil {
		fmt.Println("Error writing locales:", err)
		return
	}
}

func Install(bootloader string, de []string) error {
	packages := []string{
		"base",
		"linux",
		"linux-firmware",
		"efibootmgr",
		"lightdm",
		"lightdm-gtk-greeter",
		"sudo",
		"vim",
		"networkmanager",
		"htop",
		"firefox",
		bootloader,
	}

	if len(de) > 0 {
		packages = append(packages, de...)
	}

	args := append([]string{"-K", "/mnt"}, packages...)

	cmd := exec.Command("pacstrap", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func InstallExtraPackages() {
	// install from darkarchrepos
	cmd := exec.Command("arch-chroot", "/mnt", "pacman", "-S", "--noconfirm", "yay", "snowfetch")

	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
