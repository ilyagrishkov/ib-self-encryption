package internal

import (
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Fabric struct {
	PublicKey string
	contract  *gateway.Contract
}

type Asset struct {
	ID    string   `json:"ID"`
	Owner string   `json:"Owner"`
	CID   []string `json:"CID"`
}

func NewFabric() Fabric {
	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	wallet := gateway.NewInMemoryWallet()

	if !wallet.Exists("appUser") {
		err = populateWallet(wallet)
		if err != nil {
			log.Fatalf("Failed to populate wallet contents: %v", err)
		}
	}

	ccpPath := viper.Get("conn_config").(string)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}

	err = os.Remove("keystore")
	if err != nil {
		log.Fatalf("Failed to remove keystore: %v", err)
	}

	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}

	identity, err := wallet.Get("appUser")
	block, _ := pem.Decode([]byte(identity.(*gateway.X509Identity).Certificate()))
	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return Fabric{}
	}
	key, err := x509.MarshalPKIXPublicKey(certificate.PublicKey)
	if err != nil {
		return Fabric{}
	}

	return Fabric{
		PublicKey: hex.EncodeToString(key),
		contract:  network.GetContract("ibse"),
	}
}

func (fabric Fabric) CreateAsset(id string, CIDs []string) error {
	CIDsJson, err := json.Marshal(CIDs)
	_, err = fabric.contract.SubmitTransaction("CreateAsset", id, fabric.PublicKey, string(CIDsJson))
	if err != nil {
		log.Fatalf("Failed to Create asset: %v", err)
	}
	return nil
}

func (fabric Fabric) ReadAsset(id string) (Asset, error) {
	result, err := fabric.contract.EvaluateTransaction("ReadAsset", id)
	if err != nil {
		return Asset{}, err
	}

	var asset Asset
	err = json.Unmarshal(result, &asset)
	if err != nil {
		return Asset{}, err
	}
	return asset, nil
}

func (fabric Fabric) ReadAllAssets() ([]Asset, error) {
	result, err := fabric.contract.EvaluateTransaction("GetAllAssets")
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}
	var assets []Asset
	err = json.Unmarshal(result, &assets)
	if err != nil {
		return nil, err
	}
	return assets, nil
}

func populateWallet(wallet *gateway.Wallet) error {
	credPath := viper.Get("cred_path").(string)

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))
	return wallet.Put("appUser", identity)
}
