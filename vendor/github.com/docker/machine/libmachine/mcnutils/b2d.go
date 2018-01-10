package mcnutils

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/docker/machine/libmachine/log"
	"github.com/lambda-linux/lambda-machine-local/version"
)

const (
	defaultURL            = "https://api.github.com/repos/lambda-linux/lambda-linux-vbox/releases"
	defaultISOFilename    = "lambda-linux-vbox.iso"
	defaultVolumeIDOffset = int64(0x8028)
	versionPrefix         = "-v"
	defaultVolumeIDLength = 32
)

var (
	GithubAPIToken string
	// githubPrivate is set to false for regular builds
	githubPrivate = false
)

var (
	errGitHubAPIResponse = errors.New(`failure getting a version tag from the Github API response (are you getting rate limited by Github?)`)
)

func defaultTimeout(network, addr string) (net.Conn, error) {
	return net.Dial(network, addr)
}

func getClient() *http.Client {
	transport := http.Transport{
		DisableKeepAlives: true,
		Proxy:             http.ProxyFromEnvironment,
		Dial:              defaultTimeout,
	}

	return &http.Client{
		Transport: &transport,
	}
}

func getRequest(apiURL string) (*http.Request, error) {
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	if GithubAPIToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("token %s", GithubAPIToken))
	}

	return req, nil
}

// releaseGetter is a client that gets release information of a product and downloads it.
type releaseGetter interface {
	// filename returns filename of the product.
	filename() string
	// getReleaseTag gets a release tag from the given URL.
	getReleaseTag(apiURL string) (string, error)
	// getReleaseURL gets the latest release download URL from the given URL.
	getReleaseURL(apiURL string) (string, error)
	// download downloads a file from the given dlURL and saves it under dir.
	download(dir, file, dlURL string) error
}

// b2dReleaseGetter implements the releaseGetter interface for getting the release of Boot2Docker.
type b2dReleaseGetter struct {
	isoFilename string
}

func (b *b2dReleaseGetter) filename() string {
	if b == nil {
		return ""
	}
	return b.isoFilename
}

// getReleaseTag gets the release tag of Boot2Docker from apiURL.
func (*b2dReleaseGetter) getReleaseTag(apiURL string) (string, error) {
	if apiURL == "" {
		apiURL = defaultURL
	}

	if !version.RC() {
		// Just go straight to the convenience URL for "/latest" if we
		// are a non-release candidate version.  "/latest" won't return
		// non-RCs, so that's what we use for stable releases of
		// Machine.
		apiURL = apiURL + "/latest"
	}

	client := getClient()
	req, err := getRequest(apiURL)
	if err != nil {
		return "", err
	}
	rsp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	// If we call the API endpoint
	// "/repos/boot2docker/boot2docker/releases" without specifying
	// "/latest", we will receive a list of releases instead of a single
	// one, and we should decode accordingly.
	if version.RC() {
		var tags []struct {
			TagName string `json:"tag_name"`
		}
		if err := json.NewDecoder(rsp.Body).Decode(&tags); err != nil {
			return "", err
		}
		t := tags[0]
		if t.TagName == "" {
			return "", errGitHubAPIResponse
		}
		return t.TagName, nil
	}

	// Otherwise, we get back just one release, which we can decode to get
	// the tag.
	var t struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(rsp.Body).Decode(&t); err != nil {
		return "", err
	}
	if t.TagName == "" {
		return "", errGitHubAPIResponse
	}
	return t.TagName, nil
}

// getReleaseID gets the release ID of lambda-linux-vbox from apiURL.
func (*b2dReleaseGetter) getReleaseID(apiURL string) (string, error) {
	if apiURL == "" {
		apiURL = defaultURL
	}

	if !version.RC() {
		// Just go straight to the convenience URL for "/latest" if we
		// are a non-release candidate version.  "/latest" won't return
		// non-RCs, so that's what we use for stable releases of
		// Machine.
		apiURL = apiURL + "/latest"
	}

	client := getClient()
	req, err := getRequest(apiURL)
	if err != nil {
		return "", err
	}
	rsp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	// If we call the API endpoint
	// "/repos/boot2docker/boot2docker/releases" without specifying
	// "/latest", we will receive a list of releases instead of a single
	// one, and we should decode accordingly.
	if version.RC() {
		var ids []struct {
			ID int `json:"id"`
		}
		if err := json.NewDecoder(rsp.Body).Decode(&ids); err != nil {
			return "", err
		}
		i := ids[0]
		if i.ID == 0 {
			return "", errGitHubAPIResponse
		}
		return fmt.Sprintf("%d", i.ID), nil
	}

	// Otherwise, we get back just one release, which we can decode to get
	// the tag.
	var i struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(rsp.Body).Decode(&i); err != nil {
		return "", err
	}
	if i.ID == 0 {
		return "", errGitHubAPIResponse
	}
	return fmt.Sprintf("%d", i.ID), nil
}

