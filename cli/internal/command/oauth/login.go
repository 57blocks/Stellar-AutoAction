package oauth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/57blocks/auto-action/cli/internal/config"
	"github.com/57blocks/auto-action/cli/internal/constant"
	"github.com/57blocks/auto-action/cli/internal/pkg/errorx"
	"github.com/57blocks/auto-action/cli/internal/pkg/restyx"
	"github.com/57blocks/auto-action/cli/internal/pkg/util"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// login represents the login command
var login = &cobra.Command{
	Use:   "login",
	Short: "Login to the Stellar auto-action.",
	Long: `
Description:
  Login the Stellar auto-action based on credential path in the config.
  Or, by the --credential/-c flag, to specify the credential path.
  And will create a new credential under the path you just
  claimed and set it to config, when it's the first time.

Note:
  - You could specify other credentials by **configure** command.
`,
	Args: cobra.NoArgs,
	RunE: loginFunc,
}

func init() {
	oauthGroup.AddCommand(login)

	flagCred := constant.FlagCredential.ValStr()
	login.Flags().StringP(
		flagCred,
		"c",
		config.Vp.GetString(flagCred),
		`
The credential file for the command about to be bound.
If it's the first time, or ignored, the default path will be used.'`)

	flagEnv := constant.FlagEnvironment.ValStr()
	login.Flags().StringP(
		flagEnv,
		"e",
		config.Vp.GetString(flagEnv),
		`
The execution environment about to be bound.
If ignored, the default environment: Horizon will be used.'`)

	flagAcc := constant.FlagAccount.ValStr()
	login.Flags().StringP(
		flagAcc,
		"a",
		"",
		"name of the account")

	flagOrg := constant.FlagOrganization.ValStr()
	login.Flags().StringP(flagOrg,
		"o",
		"",
		"name of the organization")

	if err := login.MarkFlagRequired(flagAcc); err != nil {
		return
	}
	if err := login.MarkFlagRequired(flagOrg); err != nil {
		return
	}
}

type ReqLogin struct {
	Account      string `json:"account"`
	Organization string `json:"organization"`
	Password     string `json:"password"`
	Environment  string `json:"environment"`
}

func loginFunc(cmd *cobra.Command, args []string) error {
	fmt.Println("Password: ")

	pwdBytes, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return errorx.Internal(fmt.Sprintf("reading passowrd error: %s", err.Error()))
	}

	if len(pwdBytes) == 0 {
		return errorx.BadRequest("empty cryptPwd error")
	}

	key, err := loadPublicKey()
	if err != nil {
		return err
	}
	encodedPwd, err := encryptPassword(string(pwdBytes), key)
	if err != nil {
		return err
	}

	success, err := supplierLogin(encodedPwd)
	if err != nil {
		return err
	}

	return syncLogin(success)
}

func supplierLogin(pwdHash string) (*resty.Response, error) {
	URL := util.ParseReqPath(fmt.Sprintf("%s/oauth/login", config.Vp.GetString("bound_with.endpoint")))

	response, err := restyx.Client.R().
		EnableTrace().
		SetBody(ReqLogin{
			Account:      config.Vp.GetString(constant.FlagAccount.ValStr()),
			Organization: config.Vp.GetString(constant.FlagOrganization.ValStr()),
			Password:     pwdHash,
			Environment:  config.Vp.GetString(constant.FlagEnvironment.ValStr()),
		}).
		Post(URL)
	if err != nil {
		return nil, errorx.RestyError(err.Error())
	}
	if response.IsError() {
		return nil, errorx.WithRestyResp(response)
	}

	return response, nil
}

func syncLogin(resp *resty.Response) error {
	cred := new(config.Credential)
	if err := json.Unmarshal(resp.Body(), cred); err != nil {
		return errorx.Internal(fmt.Sprintf("unmarshaling json response error: %s", err.Error()))
	}

	if err := config.WriteCredential(config.Vp.GetString(constant.FlagCredential.ValStr()), cred); err != nil {
		return err
	}

	return config.SyncConfigByFlags()
}

func loadPublicKey() (*rsa.PublicKey, error) {
	pubPEMBytes, err := base64.StdEncoding.DecodeString(constant.PublicKey)
	if err != nil {
		return nil, errorx.Internal("decode public key failed")
	}

	block, _ := pem.Decode(pubPEMBytes)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errorx.Internal("failed to decode PEM block containing public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pubKey := pubKey.(type) {
	case *rsa.PublicKey:
		return pubKey, nil
	default:
		return nil, errorx.Internal("not an RSA public key")
	}
}

func encryptPassword(password string, publicKey *rsa.PublicKey) (string, error) {
	encryptedBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, []byte(password), nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}
