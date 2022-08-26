package models

import (
	"github.com/utkarsh-pro/kindli/pkg/db"
)

type Cluster struct {
	ID             uint
	Name           string
	KindConfigPath string
	VM             string
}

func ClusterPreload() {
	db.RegisterPreload(`
CREATE TABLE IF NOT EXISTS cluster (
	id INTEGER PRIMARY KEY,
	name TEXT,
	kind_config_path TEXT UNIQUE,
	vm TEXT,
	FOREIGN KEY (vm) REFERENCES vm(name)
);`)
}

func NewCluster(name, kindConfigPath, vm string) *Cluster {
	return &Cluster{
		Name:           name,
		KindConfigPath: kindConfigPath,
		VM:             vm,
	}
}

func (cluster *Cluster) Save() error {
	_, err := db.Instance().Exec(
		`INSERT INTO cluster (name, kind_config_path, vm) VALUES (?, ?, ?)`,
		cluster.Name,
		cluster.KindConfigPath,
		cluster.VM,
	)

	return err
}

func (cluster *Cluster) Delete() error {
	_, err := db.Instance().Exec(`DELETE FROM cluster WHERE name = ?`, cluster.Name)

	return err
}

func (cluster *Cluster) Exists() (bool, error) {
	var count int
	err := db.Instance().QueryRow(`SELECT COUNT(*) FROM cluster WHERE name = ? AND vm = ?`, cluster.Name, cluster.VM).Scan(&count)

	return count > 0, err
}

func (cluster *Cluster) GetByName() error {
	err := db.Instance().QueryRow(`SELECT * FROM cluster WHERE name = ?`, cluster.Name).Scan(
		&cluster.ID,
		&cluster.Name,
		&cluster.KindConfigPath,
		&cluster.VM,
	)

	return err
}

func ListCluster() ([]Cluster, error) {
	var clusters []Cluster
	rows, err := db.Instance().Query(`SELECT * FROM cluster`)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var cluster Cluster
		err := rows.Scan(
			&cluster.ID,
			&cluster.Name,
			&cluster.KindConfigPath,
			&cluster.VM,
		)

		if err != nil {
			return nil, err
		}

		clusters = append(clusters, cluster)
	}

	return clusters, nil
}

func (cluster *Cluster) AssignID() error {
	clusters, err := ListCluster()
	if err != nil {
		return err
	}

	idIdx := make([]bool, 100)
	for _, c := range clusters {
		if int(c.ID) >= len(idIdx) {
			panic("Cluster ID is too large")
		}

		idIdx[c.ID] = true
	}

	for i := range idIdx {
		if !idIdx[i] {
			cluster.ID = uint(i)
			return nil
		}
	}

	return nil
}
