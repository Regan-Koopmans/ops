package lepton

import (
	"fmt"
	"os"
	"path/filepath"
)

// Azure contains all operations for Azure
type Azure struct {
	Storage *AzureStorage
}

func (p *Azure) getArchiveName(ctx *Context) string {
	return ctx.config.CloudConfig.ImageName + ".tar.gz"
}

func (p *Azure) customizeImage(ctx *Context) (string, error) {
	imagePath := ctx.config.RunConfig.Imagename
	symlink := filepath.Join(filepath.Dir(imagePath), "disk.raw")

	if _, err := os.Lstat(symlink); err == nil {
		if err := os.Remove(symlink); err != nil {
			return "", fmt.Errorf("failed to unlink: %+v", err)
		}
	}

	err := os.Link(imagePath, symlink)
	if err != nil {
		return "", err
	}

	archPath := filepath.Join(filepath.Dir(imagePath), p.getArchiveName(ctx))
	files := []string{symlink}

	err = createArchive(archPath, files)
	if err != nil {
		return "", err
	}
	return archPath, nil
}

// BuildImage to be upload on Azure
func (p *Azure) BuildImage(ctx *Context) (string, error) {
	c := ctx.config
	err := BuildImage(*c)
	if err != nil {
		return "", err
	}

	return p.customizeImage(ctx)
}

// BuildImageWithPackage to upload on Azure
func (p *Azure) BuildImageWithPackage(ctx *Context, pkgpath string) (string, error) {
	c := ctx.config
	err := BuildImageFromPackage(pkgpath, *c)
	if err != nil {
		return "", err
	}
	return p.customizeImage(ctx)
}

// Initialize Azure related things
func (p *Azure) Initialize() error {
	p.Storage = &AzureStorage{}
	return nil
}

// CreateImage - Creates image on Azure using nanos images
func (p *Azure) CreateImage(ctx *Context) error {

	/*
	   rawdisk="MyLinuxVM.raw"
	   vhddisk="MyLinuxVM.vhd"
	   MB=$((1024*1024))
	   size=$(qemu-img info -f raw --output json "$rawdisk" | gawk 'match($0, /"virtual-size": ([0-9]+),/, val) {print val[1]}')
	   rounded_size=$((($size/$MB + 1)*$MB))
	   echo "Rounded Size = $rounded_size"
	   qemu-img resize MyLinuxVM.raw $rounded_size
	   qemu-img convert -f raw -o subformat=fixed,force_size -O vpc MyLinuxVM.raw MyLinuxVM.vhd
	*/

	// must have *.vhd extension
	// must be at least 20mb
	/* 20971520 bytes */

	fmt.Println("un-implemented")
	return nil
}

// ListImages lists images on azure
func (p *Azure) ListImages(ctx *Context) error {
	fmt.Println("un-implemented")
	return nil
}

// DeleteImage deletes image from Azure
func (p *Azure) DeleteImage(ctx *Context, imagename string) error {
	fmt.Println("un-implemented")
	return nil
}

// CreateInstance - Creates instance on azure Platform
func (p *Azure) CreateInstance(ctx *Context) error {
	fmt.Println("un-implemented")

	// https://github.com/Azure-Samples/azure-sdk-for-go-samples/blob/master/quickstarts/deploy-vm/main.go

	return nil
}

// ListInstances lists instances on Azure
func (p *Azure) ListInstances(ctx *Context) error {
	fmt.Println("un-implemented")
	return nil
}

// DeleteInstance deletes instance from Azure
func (p *Azure) DeleteInstance(ctx *Context, instancename string) error {
	fmt.Println("un-implemented")
	return nil
}

// StartInstance starts an instance in Azure
func (p *Azure) StartInstance(ctx *Context, instancename string) error {
	fmt.Println("un-implemented")

	// https://github.com/Azure-Samples/azure-sdk-for-go-samples/blob/master/compute/vm.go

	return nil
}

// StopInstance deletes instance from Azure
func (p *Azure) StopInstance(ctx *Context, instancename string) error {
	fmt.Println("un-implemented")
	return nil
}

// GetInstanceLogs gets instance related logs
func (p *Azure) GetInstanceLogs(ctx *Context, instancename string, watch bool) error {
	fmt.Println("un-implemented")
	return nil
}

// ResizeImage is not supported on azure.
func (p *Azure) ResizeImage(ctx *Context, imagename string, hbytes string) error {
	return fmt.Errorf("Operation not supported")
}
