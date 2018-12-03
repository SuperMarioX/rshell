package commands

import "fmt"

func Help() {
	fmt.Println(`Usage: <keywords> <hostgroup> <agruments>

do hostgroup cmd1; cmd2; cmd3
    --- Run cmds on hostgroup use normal user
sudo hostgroup sudo cmd1; cmd2; cmd3
    --- Run cmds on hostgroup use root which auto change from normal user
download hostgroup srcFile desDir
    --- Download srcFile from hostgroup to local desDir
upload hostgroup srcFile desDir
    --- Upload srcFile from local to hostgroup desDir

encrypt_aes cleartext_password
    --- Encrypt cleartext_password with aes 256 cfb
decrypt_aes ciphertext_password
    --- Decrypt ciphertext_password with aes 256 cfb

exit
    --- Exit rshell
?
    --- Help`)
}
