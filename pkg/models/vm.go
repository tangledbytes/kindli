package models

import (
	"database/sql"

	"github.com/utkarsh-pro/kindli/pkg/db"
)

type VM struct {
	ID             uint
	Name           string
	LimaConfigPath string
	DockerPort     int
}

func VMPreload() {
	db.RegisterPreload(`
CREATE TABLE IF NOT EXISTS vm (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT UNIQUE,
	lima_config_path TEXT UNIQUE,
	docker_port INTEGER UNIQUE
);`)
}

func NewVM(name, limaConfigPath string, dockerPort int) *VM {
	return &VM{
		Name:           name,
		LimaConfigPath: limaConfigPath,
		DockerPort:     dockerPort,
	}
}

func (vm *VM) Save() error {
	_, err := db.Instance().Exec(
		`INSERT INTO vm (name, lima_config_path, docker_port) VALUES (?, ?, ?)`,
		vm.Name,
		vm.LimaConfigPath,
		vm.DockerPort,
	)

	return err
}

func (vm *VM) Delete() error {
	_, err := db.Instance().Exec(`DELETE FROM vm WHERE name = ?`, vm.Name)

	return err
}

func (vm *VM) Exists() (bool, error) {
	var count int
	err := db.Instance().QueryRow(`SELECT COUNT(*) FROM vm WHERE name = ?`, vm.Name).Scan(&count)

	return count > 0, err
}

func (vm *VM) GetByName() error {
	err := db.Instance().QueryRow(`SELECT * FROM vm WHERE name = ?`, vm.Name).Scan(&vm.ID, &vm.Name, &vm.LimaConfigPath, &vm.DockerPort)

	if err != nil {
		return err
	}

	return nil
}

func GetMaxVMDockerPort() (int, error) {
	var maxPort sql.NullInt64
	err := db.Instance().QueryRow(`SELECT MAX(docker_port) FROM vm`).Scan(&maxPort)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}

		return 0, err
	}

	return int(maxPort.Int64), nil
}

func ListVM() ([]*VM, error) {
	var vms []*VM

	rows, err := db.Instance().Query(`SELECT * FROM vm`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var vm VM
		err := rows.Scan(&vm.ID, &vm.Name, &vm.LimaConfigPath, &vm.DockerPort)
		if err != nil {
			return nil, err
		}

		vms = append(vms, &vm)
	}

	return vms, nil
}
