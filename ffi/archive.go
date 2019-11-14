package main

import (
	"io/ioutil"
	"os"

	"github.com/mholt/archiver"
	"github.com/pkg/errors"
	kotsv1beta1 "github.com/replicatedhq/kots/kotskinds/apis/kots/v1beta1"
	"k8s.io/client-go/kubernetes/scheme"
)

func extractArchive(rootPath, fromArchivePath string) (*archiver.TarGz, error) {
	// extract the current archive to this root
	tarGz := &archiver.TarGz{
		Tar: &archiver.Tar{
			ImplicitTopLevelFolder: false,
		},
	}
	if err := tarGz.Unarchive(fromArchivePath, rootPath); err != nil {
		return nil, err
	}

	return tarGz, nil
}

func readCursorFromPath(installationFilePath string) (string, error) {
	_, err := os.Stat(installationFilePath)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", errors.Wrap(err, "failed to open file")
	}

	installationData, err := ioutil.ReadFile(installationFilePath)
	if err != nil {
		return "", errors.Wrap(err, "failed to read update installation file")
	}

	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode([]byte(installationData), nil, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to devode installation data")
	}

	installation := obj.(*kotsv1beta1.Installation)
	return installation.Spec.UpdateCursor, nil
}

func loadLicenseFromPath(expectedLicenseFile string) (*kotsv1beta1.License, error) {
	_, err := os.Stat(expectedLicenseFile)
	if err != nil {
		return nil, errors.New("find license file in archive")
	}
	licenseData, err := ioutil.ReadFile(expectedLicenseFile)
	if err != nil {
		return nil, errors.Wrap(err, "read license file")
	}

	return loadLicense(string(licenseData))
}

func loadLicense(licenseData string) (*kotsv1beta1.License, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode([]byte(licenseData), nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "decode license data")
	}

	return obj.(*kotsv1beta1.License), nil
}
