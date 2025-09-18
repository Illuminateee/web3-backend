// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol";


/**
 * @title FiatToTokenPaymentGateway
 * @dev Contract for processing fiat payments and swapping to tokens, with user-paid gas fees
 */
contract FiatToTokenPaymentGateway is Ownable, ReentrancyGuard {

    

    // Token contract
    IERC20 public token;
    
    // Payment statuses
    enum PaymentStatus { Pending, Completed, Failed, Refunded }
    
    // Payment information
    struct Payment {
        address buyer;
        address destinationWallet;
        uint256 tokenAmount;
        uint256 fiatAmount;
        uint256 timestamp;
        string gateway;
        PaymentStatus status;
        uint256 gasFundAmount;     // ETH deposited for gas
        bool gasRefunded;          // Whether remaining gas has been refunded
    }
    
    // Gas fee requirements
    uint256 public requiredGasDeposit;  // Amount of ETH required for gas fees
    
    // Mapping of payment IDs to Payment structs
    mapping(string => Payment) public payments;
    
    // Payment gateway signers (for verification)
    mapping(string => address) public gatewaySigners;
    
    // Token price in wei per smallest token unit (can be updated by owner)
    uint256 public pricePerToken;
    
    // Events
    event PaymentCreated(string paymentId, address indexed buyer, address indexed destinationWallet ,uint256 tokenAmount, uint256 fiatAmount, string gateway);
    event PaymentCompleted(string paymentId, address indexed buyer, uint256 tokenAmount);
    event PaymentFailed(string paymentId, address indexed buyer);
    event PaymentRefunded(string paymentId, address indexed buyer);
    event GasRefunded(string paymentId, address indexed buyer, uint256 amount);
    event TokenPriceUpdated(uint256 newPrice);
    event GasDepositRequirementUpdated(uint256 newRequirement);
    
    /**
     * @dev Constructor 
     * @param _token Address of the ERC20 token being sold
     * @param _pricePerToken Initial price per token in wei
     * @param _requiredGasDeposit Initial required gas deposit in wei
     */
    constructor(
        address _token, 
        uint256 _pricePerToken,
        uint256 _requiredGasDeposit
    ) Ownable(msg.sender) {
        require(_token != address(0), "Invalid token address");
        token = IERC20(_token);
        pricePerToken = _pricePerToken;
        requiredGasDeposit = _requiredGasDeposit;
        
        // Initialize gateway signers - replace with actual verification addresses in production
        gatewaySigners["midtrans"] = msg.sender;  // For testing only
        gatewaySigners["bca"] = msg.sender;       // For testing only
        gatewaySigners["stripe"] = msg.sender;    // For testing only
    }
    
    /**
     * @dev Create a payment request with gas funding and reserve tokens
     * @param paymentId Unique payment identifier from the payment gateway
     * @param tokenAmount Amount of tokens the user is purchasing
     * @param fiatAmount Amount in fiat currency (in smallest unit, e.g., cents)
     * @param gateway Payment gateway to be used (midtrans, bca, stripe)
     */
    function createPayment(
        string memory paymentId,
        uint256 tokenAmount,
        uint256 fiatAmount,
        string memory gateway,
        address destinationWallet
    ) external payable nonReentrant {
        require(payments[paymentId].buyer == address(0), "Payment ID already exists");
        require(tokenAmount > 0, "Token amount must be greater than zero");
        require(fiatAmount > 0, "Fiat amount must be greater than zero");
        require(bytes(gateway).length > 0, "Gateway must be specified");
        require(msg.value >= requiredGasDeposit, "Insufficient gas deposit");
        
        // Check if this contract has enough token allowance from the owner
        require(token.allowance(owner(), address(this)) >= tokenAmount, "Insufficient token allowance");
        
        // Create payment record with Pending status
        payments[paymentId] = Payment({
            buyer: msg.sender,
            destinationWallet: destinationWallet, //store destination wallet
            tokenAmount: tokenAmount,
            fiatAmount: fiatAmount,
            timestamp: block.timestamp,
            gateway: gateway,
            status: PaymentStatus.Pending,
            gasFundAmount: msg.value,
            gasRefunded: false
        });
        
        // Reserve tokens by transferring them from owner to this contract
        // The owner must approve this contract to spend tokens beforehand
        bool success = token.transferFrom(owner(), address(this), tokenAmount);
        require(success, "Token reservation failed");
        
        emit PaymentCreated(paymentId, msg.sender, destinationWallet, tokenAmount, fiatAmount, gateway);
    }
    
    /**
     * @dev Process a payment notification from the payment gateway
     * @param paymentId Unique payment identifier
     * @param status Payment status (1 = Complete, 2 = Failed)
     * @param signature Digital signature from the payment gateway
     */
    function processPaymentCallback(
        string memory paymentId,
        uint8 status,
        bytes memory signature
    ) external nonReentrant {
        Payment storage payment = payments[paymentId];
        
        // Ensure payment exists and is pending
        require(payment.buyer != address(0), "Payment not found");
        require(payment.status == PaymentStatus.Pending, "Payment not in pending state");
        
        // Verify the signature from the payment gateway
        bytes32 messageHash = keccak256(abi.encodePacked(paymentId, status));
        bytes32 ethSignedMessageHash = MessageHashUtils.toEthSignedMessageHash(messageHash);
        address signer = ECDSA.recover(ethSignedMessageHash, signature);
        
        require(signer == gatewaySigners[payment.gateway], "Invalid signature");
        
        if (status == 1) {
            // Payment succeeded - transfer tokens to buyer
            payment.status = PaymentStatus.Completed;
            address recipient = payment.destinationWallet;
            bool success = token.transfer(recipient, payment.tokenAmount);
            require(success, "Token transfer failed");
            emit PaymentCompleted(paymentId, recipient, payment.tokenAmount);
        } else {
            // Payment failed - return tokens to owner
            payment.status = PaymentStatus.Failed;
            bool success = token.transfer(owner(), payment.tokenAmount);
            require(success, "Token return failed");
            emit PaymentFailed(paymentId, payment.buyer);
        }
        
        // Refund any remaining gas deposit to buyer minus processing costs
        uint256 processingCost = calculateProcessingCost();
        uint256 refundAmount = payment.gasFundAmount > processingCost ? 
                               payment.gasFundAmount - processingCost : 0;
                               
        if (refundAmount > 0 && !payment.gasRefunded) {
            payment.gasRefunded = true;
            (bool refundSuccess, ) = payment.buyer.call{value: refundAmount}("");
            require(refundSuccess, "Gas refund failed");
            emit GasRefunded(paymentId, payment.buyer, refundAmount);
        }
    }
    
    /**
     * @dev Calculate the cost for processing a payment callback (gas fee)
     * @return cost in wei for processing
     */
    function calculateProcessingCost() public view returns (uint256) {
        // Fixed cost for now, could be made dynamic based on network conditions
        return requiredGasDeposit / 2; // Using half of the deposit as the processing cost
    }
    
    /**
     * @dev Allow admin to manually process refunds for payments
     * @param paymentId Payment identifier to refund
     */
    function processRefund(string memory paymentId) external onlyOwner nonReentrant {
        Payment storage payment = payments[paymentId];
        
        // Ensure payment exists and is in appropriate state for refund
        require(payment.buyer != address(0), "Payment not found");
        require(payment.status == PaymentStatus.Pending || payment.status == PaymentStatus.Completed, 
                "Payment cannot be refunded");
        
        if (payment.status == PaymentStatus.Pending) {
            // If pending, return tokens to owner
            payment.status = PaymentStatus.Refunded;
            bool success = token.transfer(owner(), payment.tokenAmount);
            require(success, "Token return failed");
        } else if (payment.status == PaymentStatus.Completed) {
            // If completed, we need to get tokens back from buyer (not implemented here)
            payment.status = PaymentStatus.Refunded;
        }
        
        // Refund any remaining gas funds if not already refunded
        if (!payment.gasRefunded && payment.gasFundAmount > 0) {
            payment.gasRefunded = true;
            (bool refundSuccess, ) = payment.buyer.call{value: payment.gasFundAmount}("");
            require(refundSuccess, "Gas refund failed");
            emit GasRefunded(paymentId, payment.buyer, payment.gasFundAmount);
        }
        
        emit PaymentRefunded(paymentId, payment.buyer);
    }
    
    /**
     * @dev Withdraw accumulated processing fees to owner
     */
    function withdrawProcessingFees() external onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "No fees to withdraw");
        (bool success, ) = owner().call{value: balance}("");
        require(success, "Fee withdrawal failed");
    }
    
    /**
     * @dev Update required gas deposit
     * @param _requiredGasDeposit New required deposit amount in wei
     */
    function updateGasDepositRequirement(uint256 _requiredGasDeposit) external onlyOwner {
        requiredGasDeposit = _requiredGasDeposit;
        emit GasDepositRequirementUpdated(_requiredGasDeposit);
    }
    
    /**
     * @dev Check payment status
     * @param paymentId Payment identifier to check
     * @return status of the payment
     */
    function getPaymentStatus(string memory paymentId) external view returns (PaymentStatus) {
        return payments[paymentId].status;
    }
    
    /**
     * @dev Update token price
     * @param _pricePerToken New price per token in wei
     */
    function updateTokenPrice(uint256 _pricePerToken) external onlyOwner {
        require(_pricePerToken > 0, "Price must be greater than zero");
        pricePerToken = _pricePerToken;
        emit TokenPriceUpdated(_pricePerToken);
    }
    
    /**
     * @dev Update a gateway signer address
     * @param gateway Payment gateway identifier
     * @param signer New signer address
     */
    function updateGatewaySigner(string memory gateway, address signer) external onlyOwner {
        require(signer != address(0), "Invalid signer address");
        gatewaySigners[gateway] = signer;
    }
    
    /**
     * @dev Calculate token amount based on fiat amount
     * @param fiatAmount Amount in fiat currency
     * @return tokenAmount equivalent in tokens
     */
    function calculateTokenAmount(uint256 fiatAmount) external view returns (uint256) {
        return (fiatAmount * 10**18) / pricePerToken;
    }
    
    /**
     * @dev Mock function to simulate payment callback for testing
     * Can only be called by the contract owner
     */
    function mockPaymentCallback(string memory paymentId, uint8 status) external onlyOwner {
        Payment storage payment = payments[paymentId];
        require(payment.buyer != address(0), "Payment not found");
        require(payment.status == PaymentStatus.Pending, "Payment not in pending state");
        
        if (status == 1) {
            // Payment succeeded - use destinationWallet instead of buyer
            payment.status = PaymentStatus.Completed;
            address recipient = payment.destinationWallet;
            bool success = token.transfer(recipient, payment.tokenAmount);
            require(success, "Token transfer failed");
            emit PaymentCompleted(paymentId, recipient, payment.tokenAmount);
        } else {
            // Payment failed
            payment.status = PaymentStatus.Failed;
            bool success = token.transfer(owner(), payment.tokenAmount);
            require(success, "Token return failed");
            emit PaymentFailed(paymentId, payment.buyer);
        }
        
        // Refund any remaining gas deposit
        if (!payment.gasRefunded && payment.gasFundAmount > 0) {
            uint256 processingCost = calculateProcessingCost();
            uint256 refundAmount = payment.gasFundAmount > processingCost ? 
                                  payment.gasFundAmount - processingCost : 0;
            
            if (refundAmount > 0) {
                payment.gasRefunded = true;
                (bool refundSuccess, ) = payment.buyer.call{value: refundAmount}("");
                require(refundSuccess, "Gas refund failed");
                emit GasRefunded(paymentId, payment.buyer, refundAmount);
            }
        }
    }
    
    // Function to receive Ether
    receive() external payable {}
}