// getISOAssetURL gets asset URL for lambda-linux-vbox.iso.
func (*b2dReleaseGetter) getISOAssetURL(releaseID string, owner string, repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/%s/assets", owner, repo, releaseID)
	client := getClient()
	req, err := getRequest(url)
	if err != nil {
		return "", err
	}
	rsp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	var assets []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	if err := json.NewDecoder(rsp.Body).Decode(&assets); err != nil {
		return "", err
	}
	for _, asset := range assets {
		if asset.Name == "lambda-linux-vbox.iso" {
			return asset.Url, nil
		}
	}
	return "", errGitHubAPIResponse
}

// getReleaseURL gets the latest release URL of Boot2Docker.
func (b *b2dReleaseGetter) getReleaseURL(apiURL string) (string, error) {
	if apiURL == "" {
		apiURL = defaultURL
	}

	// match github (enterprise) release urls:
	// https://api.github.com/repos/../../releases or
	// https://some.github.enterprise/api/v3/repos/../../releases
	re := regexp.MustCompile("(https?)://([^/]+)(/api/v3)?/repos/([^/]+)/([^/]+)/releases")
	matches := re.FindStringSubmatch(apiURL)
	if len(matches) != 6 {
		// does not match a github releases api URL
		return apiURL, nil
	}

	scheme, host, org, repo := matches[1], matches[2], matches[4], matches[5]
	if host == "api.github.com" {
		host = "github.com"
	}

	tag, err := b.getReleaseTag(apiURL)
	if err != nil {
		return "", err
	}

	log.Infof("Latest release for %s/%s/%s is %s", host, org, repo, tag)

	if githubPrivate {
		releaseid, err := b.getReleaseID(apiURL)
		if err != nil {
			return "", err
		}
		isourl, err := b.getISOAssetURL(releaseid, org, repo)
		if err != nil {
			return "", err
		}
		// curl -vLJO -H 'Accept: application/octet-stream' \
		// 'https://api.github.com/repos/ORGANIZATION/lambda-linux-vbox/releases/assets/XXXXXXX?access_token=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX'
		url := fmt.Sprintf("%s?access_token=%s", isourl, GithubAPIToken)
		return url, nil
	}

	url := fmt.Sprintf("%s://%s/%s/%s/releases/download/%s/%s", scheme, host, org, repo, tag, b.isoFilename)
	return url, nil
}

func (*b2dReleaseGetter) download(dir, file, isoURL string) error {
	u, err := url.Parse(isoURL)

	var src io.ReadCloser
	if u.Scheme == "file" || u.Scheme == "" {
		s, err := os.Open(u.Path)
		if err != nil {
			return err
		}

		src = s
	} else if githubPrivate {
		client := getClient()
		req, err := http.NewRequest("GET", isoURL, nil)
		if err != nil {
			return err
		}
		// curl -vLJO -H 'Accept: application/octet-stream' \
		// 'https://api.github.com/repos/ORGANIZATION/lambda-linux-vbox/releases/assets/XXXXXXX?access_token=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX'
		req.Header.Add("Accept", fmt.Sprintf("application/octet-stream"))
		rsp, err := client.Do(req)
		if err != nil {
			return err
		}
		src = &ReaderWithProgress{
			ReadCloser:     rsp.Body,
			out:            os.Stdout,
			expectedLength: rsp.ContentLength,
		}
	} else {
		client := getClient()
		s, err := client.Get(isoURL)
		if err != nil {
			return err
		}

		src = &ReaderWithProgress{
			ReadCloser:     s.Body,
			out:            os.Stdout,
			expectedLength: s.ContentLength,
		}
	}

	defer src.Close()

	// Download to a temp file first then rename it to avoid partial download.
	f, err := ioutil.TempFile(dir, file+".tmp")
	if err != nil {
		return err
	}

	defer func() {
		if err := removeFileIfExists(f.Name()); err != nil {
			log.Warnf("Error removing file: %s", err)
		}
	}()

	if _, err := io.Copy(f, src); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	// Dest is the final path of the boot2docker.iso file.
	dest := filepath.Join(dir, file)

	// Windows can't rename in place, so remove the old file before
	// renaming the temporary downloaded file.
	if err := removeFileIfExists(dest); err != nil {
		return err
	}

	return os.Rename(f.Name(), dest)
}

