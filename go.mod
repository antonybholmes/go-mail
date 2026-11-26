module github.com/antonybholmes/go-mail

go 1.27.4

replace github.com/antonybholmes/go-sys => ../go-sys

replace github.com/antonybholmes/go-jobs => ../go-jobs

replace github.com/antonybholmes/go-web => ../go-web

require (
	github.com/antonybholmes/go-jobs v0.0.0
	github.com/antonybholmes/go-sys v0.0.0-20261123005409-c7d66a3d2106
	github.com/aws/aws-sdk-go-v2 v1.48.0
	github.com/aws/aws-sdk-go-v2/service/sesv2 v1.79.9
	github.com/aws/aws-sdk-go-v2/service/sqs v1.53.0
	github.com/panjf2000/ants v1.3.0
	github.com/redis/go-redis/v9 v9.24.0
	github.com/segmentio/kafka-go v0.4.51
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.27 // indirect
	github.com/aws/aws-sdk-go-v2/feature/cloudfront/sign v1.12.12 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.5.12 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.8.12 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.5.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.25 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.11.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.14.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.20.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.116.0 // indirect
	github.com/aws/smithy-go v1.28.9 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.20.5 // indirect
	github.com/matoous/go-nanoid/v2 v2.1.0 // indirect
	github.com/mattn/go-colorable v0.1.16 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/pierrec/lz4/v4 v4.1.34 // indirect
	github.com/richardlehane/mscfb v1.0.10 // indirect
	github.com/richardlehane/msoleps v1.0.6 // indirect
	github.com/rs/zerolog v1.35.1 // indirect
	github.com/tiendc/go-deepcopy v1.7.2 // indirect
	github.com/xdg-go/scram v1.4.0 // indirect
	github.com/xuri/efp v0.0.3-0.20260924042702-217b81952dec // indirect
	github.com/xuri/excelize/v2 v2.11.1 // indirect
	github.com/xuri/nfp v0.0.2-0.20250530014748-2ddeb826f9a9 // indirect
	go.uber.org/atomic v1.12.0 // indirect
	golang.org/x/crypto v0.59.0 // indirect
	golang.org/x/net v0.63.0 // indirect
	golang.org/x/sys v0.50.0 // indirect
	golang.org/x/text v0.44.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)
