# Web3 Token Sale Backend

A comprehensive Go-based backend service for Web3 token sales, featuring blockchain integration, payment processing, wallet management, and user authentication.

## 🚀 Features

### Core Functionality
- **Token Sales & Swapping**: Complete token purchase and swap functionality
- **Multi-Payment Support**: Fiat-to-crypto payments via Transak and Midtrans
- **Blockchain Integration**: Direct Ethereum blockchain interaction with smart contracts
- **Wallet Management**: HD wallet creation, import, and balance tracking
- **Uniswap Integration**: Automated token swapping through Uniswap V3

### Security & Authentication
- **JWT Authentication**: Secure user session management
- **2FA Support**: Two-factor authentication with recovery options
- **OTP Verification**: Email-based OTP for secure operations
- **Token Blacklisting**: Session invalidation and security controls
- **Recovery System**: Account and wallet recovery mechanisms

### Advanced Features
- **Real-time Price Feeds**: Live token pricing and market data
- **Transaction Monitoring**: Comprehensive transaction tracking and status updates
- **Activity Logging**: Detailed audit trails for all user activities
- **Auto-swap Functionality**: Automated token conversion workflows
- **Smart Contract Bindings**: Generated Go bindings for contract interactions

## 🏗️ Architecture

```
web3-tokensale-be/
├── cmd/                    # Application entry point
├── internal/
│   ├── api/               # HTTP API layer
│   │   ├── auth/          # Authentication handlers
│   │   ├── handlers/      # Business logic handlers
│   │   ├── middleware/    # HTTP middleware
│   │   └── routes/        # Route definitions
│   ├── blockchain/        # Blockchain interaction layer
│   ├── config/           # Configuration management
│   ├── database/         # Database connections and models
│   ├── models/           # Data models
│   └── services/         # Business logic services
├── pkg/                  # Shared packages
└── bindings/            # Smart contract bindings
```

## 🛠️ Technology Stack

- **Language**: Go 1.23.1
- **Web Framework**: Gin
- **Database**: PostgreSQL with GORM
- **Blockchain**: Ethereum (go-ethereum)
- **Authentication**: JWT with 2FA support
- **Payment Processing**: Transak, Midtrans
- **Smart Contracts**: OpenZeppelin, Solidity
- **Email Service**: SendGrid

## 📋 Prerequisites

- Go 1.23.1 or higher
- PostgreSQL database
- Ethereum node access (Infura, Alchemy, or local node)
- Environment variables configured (see Configuration section)

## ⚙️ Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/Illuminateee/web3-backend-go.git
   cd web3-tokensale-be
   ```

2. **Install Go dependencies**
   ```bash
   go mod download
   ```

3. **Install smart contract dependencies**
   ```bash
   npm install
   ```

4. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

5. **Run database migrations**
   ```bash
   # Ensure PostgreSQL is running and configured
   go run cmd/main.go
   ```

## 🔧 Configuration

Create a `.env` file in the root directory with the following variables:

### Server Configuration
```env
PORT=8080
ENVIRONMENT=development
```

### Database Configuration
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=web3_tokensale
DB_SSL_MODE=disable
```

### Blockchain Configuration
```env
ETHEREUM_RPC=https://mainnet.infura.io/v3/your-project-id
PRIVATE_KEY=your_private_key
WALLET_PRIVATE_KEY=your_wallet_private_key
TOKEN_ADDRESS=0x...
PAYMENT_GATEWAY_ADDRESS=0x...
STABLECOIN_ADDRESS=0x...
WETH_ADDRESS=0x...
UNISWAP_ROUTER_ADDRESS=0x...
```

### Authentication
```env
JWT_SECRET=your_jwt_secret
JWT_EXPIRATION=24h
```

### Email Service
```env
SENDGRID_API_KEY=your_sendgrid_api_key
FROM_EMAIL=noreply@yourdomain.com
FROM_NAME=Your App Name
APP_URL=https://yourdomain.com
```

### Payment Services
```env
TRANSAK_API_KEY=your_transak_api_key
TRANSAK_SECRET_KEY=your_transak_secret_key
TRANSAK_BASE_URL=https://api.transak.com
```

## 🚀 Running the Application

### Development
```bash
go run cmd/main.go
```

### Production Build
```bash
go build -o bin/web3-tokensale-be cmd/main.go
./bin/web3-tokensale-be
```

The server will start on the configured port (default: 8080) and display connection information.

## 📚 API Documentation

### Authentication Endpoints
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `POST /api/auth/logout` - User logout
- `POST /api/auth/refresh` - Token refresh

### Wallet Management
- `POST /api/wallet/create` - Create new wallet
- `POST /api/wallet/import` - Import existing wallet
- `GET /api/wallet/balance` - Get wallet balance

### Token Operations
- `POST /api/tokens/swap` - Swap tokens
- `GET /api/tokens/price` - Get token prices
- `GET /api/transactions` - Get transaction history

### Payment Processing
- `POST /api/payments/transak/order` - Create Transak order
- `POST /api/payments/fiat-to-token` - Fiat to token conversion
- `POST /api/payments/webhook` - Payment webhooks

### 2FA Management
- `POST /api/2fa/setup` - Setup 2FA
- `POST /api/2fa/verify` - Verify 2FA code
- `POST /api/2fa/disable` - Disable 2FA

## 🧪 Testing

Run the test suite:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

## 🔒 Security Features

- **JWT Token Management**: Secure authentication with token blacklisting
- **2FA Integration**: TOTP-based two-factor authentication
- **OTP Verification**: Email-based one-time passwords
- **Request Validation**: Input sanitization and validation
- **Rate Limiting**: API rate limiting middleware
- **CORS Protection**: Cross-origin request security
- **Recovery Systems**: Account and wallet recovery mechanisms

## 🔗 Smart Contract Integration

The application integrates with several smart contracts:
- **ERC20 Token**: Custom token implementation
- **Payment Gateway**: Fiat-to-crypto payment processing
- **Uniswap Router**: Automated token swapping
- **Test Contracts**: Development and testing utilities

Contract bindings are automatically generated and stored in the `bindings/` directory.

## 📊 Database Schema

Key database models:
- **Users**: User accounts and authentication
- **Wallets**: HD wallet management
- **Transactions**: Transaction tracking and history
- **Activity Logs**: Audit trail for user actions
- **Transak Orders**: Payment processing records

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

For support and questions:
- Create an issue in the GitHub repository
- Contact the development team
- Check the documentation for common solutions

## 🔄 Changelog

### Latest Updates
- Enhanced Uniswap V3 integration
- Improved wallet recovery system
- Advanced 2FA functionality
- Comprehensive activity logging
- Payment gateway optimizations

---

**Note**: This is a production-ready Web3 backend service. Ensure all environment variables are properly configured and security measures are in place before deployment.