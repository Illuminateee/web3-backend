// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package fiattotokenpaymentgateway

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// FiatToTokenPaymentGatewayMetaData contains all meta data concerning the FiatToTokenPaymentGateway contract.
var FiatToTokenPaymentGatewayMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_pricePerToken\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_requiredGasDeposit\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newRequirement\",\"type\":\"uint256\"}],\"name\":\"GasDepositRequirementUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"paymentId\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"buyer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"GasRefunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"paymentId\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"buyer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenAmount\",\"type\":\"uint256\"}],\"name\":\"PaymentCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"paymentId\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"buyer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"destinationWallet\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fiatAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"gateway\",\"type\":\"string\"}],\"name\":\"PaymentCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"paymentId\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"buyer\",\"type\":\"address\"}],\"name\":\"PaymentFailed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"paymentId\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"buyer\",\"type\":\"address\"}],\"name\":\"PaymentRefunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newPrice\",\"type\":\"uint256\"}],\"name\":\"TokenPriceUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"calculateProcessingCost\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"fiatAmount\",\"type\":\"uint256\"}],\"name\":\"calculateTokenAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"paymentId\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"tokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fiatAmount\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"gateway\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"destinationWallet\",\"type\":\"address\"}],\"name\":\"createPayment\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"gatewaySigners\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"paymentId\",\"type\":\"string\"}],\"name\":\"getPaymentStatus\",\"outputs\":[{\"internalType\":\"enumFiatToTokenPaymentGateway.PaymentStatus\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"paymentId\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"mockPaymentCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"payments\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"buyer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"destinationWallet\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fiatAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"gateway\",\"type\":\"string\"},{\"internalType\":\"enumFiatToTokenPaymentGateway.PaymentStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasFundAmount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"gasRefunded\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pricePerToken\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"paymentId\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"processPaymentCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"paymentId\",\"type\":\"string\"}],\"name\":\"processRefund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requiredGasDeposit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requiredGasDeposit\",\"type\":\"uint256\"}],\"name\":\"updateGasDepositRequirement\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"gateway\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"updateGatewaySigner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_pricePerToken\",\"type\":\"uint256\"}],\"name\":\"updateTokenPrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawProcessingFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x608060405234801562000010575f80fd5b506040516200251f3803806200251f8339810160408190526200003391620001bf565b33806200005a57604051631e4fbdf760e01b81525f60048201526024015b60405180910390fd5b620000658162000170565b50600180556001600160a01b038316620000c25760405162461bcd60e51b815260206004820152601560248201527f496e76616c696420746f6b656e20616464726573730000000000000000000000604482015260640162000051565b600280546001600160a01b03949094166001600160a01b03199485161790556006918255600390815560408051676d69647472616e7360c01b8152600560088201819052825160289281900392909201822080548716339081179091556262636160e81b8352938201819052825160239281900392909201822080548716851790556573747269706560d01b8252938101939093525191829003602601909120805490921617905562000200565b5f80546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b5f805f60608486031215620001d2575f80fd5b83516001600160a01b0381168114620001e9575f80fd5b602085015160409095015190969495509392505050565b612311806200020e5f395ff3fe608060405260043610610113575f3560e01c80637b1b1de61161009d578063ee3efd5811610062578063ee3efd58146102da578063f2fde38b146102f9578063f946871314610318578063fc0c546a1461034c578063fd2240d31461036b575f80fd5b80637b1b1de6146102535780638da5cb5b146102765780638ff1148814610292578063a24bcf46146102a6578063e48770df146102c5575f80fd5b80634e981b68116100e35780634e981b681461019157806352c3cc2e146101ee5780635bb8e4d11461020d578063676c0d7714610220578063715018a61461023f575f80fd5b806303ff03991461011e5780630ac6e1311461013f57806316692d471461015e57806334eec86a1461017d575f80fd5b3661011a57005b5f80fd5b348015610129575f80fd5b5061013d610138366004611c4a565b610397565b005b34801561014a575f80fd5b5061013d610159366004611c95565b6107b1565b348015610169575f80fd5b5061013d610178366004611d2d565b610ca2565b348015610188575f80fd5b5061013d610d3d565b34801561019c575f80fd5b506101d16101ab366004611d6f565b80516020818301810180516005825292820191909301209152546001600160a01b031681565b6040516001600160a01b0390911681526020015b60405180910390f35b3480156101f9575f80fd5b5061013d610208366004611da9565b610e25565b61013d61021b366004611dc0565b610e69565b34801561022b575f80fd5b5061013d61023a366004611da9565b61138a565b34801561024a575f80fd5b5061013d611416565b34801561025e575f80fd5b5061026860065481565b6040519081526020016101e5565b348015610281575f80fd5b505f546001600160a01b03166101d1565b34801561029d575f80fd5b50610268611429565b3480156102b1575f80fd5b506102686102c0366004611da9565b61143e565b3480156102d0575f80fd5b5061026860035481565b3480156102e5575f80fd5b5061013d6102f4366004611d6f565b611465565b348015610304575f80fd5b5061013d610313366004611e42565b611798565b348015610323575f80fd5b50610337610332366004611d6f565b6117d2565b6040516101e599989796959493929190611edc565b348015610357575f80fd5b506002546101d1906001600160a01b031681565b348015610376575f80fd5b5061038a610385366004611d6f565b6118c6565b6040516101e59190611f4a565b61039f6118f3565b5f6004836040516103b09190611f58565b90815260405190819003602001902080549091506001600160a01b03166103f25760405162461bcd60e51b81526004016103e990611f73565b60405180910390fd5b5f600682015460ff16600381111561040c5761040c611ea8565b146104595760405162461bcd60e51b815260206004820152601c60248201527f5061796d656e74206e6f7420696e2070656e64696e672073746174650000000060448201526064016103e9565b8160ff1660010361058c5760068101805460ff19166001908117909155810154600280549083015460405163a9059cbb60e01b81526001600160a01b0393841660048201819052602482019290925290925f92169063a9059cbb906044016020604051808303815f875af11580156104d3573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906104f79190611f9e565b90508061053e5760405162461bcd60e51b8152602060048201526015602482015274151bdad95b881d1c985b9cd9995c8819985a5b1959605a1b60448201526064016103e9565b816001600160a01b03167f13d1118d2532a9b21c3e26bd2ba1b959263d90ea306c910ff6339e405192f1f286856002015460405161057d929190611fbd565b60405180910390a2505061068f565b60068101805460ff19166002908117909155545f805490916001600160a01b039081169163a9059cbb911660028501546040516001600160e01b031960e085901b1681526001600160a01b03909216600483015260248201526044016020604051808303815f875af1158015610604573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906106289190611f9e565b9050806106475760405162461bcd60e51b81526004016103e990611fde565b81546040516001600160a01b03909116907fd467a4459f62e0868cf38768c34059712f66ac7d9f3897703192a4a7e8a9db2b9061068590879061200b565b60405180910390a2505b600881015460ff161580156106a757505f8160070154115b156107ac575f6106b5611429565b90505f818360070154116106c9575f6106d9565b8183600701546106d99190612031565b905080156107a95760088301805460ff1916600117905582546040515f916001600160a01b03169083908381818185875af1925050503d805f8114610739576040519150601f19603f3d011682016040523d82523d5f602084013e61073e565b606091505b505090508061075f5760405162461bcd60e51b81526004016103e990612044565b83546040516001600160a01b03909116907f89092cc5430d6dbac7d084b6205265037b6589502bc0cedd66b2f46eefb79e469061079f9089908690611fbd565b60405180910390a2505b50505b505050565b6107b961191f565b5f6004846040516107ca9190611f58565b90815260405190819003602001902080549091506001600160a01b03166108035760405162461bcd60e51b81526004016103e990611f73565b5f600682015460ff16600381111561081d5761081d611ea8565b1461086a5760405162461bcd60e51b815260206004820152601c60248201527f5061796d656e74206e6f7420696e2070656e64696e672073746174650000000060448201526064016103e9565b5f848460405160200161087e92919061206f565b6040516020818303038152906040528051906020012090505f6108cd827f19457468657265756d205369676e6564204d6573736167653a0a3332000000005f908152601c91909152603c902090565b90505f6108da8286611949565b90506005846005016040516108ef91906120d8565b908152604051908190036020019020546001600160a01b0382811691161461094d5760405162461bcd60e51b8152602060048201526011602482015270496e76616c6964207369676e617475726560781b60448201526064016103e9565b8560ff16600103610a805760068401805460ff19166001908117909155840154600280549086015460405163a9059cbb60e01b81526001600160a01b0393841660048201819052602482019290925290925f92169063a9059cbb906044016020604051808303815f875af11580156109c7573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906109eb9190611f9e565b905080610a325760405162461bcd60e51b8152602060048201526015602482015274151bdad95b881d1c985b9cd9995c8819985a5b1959605a1b60448201526064016103e9565b816001600160a01b03167f13d1118d2532a9b21c3e26bd2ba1b959263d90ea306c910ff6339e405192f1f28a8860020154604051610a71929190611fbd565b60405180910390a25050610b83565b60068401805460ff19166002908117909155545f805490916001600160a01b039081169163a9059cbb911660028801546040516001600160e01b031960e085901b1681526001600160a01b03909216600483015260248201526044016020604051808303815f875af1158015610af8573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610b1c9190611f9e565b905080610b3b5760405162461bcd60e51b81526004016103e990611fde565b84546040516001600160a01b03909116907fd467a4459f62e0868cf38768c34059712f66ac7d9f3897703192a4a7e8a9db2b90610b79908b9061200b565b60405180910390a2505b5f610b8c611429565b90505f81866007015411610ba0575f610bb0565b818660070154610bb09190612031565b90505f81118015610bc65750600886015460ff16155b15610c935760088601805460ff1916600117905585546040515f916001600160a01b03169083908381818185875af1925050503d805f8114610c23576040519150601f19603f3d011682016040523d82523d5f602084013e610c28565b606091505b5050905080610c495760405162461bcd60e51b81526004016103e990612044565b86546040516001600160a01b03909116907f89092cc5430d6dbac7d084b6205265037b6589502bc0cedd66b2f46eefb79e4690610c89908d908690611fbd565b60405180910390a2505b5050505050506107ac60018055565b610caa6118f3565b6001600160a01b038116610cf95760405162461bcd60e51b8152602060048201526016602482015275496e76616c6964207369676e6572206164647265737360501b60448201526064016103e9565b80600583604051610d0a9190611f58565b90815260405190819003602001902080546001600160a01b03929092166001600160a01b03199092169190911790555050565b610d456118f3565b4780610d895760405162461bcd60e51b81526020600482015260136024820152724e6f206665657320746f20776974686472617760681b60448201526064016103e9565b5f80546040516001600160a01b039091169083908381818185875af1925050503d805f8114610dd3576040519150601f19603f3d011682016040523d82523d5f602084013e610dd8565b606091505b5050905080610e215760405162461bcd60e51b8152602060048201526015602482015274119959481dda5d1a191c985dd85b0819985a5b1959605a1b60448201526064016103e9565b5050565b610e2d6118f3565b60038190556040518181527f6ba7f9a0367ca0f8ad8598911d0fa547a28ea788a8b3c0607092d5455c2b583e906020015b60405180910390a150565b610e7161191f565b5f6001600160a01b0316600486604051610e8b9190611f58565b908152604051908190036020019020546001600160a01b031614610ef15760405162461bcd60e51b815260206004820152601960248201527f5061796d656e7420494420616c7265616479206578697374730000000000000060448201526064016103e9565b5f8411610f4f5760405162461bcd60e51b815260206004820152602660248201527f546f6b656e20616d6f756e74206d7573742062652067726561746572207468616044820152656e207a65726f60d01b60648201526084016103e9565b5f8311610fac5760405162461bcd60e51b815260206004820152602560248201527f4669617420616d6f756e74206d7573742062652067726561746572207468616e604482015264207a65726f60d81b60648201526084016103e9565b5f825111610ffc5760405162461bcd60e51b815260206004820152601960248201527f47617465776179206d757374206265207370656369666965640000000000000060448201526064016103e9565b60035434101561104e5760405162461bcd60e51b815260206004820152601860248201527f496e73756666696369656e7420676173206465706f736974000000000000000060448201526064016103e9565b60025484906001600160a01b031663dd62ed3e6110725f546001600160a01b031690565b6040516001600160e01b031960e084901b1681526001600160a01b039091166004820152306024820152604401602060405180830381865afa1580156110ba573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906110de919061214a565b101561112c5760405162461bcd60e51b815260206004820152601c60248201527f496e73756666696369656e7420746f6b656e20616c6c6f77616e63650000000060448201526064016103e9565b60408051610120810182523381526001600160a01b03831660208201529081018590526060810184905242608082015260a0810183905260c081015f81523460208201525f60409182015251600490611186908890611f58565b9081526040805160209281900383019020835181546001600160a01b03199081166001600160a01b03928316178355938501516001830180549095169116179092558201516002820155606082015160038201556080820151600482015560a082015160058201906111f890826121ae565b5060c082015160068201805460ff1916600183600381111561121c5761121c611ea8565b021790555060e08201516007820155610100909101516008909101805491151560ff199092169190911790556002545f906001600160a01b03166323b872dd61126c5f546001600160a01b031690565b6040516001600160e01b031960e084901b1681526001600160a01b039091166004820152306024820152604481018890526064016020604051808303815f875af11580156112bc573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906112e09190611f9e565b90508061132f5760405162461bcd60e51b815260206004820152601860248201527f546f6b656e207265736572766174696f6e206661696c6564000000000000000060448201526064016103e9565b816001600160a01b0316336001600160a01b03167f73e020d64bf5e24a99ceab0ef0f92283d5e48bd8700e072769bf844d6f0cf34e88888888604051611378949392919061226a565b60405180910390a3506107a960018055565b6113926118f3565b5f81116113e15760405162461bcd60e51b815260206004820152601f60248201527f5072696365206d7573742062652067726561746572207468616e207a65726f0060448201526064016103e9565b60068190556040518181527f6b361c807733cc94a3b7bda39c1c94dcc0511e91aa47b0f7f918068e9397689090602001610e5e565b61141e6118f3565b6114275f611971565b565b5f600260035461143991906122a5565b905090565b6006545f9061145583670de0b6b3a76400006122c4565b61145f91906122a5565b92915050565b61146d6118f3565b61147561191f565b5f6004826040516114869190611f58565b90815260405190819003602001902080549091506001600160a01b03166114bf5760405162461bcd60e51b81526004016103e990611f73565b5f600682015460ff1660038111156114d9576114d9611ea8565b14806114fd57506001600682015460ff1660038111156114fb576114fb611ea8565b145b6115495760405162461bcd60e51b815260206004820152601a60248201527f5061796d656e742063616e6e6f7420626520726566756e64656400000000000060448201526064016103e9565b5f600682015460ff16600381111561156357611563611ea8565b036116285760068101805460ff191660031790556002545f805490916001600160a01b039081169163a9059cbb911660028501546040516001600160e01b031960e085901b1681526001600160a01b03909216600483015260248201526044016020604051808303815f875af11580156115df573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906116039190611f9e565b9050806116225760405162461bcd60e51b81526004016103e990611fde565b50611658565b6001600682015460ff16600381111561164357611643611ea8565b036116585760068101805460ff191660031790555b600881015460ff1615801561167057505f8160070154115b156117455760088101805460ff19166001179055805460078201546040515f926001600160a01b031691908381818185875af1925050503d805f81146116d1576040519150601f19603f3d011682016040523d82523d5f602084013e6116d6565b606091505b50509050806116f75760405162461bcd60e51b81526004016103e990612044565b815460078301546040516001600160a01b03909216917f89092cc5430d6dbac7d084b6205265037b6589502bc0cedd66b2f46eefb79e469161173b91879190611fbd565b60405180910390a2505b80546040516001600160a01b03909116907fb7309985c95c87cf571a649294c1e24320da430df8238d2ec65c732619f8e53b9061178390859061200b565b60405180910390a25061179560018055565b50565b6117a06118f3565b6001600160a01b0381166117c957604051631e4fbdf760e01b81525f60048201526024016103e9565b61179581611971565b805180820160209081018051600480835293830192909401919091209290528154600183015460028401546003850154938501546005860180546001600160a01b0395861697949095169592949293919261182c906120a0565b80601f0160208091040260200160405190810160405280929190818152602001828054611858906120a0565b80156118a35780601f1061187a576101008083540402835291602001916118a3565b820191905f5260205f20905b81548152906001019060200180831161188657829003601f168201915b5050505060068301546007840154600890940154929360ff918216939092501689565b5f6004826040516118d79190611f58565b9081526040519081900360200190206006015460ff1692915050565b5f546001600160a01b031633146114275760405163118cdaa760e01b81523360048201526024016103e9565b60026001540361194257604051633ee5aeb560e01b815260040160405180910390fd5b6002600155565b5f805f8061195786866119c0565b9250925092506119678282611a09565b5090949350505050565b5f80546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b5f805f83516041036119f7576020840151604085015160608601515f1a6119e988828585611ac1565b955095509550505050611a02565b505081515f91506002905b9250925092565b5f826003811115611a1c57611a1c611ea8565b03611a25575050565b6001826003811115611a3957611a39611ea8565b03611a575760405163f645eedf60e01b815260040160405180910390fd5b6002826003811115611a6b57611a6b611ea8565b03611a8c5760405163fce698f760e01b8152600481018290526024016103e9565b6003826003811115611aa057611aa0611ea8565b03610e21576040516335e2f38360e21b8152600481018290526024016103e9565b5f80807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0841115611afa57505f91506003905082611b7f565b604080515f808252602082018084528a905260ff891692820192909252606081018790526080810186905260019060a0016020604051602081039080840390855afa158015611b4b573d5f803e3d5ffd5b5050604051601f1901519150506001600160a01b038116611b7657505f925060019150829050611b7f565b92505f91508190505b9450945094915050565b634e487b7160e01b5f52604160045260245ffd5b5f67ffffffffffffffff80841115611bb757611bb7611b89565b604051601f8501601f19908116603f01168101908282118183101715611bdf57611bdf611b89565b81604052809350858152868686011115611bf7575f80fd5b858560208301375f602087830101525050509392505050565b5f82601f830112611c1f575f80fd5b611c2e83833560208501611b9d565b9392505050565b803560ff81168114611c45575f80fd5b919050565b5f8060408385031215611c5b575f80fd5b823567ffffffffffffffff811115611c71575f80fd5b611c7d85828601611c10565b925050611c8c60208401611c35565b90509250929050565b5f805f60608486031215611ca7575f80fd5b833567ffffffffffffffff80821115611cbe575f80fd5b611cca87838801611c10565b9450611cd860208701611c35565b93506040860135915080821115611ced575f80fd5b508401601f81018613611cfe575f80fd5b611d0d86823560208401611b9d565b9150509250925092565b80356001600160a01b0381168114611c45575f80fd5b5f8060408385031215611d3e575f80fd5b823567ffffffffffffffff811115611d54575f80fd5b611d6085828601611c10565b925050611c8c60208401611d17565b5f60208284031215611d7f575f80fd5b813567ffffffffffffffff811115611d95575f80fd5b611da184828501611c10565b949350505050565b5f60208284031215611db9575f80fd5b5035919050565b5f805f805f60a08688031215611dd4575f80fd5b853567ffffffffffffffff80821115611deb575f80fd5b611df789838a01611c10565b965060208801359550604088013594506060880135915080821115611e1a575f80fd5b50611e2788828901611c10565b925050611e3660808701611d17565b90509295509295909350565b5f60208284031215611e52575f80fd5b611c2e82611d17565b5f5b83811015611e75578181015183820152602001611e5d565b50505f910152565b5f8151808452611e94816020860160208601611e5b565b601f01601f19169290920160200192915050565b634e487b7160e01b5f52602160045260245ffd5b60048110611ed857634e487b7160e01b5f52602160045260245ffd5b9052565b6001600160a01b038a811682528916602082015260408101889052606081018790526080810186905261012060a082018190525f90611f1d83820188611e7d565b915050611f2d60c0830186611ebc565b8360e08301528215156101008301529a9950505050505050505050565b6020810161145f8284611ebc565b5f8251611f69818460208701611e5b565b9190910192915050565b60208082526011908201527014185e5b595b9d081b9bdd08199bdd5b99607a1b604082015260600190565b5f60208284031215611fae575f80fd5b81518015158114611c2e575f80fd5b604081525f611fcf6040830185611e7d565b90508260208301529392505050565b602080825260139082015272151bdad95b881c995d1d5c9b8819985a5b1959606a1b604082015260600190565b602081525f611c2e6020830184611e7d565b634e487b7160e01b5f52601160045260245ffd5b8181038181111561145f5761145f61201d565b60208082526011908201527011d85cc81c99599d5b990819985a5b1959607a1b604082015260600190565b5f8351612080818460208801611e5b565b60f89390931b6001600160f81b0319169190920190815260010192915050565b600181811c908216806120b457607f821691505b6020821081036120d257634e487b7160e01b5f52602260045260245ffd5b50919050565b5f8083546120e5816120a0565b600182811680156120fd57600181146121125761213e565b60ff198416875282151583028701945061213e565b875f526020805f205f5b858110156121355781548a82015290840190820161211c565b50505082870194505b50929695505050505050565b5f6020828403121561215a575f80fd5b5051919050565b601f8211156107ac575f81815260208120601f850160051c810160208610156121875750805b601f850160051c820191505b818110156121a657828155600101612193565b505050505050565b815167ffffffffffffffff8111156121c8576121c8611b89565b6121dc816121d684546120a0565b84612161565b602080601f83116001811461220f575f84156121f85750858301515b5f19600386901b1c1916600185901b1785556121a6565b5f85815260208120601f198616915b8281101561223d5788860151825594840194600190910190840161221e565b508582101561225a57878501515f19600388901b60f8161c191681555b5050505050600190811b01905550565b608081525f61227c6080830187611e7d565b856020840152846040840152828103606084015261229a8185611e7d565b979650505050505050565b5f826122bf57634e487b7160e01b5f52601260045260245ffd5b500490565b808202811582820484141761145f5761145f61201d56fea2646970667358221220f83782cb1903c2039f3b38a4375f61197787cb3fbe9e1b835f0b0eb6c3df237264736f6c63430008140033",
}

