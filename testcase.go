package aoj

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	// Endpoint of testcase_header
	TESTCASE_HEADER_ENDPOINT = BASE_ENDPOINT + "/testcase_header.jsp"

	// Endpoint of testcase
	TESTCASE_ENDPOINT = BASE_ENDPOINT + "/testcase.jsp"
)

// Testcase represents the test cases of an AOJ problem.
type Testcase struct {
	testcase
}

type testcase struct {
	ID          string
	Available   int
	Input       []caseFileInfo
	Output      []caseFileInfo
	CaseMapping []string `json:"case_mapping"`
}

type caseFileInfo struct {
	Name string
	Size int64
}

// GetTestcase returns the Testcase for problemID.
func GetTestcase(problemID string) (*Testcase, error) {
	return getTestcase(problemID, true)
}

func GetTestcaseIgnoringCache(problemID string) (*Testcase, error) {
	return getTestcase(problemID, false)
}

func getTestcase(prob string, useCache bool) (*Testcase, error) {
	hr, err := headerReader(prob, useCache)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch testcase_header for '%s': %v", prob, err)
	}
	defer hr.Close()

	b, err := ioutil.ReadAll(hr)
	if err != nil {
		return nil, fmt.Errorf("failed to read testcase_header for '%s': %v", prob, err)
	}

	t := &Testcase{}
	err = json.Unmarshal(b, t)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal testcase_header '%s': %v", string(b), err)
	}
	logdbg("testcase header: '%s'", t)

	if t.Length() == 0 {
		t.removeHeaderCache()
		return nil, fmt.Errorf("Testcase is not available for '%s'", prob)
	}

	return t, nil
}

func (cfi caseFileInfo) String() string {
	return fmt.Sprintf(`{"name":"%s","size":%d}`, cfi.Name, cfi.Size)
}

func (t *Testcase) String() string {
	s := fmt.Sprintf(`{"id":"%s","available":%d,"input":[`, t.ID, t.Available)
	sp := ""
	for _, c := range t.Input {
		s += fmt.Sprintf("%s%v", sp, c)
		sp = ","
	}
	s += `],"output":[`
	sp = ""
	for _, c := range t.Output {
		s += fmt.Sprintf("%s%v", sp, c)
		sp = ","
	}
	sp = ""
	s += `],"case_mapping":[`
	for _, c := range t.CaseMapping {
		s += fmt.Sprintf(`%s"%s"`, sp, c)
		sp = ","
	}
	s += `]}`
	return s
}

func (t *Testcase) Length() int {
	if t.Available > 0 {
		return len(t.Input)
	} else {
		return 0
	}
}

func (t *Testcase) CaseInput(idx int) (io.ReadCloser, error) {
	return t.caseReader(idx, caseTypeInput, true)
}

func (t *Testcase) CaseOutput(idx int) (io.ReadCloser, error) {
	return t.caseReader(idx, caseTypeOutput, true)
}

type caseType string

const (
	caseTypeInput  caseType = "in"
	caseTypeOutput caseType = "out"
)

