package credentials

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/denisbrodbeck/machineid"
	"github.com/thamaji/files"
	"github.com/thamaji/lazycrypto"
)

const appID = "github.com/kurusugawa-computer/ace"

type Credentials struct {
	OpenAIAPIKey string
}

type credentialsFile struct {
	OpenAIAPIKey string `json:"openai_api_key"`
}

// Credentialsを暗号化してJSON形式で保存する
func Save(appName string, credentials *Credentials) error {
	passphrase, err := machineid.ProtectedID(appID)
	if err != nil {
		return err
	}

	openAIAPIKey, err := lazycrypto.EncryptToString([]byte(passphrase), []byte(credentials.OpenAIAPIKey))
	if err != nil {
		return err
	}

	credentialsFile := credentialsFile{
		OpenAIAPIKey: string(openAIAPIKey),
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	path := filepath.Join(cacheDir, appName, "credentials.json")
	f, err := files.OpenFileWriter(path)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	err = enc.Encode(credentialsFile)
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}
	if err != nil {
		return err
	}

	return nil
}

// 暗号化されているCredentialsファイルを読み込んで複合化する
func Load(appName string) (*Credentials, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	credentialsFile := credentialsFile{}

	path := filepath.Join(cacheDir, appName, "credentials.json")
	f, err := files.OpenFileReader(path)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(f).Decode(&credentialsFile)
	f.Close()
	if err != nil {
		return nil, err
	}

	passphrase, err := machineid.ProtectedID(appID)
	if err != nil {
		return nil, err
	}

	openAIAPIKey, err := lazycrypto.DecryptString([]byte(passphrase), credentialsFile.OpenAIAPIKey)
	if err != nil {
		return nil, err
	}

	credentials := &Credentials{
		OpenAIAPIKey: string(openAIAPIKey),
	}

	return credentials, nil
}