// FiatToTokenPaymentGatewayABI is the input ABI used to generate the binding from.
// Deprecated: Use FiatToTokenPaymentGatewayMetaData.ABI instead.
var FiatToTokenPaymentGatewayABI = FiatToTokenPaymentGatewayMetaData.ABI

// FiatToTokenPaymentGatewayBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use FiatToTokenPaymentGatewayMetaData.Bin instead.
var FiatToTokenPaymentGatewayBin = FiatToTokenPaymentGatewayMetaData.Bin

// DeployFiatToTokenPaymentGateway deploys a new Ethereum contract, binding an instance of FiatToTokenPaymentGateway to it.
func DeployFiatToTokenPaymentGateway(auth *bind.TransactOpts, backend bind.ContractBackend, _token common.Address, _pricePerToken *big.Int, _requiredGasDeposit *big.Int) (common.Address, *types.Transaction, *FiatToTokenPaymentGateway, error) {
	parsed, err := FiatToTokenPaymentGatewayMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FiatToTokenPaymentGatewayBin), backend, _token, _pricePerToken, _requiredGasDeposit)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FiatToTokenPaymentGateway{FiatToTokenPaymentGatewayCaller: FiatToTokenPaymentGatewayCaller{contract: contract}, FiatToTokenPaymentGatewayTransactor: FiatToTokenPaymentGatewayTransactor{contract: contract}, FiatToTokenPaymentGatewayFilterer: FiatToTokenPaymentGatewayFilterer{contract: contract}}, nil
}

