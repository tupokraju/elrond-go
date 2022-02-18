package factory

import (
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go/common"
	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/errors"
	"github.com/ElrondNetwork/elrond-go/heartbeat/processor"
	"github.com/ElrondNetwork/elrond-go/heartbeat/sender"
)

// ArgHeartbeatV2ComponentsFactory represents the argument for the heartbeat v2 components factory
type ArgHeartbeatV2ComponentsFactory struct {
	Config             config.Config
	Prefs              config.Preferences
	AppVersion         string
	BoostrapComponents BootstrapComponentsHolder
	CoreComponents     CoreComponentsHolder
	DataComponents     DataComponentsHolder
	NetworkComponents  NetworkComponentsHolder
	CryptoComponents   CryptoComponentsHolder
	ProcessComponents  ProcessComponentsHolder
}

type heartbeatV2ComponentsFactory struct {
	config             config.Config
	prefs              config.Preferences
	version            string
	boostrapComponents BootstrapComponentsHolder
	coreComponents     CoreComponentsHolder
	dataComponents     DataComponentsHolder
	networkComponents  NetworkComponentsHolder
	cryptoComponents   CryptoComponentsHolder
	processComponents  ProcessComponentsHolder
}

type heartbeatV2Components struct {
	sender    HeartbeatV2Sender
	processor PeerAuthenticationRequestsProcessor
}

// NewHeartbeatV2ComponentsFactory creates a new instance of heartbeatV2ComponentsFactory
func NewHeartbeatV2ComponentsFactory(args ArgHeartbeatV2ComponentsFactory) (*heartbeatV2ComponentsFactory, error) {
	err := checkHeartbeatV2FactoryArgs(args)
	if err != nil {
		return nil, err
	}

	return &heartbeatV2ComponentsFactory{
		config:             args.Config,
		prefs:              args.Prefs,
		version:            args.AppVersion,
		boostrapComponents: args.BoostrapComponents,
		coreComponents:     args.CoreComponents,
		dataComponents:     args.DataComponents,
		networkComponents:  args.NetworkComponents,
		cryptoComponents:   args.CryptoComponents,
		processComponents:  args.ProcessComponents,
	}, nil
}

func checkHeartbeatV2FactoryArgs(args ArgHeartbeatV2ComponentsFactory) error {
	if check.IfNil(args.BoostrapComponents) {
		return errors.ErrNilBootstrapComponentsHolder
	}
	if check.IfNil(args.CoreComponents) {
		return errors.ErrNilCoreComponentsHolder
	}
	if check.IfNil(args.DataComponents) {
		return errors.ErrNilDataComponentsHolder
	}
	if check.IfNil(args.NetworkComponents) {
		return errors.ErrNilNetworkComponentsHolder
	}
	if check.IfNil(args.CryptoComponents) {
		return errors.ErrNilCryptoComponentsHolder
	}
	if check.IfNil(args.ProcessComponents) {
		return errors.ErrNilProcessComponentsHolder
	}

	return nil
}

// Create creates the heartbeatV2 components
func (hcf *heartbeatV2ComponentsFactory) Create() (*heartbeatV2Components, error) {
	peerSubType := core.RegularPeer
	if hcf.prefs.Preferences.FullArchive {
		peerSubType = core.FullHistoryObserver
	}

	cfg := hcf.config.HeartbeatV2

	argsSender := sender.ArgSender{
		Messenger:                          hcf.networkComponents.NetworkMessenger(),
		Marshaller:                         hcf.coreComponents.InternalMarshalizer(),
		PeerAuthenticationTopic:            common.PeerAuthenticationTopic,
		HeartbeatTopic:                     common.HeartbeatV2Topic,
		PeerAuthenticationTimeBetweenSends: time.Second * time.Duration(cfg.PeerAuthenticationTimeBetweenSendsInSec),
		PeerAuthenticationTimeBetweenSendsWhenError: time.Second * time.Duration(cfg.PeerAuthenticationTimeBetweenSendsWhenErrorInSec),
		PeerAuthenticationThresholdBetweenSends:     cfg.PeerAuthenticationThresholdBetweenSends,
		HeartbeatTimeBetweenSends:                   time.Second * time.Duration(cfg.HeartbeatTimeBetweenSendsInSec),
		HeartbeatTimeBetweenSendsWhenError:          time.Second * time.Duration(cfg.HeartbeatTimeBetweenSendsWhenErrorInSec),
		HeartbeatThresholdBetweenSends:              cfg.HeartbeatThresholdBetweenSends,
		VersionNumber:                               hcf.version,
		NodeDisplayName:                             hcf.prefs.Preferences.NodeDisplayName,
		Identity:                                    hcf.prefs.Preferences.Identity,
		PeerSubType:                                 peerSubType,
		CurrentBlockProvider:                        hcf.dataComponents.Blockchain(),
		PeerSignatureHandler:                        hcf.cryptoComponents.PeerSignatureHandler(),
		PrivateKey:                                  hcf.cryptoComponents.PrivateKey(),
		RedundancyHandler:                           hcf.processComponents.NodeRedundancyHandler(),
	}
	heartbeatV2Sender, err := sender.NewSender(argsSender)
	if err != nil {
		return nil, err
	}

	epochBootstrapParams := hcf.boostrapComponents.EpochBootstrapParams()
	argsProcessor := processor.ArgPeerAuthenticationRequestsProcessor{
		RequestHandler:           hcf.processComponents.RequestHandler(),
		NodesCoordinator:         hcf.processComponents.NodesCoordinator(),
		PeerAuthenticationPool:   hcf.dataComponents.Datapool().PeerAuthentications(),
		ShardId:                  epochBootstrapParams.SelfShardID(),
		Epoch:                    epochBootstrapParams.Epoch(),
		MessagesInChunk:          uint32(cfg.MaxNumOfPeerAuthenticationInResponse),
		MinPeersThreshold:        cfg.MinPeersThreshold,
		DelayBetweenRequests:     time.Second * time.Duration(cfg.DelayBetweenRequestsInSec),
		MaxTimeout:               time.Second * time.Duration(cfg.MaxTimeoutInSec),
		MaxMissingKeysInResponse: cfg.MaxMissingKeysInResponse,
	}
	paRequestsProcessor, err := processor.NewPeerAuthenticationRequestsProcessor(argsProcessor)
	if err != nil {
		return nil, err
	}

	return &heartbeatV2Components{
		sender:    heartbeatV2Sender,
		processor: paRequestsProcessor,
	}, nil
}

// Close closes the heartbeat components
func (hc *heartbeatV2Components) Close() error {
	log.Debug("calling close on heartbeatV2 components")

	if !check.IfNil(hc.sender) {
		log.LogIfError(hc.sender.Close())
	}

	if !check.IfNil(hc.processor) {
		log.LogIfError(hc.processor.Close())
	}

	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (hcf *heartbeatV2ComponentsFactory) IsInterfaceNil() bool {
	return hcf == nil
}
