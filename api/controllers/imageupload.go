package controllers

import (
	"context"
	"net/http"
	"pi-gravity/internal/config"
	"pi-gravity/pkg/tracer"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UploadImage(config *config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg, err := awsconfig.LoadDefaultConfig(
			c,
			awsconfig.WithRegion("auto"), // Specify the region, if required
			awsconfig.WithCredentialsProvider(
				aws.NewCredentialsCache(
					credentials.NewStaticCredentialsProvider("yourAccessKeyId", "yourSecretAccessKey", ""),
				),
			),
			// Add endpoint resolver for Cloudflare R2
			awsconfig.WithEndpointResolverWithOptions(
				aws.EndpointResolverWithOptionsFunc(
					func(service, region string, options ...interface{}) (aws.Endpoint, error) {
						return aws.Endpoint{
								URL:           "https://<your-r2-endpoint>",
								SigningRegion: "auto",
							},
							nil
					},
				),
			),
		)

		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				gin.H{"error": tracer.TraceError(err).Error()},
			)
			return
		}

		s3Client := s3.NewFromConfig(cfg)

		file, err := c.FormFile("file")
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				gin.H{"error": tracer.TraceError(err).Error()},
			)
			return
		}

		// Open the file
		openedFile, err := file.Open()
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				gin.H{"error": tracer.TraceError(err).Error()},
			)
			return
		}
		defer openedFile.Close()

		// Upload the file to Cloudflare R2
		_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String("yourBucketName"), // Replace with your bucket name
			Key:    aws.String(file.Filename),
			Body:   openedFile,
		})

		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				gin.H{"error": tracer.TraceError(err).Error()},
			)
		}

		c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
	}
}
