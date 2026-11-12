package login

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

// selfSigned encodes a login request whose multiplayer token the client signed
// itself, stating authType, as a client not logged in with XBOX Live sends.
func selfSigned(t *testing.T, authType uint8, id IdentityData) []byte {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	keyData := MarshalPublicKey(&key.PublicKey)
	signer, _ := jose.NewSigner(jose.SigningKey{Key: key, Algorithm: jose.ES384}, &jose.SignerOptions{
		ExtraHeaders: map[jose.HeaderKey]any{"x5u": keyData},
	})
	claims := jwt.Claims{
		Expiry:    jwt.NewNumericDate(time.Now().Add(time.Hour)),
		NotBefore: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		Audience:  jwt.Audience{"api://auth-minecraft-services/multiplayer"},
	}
	req := &request{AuthenticationType: authType, Certificate: certificate{Chain: chain{""}}}
	req.Token, _ = jwt.Signed(signer).Claims(tokenClaims{
		Claims:          claims,
		ClientPublicKey: keyData,
		XUID:            id.XUID,
		DisplayName:     id.DisplayName,
		Identity:        id.Identity,
	}).Serialize()
	req.RawToken, _ = jwt.Signed(signer).Claims(ClientData{
		GameVersion:       "1.21.0",
		DeviceID:          "a4f365bb1e04459bbe4cb3cbf9d546e0",
		DeviceOS:          7,
		LanguageCode:      "ko_KR",
		ServerAddress:     "127.0.0.1:19132",
		SkinID:            "s",
		SkinImageWidth:    1,
		SkinImageHeight:   1,
		SkinData:          "AAAAAA==",
		SkinResourcePatch: "e30=",
	}).Serialize()
	return encodeRequest(req)
}

// verifier verifies tokens signed by a key no client holds, like the real one.
func verifier(t *testing.T) *oidc.IDTokenVerifier {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return oidc.NewVerifier("https://authorization.franchise.minecraft-services.net", &oidc.StaticKeySet{PublicKeys: []crypto.PublicKey{&key.PublicKey}}, &oidc.Config{SkipClientIDCheck: true})
}

func TestParseGuestNotAuthenticated(t *testing.T) {
	// A guest claiming an Xbox account's XUID and no name: nothing of it verifies.
	data := selfSigned(t, 1, IdentityData{XUID: "2535400000000001"})
	id, _, res, err := Parse(data, verifier(t))
	if err != nil {
		t.Fatalf("guest request refused by Parse: %v", err)
	}
	if res.XBOXLiveAuthenticated {
		t.Fatal("guest request reported as authenticated")
	}
	if res.AuthenticationType != 1 || res.VerifyErr == nil {
		t.Fatalf("AuthenticationType %d, VerifyErr %v", res.AuthenticationType, res.VerifyErr)
	}
	if res.IdentityErr == nil {
		t.Fatalf("empty name of %+v passed validation", id)
	}
}

func TestParseUnverifiedXboxRefused(t *testing.T) {
	// A client stating it is logged in with XBOX Live whose token does not verify.
	data := selfSigned(t, 0, IdentityData{XUID: "2535400000000001", DisplayName: "Steve"})
	if _, _, _, err := Parse(data, verifier(t)); err == nil {
		t.Fatal("an XBOX Live request whose token does not verify was parsed")
	}
}

func TestParseSelfSignedWithoutVerifier(t *testing.T) {
	// Authentication disabled: a valid self-signed identity parses as before.
	data := selfSigned(t, 2, IdentityData{DisplayName: "Steve"})
	id, _, res, err := Parse(data, nil)
	if err != nil || res.IdentityErr != nil || res.XBOXLiveAuthenticated || id.DisplayName != "Steve" {
		t.Fatalf("id %+v, res %+v, err %v", id, res, err)
	}
}
