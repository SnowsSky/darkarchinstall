package installer

import (
	"darkarchinstall/fs"
	"darkarchinstall/types"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

var RootPartition string = ""
var EFIPartition string = ""
var SwapPartition string = ""

// colors
var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

func Setup(disk *string, bootloader string, de []string, timezone *string, locale string, keymap string, hostname string, rootpasswd string, session_manager string, accounts []types.Account) {
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
	fmt.Println(Blue + "==>" + Reset + " Formating Partitions...")
	// format disk
	err := fs.FormatDisk(RootPartition, EFIPartition, SwapPartition)
	if err != nil {
		fmt.Println(Red+"==> ERROR"+Reset, err)
		os.Exit(1)
	}

	// mount partitions
	fmt.Println(Blue + "==>" + Reset + "Mouting Partitions...")
	err = fs.MountPartitions(RootPartition, EFIPartition, SwapPartition)
	if err != nil {
		fmt.Println(Red+"==> ERROR"+Reset, err)
		os.Exit(1)
	}

	fmt.Println(Blue + "==>" + Reset + "Installing Base System...")
	err = InstallBase(bootloader)
	if err != nil {
		fmt.Println(Red+"==> ERROR"+Reset, err)
		return
	}

	cmd := exec.Command("genfstab", "-U", "/mnt", ">>", "/mnt/etc/fstab")

	fmt.Println(Blue + "==>" + Reset + " Installing Full Desktop Environment...")
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(Red+"==> ERROR"+Reset, err)
		os.Exit(1)
	}
	err = InstallFullDE(de)
	if err != nil {
		fmt.Println(Red+"==> ERROR"+Reset, err)
		return
	}
	err = InstallSessionManager(session_manager)
	if err != nil {
		fmt.Println(Red+"==> ERROR"+Reset, err)
		return
	}

	//after chroot
	fmt.Println(Blue + "==>" + Reset + " System Configuration...")
	err = AddDarkArchRepos()
	if err != nil {
		fmt.Println(Red+"==> ERROR"+Reset, err)
		os.Exit(1)
	}
	EditOSRelease()
	SetTime(timezone)
	SetLocalisation(locale)
	SetKeymap(keymap)
	SetHostname(hostname)
	// set root password
	cmd = exec.Command("arch-chroot", "/mnt", "chpasswd")
	cmd.Stdin = strings.NewReader("root:" + rootpasswd + "\n")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(Red+"==> ERROR"+Reset, err)
		os.Exit(1)
	}
	err = SetupAccounts(accounts)
	if err != nil {
		fmt.Println(Red+"==> ERROR"+Reset, err)
		os.Exit(1)
	}

	fmt.Println(Blue + "==>" + Reset + " Installing Bootloader...")
	err = SetupBootloader(bootloader, disk)
	if err != nil {
		fmt.Println(Red+"==> ERROR"+Reset, err)
		return
	}
	fmt.Println(Blue + "==>" + Reset + " Installing Extra Feature...")
	err = InstallBlackArchRepos()
	if err != nil {
		fmt.Println(Red+"==> ERROR"+Reset, err)
		os.Exit(1)
	}
	InstallExtraPackages()
	fmt.Println(Blue + "==>" + Reset + " Enabling Services...")
	err = EnableServices(session_manager)
	if err != nil {
		fmt.Println(Red+"==> ERROR"+Reset, err)
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

	if _, err := f.WriteString("\n[darkarch]\nSigLevel = Optional TrustAll\nServer = https://raw.githubusercontent.com/darkarchlinux/DarkArchPackages/main/main/binaries/$arch\n#[darkarch-unstable]\n#SigLevel = Optional TrustAll\n#Server = https://raw.githubusercontent.com/darkarchlinux/DarkArchPackages/main/unstable/binaries/$arch\n"); err != nil {
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

func EnableServices(session_manager string) error {
	cmd := exec.Command("arch-chroot", "/mnt", "systemctl", "enable", "NetworkManager")
	err := cmd.Run()
	if err != nil {
		return err
	}
	cmd = exec.Command("arch-chroot", "/mnt", "systemctl", "enable", session_manager)
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

func SetupBootloader(bootloader string, disk *string) error {
	switch bootloader {
	case "grub":
		if EFIPartition == "" {
			cmd := exec.Command("arch-chroot", "/mnt", "grub-install", "--target=i386-pc", *disk)

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
		} else {
			cmd := exec.Command("arch-chroot", "/mnt", "grub-install", "--target=x86_64-efi", "--efi-directory=/boot", "--bootloader-id=GRUB")

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
		cmd.Stdout = os.Stdout
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
	result := strings.Split(locale, " ")[0]
	f, err := os.OpenFile("/mnt/etc/locale.gen", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(locale + "\n"); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	cmd := exec.Command("arch-chroot", "/mnt", "locale-gen")

	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	data := []byte("LANG=" + result)
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

func InstallFullDE(de []string) error {
	if slices.Contains(de, "xfce4") {
		cmd := exec.Command("arch-chroot", "/mnt", "pacman", "-S", "--noconfirm", "xfce4", "xfce4-goodies", "mugshot", "gvfs")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	if slices.Contains(de, "plasma") {
		cmd := exec.Command("arch-chroot", "/mnt", "pacman", "-S", "--noconfirm", "plasma", "konsole", "packagekit")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	if slices.Contains(de, "gnome") {
		cmd := exec.Command("arch-chroot", "/mnt", "pacman", "-S", "--noconfirm", "gnome", "gnome-terminal", "gnome-browser-connector", "packagekit", "gnome-tweaks")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func InstallSessionManager(session_manager string) error {
	switch session_manager {
	case "lightdm":
		cmd := exec.Command("arch-chroot", "/mnt", "pacman", "-S", "--noconfirm", "lightdm", "lightdm-gtk-greeter")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			return err
		}
	case "sddm":
		cmd := exec.Command("arch-chroot", "/mnt", "pacman", "-S", "--noconfirm", "sddm")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			return err
		}
	case "gdm":
		cmd := exec.Command("arch-chroot", "/mnt", "pacman", "-S", "--noconfirm", "gdm")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func InstallBase(bootloader string) error {
	packages := []string{
		"base",
		"linux-hardened",
		"linux-hardened-headers",
		"linux-firmware",
		"efibootmgr",
		"sudo",
		"vim",
		"networkmanager",
		"htop",
		"firefox",
		bootloader,
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
