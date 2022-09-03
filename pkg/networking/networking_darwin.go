package networking

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/utkarsh-pro/kindli/pkg/models"
	"github.com/utkarsh-pro/kindli/pkg/sh"
	"github.com/utkarsh-pro/kindli/pkg/utils"
)

const disableIPv6 = true

func Setup(vmName string) error {
	logrus.Info("Setting up inside the VM...")
	if err := setupPacketRoutingInsideVM(vmName); err != nil {
		return fmt.Errorf("failed to setup packet routing inside VM: %s", err)
	}
	logrus.Info("✅ Completed setup inside the VM")

	logrus.Info("Setting up on the host...")
	if err := setupPacketRoutingOnHost(vmName); err != nil {
		return fmt.Errorf("failed to setup packet routing on host: %s", err)
	}
	logrus.Info("✅ Completed setup on the host")

	logrus.Info("Waiting for SIGINT (Ctr + C) to cleanup...")

	utils.SigIntHandler(func() {
		if err := Cleanup(vmName); err != nil {
			logrus.Error("failed to cleanup networking: ", err)
			logrus.Warn("please cleanup networking manually - `kindli network cleanup --vm-name <vm-name>`")
		}
	})

	return nil
}

func Cleanup(vmName string) error {
	vm := models.NewVM(vmName, "", 0)
	if err := vm.GetByName(); err != nil {
		return fmt.Errorf("failed to get VM by name: %w", err)
	}

	limaVMIPv4 := vm.GetVMIPv4()
	ipv4Subnetprefix, err := GetIPv4SubnetPrefix("kind")
	if err != nil {
		return fmt.Errorf("failed to get IPv4 subnet prefix: %w", err)
	}
	if err := sh.RunSilent(fmt.Sprintf("sudo route -nv delete -net %s %s", ipv4Subnetprefix, limaVMIPv4)); err != nil {
		return fmt.Errorf("failed to cleanup route from system to VM: %s", err)
	}

	if !disableIPv6 {
		limaVMIPv6 := vm.GetVMIPv6()
		ipv6Subnetprefix, err := GetIPv6SubnetPrefix("kind")
		if err != nil {
			return fmt.Errorf("failed to get IPv6 subnet prefix: %w", err)
		}
		if err := sh.RunSilent(fmt.Sprintf("sudo route -nv delete -inet6 %s:: %s", ipv6Subnetprefix, limaVMIPv6)); err != nil {
			return fmt.Errorf("failed to cleanup route from system to VM: %s", err)
		}
	}

	logrus.Info("✅ Completed cleanup")
	return nil
}

func setupPacketRoutingInsideVM(vmName string) error {
	kindIf, err := sh.RunIO(fmt.Sprintf("limactl shell %s -- ip -o link show | awk -F': ' '{print $2}' | grep 'br-'", vmName))
	if err != nil {
		return fmt.Errorf("failed to get kind network interface name: %s", err)
	}

	hostIf := "lima0"

	// Forward all the packets that are coming from the host interface to the kind network interface
	rule := NewIPTable().
		Sudo().
		Table("filter").
		Specification(fmt.Sprintf(
			"-4 -p tcp -s 192.168.105.1 -d 172.18.0.0/16 -j ACCEPT -i %s -o %s",
			hostIf,
			trim(kindIf),
		))

	if err := sh.Run(fmt.Sprintf("limactl shell %s -- %s", vmName, rule.Command("-C FORWARD").String())); err != nil {
		if err := sh.Run(fmt.Sprintf("limactl shell %s -- %s", vmName, rule.Command("-A FORWARD").String())); err != nil {
			return fmt.Errorf("failed to setup route from VM network interface to kind network interface: %w", err)
		}
	}

	return nil
}

func setupPacketRoutingOnHost(vmName string) error {
	vm := models.NewVM(vmName, "", 0)
	if err := vm.GetByName(); err != nil {
		return fmt.Errorf("failed to get VM by name: %w", err)
	}

	// IPv4 Routing
	limaVMIPv4 := vm.GetVMIPv4()
	ipv4Subnetprefix, err := GetIPv4SubnetPrefix("kind")
	if err != nil {
		return fmt.Errorf("failed to get IPv4 subnet prefix: %w", err)
	}
	if err := sh.RunSilent(fmt.Sprintf("sudo route -nv add -net %s %s", ipv4Subnetprefix, limaVMIPv4)); err != nil {
		return fmt.Errorf("failed to setup route from system to VM: %s", err)
	}

	// IPv6 Routing
	if !disableIPv6 {
		limaVMIPv6 := vm.GetVMIPv6()
		ipv6Subnetprefix, err := GetIPv6SubnetPrefix("kind")
		if err != nil {
			return fmt.Errorf("failed to get IPv6 subnet prefix: %w", err)
		}
		if err := sh.RunSilent(fmt.Sprintf("sudo route -nv add -inet6 %s:: %s", ipv6Subnetprefix, limaVMIPv6)); err != nil {
			return fmt.Errorf("failed to setup route from system to VM: %s", err)
		}
	}

	return nil
}

func trim(data []byte) string {
	return strings.Trim(string(data), " \n")
}