func headerIsCached(prob string) bool {
	path := headerCachePath(prob)

	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func headerReader(prob string, useCache bool) (io.ReadCloser, error) {
	if !useCache || !headerIsCached(prob) {
		logdbg("headerReader: fetching test case header for '%s'", prob)
		if err := fetchHeader(prob); err != nil {
			return nil, err
		}
		logdbg("headerReader: fetched test case header for '%s'", prob)
	}

	path := headerCachePath(prob)
	logdbg("headerReader: path: %s", path)

	return os.Open(path)
}

func fetchHeader(prob string) error {
	dir := headerCacheDir()
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return fmt.Errorf("Failed to create cache dir: %s: %v", dir, err)
	}

	url := fmt.Sprintf("%s?id=%s", TESTCASE_HEADER_ENDPOINT, prob)
	c := &http.Client{}
	resp, err := c.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("http status code=%d", resp.StatusCode)
	}

	path := headerCachePath(prob)
	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("failed to open %s: %v", path, err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (t *Testcase) removeHeaderCache() error {
	if !headerIsCached(t.ID) {
		return nil
	}
	path := headerCachePath(t.ID)
	return os.Remove(path)
}

func (t *Testcase) cfi(idx int, ct caseType) (*caseFileInfo, error) {
	cfis := t.Input
	if ct == caseTypeOutput {
		cfis = t.Output
	}

	if idx >= len(cfis) {
		return nil, fmt.Errorf("out of range")
	}

	return &cfis[idx], nil
}

func cacheDir() string {
	home, ok := os.LookupEnv("HOME")
	if !ok || home == "" {
		logfatal("HOME is not set")
	}
	return fmt.Sprintf("%s/%s", home, CACHE_DIR)
}

func headerCacheDir() string {
	return fmt.Sprintf("%s/testcase_header", cacheDir())
}

func headerCachePath(prob string) string {
	return fmt.Sprintf("%s/%s", headerCacheDir(), prob)
}

func (t *Testcase) caseCacheDir() string {
	return fmt.Sprintf("%s/testcase/%s", cacheDir(), t.ID)
}

// Returns cache file path.
func (t *Testcase) caseCachePath(idx int, ct caseType) (string, error) {
	cfi, err := t.cfi(idx, ct)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", t.caseCacheDir(), cfi.Name), nil
}

func (t *Testcase) caseCacheSize(idx int, ct caseType) (int64, error) {
	cfi, err := t.cfi(idx, ct)
	if err != nil {
		return 0, err
	}

	return cfi.Size, nil
}

func (t *Testcase) caseIsCached(idx int, ct caseType) bool {
	path, err := t.caseCachePath(idx, ct)
	if err != nil {
		logdbg("failed to determine cache path for %sput test case[%d]: %v", ct, idx, err)
		return false
	}

	cfi, err := t.cfi(idx, ct)
	if err != nil {
		logdbg("failed to get cfi for %sput test case[%d]: %v", ct, idx, err)
		return false
	}

	st, err := os.Stat(path)
	if err != nil {
		logdbg("failed to stat cache file for %sput test case[%d]: %s: %v", ct, idx, path, err)
		return false
	}
	if st.Size() != cfi.Size {
		logdbg("cache file for %sput test case[%d] has wrong size: %s: %v", ct, idx, path, err)
		return false
	}
	logdbg("found valid cache file for %sput test case[%d]: %s", ct, idx, path)
	return true
}

func (t *Testcase) caseReader(idx int, ct caseType, useCache bool) (io.ReadCloser, error) {
	if !useCache || !t.caseIsCached(idx, ct) {
		logdbg("caseReader: fetching %sput test case[%d] for %s", ct, idx, t.ID)
		if err := t.fetchCase(idx, ct); err != nil {
			return nil, err
		}
		logdbg("caseReader: fetched %sput test case[%d] for %s", ct, idx, t.ID)
	}
	logdbg("caseReader: %sput test case[%d] for %s is cached", ct, idx, t.ID)

	path, err := t.caseCachePath(idx, ct)
	if err != nil {
		return nil, err
	}
	logdbg("caseReader: path: %s", path)

	return os.Open(path)
}

func (t *Testcase) fetchCase(idx int, ct caseType) error {
	dir := t.caseCacheDir()
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return fmt.Errorf("Failed to create cache dir: %s: %v", dir, err)
	}

	url := fmt.Sprintf("%s?id=%s&case=%d&type=%s", TESTCASE_ENDPOINT, t.ID, idx+1, ct)
	logdbg("fetchCase: url: %s", url)
	c := &http.Client{}
	resp, err := c.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("http status code=%d", resp.StatusCode)
	}

	path, err := t.caseCachePath(idx, ct)
	if err != nil {
		return err
	}

	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("failed to open %s: %v", path, err)
	}
	defer out.Close()

	sz, err := io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	esz, err := t.caseCacheSize(idx, ct)
	if err != nil {
		return err
	}

	if sz != esz {
		return fmt.Errorf("fetched file size mismatch: %s: expected=%d, received=%d", path, esz, sz)
	}

	return nil
}

func (t *Testcase) FetchAllCases() error {
	n := len(t.Input)
	if on := len(t.Output); on < n {
		n = on
	}

	for i := 0; i < n; i++ {
		for _, ct := range []caseType{caseTypeInput, caseTypeOutput} {
			err := t.fetchCase(i, ct)
			if err != nil {
				return fmt.Errorf("failed to fetch the %s testcase[%d]: %v", ct, i, err)
			}
		}
	}
	return nil
}
