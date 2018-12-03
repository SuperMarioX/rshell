package sftp

import (
	"github.com/luckywinds/rshell/modes/client"
	"github.com/pkg/sftp"
	"io"
	"os"
	"path"
)

func Upload(groupname, host string, port int, user, pass, keyname, passphrase string, timeout int, ciphers []string, srcFilePath, desDirPath string) error {
	var (
		session *sftp.Client
		err     error
	)

	c, err := client.New(groupname, host, port, user, pass, keyname, passphrase, timeout, ciphers)
	if err != nil {
		return err
	}

	session, err = sftp.NewClient(c)
	if err != nil {
		return err
	}
	defer session.Close()

	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	var desFileName = path.Base(srcFilePath)
	desFile, err := session.Create(path.Join(desDirPath, desFileName))
	if err != nil {
		return err
	}
	defer desFile.Close()

	_, err = io.Copy(desFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func Download(groupname, host string, port int, user, pass, keyname, passphrase string, timeout int, ciphers []string, srcFilePath, desDirPath string) error {
	var (
		session *sftp.Client
		err     error
	)

	c, err := client.New(groupname, host, port, user, pass, keyname, passphrase, timeout, ciphers)
	if err != nil {
		return err
	}

	session, err = sftp.NewClient(c)
	if err != nil {
		return err
	}
	defer session.Close()

	var desFileName = path.Base(srcFilePath)

	if err = os.Mkdir(desDirPath, os.ModeDir); err != nil {
		if os.IsNotExist(err) {
			return err
		}
	}
	if err = os.Mkdir(path.Join(desDirPath, host), os.ModeDir); err != nil {
		if os.IsNotExist(err) {
			return err
		}
	}

	desFile, err := os.Create(path.Join(path.Join(desDirPath, host), desFileName))
	if err != nil {
		return err
	}
	defer desFile.Close()

	srcFile, err := session.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer desFile.Close()

	_, err = io.Copy(desFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}