// iso is an ISO volume.
type iso interface {
	// path returns the path of the ISO.
	path() string
	// exists reports whether the ISO exists.
	exists() bool
	// version returns version information of the ISO.
	version() (string, error)
}

// b2dISO represents a Boot2Docker ISO. It implements the ISO interface.
type b2dISO struct {
	// path of Boot2Docker ISO
	commonIsoPath string

	// offset and length of ISO volume ID
	// cf. http://serverfault.com/questions/361474/is-there-a-way-to-change-a-iso-files-volume-id-from-the-command-line
	volumeIDOffset int64
	volumeIDLength int
}

func (b *b2dISO) path() string {
	if b == nil {
		return ""
	}
	return b.commonIsoPath
}

func (b *b2dISO) exists() bool {
	if b == nil {
		return false
	}

	_, err := os.Stat(b.commonIsoPath)
	return !os.IsNotExist(err)
}

// version scans the volume ID in b and returns its version tag.
func (b *b2dISO) version() (string, error) {
	if b == nil {
		return "", nil
	}

	iso, err := os.Open(b.commonIsoPath)
	if err != nil {
		return "", err
	}
	defer iso.Close()

	isoMetadata := make([]byte, b.volumeIDLength)
	_, err = iso.ReadAt(isoMetadata, b.volumeIDOffset)
	if err != nil {
		return "", err
	}

	fullVersion := string(isoMetadata)

	versionIndex := strings.Index(fullVersion, versionPrefix)
	if versionIndex == -1 {
		return "", fmt.Errorf("Did not find prefix %q in version string", versionPrefix)
	}

	// Original magic file string looks similar to this: "Boot2Docker-v0.1.0              "
	// This will return "v0.1.0" given the above string
	vers := strings.TrimSpace(fullVersion)[versionIndex+1:]

	log.Debug("local Lambda Linux VirtualBox ISO version: ", vers)
	return vers, nil
}

func removeFileIfExists(name string) error {
	if _, err := os.Stat(name); err == nil {
		if err := os.Remove(name); err != nil {
			return fmt.Errorf("Error removing temporary download file: %s", err)
		}
	}
	return nil
}

type B2dUtils struct {
	releaseGetter
	iso
	storePath    string
	imgCachePath string
}

func NewB2dUtils(storePath string) *B2dUtils {
	imgCachePath := filepath.Join(storePath, "cache")

	return &B2dUtils{
		releaseGetter: &b2dReleaseGetter{isoFilename: defaultISOFilename},
		iso: &b2dISO{
			commonIsoPath:  filepath.Join(imgCachePath, defaultISOFilename),
			volumeIDOffset: defaultVolumeIDOffset,
			volumeIDLength: defaultVolumeIDLength,
		},
		storePath:    storePath,
		imgCachePath: imgCachePath,
	}
}

// DownloadISO downloads boot2docker ISO image for the given tag and save it at dest.
func (b *B2dUtils) DownloadISO(dir, file, isoURL string) error {
	log.Infof("We will now attempt to automatically download the latest release. This can take some time depending on the speed of your internet connection")
	log.Infof("You can also download the latest lambda-linux-vbox.iso from https://github.com/lambda-linux/lambda-linux-vbox/releases and copy it to %s", b.path())
	log.Infof("Downloading %s from %s...", b.path(), isoURL)
	return b.download(dir, file, isoURL)
}

type ReaderWithProgress struct {
	io.ReadCloser
	out                io.Writer
	bytesTransferred   int64
	expectedLength     int64
	nextPercentToPrint int64
}

func (r *ReaderWithProgress) Read(p []byte) (int, error) {
	n, err := r.ReadCloser.Read(p)

	if n > 0 {
		r.bytesTransferred += int64(n)
		percentage := r.bytesTransferred * 100 / r.expectedLength

		for percentage >= r.nextPercentToPrint {
			if r.nextPercentToPrint%10 == 0 && r.nextPercentToPrint != 0 {
				if r.nextPercentToPrint == 100 {
					fmt.Fprintf(r.out, "....%d%% (done) \n", r.nextPercentToPrint)
				} else {
					fmt.Fprintf(r.out, "....%d%%  (downloaded) \n", r.nextPercentToPrint)
				}
			}
			r.nextPercentToPrint += 2
		}
	}

	return n, err
}

func (r *ReaderWithProgress) Close() error {
	return r.ReadCloser.Close()
}

func (b *B2dUtils) DownloadLatestBoot2Docker(apiURL string) error {
	latestReleaseURL, err := b.getReleaseURL(apiURL)
	if err != nil {
		return err
	}

	return b.DownloadISOFromURL(latestReleaseURL)
}

