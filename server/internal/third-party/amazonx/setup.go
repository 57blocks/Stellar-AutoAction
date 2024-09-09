package amazonx

import (
	"context"
	"sync"

	configx "github.com/57blocks/auto-action/server/internal/config"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

var (
	once sync.Once

	amazonConfig         aws.Config
	secretManagerClient  *secretsmanager.Client
	lambdaClient         *lambda.Client
	schedulerClient      *scheduler.Client
	cloudWatchLogsClient *cloudwatchlogs.Client
)

func Setup() error {
	var err error

	once.Do(func() {
		amazonConfig, err = config.LoadDefaultConfig(
			context.TODO(),
			config.WithRegion(configx.GlobalConfig.Region),
		)
		if err != nil {
			err = errorx.AmazonConfig(err.Error())
			return
		}

		lambdaClient = lambda.NewFromConfig(amazonConfig)
		schedulerClient = scheduler.NewFromConfig(amazonConfig)
		cloudWatchLogsClient = cloudwatchlogs.NewFromConfig(amazonConfig)
		secretManagerClient = secretsmanager.NewFromConfig(amazonConfig)

		Conductor = buildAmazonConductor(
			withConfig(amazonConfig),
			withSecretManagerClient(secretManagerClient),
			withLambdaClient(lambdaClient),
			withSchedulerClient(schedulerClient),
			withCloudWatchLogsClient(cloudWatchLogsClient),
		)
	})

	return err
}
