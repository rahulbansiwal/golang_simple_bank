package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetrickey []byte
}

func NewPasetoMaker(symmetrickey string) (Maker, error) {
	if len(symmetrickey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid symmetric key size")
	}
	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetrickey: []byte(symmetrickey),
	}
	return maker, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, *Payload,error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "",payload, err
	}
	token,err:= maker.paseto.Encrypt(maker.symmetrickey, payload, nil)
	return token,payload,err
}
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetrickey, payload, nil)
	if err != nil {
		return nil, err
	}
	if payload.ExpiredAt.Before(time.Now()) {

		return nil, fmt.Errorf("payload is expired")
	}
	return payload, nil
}