func (b *B2dUtils) DownloadISOFromURL(latestReleaseURL string) error {
	return b.DownloadISO(b.imgCachePath, b.filename(), latestReleaseURL)
}

func (b *B2dUtils) UpdateISOCache(isoURL string) error {
	// recreate the cache dir if it has been manually deleted
	if _, err := os.Stat(b.imgCachePath); os.IsNotExist(err) {
		log.Infof("Image cache directory does not exist, creating it at %s...", b.imgCachePath)
		if err := os.MkdirAll(b.imgCachePath, 0700); err != nil {
			return err
		}
	}

	exists := b.exists()

	if isoURL != "" {
		if exists {
			// Warn that the b2d iso won't be updated if isoURL is set
			log.Warnf("Lambda Linux VirtualBox URL was explicitly set to %q at create time, so Lambda Machine Local cannot upgrade this machine to the latest version.", isoURL)
		}
		// Non-default B2D are not cached
		return nil
	}

	if !exists {
		log.Info("No default Lambda Linux VirtualBox ISO found locally, downloading the latest release...")
		return b.DownloadLatestBoot2Docker("")
	}

	latest := b.isLatest()
	if !latest {
		log.Info("Default Lambda Linux VirtualBox ISO is out-of-date, downloading the latest release...")
		return b.DownloadLatestBoot2Docker("")
	}

	return nil
}

func (b *B2dUtils) CopyIsoToMachineDir(isoURL, machineName string) error {
	if err := b.UpdateISOCache(isoURL); err != nil {
		return err
	}

	// TODO: This is a bit off-color.
	machineDir := filepath.Join(b.storePath, "machines", machineName)
	machineIsoPath := filepath.Join(machineDir, b.filename())

	// By default just copy the existing "cached" iso to the machine's directory...
	if isoURL == "" {
		log.Infof("Copying %s to %s...", b.path(), machineIsoPath)
		return CopyFile(b.path(), machineIsoPath)
	}

	// if ISO is specified, check if it matches a github releases url or fallback to a direct download
	downloadURL, err := b.getReleaseURL(isoURL)
	if err != nil {
		return err
	}

	return b.DownloadISO(machineDir, b.filename(), downloadURL)
}

// isLatest checks the latest release tag and
// reports whether the local ISO cache is the latest version.
//
// It returns false if failing to get the local ISO version
// and true if failing to fetch the latest release tag.
func (b *B2dUtils) isLatest() bool {
	localVer, err := b.version()
	if err != nil {
		log.Warn("Unable to get the local Lambda Linux VirtualBox ISO version: ", err)
		return false
	}

	latestVer, err := b.getReleaseTag("")
	if err != nil {
		log.Warn("Unable to get the latest Lambda Linux VirtualBox ISO release version: ", err)
		return true
	}

	return localVer == latestVer
}

// MakeDiskImage makes a boot2docker VM disk image.
// See https://github.com/boot2docker/boot2docker/blob/master/rootfs/rootfs/etc/rc.d/automount
// Ext4 label can be a maximum of 16 characters. We use 'll-vbox-data'
func MakeDiskImage(publicSSHKeyPath string) (*bytes.Buffer, error) {
	magicString := "lambda-linux-vbox-disk, please format-me"

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	// magicString first so the automount script knows to format the disk
	file := &tar.Header{Name: magicString, Size: int64(len(magicString))}

	log.Debug("Writing magic tar header")

	if err := tw.WriteHeader(file); err != nil {
		return nil, err
	}

	if _, err := tw.Write([]byte(magicString)); err != nil {
		return nil, err
	}

	// .ssh/key.pub => authorized_keys
	file = &tar.Header{Name: ".ssh", Typeflag: tar.TypeDir, Mode: 0700}
	if err := tw.WriteHeader(file); err != nil {
		return nil, err
	}

	log.Debug("Writing SSH key tar header")

	pubKey, err := ioutil.ReadFile(publicSSHKeyPath)
	if err != nil {
		return nil, err
	}

	file = &tar.Header{Name: ".ssh/authorized_keys", Size: int64(len(pubKey)), Mode: 0644}
	if err := tw.WriteHeader(file); err != nil {
		return nil, err
	}

	if _, err := tw.Write([]byte(pubKey)); err != nil {
		return nil, err
	}

	file = &tar.Header{Name: ".ssh/authorized_keys2", Size: int64(len(pubKey)), Mode: 0644}
	if err := tw.WriteHeader(file); err != nil {
		return nil, err
	}

	if _, err := tw.Write([]byte(pubKey)); err != nil {
		return nil, err
	}

	if err := tw.Close(); err != nil {
		return nil, err
	}

	return buf, nil
}
