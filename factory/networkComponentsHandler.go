package factory

import (
	"context"
	"sync"

	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/p2p"
	"github.com/ElrondNetwork/elrond-go/process"
)

var _ ComponentHandler = (*managedNetworkComponents)(nil)
var _ NetworkComponentsHolder = (*managedNetworkComponents)(nil)
var _ NetworkComponentsHandler = (*managedNetworkComponents)(nil)

// NetworkComponentsHandlerArgs holds the arguments to create a network component handler instance
type NetworkComponentsHandlerArgs struct {
	P2pConfig     config.P2PConfig
	MainConfig    config.Config
	StatusHandler core.AppStatusHandler
}

// managedNetworkComponents creates the data components handler that can create, close and access the data components
type managedNetworkComponents struct {
	*networkComponents
	networkComponentsFactory *networkComponentsFactory
	cancelFunc               func()
	mutNetworkComponents     sync.RWMutex
}

// NewManagedNetworkComponents creates a new data components handler
func NewManagedNetworkComponents(args NetworkComponentsHandlerArgs) (*managedNetworkComponents, error) {
	ncf, err := NewNetworkComponentsFactory(
		args.P2pConfig,
		args.MainConfig,
		args.StatusHandler,
	)
	if err != nil {
		return nil, err
	}

	return &managedNetworkComponents{
		networkComponents:        nil,
		networkComponentsFactory: ncf,
	}, nil
}

// Create creates the network components
func (mnc *managedNetworkComponents) Create() error {
	nc, err := mnc.networkComponentsFactory.Create()
	if err != nil {
		return err
	}

	mnc.mutNetworkComponents.Lock()
	mnc.networkComponents = nc
	_, mnc.cancelFunc = context.WithCancel(context.Background())
	mnc.mutNetworkComponents.Unlock()

	return nil
}

// Close closes the network components
func (mnc *managedNetworkComponents) Close() error {
	mnc.mutNetworkComponents.Lock()
	defer mnc.mutNetworkComponents.Unlock()

	mnc.cancelFunc()
	err := mnc.netMessenger.Close()
	if err != nil {
		return err
	}

	mnc.cancelFunc = nil
	mnc.networkComponents = nil

	return nil
}

// NetworkMessenger returns the p2p messenger
func (mnc *managedNetworkComponents) NetworkMessenger() p2p.Messenger {
	mnc.mutNetworkComponents.RLock()
	defer mnc.mutNetworkComponents.RUnlock()

	if mnc.networkComponents == nil {
		return nil
	}

	return mnc.netMessenger
}

// InputAntiFloodHandler returns the input p2p anti-flood handler
func (mnc *managedNetworkComponents) InputAntiFloodHandler() P2PAntifloodHandler {
	mnc.mutNetworkComponents.RLock()
	defer mnc.mutNetworkComponents.RUnlock()

	if mnc.networkComponents == nil {
		return nil
	}

	return mnc.inputAntifloodHandler
}

// OutputAntiFloodHandler returns the output p2p anti-flood handler
func (mnc *managedNetworkComponents) OutputAntiFloodHandler() P2PAntifloodHandler {
	mnc.mutNetworkComponents.RLock()
	defer mnc.mutNetworkComponents.RUnlock()

	if mnc.networkComponents == nil {
		return nil
	}

	return mnc.outputAntifloodHandler
}

// PeerBlackListHandler returns the blacklist handler
func (mnc *managedNetworkComponents) PeerBlackListHandler() process.BlackListHandler {
	mnc.mutNetworkComponents.RLock()
	defer mnc.mutNetworkComponents.RUnlock()

	if mnc.networkComponents == nil {
		return nil
	}

	return mnc.networkComponents.peerBlackListHandler
}

// IsInterfaceNil returns true if the interface is nil
func (mnc *managedNetworkComponents) IsInterfaceNil() bool {
	return mnc == nil
}
