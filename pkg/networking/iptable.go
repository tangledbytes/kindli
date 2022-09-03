package networking

import "strings"

// IPTable is a wrapper around iptables command
type IPTable struct {
	sudo          bool
	table         string
	cmd           string
	specification string
}

func NewIPTable() *IPTable {
	return &IPTable{}
}

func (ipt *IPTable) Sudo() *IPTable {
	ipt.sudo = true
	return ipt
}

func (ipt *IPTable) Table(table string) *IPTable {
	ipt.table = table
	return ipt
}

func (ipt *IPTable) Command(cmd string) *IPTable {
	ipt.cmd = cmd
	return ipt
}

func (ipt *IPTable) Specification(specification string) *IPTable {
	ipt.specification = specification
	return ipt
}

func (ipt *IPTable) String() string {
	var sb strings.Builder
	if ipt.sudo {
		sb.WriteString("sudo ")
	}

	sb.WriteString("iptables ")

	if ipt.table != "" {
		sb.WriteString("-t ")
		sb.WriteString(ipt.table)
		sb.WriteString(" ")
	}

	sb.WriteString(ipt.cmd)
	sb.WriteString(" ")
	sb.WriteString(ipt.specification)

	return sb.String()
}
