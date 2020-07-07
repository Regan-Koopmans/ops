package lepton

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

// AzureStorage provides Azure storage related operations
type AzureStorage struct{}

// CopyToBucket copies archive to bucket
func (az *AzureStorage) CopyToBucket(config *Config, archPath string) error {

	vhdxPath := "/tmp/" + config.CloudConfig.ImageName + ".vhdx"

	vhdxPath = strings.ReplaceAll(vhdxPath, "-image", "")

	// this is probably just for hyper-v not azure
	args := []string{
		"convert", "-f", "raw",
		"-O", "vhdx", "-o", "subformat=dynamic",
		archPath, vhdxPath,
	}

	cmd := exec.Command("qemu-img", args...)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT"), os.Getenv("AZURE_STORAGE_ACCESS_KEY")
	if len(accountName) == 0 || len(accountKey) == 0 {
		log.Fatal("Either the AZURE_STORAGE_ACCOUNT or AZURE_STORAGE_ACCESS_KEY environment variable is not set")
	}

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	containerName := "quickstart-nanos"

	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))

	containerURL := azblob.NewContainerURL(*URL, p)

	fmt.Printf("Creating a container named %s\n", containerName)
	ctx := context.Background()
	_, err = containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	if err != nil {
		fmt.Println(err)
	}

	blobURL := containerURL.NewBlockBlobURL(config.CloudConfig.ImageName + ".vhdx")

	file, err := os.Open(vhdxPath)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Uploading the file with blob name: %s\n", vhdxPath)
	_, err = azblob.UploadFileToBlockBlob(ctx, file, blobURL, azblob.UploadToBlockBlobOptions{
		BlockSize:   4 * 1024 * 1024,
		Parallelism: 16})
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

// DeleteFromBucket deletes key from config's bucket
func (az *AzureStorage) DeleteFromBucket(config *Config, key string) error {

	fmt.Println("un-implemented")

	return nil
}
