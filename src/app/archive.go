package app

import (
	"os"
	"path"
)

const archive_name string = "pwman_archive.bz2"

func LookForArchive() (bool, error) {
	exec_path, err := os.Executable()
	if err != nil {
		return false, err
	}
	exec_path = path.Dir(exec_path)
	archive_path := exec_path + "/" + archive_name
	return CheckIfExists(archive_path)
}

func CheckIfExists(file_name string) (bool, error) {
	if _, err := os.Stat(file_name); !os.IsNotExist(err) {
		return true, err
	}
	return false, nil
}

func (ar *PWMan_Archive) GetPhysicalArchivePath() (string, error) {
	if exec_path, err := os.Executable(); err != nil {
		ar.Location = ""
		return "", err
	} else {
		exec_path = path.Dir(exec_path)
		ar.Location = exec_path + "/" + archive_name
		return exec_path + "/" + archive_name, nil
	}
}