// FiatToTokenPaymentGateway is an auto generated Go binding around an Ethereum contract.
type FiatToTokenPaymentGateway struct {
	FiatToTokenPaymentGatewayCaller     // Read-only binding to the contract
	FiatToTokenPaymentGatewayTransactor // Write-only binding to the contract
	FiatToTokenPaymentGatewayFilterer   // Log filterer for contract events
}

// FiatToTokenPaymentGatewayCaller is an auto generated read-only Go binding around an Ethereum contract.
type FiatToTokenPaymentGatewayCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FiatToTokenPaymentGatewayTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FiatToTokenPaymentGatewayTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FiatToTokenPaymentGatewayFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FiatToTokenPaymentGatewayFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FiatToTokenPaymentGatewaySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FiatToTokenPaymentGatewaySession struct {
	Contract     *FiatToTokenPaymentGateway // Generic contract binding to set the session for
	CallOpts     bind.CallOpts              // Call options to use throughout this session
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// FiatToTokenPaymentGatewayCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FiatToTokenPaymentGatewayCallerSession struct {
	Contract *FiatToTokenPaymentGatewayCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                    // Call options to use throughout this session
}

// FiatToTokenPaymentGatewayTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FiatToTokenPaymentGatewayTransactorSession struct {
	Contract     *FiatToTokenPaymentGatewayTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                    // Transaction auth options to use throughout this session
}

// FiatToTokenPaymentGatewayRaw is an auto generated low-level Go binding around an Ethereum contract.
type FiatToTokenPaymentGatewayRaw struct {
	Contract *FiatToTokenPaymentGateway // Generic contract binding to access the raw methods on
}

// FiatToTokenPaymentGatewayCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FiatToTokenPaymentGatewayCallerRaw struct {
	Contract *FiatToTokenPaymentGatewayCaller // Generic read-only contract binding to access the raw methods on
}

// FiatToTokenPaymentGatewayTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FiatToTokenPaymentGatewayTransactorRaw struct {
	Contract *FiatToTokenPaymentGatewayTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFiatToTokenPaymentGateway creates a new instance of FiatToTokenPaymentGateway, bound to a specific deployed contract.
