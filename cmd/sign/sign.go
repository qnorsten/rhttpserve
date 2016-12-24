package sign

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/brandur/rserve/cmd"
	"github.com/brandur/rserve/common"
	"github.com/joeshaw/envdecode"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ed25519"
)

var (
	curl bool
)

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: `Creates a shareable link.`,
	Long: `
rclone cat sends any files to standard output.

You can use it like this to output a single file

    rclone cat remote:path/to/file

Or like this to output any file in dir or subdirectories.

    rclone cat remote:path/to/dir

Or like this to output any .txt files in dir or subdirectories.

    rclone --include "*.txt" cat remote:path/to/dir
`,
	Run: func(command *cobra.Command, args []string) {
		cmd.CheckArgs(1, 1, command, args)

		var conf Config
		err := envdecode.Decode(&conf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			os.Exit(1)
		}

		privateKey, err := base64.URLEncoding.DecodeString(conf.PrivateKey)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			os.Exit(1)
		}

		generator := URLGenerator{
			Host:       conf.Host,
			PrivateKey: ed25519.PrivateKey(privateKey),
		}

		// Maybe make this configurable at some point.
		expiresAt := time.Now().Add(48 * time.Hour)
		url, err := generator.Generate(args[0], expiresAt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			os.Exit(1)
		}

		if curl {
			filename := filepath.Base(args[0])
			fmt.Printf("curl -o '%s' '%s'\n", filename, url)
		} else {
			fmt.Printf("%s\n", url)
		}
	},
}

func init() {
	cmd.Root.AddCommand(signCmd)
	signCmd.Flags().BoolVar(&curl, "curl", false, "Output as cURL command")
}

type Config struct {
	Host       string `env:"RSERVE_HOST,required"`
	PrivateKey string `env:"RSERVE_PRIVATE_KEY,required"`
}

type URLGenerator struct {
	Host       string
	PrivateKey ed25519.PrivateKey
}

func (s *URLGenerator) Generate(path string, expiresAt time.Time) (string, error) {
	scheme := "https"
	if s.Host == "localhost" || strings.HasPrefix(s.Host, "localhost:") {
		scheme = "http"
	}

	u := url.URL{
		Host:   s.Host,
		Path:   path,
		Scheme: scheme,
	}

	message := common.Message(path, expiresAt.Unix())
	if cmd.Verbose {
		log.Printf("Message: %v", string(message))
	}

	signature := ed25519.Sign(s.PrivateKey, message)

	u.RawQuery = fmt.Sprintf("expires_at=%v&signature=%v",
		expiresAt.Unix(), base64.URLEncoding.EncodeToString(signature))

	return u.String(), nil
}