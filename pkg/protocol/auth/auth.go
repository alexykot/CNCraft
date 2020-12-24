package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/alexykot/cncraft/pkg/protocol/auth/mojang"
)

// A handles player authentication before they are allowed to join the server.
type A interface {
	BootstrapUser(userID uuid.UUID, username string) error
	GetUserPubkey(userID uuid.UUID) []byte
	GetUserVerifyToken(userID uuid.UUID) []byte
	DecryptUserSharedSecret(userID uuid.UUID, encSharedSecret []byte) ([]byte, error)
	RunMojangSessionAuth(userID uuid.UUID, sharedSecret []byte) (*mojang.AuthResponse, error)
}

var authRunner auther // the singleton

const verifyTokenLength = 16

// allowed length of the shared secret corresponding to three types of AES encryption
const secretAES128 = 16
const secretAES192 = 24
const secretAES256 = 32

type auther struct {
	// DEBT With auther being a singleton this is global state. Hmm.
	stagingUsers map[uuid.UUID]stagingUser
	mu           sync.Mutex
}

type stagingUser struct {
	username      string
	secretCrypter *mojang.RSACrypter
	verifyToken   []byte
}

func init() {
	authRunner = auther{
		stagingUsers: make(map[uuid.UUID]stagingUser),
	}
}

func GetAuther() A {
	return &authRunner
}

func (a *auther) BootstrapUser(userID uuid.UUID, username string) error {
	crypter, err := mojang.NewRSACrypter()
	if err != nil {
		return fmt.Errorf("failed to create new crypter: %w", err)
	}

	token := make([]byte, verifyTokenLength)
	if resLen, err := rand.Read(token); err != nil {
		return fmt.Errorf("failed to generate random token: %w", err)
	} else if resLen != verifyTokenLength {
		return fmt.Errorf("generated random token is incorrect size: %w", err)
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if _, ok := a.stagingUsers[userID]; ok {
		return fmt.Errorf("user with userID %s already exists", userID.String())
	}

	a.stagingUsers[userID] = stagingUser{
		username:      username,
		secretCrypter: crypter,
		verifyToken:   token,
	}

	return nil
}

func (a *auther) GetUserPubkey(userID uuid.UUID) []byte {
	user, ok := a.stagingUsers[userID]
	if !ok {
		return nil
	}
	return user.secretCrypter.GetPubKey()
}

func (a *auther) GetUserVerifyToken(userID uuid.UUID) []byte {
	user, ok := a.stagingUsers[userID]
	if !ok {
		return nil
	}
	return user.verifyToken
}

// SetUserSharedSecret decrypts and saves shared secret for given user.
func (a *auther) DecryptUserSharedSecret(userID uuid.UUID, encSharedSecret []byte) ([]byte, error) {
	user, ok := a.stagingUsers[userID]
	if !ok {
		return nil, fmt.Errorf("stagingUser %s not found", userID.String())
	}

	sharedSecret, err := user.secretCrypter.Decrypt(encSharedSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt shared secret: %w", err)
	}
	if len(sharedSecret) != secretAES128 && len(sharedSecret) != secretAES192 && len(sharedSecret) != secretAES256 {
		return nil, fmt.Errorf("received shared secret of unexpected length %d", len(sharedSecret))
	}

	return sharedSecret, nil
}

func (a *auther) RunMojangSessionAuth(userID uuid.UUID, sharedSecret []byte) (*mojang.AuthResponse, error) {
	user, ok := a.stagingUsers[userID]
	if !ok {
		return nil, fmt.Errorf("stagingUser %s not found", userID.String())
	}

	auth, err := mojang.RunMojangSessionAuth(user.username, user.secretCrypter.GetPubKey(), sharedSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to run Mojang authentication: %w", err)
	}

	return auth, nil
}

func (a *auther) LoginSuccess(userID uuid.UUID) error {
	return errors.New("unimplemented")
}