func NewFiatToTokenPaymentGateway(address common.Address, backend bind.ContractBackend) (*FiatToTokenPaymentGateway, error) {
	contract, err := bindFiatToTokenPaymentGateway(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FiatToTokenPaymentGateway{FiatToTokenPaymentGatewayCaller: FiatToTokenPaymentGatewayCaller{contract: contract}, FiatToTokenPaymentGatewayTransactor: FiatToTokenPaymentGatewayTransactor{contract: contract}, FiatToTokenPaymentGatewayFilterer: FiatToTokenPaymentGatewayFilterer{contract: contract}}, nil
}

// NewFiatToTokenPaymentGatewayCaller creates a new read-only instance of FiatToTokenPaymentGateway, bound to a specific deployed contract.
func NewFiatToTokenPaymentGatewayCaller(address common.Address, caller bind.ContractCaller) (*FiatToTokenPaymentGatewayCaller, error) {
	contract, err := bindFiatToTokenPaymentGateway(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FiatToTokenPaymentGatewayCaller{contract: contract}, nil
}

// NewFiatToTokenPaymentGatewayTransactor creates a new write-only instance of FiatToTokenPaymentGateway, bound to a specific deployed contract.
func NewFiatToTokenPaymentGatewayTransactor(address common.Address, transactor bind.ContractTransactor) (*FiatToTokenPaymentGatewayTransactor, error) {
	contract, err := bindFiatToTokenPaymentGateway(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FiatToTokenPaymentGatewayTransactor{contract: contract}, nil
}

// NewFiatToTokenPaymentGatewayFilterer creates a new log filterer instance of FiatToTokenPaymentGateway, bound to a specific deployed contract.
func NewFiatToTokenPaymentGatewayFilterer(address common.Address, filterer bind.ContractFilterer) (*FiatToTokenPaymentGatewayFilterer, error) {
	contract, err := bindFiatToTokenPaymentGateway(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FiatToTokenPaymentGatewayFilterer{contract: contract}, nil
}

// bindFiatToTokenPaymentGateway binds a generic wrapper to an already deployed contract.
func bindFiatToTokenPaymentGateway(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FiatToTokenPaymentGatewayMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FiatToTokenPaymentGateway.Contract.FiatToTokenPaymentGatewayCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.Contract.FiatToTokenPaymentGatewayTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.Contract.FiatToTokenPaymentGatewayTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FiatToTokenPaymentGateway.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.Contract.contract.Transact(opts, method, params...)
}

// CalculateProcessingCost is a free data retrieval call binding the contract method 0x8ff11488.
//
// Solidity: function calculateProcessingCost() view returns(uint256)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayCaller) CalculateProcessingCost(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FiatToTokenPaymentGateway.contract.Call(opts, &out, "calculateProcessingCost")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateProcessingCost is a free data retrieval call binding the contract method 0x8ff11488.
//
// Solidity: function calculateProcessingCost() view returns(uint256)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewaySession) CalculateProcessingCost() (*big.Int, error) {
	return _FiatToTokenPaymentGateway.Contract.CalculateProcessingCost(&_FiatToTokenPaymentGateway.CallOpts)
}

// CalculateProcessingCost is a free data retrieval call binding the contract method 0x8ff11488.
//
// Solidity: function calculateProcessingCost() view returns(uint256)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayCallerSession) CalculateProcessingCost() (*big.Int, error) {
	return _FiatToTokenPaymentGateway.Contract.CalculateProcessingCost(&_FiatToTokenPaymentGateway.CallOpts)
}

// CalculateTokenAmount is a free data retrieval call binding the contract method 0xa24bcf46.
//
// Solidity: function calculateTokenAmount(uint256 fiatAmount) view returns(uint256)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayCaller) CalculateTokenAmount(opts *bind.CallOpts, fiatAmount *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _FiatToTokenPaymentGateway.contract.Call(opts, &out, "calculateTokenAmount", fiatAmount)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateTokenAmount is a free data retrieval call binding the contract method 0xa24bcf46.
//
// Solidity: function calculateTokenAmount(uint256 fiatAmount) view returns(uint256)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewaySession) CalculateTokenAmount(fiatAmount *big.Int) (*big.Int, error) {
	return _FiatToTokenPaymentGateway.Contract.CalculateTokenAmount(&_FiatToTokenPaymentGateway.CallOpts, fiatAmount)
}

// CalculateTokenAmount is a free data retrieval call binding the contract method 0xa24bcf46.
//
// Solidity: function calculateTokenAmount(uint256 fiatAmount) view returns(uint256)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayCallerSession) CalculateTokenAmount(fiatAmount *big.Int) (*big.Int, error) {
	return _FiatToTokenPaymentGateway.Contract.CalculateTokenAmount(&_FiatToTokenPaymentGateway.CallOpts, fiatAmount)
}

// GatewaySigners is a free data retrieval call binding the contract method 0x4e981b68.
//
// Solidity: function gatewaySigners(string ) view returns(address)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayCaller) GatewaySigners(opts *bind.CallOpts, arg0 string) (common.Address, error) {
	var out []interface{}
	err := _FiatToTokenPaymentGateway.contract.Call(opts, &out, "gatewaySigners", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GatewaySigners is a free data retrieval call binding the contract method 0x4e981b68.
//
// Solidity: function gatewaySigners(string ) view returns(address)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewaySession) GatewaySigners(arg0 string) (common.Address, error) {
	return _FiatToTokenPaymentGateway.Contract.GatewaySigners(&_FiatToTokenPaymentGateway.CallOpts, arg0)
}

// GatewaySigners is a free data retrieval call binding the contract method 0x4e981b68.
//
// Solidity: function gatewaySigners(string ) view returns(address)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayCallerSession) GatewaySigners(arg0 string) (common.Address, error) {
	return _FiatToTokenPaymentGateway.Contract.GatewaySigners(&_FiatToTokenPaymentGateway.CallOpts, arg0)
}

// GetPaymentStatus is a free data retrieval call binding the contract method 0xfd2240d3.
//
// Solidity: function getPaymentStatus(string paymentId) view returns(uint8)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayCaller) GetPaymentStatus(opts *bind.CallOpts, paymentId string) (uint8, error) {
	var out []interface{}
	err := _FiatToTokenPaymentGateway.contract.Call(opts, &out, "getPaymentStatus", paymentId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetPaymentStatus is a free data retrieval call binding the contract method 0xfd2240d3.
//
// Solidity: function getPaymentStatus(string paymentId) view returns(uint8)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewaySession) GetPaymentStatus(paymentId string) (uint8, error) {
	return _FiatToTokenPaymentGateway.Contract.GetPaymentStatus(&_FiatToTokenPaymentGateway.CallOpts, paymentId)
}

// GetPaymentStatus is a free data retrieval call binding the contract method 0xfd2240d3.
//
// Solidity: function getPaymentStatus(string paymentId) view returns(uint8)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayCallerSession) GetPaymentStatus(paymentId string) (uint8, error) {
	return _FiatToTokenPaymentGateway.Contract.GetPaymentStatus(&_FiatToTokenPaymentGateway.CallOpts, paymentId)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FiatToTokenPaymentGateway.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewaySession) Owner() (common.Address, error) {
	return _FiatToTokenPaymentGateway.Contract.Owner(&_FiatToTokenPaymentGateway.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayCallerSession) Owner() (common.Address, error) {
	return _FiatToTokenPaymentGateway.Contract.Owner(&_FiatToTokenPaymentGateway.CallOpts)
}

// Payments is a free data retrieval call binding the contract method 0xf9468713.
//
// Solidity: function payments(string ) view returns(address buyer, address destinationWallet, uint256 tokenAmount, uint256 fiatAmount, uint256 timestamp, string gateway, uint8 status, uint256 gasFundAmount, bool gasRefunded)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayCaller) Payments(opts *bind.CallOpts, arg0 string) (struct {
	Buyer             common.Address
	DestinationWallet common.Address
	TokenAmount       *big.Int
	FiatAmount        *big.Int
	Timestamp         *big.Int
	Gateway           string
	Status            uint8
	GasFundAmount     *big.Int
	GasRefunded       bool
}, error) {
	var out []interface{}
	err := _FiatToTokenPaymentGateway.contract.Call(opts, &out, "payments", arg0)

	outstruct := new(struct {
		Buyer             common.Address
		DestinationWallet common.Address
		TokenAmount       *big.Int
		FiatAmount        *big.Int
		Timestamp         *big.Int
		Gateway           string
		Status            uint8
		GasFundAmount     *big.Int
		GasRefunded       bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Buyer = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.DestinationWallet = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.TokenAmount = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.FiatAmount = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Timestamp = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.Gateway = *abi.ConvertType(out[5], new(string)).(*string)
	outstruct.Status = *abi.ConvertType(out[6], new(uint8)).(*uint8)
	outstruct.GasFundAmount = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)
	outstruct.GasRefunded = *abi.ConvertType(out[8], new(bool)).(*bool)

	return *outstruct, err

}

// Payments is a free data retrieval call binding the contract method 0xf9468713.
//
// Solidity: function payments(string ) view returns(address buyer, address destinationWallet, uint256 tokenAmount, uint256 fiatAmount, uint256 timestamp, string gateway, uint8 status, uint256 gasFundAmount, bool gasRefunded)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewaySession) Payments(arg0 string) (struct {
	Buyer             common.Address
	DestinationWallet common.Address
	TokenAmount       *big.Int
	FiatAmount        *big.Int
	Timestamp         *big.Int
	Gateway           string
	Status            uint8
	GasFundAmount     *big.Int
	GasRefunded       bool
}, error) {
	return _FiatToTokenPaymentGateway.Contract.Payments(&_FiatToTokenPaymentGateway.CallOpts, arg0)
}

// Payments is a free data retrieval call binding the contract method 0xf9468713.
//
// Solidity: function payments(string ) view returns(address buyer, address destinationWallet, uint256 tokenAmount, uint256 fiatAmount, uint256 timestamp, string gateway, uint8 status, uint256 gasFundAmount, bool gasRefunded)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayCallerSession) Payments(arg0 string) (struct {
	Buyer             common.Address
	DestinationWallet common.Address
	TokenAmount       *big.Int
	FiatAmount        *big.Int
	Timestamp         *big.Int
	Gateway           string
	Status            uint8
	GasFundAmount     *big.Int
	GasRefunded       bool
}, error) {
	return _FiatToTokenPaymentGateway.Contract.Payments(&_FiatToTokenPaymentGateway.CallOpts, arg0)
}

