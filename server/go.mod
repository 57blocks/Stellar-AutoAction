module github.com/57blocks/auto-action/server

go 1.22.2

require (
	github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/aws/aws-sdk-go-v2 v1.30.5
	github.com/aws/aws-sdk-go-v2/config v1.27.30
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.38.0
	github.com/aws/aws-sdk-go-v2/service/iam v1.35.2
	github.com/aws/aws-sdk-go-v2/service/lambda v1.58.1
	github.com/aws/aws-sdk-go-v2/service/scheduler v1.10.4
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.32.7
	github.com/aws/smithy-go v1.20.4
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.10.0
	github.com/go-resty/resty/v2 v2.14.0
	github.com/golang-migrate/migrate/v4 v4.17.1
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.3
	github.com/jarcoal/httpmock v1.3.1
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	github.com/spf13/viper v1.19.0
	github.com/stellar/go v0.0.0-20240906221814-1cd3bbe73531
	github.com/stretchr/testify v1.9.0
	go.uber.org/zap v1.21.0
	golang.org/x/crypto v0.26.0
	gorm.io/driver/postgres v1.5.9
	gorm.io/gorm v1.25.11
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.4 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.29 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.12 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.17 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.17 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.18 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.22.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.26.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.30.5 // indirect
	github.com/bytedance/sonic v1.12.1 // indirect
	github.com/bytedance/sonic/loader v0.2.0 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.5 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-chi/chi v4.1.2+incompatible // indirect
	github.com/go-errors/errors v1.5.1 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.22.0 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/gorilla/schema v1.2.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/manucorporat/sse v0.0.0-20160126180136-ee05b128a739 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/segmentio/go-loggly v0.5.1-0.20171222203950-eb91657e62b2 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stellar/go-xdr v0.0.0-20231122183749-b53fb00bcac2 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/arch v0.9.0 // indirect
	golang.org/x/exp v0.0.0-20231006140011-7918f672742d // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
