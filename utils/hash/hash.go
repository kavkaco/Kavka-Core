package hash

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

type HashParams struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

var DefaultHashParams = &HashParams{
	memory:      64 * 1024,
	iterations:  10,
	parallelism: 2,
	saltLength:  16,
	keyLength:   32,
}

type HashManager struct {
	params *HashParams
}

func NewHashManager(params *HashParams) *HashManager {
	return &HashManager{params}
}

func (h *HashManager) HashPassword(password string) (hashedPassword string, err error) {
	salt, err := h.generateRandomBytes(h.params.saltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, h.params.iterations, h.params.memory, h.params.parallelism, h.params.keyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	hashedPassword = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, h.params.memory, h.params.iterations, h.params.parallelism, b64Salt, b64Hash)

	return hashedPassword, nil
}

func (h *HashManager) CheckPasswordHash(password, hashedPassword string) bool {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	p, salt, hash, err := h.decodeHash(hashedPassword)
	if err != nil {
		return false
	}

	otherHash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	return subtle.ConstantTimeCompare(hash, otherHash) == 1
}

func (h *HashManager) generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (h *HashManager) decodeHash(encodedHash string) (params *HashParams, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	p := &HashParams{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
}