// PricePerToken is a free data retrieval call binding the contract method 0x7b1b1de6.
//
// Solidity: function pricePerToken() view returns(uint256)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayCaller) PricePerToken(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FiatToTokenPaymentGateway.contract.Call(opts, &out, "pricePerToken")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PricePerToken is a free data retrieval call binding the contract method 0x7b1b1de6.
//
// Solidity: function pricePerToken() view returns(uint256)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewaySession) PricePerToken() (*big.Int, error) {
	return _FiatToTokenPaymentGateway.Contract.PricePerToken(&_FiatToTokenPaymentGateway.CallOpts)
}

// PricePerToken is a free data retrieval call binding the contract method 0x7b1b1de6.
//
// Solidity: function pricePerToken() view returns(uint256)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayCallerSession) PricePerToken() (*big.Int, error) {
	return _FiatToTokenPaymentGateway.Contract.PricePerToken(&_FiatToTokenPaymentGateway.CallOpts)
}

// RequiredGasDeposit is a free data retrieval call binding the contract method 0xe48770df.
//
// Solidity: function requiredGasDeposit() view returns(uint256)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayCaller) RequiredGasDeposit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FiatToTokenPaymentGateway.contract.Call(opts, &out, "requiredGasDeposit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RequiredGasDeposit is a free data retrieval call binding the contract method 0xe48770df.
//
// Solidity: function requiredGasDeposit() view returns(uint256)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewaySession) RequiredGasDeposit() (*big.Int, error) {
	return _FiatToTokenPaymentGateway.Contract.RequiredGasDeposit(&_FiatToTokenPaymentGateway.CallOpts)
}

// RequiredGasDeposit is a free data retrieval call binding the contract method 0xe48770df.
//
// Solidity: function requiredGasDeposit() view returns(uint256)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayCallerSession) RequiredGasDeposit() (*big.Int, error) {
	return _FiatToTokenPaymentGateway.Contract.RequiredGasDeposit(&_FiatToTokenPaymentGateway.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FiatToTokenPaymentGateway.contract.Call(opts, &out, "token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewaySession) Token() (common.Address, error) {
	return _FiatToTokenPaymentGateway.Contract.Token(&_FiatToTokenPaymentGateway.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayCallerSession) Token() (common.Address, error) {
	return _FiatToTokenPaymentGateway.Contract.Token(&_FiatToTokenPaymentGateway.CallOpts)
}

// CreatePayment is a paid mutator transaction binding the contract method 0x5bb8e4d1.
//
// Solidity: function createPayment(string paymentId, uint256 tokenAmount, uint256 fiatAmount, string gateway, address destinationWallet) payable returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayTransactor) CreatePayment(opts *bind.TransactOpts, paymentId string, tokenAmount *big.Int, fiatAmount *big.Int, gateway string, destinationWallet common.Address) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.contract.Transact(opts, "createPayment", paymentId, tokenAmount, fiatAmount, gateway, destinationWallet)
}

// CreatePayment is a paid mutator transaction binding the contract method 0x5bb8e4d1.
//
// Solidity: function createPayment(string paymentId, uint256 tokenAmount, uint256 fiatAmount, string gateway, address destinationWallet) payable returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewaySession) CreatePayment(paymentId string, tokenAmount *big.Int, fiatAmount *big.Int, gateway string, destinationWallet common.Address) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.Contract.CreatePayment(&_FiatToTokenPaymentGateway.TransactOpts, paymentId, tokenAmount, fiatAmount, gateway, destinationWallet)
}

// CreatePayment is a paid mutator transaction binding the contract method 0x5bb8e4d1.
//
// Solidity: function createPayment(string paymentId, uint256 tokenAmount, uint256 fiatAmount, string gateway, address destinationWallet) payable returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayTransactorSession) CreatePayment(paymentId string, tokenAmount *big.Int, fiatAmount *big.Int, gateway string, destinationWallet common.Address) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.Contract.CreatePayment(&_FiatToTokenPaymentGateway.TransactOpts, paymentId, tokenAmount, fiatAmount, gateway, destinationWallet)
}

// MockPaymentCallback is a paid mutator transaction binding the contract method 0x03ff0399.
//
// Solidity: function mockPaymentCallback(string paymentId, uint8 status) returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayTransactor) MockPaymentCallback(opts *bind.TransactOpts, paymentId string, status uint8) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.contract.Transact(opts, "mockPaymentCallback", paymentId, status)
}

// MockPaymentCallback is a paid mutator transaction binding the contract method 0x03ff0399.
//
// Solidity: function mockPaymentCallback(string paymentId, uint8 status) returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewaySession) MockPaymentCallback(paymentId string, status uint8) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.Contract.MockPaymentCallback(&_FiatToTokenPaymentGateway.TransactOpts, paymentId, status)
}

// MockPaymentCallback is a paid mutator transaction binding the contract method 0x03ff0399.
//
// Solidity: function mockPaymentCallback(string paymentId, uint8 status) returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayTransactorSession) MockPaymentCallback(paymentId string, status uint8) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.Contract.MockPaymentCallback(&_FiatToTokenPaymentGateway.TransactOpts, paymentId, status)
}

// ProcessPaymentCallback is a paid mutator transaction binding the contract method 0x0ac6e131.
//
// Solidity: function processPaymentCallback(string paymentId, uint8 status, bytes signature) returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayTransactor) ProcessPaymentCallback(opts *bind.TransactOpts, paymentId string, status uint8, signature []byte) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.contract.Transact(opts, "processPaymentCallback", paymentId, status, signature)
}

// ProcessPaymentCallback is a paid mutator transaction binding the contract method 0x0ac6e131.
//
// Solidity: function processPaymentCallback(string paymentId, uint8 status, bytes signature) returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewaySession) ProcessPaymentCallback(paymentId string, status uint8, signature []byte) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.Contract.ProcessPaymentCallback(&_FiatToTokenPaymentGateway.TransactOpts, paymentId, status, signature)
}

// ProcessPaymentCallback is a paid mutator transaction binding the contract method 0x0ac6e131.
//
// Solidity: function processPaymentCallback(string paymentId, uint8 status, bytes signature) returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayTransactorSession) ProcessPaymentCallback(paymentId string, status uint8, signature []byte) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.Contract.ProcessPaymentCallback(&_FiatToTokenPaymentGateway.TransactOpts, paymentId, status, signature)
}

// ProcessRefund is a paid mutator transaction binding the contract method 0xee3efd58.
//
// Solidity: function processRefund(string paymentId) returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayTransactor) ProcessRefund(opts *bind.TransactOpts, paymentId string) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.contract.Transact(opts, "processRefund", paymentId)
}

// ProcessRefund is a paid mutator transaction binding the contract method 0xee3efd58.
//
// Solidity: function processRefund(string paymentId) returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewaySession) ProcessRefund(paymentId string) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.Contract.ProcessRefund(&_FiatToTokenPaymentGateway.TransactOpts, paymentId)
}

// ProcessRefund is a paid mutator transaction binding the contract method 0xee3efd58.
//
// Solidity: function processRefund(string paymentId) returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayTransactorSession) ProcessRefund(paymentId string) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.Contract.ProcessRefund(&_FiatToTokenPaymentGateway.TransactOpts, paymentId)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewaySession) RenounceOwnership() (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.Contract.RenounceOwnership(&_FiatToTokenPaymentGateway.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.Contract.RenounceOwnership(&_FiatToTokenPaymentGateway.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewaySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.Contract.TransferOwnership(&_FiatToTokenPaymentGateway.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.Contract.TransferOwnership(&_FiatToTokenPaymentGateway.TransactOpts, newOwner)
}

// UpdateGasDepositRequirement is a paid mutator transaction binding the contract method 0x52c3cc2e.
//
// Solidity: function updateGasDepositRequirement(uint256 _requiredGasDeposit) returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayTransactor) UpdateGasDepositRequirement(opts *bind.TransactOpts, _requiredGasDeposit *big.Int) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.contract.Transact(opts, "updateGasDepositRequirement", _requiredGasDeposit)
}

// UpdateGasDepositRequirement is a paid mutator transaction binding the contract method 0x52c3cc2e.
//
// Solidity: function updateGasDepositRequirement(uint256 _requiredGasDeposit) returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewaySession) UpdateGasDepositRequirement(_requiredGasDeposit *big.Int) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.Contract.UpdateGasDepositRequirement(&_FiatToTokenPaymentGateway.TransactOpts, _requiredGasDeposit)
}

// UpdateGasDepositRequirement is a paid mutator transaction binding the contract method 0x52c3cc2e.
//
// Solidity: function updateGasDepositRequirement(uint256 _requiredGasDeposit) returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayTransactorSession) UpdateGasDepositRequirement(_requiredGasDeposit *big.Int) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.Contract.UpdateGasDepositRequirement(&_FiatToTokenPaymentGateway.TransactOpts, _requiredGasDeposit)
}

