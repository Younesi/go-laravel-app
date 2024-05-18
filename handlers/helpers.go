package handlers

import (
	"github.com/younesi/atlas"
)

func (h *Handlers) encrypt(text string) (string, error) {
	enc := atlas.Encryption{Key: []byte(h.App.EncryptionKey)}

	encrypted, err := enc.Encrypt(text)
	if err != nil {
		return "", err
	}

	return encrypted, nil
}

func (h *Handlers) decrypt(text string) (string, error) {
	decrypt := atlas.Encryption{Key: []byte(h.App.EncryptionKey)}

	decrypted, err := decrypt.Decrypt(text)
	if err != nil {
		return "", err
	}

	return decrypted, nil
}
