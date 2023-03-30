package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/calvindc/Web3RpcHub/internal/repository"
	"github.com/gorilla/securecookie"
)

// LoadOrCreateCookieSecrets either parses the bytes from $repo/web/cookie-secret or creates a new file with suitable keys in it
func LoadOrCreateCookieSecrets(repo repository.Interface) ([]securecookie.Codec, error) {
	secretPath := repo.GetPath("web", "cookie-secret")
	err := os.MkdirAll(filepath.Dir(secretPath), 0700)
	if err != nil && !os.IsExist(err) {
		return nil, fmt.Errorf("failed to create folder for cookie secret: %w", err)
	}

	// 加载本地已经存在的
	secrets, err := ioutil.ReadFile(secretPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load cookie secrets: %w", err)
		}

		// create new keys, save them and return the codec
		hashKey := securecookie.GenerateRandomKey(32)
		blockKey := securecookie.GenerateRandomKey(32)

		data := append(hashKey, blockKey...)
		err = ioutil.WriteFile(secretPath, data, 0600)
		if err != nil {
			return nil, err
		}
		sc := securecookie.CodecsFromPairs(hashKey, blockKey)
		return sc, nil
	}

	// secrets should contain multiple of 64byte (to enable key rotation as supported by gorilla)
	if n := len(secrets); n%64 != 0 {
		return nil, fmt.Errorf("expected multiple of 64bytes in cookie secret file but got: %d", n)
	}

	// range over the secrets []byte in chunks of 64 bytes
	// and slice it into 32byte pairs
	var pairs [][]byte

	// the increment/next part (which usually is i++)
	// is the multiple comma assigment (a,b = b+1,a-1)
	// so chunk is the next 64 bytes and then it slices of the first 64 bytes of secrets for the next iteration
	for chunk := secrets[:64]; len(secrets) >= 64; chunk, secrets = secrets[:64], secrets[64:] {
		pairs = append(pairs,
			chunk[0:32],  // hash key
			chunk[32:64], // block key
		)
	}

	sc := securecookie.CodecsFromPairs(pairs...)
	return sc, nil
}