// UpdateGatewaySigner is a paid mutator transaction binding the contract method 0x16692d47.
//
// Solidity: function updateGatewaySigner(string gateway, address signer) returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayTransactor) UpdateGatewaySigner(opts *bind.TransactOpts, gateway string, signer common.Address) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.contract.Transact(opts, "updateGatewaySigner", gateway, signer)
}

// UpdateGatewaySigner is a paid mutator transaction binding the contract method 0x16692d47.
//
// Solidity: function updateGatewaySigner(string gateway, address signer) returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewaySession) UpdateGatewaySigner(gateway string, signer common.Address) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.Contract.UpdateGatewaySigner(&_FiatToTokenPaymentGateway.TransactOpts, gateway, signer)
}

// UpdateGatewaySigner is a paid mutator transaction binding the contract method 0x16692d47.
//
// Solidity: function updateGatewaySigner(string gateway, address signer) returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayTransactorSession) UpdateGatewaySigner(gateway string, signer common.Address) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.Contract.UpdateGatewaySigner(&_FiatToTokenPaymentGateway.TransactOpts, gateway, signer)
}

// UpdateTokenPrice is a paid mutator transaction binding the contract method 0x676c0d77.
//
// Solidity: function updateTokenPrice(uint256 _pricePerToken) returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayTransactor) UpdateTokenPrice(opts *bind.TransactOpts, _pricePerToken *big.Int) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.contract.Transact(opts, "updateTokenPrice", _pricePerToken)
}

// UpdateTokenPrice is a paid mutator transaction binding the contract method 0x676c0d77.
//
// Solidity: function updateTokenPrice(uint256 _pricePerToken) returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewaySession) UpdateTokenPrice(_pricePerToken *big.Int) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.Contract.UpdateTokenPrice(&_FiatToTokenPaymentGateway.TransactOpts, _pricePerToken)
}

// UpdateTokenPrice is a paid mutator transaction binding the contract method 0x676c0d77.
//
// Solidity: function updateTokenPrice(uint256 _pricePerToken) returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayTransactorSession) UpdateTokenPrice(_pricePerToken *big.Int) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.Contract.UpdateTokenPrice(&_FiatToTokenPaymentGateway.TransactOpts, _pricePerToken)
}

// WithdrawProcessingFees is a paid mutator transaction binding the contract method 0x34eec86a.
//
// Solidity: function withdrawProcessingFees() returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayTransactor) WithdrawProcessingFees(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.contract.Transact(opts, "withdrawProcessingFees")
}

// WithdrawProcessingFees is a paid mutator transaction binding the contract method 0x34eec86a.
//
// Solidity: function withdrawProcessingFees() returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewaySession) WithdrawProcessingFees() (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.Contract.WithdrawProcessingFees(&_FiatToTokenPaymentGateway.TransactOpts)
}

// WithdrawProcessingFees is a paid mutator transaction binding the contract method 0x34eec86a.
//
// Solidity: function withdrawProcessingFees() returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayTransactorSession) WithdrawProcessingFees() (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.Contract.WithdrawProcessingFees(&_FiatToTokenPaymentGateway.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewaySession) Receive() (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.Contract.Receive(&_FiatToTokenPaymentGateway.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayTransactorSession) Receive() (*types.Transaction, error) {
	return _FiatToTokenPaymentGateway.Contract.Receive(&_FiatToTokenPaymentGateway.TransactOpts)
}

