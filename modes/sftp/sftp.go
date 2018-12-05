package sftp

import (
	"github.com/luckywinds/rshell/modes/client"
	"github.com/pkg/sftp"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func Upload(groupname, host string, port int, user, pass, keyname, passphrase string, timeout int, ciphers []string, srcFilePath, desDirPath string) ([]string, error) {
	var (
		session *sftp.Client
		err     error
	)

	c, err := client.New(groupname, host, port, user, pass, keyname, passphrase, timeout, ciphers)
	if err != nil {
		return nil, err
	}

	session, err = sftp.NewClient(c)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	srcFiles, err := filepath.Glob(srcFilePath)
	if err != nil {
		return nil, err
	}
	if srcFiles != nil {
		for _, sf := range srcFiles {
			srcFile, err := os.Open(sf)
			if err != nil {
				return nil, err
			}
			defer srcFile.Close()

			var desFileName string
			if runtime.GOOS == "windows" {
				desFileName = path.Base(strings.Replace(srcFile.Name(), "\\", "/", -1))
			} else {
				desFileName = path.Base(srcFile.Name())
			}
			desFile, err := session.Create(path.Join(desDirPath, desFileName))
			if err != nil {
				return nil, err
			}
			defer desFile.Close()

			_, err = io.Copy(desFile, srcFile)
			if err != nil {
				return nil, err
			}
		}
		return srcFiles, nil
	} else {
		return []string{}, nil
	}
}

func Download(groupname, host string, port int, user, pass, keyname, passphrase string, timeout int, ciphers []string, srcFilePath, desDirPath string) ([]string, error) {
	var (
		session *sftp.Client
		err     error
	)

	c, err := client.New(groupname, host, port, user, pass, keyname, passphrase, timeout, ciphers)
	if err != nil {
		return nil, err
	}

	session, err = sftp.NewClient(c)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	if err = os.Mkdir(desDirPath, os.ModeDir); err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
	}
	if err = os.Mkdir(path.Join(desDirPath, groupname), os.ModeDir); err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
	}
	if err = os.Mkdir(path.Join(path.Join(desDirPath, groupname), host), os.ModeDir); err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
	}

	srcFiles, err := session.Glob(srcFilePath)
	if err != nil {
		return nil, err
	}
	if srcFiles != nil {
		for _, sf := range srcFiles {
			var desFileName = path.Base(sf)
			srcFile, err := session.Open(sf)
			if err != nil {
				return nil, err
			} else {
				desFile, err := os.Create(path.Join(path.Join(path.Join(desDirPath, groupname), host), desFileName))
				if err != nil {
					return nil, err
				}
				defer desFile.Close()

				_, err = io.Copy(desFile, srcFile)
				if err != nil {
					return nil, err
				}
			}
			defer srcFile.Close()
		}
		return srcFiles, nil
	} else {
		return []string{}, nil
	}
}
