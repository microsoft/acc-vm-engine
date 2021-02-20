package main

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	genCmd := NewGenerateCmd()

	gen := kingpin.Command("generate", "Generate VM template.")
	gen.Flag("config.file", "Configuration file path.").Short('c').Required().StringVar(&genCmd.configFile)
	gen.Flag("output.directory", "Output directory.").Short('o').StringVar(&genCmd.outputDir)
	gen.Flag("ssh.public.key", "SSH public key file path.").Short('k').StringsVar(&genCmd.sshPubKeys)
	gen.Flag("disable-snp","Disable SNP.").Short('s').Default("false").BoolVar(&genCmd.Disable_SNP)
	gen.Flag("diskRP-version","Disk RP version 1 or 2.").Short('d').StringVar(&genCmd.DiskRP_ver)

	switch kingpin.Parse() {
	case "generate":
		if err := genCmd.Run(); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	}
}