// FiatToTokenPaymentGatewayGasDepositRequirementUpdatedIterator is returned from FilterGasDepositRequirementUpdated and is used to iterate over the raw logs and unpacked data for GasDepositRequirementUpdated events raised by the FiatToTokenPaymentGateway contract.
type FiatToTokenPaymentGatewayGasDepositRequirementUpdatedIterator struct {
	Event *FiatToTokenPaymentGatewayGasDepositRequirementUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FiatToTokenPaymentGatewayGasDepositRequirementUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FiatToTokenPaymentGatewayGasDepositRequirementUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FiatToTokenPaymentGatewayGasDepositRequirementUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FiatToTokenPaymentGatewayGasDepositRequirementUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FiatToTokenPaymentGatewayGasDepositRequirementUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FiatToTokenPaymentGatewayGasDepositRequirementUpdated represents a GasDepositRequirementUpdated event raised by the FiatToTokenPaymentGateway contract.
type FiatToTokenPaymentGatewayGasDepositRequirementUpdated struct {
	NewRequirement *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterGasDepositRequirementUpdated is a free log retrieval operation binding the contract event 0x6ba7f9a0367ca0f8ad8598911d0fa547a28ea788a8b3c0607092d5455c2b583e.
//
// Solidity: event GasDepositRequirementUpdated(uint256 newRequirement)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayFilterer) FilterGasDepositRequirementUpdated(opts *bind.FilterOpts) (*FiatToTokenPaymentGatewayGasDepositRequirementUpdatedIterator, error) {

	logs, sub, err := _FiatToTokenPaymentGateway.contract.FilterLogs(opts, "GasDepositRequirementUpdated")
	if err != nil {
		return nil, err
	}
	return &FiatToTokenPaymentGatewayGasDepositRequirementUpdatedIterator{contract: _FiatToTokenPaymentGateway.contract, event: "GasDepositRequirementUpdated", logs: logs, sub: sub}, nil
}

// WatchGasDepositRequirementUpdated is a free log subscription operation binding the contract event 0x6ba7f9a0367ca0f8ad8598911d0fa547a28ea788a8b3c0607092d5455c2b583e.
//
// Solidity: event GasDepositRequirementUpdated(uint256 newRequirement)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayFilterer) WatchGasDepositRequirementUpdated(opts *bind.WatchOpts, sink chan<- *FiatToTokenPaymentGatewayGasDepositRequirementUpdated) (event.Subscription, error) {

	logs, sub, err := _FiatToTokenPaymentGateway.contract.WatchLogs(opts, "GasDepositRequirementUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FiatToTokenPaymentGatewayGasDepositRequirementUpdated)
				if err := _FiatToTokenPaymentGateway.contract.UnpackLog(event, "GasDepositRequirementUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseGasDepositRequirementUpdated is a log parse operation binding the contract event 0x6ba7f9a0367ca0f8ad8598911d0fa547a28ea788a8b3c0607092d5455c2b583e.
//
// Solidity: event GasDepositRequirementUpdated(uint256 newRequirement)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayFilterer) ParseGasDepositRequirementUpdated(log types.Log) (*FiatToTokenPaymentGatewayGasDepositRequirementUpdated, error) {
	event := new(FiatToTokenPaymentGatewayGasDepositRequirementUpdated)
	if err := _FiatToTokenPaymentGateway.contract.UnpackLog(event, "GasDepositRequirementUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FiatToTokenPaymentGatewayGasRefundedIterator is returned from FilterGasRefunded and is used to iterate over the raw logs and unpacked data for GasRefunded events raised by the FiatToTokenPaymentGateway contract.
type FiatToTokenPaymentGatewayGasRefundedIterator struct {
	Event *FiatToTokenPaymentGatewayGasRefunded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FiatToTokenPaymentGatewayGasRefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FiatToTokenPaymentGatewayGasRefunded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FiatToTokenPaymentGatewayGasRefunded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FiatToTokenPaymentGatewayGasRefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FiatToTokenPaymentGatewayGasRefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FiatToTokenPaymentGatewayGasRefunded represents a GasRefunded event raised by the FiatToTokenPaymentGateway contract.
type FiatToTokenPaymentGatewayGasRefunded struct {
	PaymentId string
	Buyer     common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterGasRefunded is a free log retrieval operation binding the contract event 0x89092cc5430d6dbac7d084b6205265037b6589502bc0cedd66b2f46eefb79e46.
//
// Solidity: event GasRefunded(string paymentId, address indexed buyer, uint256 amount)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayFilterer) FilterGasRefunded(opts *bind.FilterOpts, buyer []common.Address) (*FiatToTokenPaymentGatewayGasRefundedIterator, error) {

	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}

	logs, sub, err := _FiatToTokenPaymentGateway.contract.FilterLogs(opts, "GasRefunded", buyerRule)
	if err != nil {
		return nil, err
	}
	return &FiatToTokenPaymentGatewayGasRefundedIterator{contract: _FiatToTokenPaymentGateway.contract, event: "GasRefunded", logs: logs, sub: sub}, nil
}

// WatchGasRefunded is a free log subscription operation binding the contract event 0x89092cc5430d6dbac7d084b6205265037b6589502bc0cedd66b2f46eefb79e46.
//
// Solidity: event GasRefunded(string paymentId, address indexed buyer, uint256 amount)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayFilterer) WatchGasRefunded(opts *bind.WatchOpts, sink chan<- *FiatToTokenPaymentGatewayGasRefunded, buyer []common.Address) (event.Subscription, error) {

	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}

	logs, sub, err := _FiatToTokenPaymentGateway.contract.WatchLogs(opts, "GasRefunded", buyerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FiatToTokenPaymentGatewayGasRefunded)
				if err := _FiatToTokenPaymentGateway.contract.UnpackLog(event, "GasRefunded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseGasRefunded is a log parse operation binding the contract event 0x89092cc5430d6dbac7d084b6205265037b6589502bc0cedd66b2f46eefb79e46.
//
// Solidity: event GasRefunded(string paymentId, address indexed buyer, uint256 amount)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayFilterer) ParseGasRefunded(log types.Log) (*FiatToTokenPaymentGatewayGasRefunded, error) {
	event := new(FiatToTokenPaymentGatewayGasRefunded)
	if err := _FiatToTokenPaymentGateway.contract.UnpackLog(event, "GasRefunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FiatToTokenPaymentGatewayOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the FiatToTokenPaymentGateway contract.
type FiatToTokenPaymentGatewayOwnershipTransferredIterator struct {
	Event *FiatToTokenPaymentGatewayOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FiatToTokenPaymentGatewayOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FiatToTokenPaymentGatewayOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FiatToTokenPaymentGatewayOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FiatToTokenPaymentGatewayOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FiatToTokenPaymentGatewayOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FiatToTokenPaymentGatewayOwnershipTransferred represents a OwnershipTransferred event raised by the FiatToTokenPaymentGateway contract.
type FiatToTokenPaymentGatewayOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*FiatToTokenPaymentGatewayOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _FiatToTokenPaymentGateway.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &FiatToTokenPaymentGatewayOwnershipTransferredIterator{contract: _FiatToTokenPaymentGateway.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FiatToTokenPaymentGatewayOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _FiatToTokenPaymentGateway.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FiatToTokenPaymentGatewayOwnershipTransferred)
				if err := _FiatToTokenPaymentGateway.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayFilterer) ParseOwnershipTransferred(log types.Log) (*FiatToTokenPaymentGatewayOwnershipTransferred, error) {
	event := new(FiatToTokenPaymentGatewayOwnershipTransferred)
	if err := _FiatToTokenPaymentGateway.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FiatToTokenPaymentGatewayPaymentCompletedIterator is returned from FilterPaymentCompleted and is used to iterate over the raw logs and unpacked data for PaymentCompleted events raised by the FiatToTokenPaymentGateway contract.
type FiatToTokenPaymentGatewayPaymentCompletedIterator struct {
	Event *FiatToTokenPaymentGatewayPaymentCompleted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FiatToTokenPaymentGatewayPaymentCompletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FiatToTokenPaymentGatewayPaymentCompleted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FiatToTokenPaymentGatewayPaymentCompleted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FiatToTokenPaymentGatewayPaymentCompletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FiatToTokenPaymentGatewayPaymentCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FiatToTokenPaymentGatewayPaymentCompleted represents a PaymentCompleted event raised by the FiatToTokenPaymentGateway contract.
type FiatToTokenPaymentGatewayPaymentCompleted struct {
	PaymentId   string
	Buyer       common.Address
	TokenAmount *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPaymentCompleted is a free log retrieval operation binding the contract event 0x13d1118d2532a9b21c3e26bd2ba1b959263d90ea306c910ff6339e405192f1f2.
//
// Solidity: event PaymentCompleted(string paymentId, address indexed buyer, uint256 tokenAmount)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayFilterer) FilterPaymentCompleted(opts *bind.FilterOpts, buyer []common.Address) (*FiatToTokenPaymentGatewayPaymentCompletedIterator, error) {

	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}

	logs, sub, err := _FiatToTokenPaymentGateway.contract.FilterLogs(opts, "PaymentCompleted", buyerRule)
	if err != nil {
		return nil, err
	}
	return &FiatToTokenPaymentGatewayPaymentCompletedIterator{contract: _FiatToTokenPaymentGateway.contract, event: "PaymentCompleted", logs: logs, sub: sub}, nil
}

// WatchPaymentCompleted is a free log subscription operation binding the contract event 0x13d1118d2532a9b21c3e26bd2ba1b959263d90ea306c910ff6339e405192f1f2.
//
// Solidity: event PaymentCompleted(string paymentId, address indexed buyer, uint256 tokenAmount)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayFilterer) WatchPaymentCompleted(opts *bind.WatchOpts, sink chan<- *FiatToTokenPaymentGatewayPaymentCompleted, buyer []common.Address) (event.Subscription, error) {

	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}

	logs, sub, err := _FiatToTokenPaymentGateway.contract.WatchLogs(opts, "PaymentCompleted", buyerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FiatToTokenPaymentGatewayPaymentCompleted)
				if err := _FiatToTokenPaymentGateway.contract.UnpackLog(event, "PaymentCompleted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaymentCompleted is a log parse operation binding the contract event 0x13d1118d2532a9b21c3e26bd2ba1b959263d90ea306c910ff6339e405192f1f2.
//
// Solidity: event PaymentCompleted(string paymentId, address indexed buyer, uint256 tokenAmount)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayFilterer) ParsePaymentCompleted(log types.Log) (*FiatToTokenPaymentGatewayPaymentCompleted, error) {
	event := new(FiatToTokenPaymentGatewayPaymentCompleted)
	if err := _FiatToTokenPaymentGateway.contract.UnpackLog(event, "PaymentCompleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FiatToTokenPaymentGatewayPaymentCreatedIterator is returned from FilterPaymentCreated and is used to iterate over the raw logs and unpacked data for PaymentCreated events raised by the FiatToTokenPaymentGateway contract.
type FiatToTokenPaymentGatewayPaymentCreatedIterator struct {
	Event *FiatToTokenPaymentGatewayPaymentCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FiatToTokenPaymentGatewayPaymentCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FiatToTokenPaymentGatewayPaymentCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FiatToTokenPaymentGatewayPaymentCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FiatToTokenPaymentGatewayPaymentCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FiatToTokenPaymentGatewayPaymentCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FiatToTokenPaymentGatewayPaymentCreated represents a PaymentCreated event raised by the FiatToTokenPaymentGateway contract.
type FiatToTokenPaymentGatewayPaymentCreated struct {
	PaymentId         string
	Buyer             common.Address
	DestinationWallet common.Address
	TokenAmount       *big.Int
	FiatAmount        *big.Int
	Gateway           string
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterPaymentCreated is a free log retrieval operation binding the contract event 0x73e020d64bf5e24a99ceab0ef0f92283d5e48bd8700e072769bf844d6f0cf34e.
//
// Solidity: event PaymentCreated(string paymentId, address indexed buyer, address indexed destinationWallet, uint256 tokenAmount, uint256 fiatAmount, string gateway)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayFilterer) FilterPaymentCreated(opts *bind.FilterOpts, buyer []common.Address, destinationWallet []common.Address) (*FiatToTokenPaymentGatewayPaymentCreatedIterator, error) {

	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}
	var destinationWalletRule []interface{}
	for _, destinationWalletItem := range destinationWallet {
		destinationWalletRule = append(destinationWalletRule, destinationWalletItem)
	}

	logs, sub, err := _FiatToTokenPaymentGateway.contract.FilterLogs(opts, "PaymentCreated", buyerRule, destinationWalletRule)
	if err != nil {
		return nil, err
	}
	return &FiatToTokenPaymentGatewayPaymentCreatedIterator{contract: _FiatToTokenPaymentGateway.contract, event: "PaymentCreated", logs: logs, sub: sub}, nil
}

// WatchPaymentCreated is a free log subscription operation binding the contract event 0x73e020d64bf5e24a99ceab0ef0f92283d5e48bd8700e072769bf844d6f0cf34e.
//
// Solidity: event PaymentCreated(string paymentId, address indexed buyer, address indexed destinationWallet, uint256 tokenAmount, uint256 fiatAmount, string gateway)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayFilterer) WatchPaymentCreated(opts *bind.WatchOpts, sink chan<- *FiatToTokenPaymentGatewayPaymentCreated, buyer []common.Address, destinationWallet []common.Address) (event.Subscription, error) {

	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}
	var destinationWalletRule []interface{}
	for _, destinationWalletItem := range destinationWallet {
		destinationWalletRule = append(destinationWalletRule, destinationWalletItem)
	}

	logs, sub, err := _FiatToTokenPaymentGateway.contract.WatchLogs(opts, "PaymentCreated", buyerRule, destinationWalletRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FiatToTokenPaymentGatewayPaymentCreated)
				if err := _FiatToTokenPaymentGateway.contract.UnpackLog(event, "PaymentCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaymentCreated is a log parse operation binding the contract event 0x73e020d64bf5e24a99ceab0ef0f92283d5e48bd8700e072769bf844d6f0cf34e.
//
// Solidity: event PaymentCreated(string paymentId, address indexed buyer, address indexed destinationWallet, uint256 tokenAmount, uint256 fiatAmount, string gateway)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayFilterer) ParsePaymentCreated(log types.Log) (*FiatToTokenPaymentGatewayPaymentCreated, error) {
	event := new(FiatToTokenPaymentGatewayPaymentCreated)
	if err := _FiatToTokenPaymentGateway.contract.UnpackLog(event, "PaymentCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FiatToTokenPaymentGatewayPaymentFailedIterator is returned from FilterPaymentFailed and is used to iterate over the raw logs and unpacked data for PaymentFailed events raised by the FiatToTokenPaymentGateway contract.
type FiatToTokenPaymentGatewayPaymentFailedIterator struct {
	Event *FiatToTokenPaymentGatewayPaymentFailed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FiatToTokenPaymentGatewayPaymentFailedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FiatToTokenPaymentGatewayPaymentFailed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FiatToTokenPaymentGatewayPaymentFailed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FiatToTokenPaymentGatewayPaymentFailedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FiatToTokenPaymentGatewayPaymentFailedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FiatToTokenPaymentGatewayPaymentFailed represents a PaymentFailed event raised by the FiatToTokenPaymentGateway contract.
type FiatToTokenPaymentGatewayPaymentFailed struct {
	PaymentId string
	Buyer     common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterPaymentFailed is a free log retrieval operation binding the contract event 0xd467a4459f62e0868cf38768c34059712f66ac7d9f3897703192a4a7e8a9db2b.
//
// Solidity: event PaymentFailed(string paymentId, address indexed buyer)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayFilterer) FilterPaymentFailed(opts *bind.FilterOpts, buyer []common.Address) (*FiatToTokenPaymentGatewayPaymentFailedIterator, error) {

	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}

	logs, sub, err := _FiatToTokenPaymentGateway.contract.FilterLogs(opts, "PaymentFailed", buyerRule)
	if err != nil {
		return nil, err
	}
	return &FiatToTokenPaymentGatewayPaymentFailedIterator{contract: _FiatToTokenPaymentGateway.contract, event: "PaymentFailed", logs: logs, sub: sub}, nil
}

// WatchPaymentFailed is a free log subscription operation binding the contract event 0xd467a4459f62e0868cf38768c34059712f66ac7d9f3897703192a4a7e8a9db2b.
//
// Solidity: event PaymentFailed(string paymentId, address indexed buyer)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayFilterer) WatchPaymentFailed(opts *bind.WatchOpts, sink chan<- *FiatToTokenPaymentGatewayPaymentFailed, buyer []common.Address) (event.Subscription, error) {

	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}

	logs, sub, err := _FiatToTokenPaymentGateway.contract.WatchLogs(opts, "PaymentFailed", buyerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FiatToTokenPaymentGatewayPaymentFailed)
				if err := _FiatToTokenPaymentGateway.contract.UnpackLog(event, "PaymentFailed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaymentFailed is a log parse operation binding the contract event 0xd467a4459f62e0868cf38768c34059712f66ac7d9f3897703192a4a7e8a9db2b.
//
// Solidity: event PaymentFailed(string paymentId, address indexed buyer)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayFilterer) ParsePaymentFailed(log types.Log) (*FiatToTokenPaymentGatewayPaymentFailed, error) {
	event := new(FiatToTokenPaymentGatewayPaymentFailed)
	if err := _FiatToTokenPaymentGateway.contract.UnpackLog(event, "PaymentFailed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FiatToTokenPaymentGatewayPaymentRefundedIterator is returned from FilterPaymentRefunded and is used to iterate over the raw logs and unpacked data for PaymentRefunded events raised by the FiatToTokenPaymentGateway contract.
type FiatToTokenPaymentGatewayPaymentRefundedIterator struct {
	Event *FiatToTokenPaymentGatewayPaymentRefunded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FiatToTokenPaymentGatewayPaymentRefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FiatToTokenPaymentGatewayPaymentRefunded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FiatToTokenPaymentGatewayPaymentRefunded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FiatToTokenPaymentGatewayPaymentRefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FiatToTokenPaymentGatewayPaymentRefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FiatToTokenPaymentGatewayPaymentRefunded represents a PaymentRefunded event raised by the FiatToTokenPaymentGateway contract.
type FiatToTokenPaymentGatewayPaymentRefunded struct {
	PaymentId string
	Buyer     common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterPaymentRefunded is a free log retrieval operation binding the contract event 0xb7309985c95c87cf571a649294c1e24320da430df8238d2ec65c732619f8e53b.
//
// Solidity: event PaymentRefunded(string paymentId, address indexed buyer)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayFilterer) FilterPaymentRefunded(opts *bind.FilterOpts, buyer []common.Address) (*FiatToTokenPaymentGatewayPaymentRefundedIterator, error) {

	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}

	logs, sub, err := _FiatToTokenPaymentGateway.contract.FilterLogs(opts, "PaymentRefunded", buyerRule)
	if err != nil {
		return nil, err
	}
	return &FiatToTokenPaymentGatewayPaymentRefundedIterator{contract: _FiatToTokenPaymentGateway.contract, event: "PaymentRefunded", logs: logs, sub: sub}, nil
}

// WatchPaymentRefunded is a free log subscription operation binding the contract event 0xb7309985c95c87cf571a649294c1e24320da430df8238d2ec65c732619f8e53b.
//
// Solidity: event PaymentRefunded(string paymentId, address indexed buyer)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayFilterer) WatchPaymentRefunded(opts *bind.WatchOpts, sink chan<- *FiatToTokenPaymentGatewayPaymentRefunded, buyer []common.Address) (event.Subscription, error) {

	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}

	logs, sub, err := _FiatToTokenPaymentGateway.contract.WatchLogs(opts, "PaymentRefunded", buyerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FiatToTokenPaymentGatewayPaymentRefunded)
				if err := _FiatToTokenPaymentGateway.contract.UnpackLog(event, "PaymentRefunded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaymentRefunded is a log parse operation binding the contract event 0xb7309985c95c87cf571a649294c1e24320da430df8238d2ec65c732619f8e53b.
//
// Solidity: event PaymentRefunded(string paymentId, address indexed buyer)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayFilterer) ParsePaymentRefunded(log types.Log) (*FiatToTokenPaymentGatewayPaymentRefunded, error) {
	event := new(FiatToTokenPaymentGatewayPaymentRefunded)
	if err := _FiatToTokenPaymentGateway.contract.UnpackLog(event, "PaymentRefunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FiatToTokenPaymentGatewayTokenPriceUpdatedIterator is returned from FilterTokenPriceUpdated and is used to iterate over the raw logs and unpacked data for TokenPriceUpdated events raised by the FiatToTokenPaymentGateway contract.
type FiatToTokenPaymentGatewayTokenPriceUpdatedIterator struct {
	Event *FiatToTokenPaymentGatewayTokenPriceUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FiatToTokenPaymentGatewayTokenPriceUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FiatToTokenPaymentGatewayTokenPriceUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FiatToTokenPaymentGatewayTokenPriceUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FiatToTokenPaymentGatewayTokenPriceUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FiatToTokenPaymentGatewayTokenPriceUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FiatToTokenPaymentGatewayTokenPriceUpdated represents a TokenPriceUpdated event raised by the FiatToTokenPaymentGateway contract.
type FiatToTokenPaymentGatewayTokenPriceUpdated struct {
	NewPrice *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTokenPriceUpdated is a free log retrieval operation binding the contract event 0x6b361c807733cc94a3b7bda39c1c94dcc0511e91aa47b0f7f918068e93976890.
//
// Solidity: event TokenPriceUpdated(uint256 newPrice)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayFilterer) FilterTokenPriceUpdated(opts *bind.FilterOpts) (*FiatToTokenPaymentGatewayTokenPriceUpdatedIterator, error) {

	logs, sub, err := _FiatToTokenPaymentGateway.contract.FilterLogs(opts, "TokenPriceUpdated")
	if err != nil {
		return nil, err
	}
	return &FiatToTokenPaymentGatewayTokenPriceUpdatedIterator{contract: _FiatToTokenPaymentGateway.contract, event: "TokenPriceUpdated", logs: logs, sub: sub}, nil
}

// WatchTokenPriceUpdated is a free log subscription operation binding the contract event 0x6b361c807733cc94a3b7bda39c1c94dcc0511e91aa47b0f7f918068e93976890.
//
// Solidity: event TokenPriceUpdated(uint256 newPrice)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayFilterer) WatchTokenPriceUpdated(opts *bind.WatchOpts, sink chan<- *FiatToTokenPaymentGatewayTokenPriceUpdated) (event.Subscription, error) {

	logs, sub, err := _FiatToTokenPaymentGateway.contract.WatchLogs(opts, "TokenPriceUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FiatToTokenPaymentGatewayTokenPriceUpdated)
				if err := _FiatToTokenPaymentGateway.contract.UnpackLog(event, "TokenPriceUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTokenPriceUpdated is a log parse operation binding the contract event 0x6b361c807733cc94a3b7bda39c1c94dcc0511e91aa47b0f7f918068e93976890.
//
// Solidity: event TokenPriceUpdated(uint256 newPrice)
func (_FiatToTokenPaymentGateway *FiatToTokenPaymentGatewayFilterer) ParseTokenPriceUpdated(log types.Log) (*FiatToTokenPaymentGatewayTokenPriceUpdated, error) {
	event := new(FiatToTokenPaymentGatewayTokenPriceUpdated)
	if err := _FiatToTokenPaymentGateway.contract.UnpackLog(event, "TokenPriceUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
