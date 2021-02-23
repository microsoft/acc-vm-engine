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

	switch kingpin.Parse() {
	case "generate":
		if err := genCmd.Run(); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	}
}
