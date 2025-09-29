# 📊 Stock Investment Analyst Hub

A full-stack application that aggregates analyst recommendations and provides AI-powered investment insights. Built with Go backend and Vue.js frontend.

![Stock Investment Hub](https://via.placeholder.com/800x400/4F46E5/FFFFFF?text=Stock+Investment+Hub)

## 🚀 Features

- **Real-time Data Loading**: Progressive loading of analyst recommendations with live progress tracking
- **Interactive Dashboard**: Beautiful cards showing recommendation statistics with animated counters
- **Advanced Table View**: Sortable, filterable table with pagination and search capabilities
- **Visual Heatmap**: Interactive recommendation heatmap by company and brokerage
- **AI Investment Algorithm**: Smart scoring system analyzing 6 key factors to recommend best stocks
- **Responsive Design**: Professional financial dashboard that works on all devices

## ��️ Architecture

```markdown:README.md
<code_block_to_apply_changes_from>
```
┌─────────────────┐    HTTP API    ┌─────────────────┐
│   Vue.js 3      │ ◄─────────────► │   Go Backend    │
│   Frontend      │                 │                 │
│                 │                 │                 │
│ • TypeScript    │                 │ • Gin/Gorilla   │
│ • Pinia Store   │                 │ • PostgreSQL    │
│ • TailwindCSS   │                 │ • REST API      │
│ • Vite          │                 │ • Data Fetching │
└─────────────────┘                 └─────────────────┘
```

## 📦 Project Structure

```
stock-investment/
├── backend/                 # Go backend application
│   ├── connection/         # Database connection logic
│   ├── server/            # HTTP server and handlers
│   ├── service/           # Business logic and API services
│   ├── main.go           # Application entry point
│   └── .env.example      # Environment variables template
├── frontend/frontend/     # Vue.js frontend application
│   ├── src/
│   │   ├── components/   # Reusable Vue components
│   │   ├── views/        # Page components
│   │   ├── stores/       # Pinia state management
│   │   ├── services/     # API service layer
│   │   ├── utils/        # Utility functions and algorithms
│   │   └── types/        # TypeScript type definitions
│   ├── .env.example      # Environment variables template
│   └── package.json
└── README.md             # This file
```

## 🚀 Quick Start

### Prerequisites

- **Go 1.19+**
- **Node.js 18+**
- **PostgreSQL 12+**
- **Git**

### 1. Clone the Repository

```bash
git clone <repository-url>
cd stock-investment
```

### 2. Set Up Backend

```bash
cd backend

# Copy environment template
cp .env.example .env

# Edit .env with your database credentials and API settings
nano .env

# Install dependencies
go mod download

# Run the backend
go run main.go -api -port=8080
```

### 3. Set Up Frontend

```bash
cd frontend/frontend

# Copy environment template
cp .env.example .env

# Edit .env with your API settings
nano .env

# Install dependencies
npm install

# Start development server
npm run dev
```

### 4. Set Up Database

Create the required PostgreSQL tables:

```sql
-- Companies table
CREATE TABLE company (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
  ticker VARCHAR(10) UNIQUE NOT NULL, 
  name VARCHAR(255) NOT NULL, 
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

-- Brokerages table
CREATE TABLE brokerage (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) UNIQUE NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

-- Analyst recommendations table
CREATE TABLE analyst_recommendation(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  company_id UUID REFERENCES company(id) ON DELETE CASCADE NOT NULL,
  brokerage_id UUID REFERENCES brokerage(id) ON DELETE SET NULL,
  target_from DECIMAL(10,2),
  target_to DECIMAL(10,2),
  rating_from VARCHAR(50),
  rating_to VARCHAR(50),
  action VARCHAR(100) NOT NULL,
  time TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

-- Create indexes
CREATE INDEX company_ticker_idx ON company(ticker);
CREATE INDEX brokerage_name_idx ON brokerage(name);
CREATE INDEX analyst_recommendation_company_time ON analyst_recommendation(company_id, time DESC);
CREATE INDEX analyst_recommendation_time ON analyst_recommendation(time DESC);
CREATE INDEX analyst_recommendation_brokerage_time ON analyst_recommendation(brokerage_id, time DESC);
```

## 🔧 Configuration

### Backend Environment Variables

See `backend/.env.example` for all available options.

### Frontend Environment Variables

See `frontend/frontend/.env.example` for configuration options.

## 📖 Usage

### 1. Data Loading
- Start the backend to fetch data from external API
- Data loads progressively (50 items at a time)
- View real-time progress on the frontend dashboard

### 2. Navigation
- **Home**: Dashboard with stats and AI investment recommendations
- **Recommendations**: Table and heatmap views of all analyst recommendations
- **Companies**: Browse all companies with recommendations
- **Brokerages**: View all brokerage firms and their activity

### 3. Investment Algorithm

The AI investment algorithm analyzes 6 key factors:

| Factor | Weight | Description |
|--------|--------|-------------|
| **Recommendation Volume** | 15% | Number of analyst recommendations |
| **Upgrade Ratio** | 25% | Percentage of positive rating changes |
| **Target Upside** | 20% | Potential price appreciation |
| **Analyst Consensus** | 20% | Agreement among analysts |
| **Recency** | 10% | How recent the recommendations are |
| **Brokerage Quality** | 10% | Reputation of recommending firms |

## 🛠️ Development

### Backend Commands

```bash
# Fetch data from external API
go run main.go -fetch

# Start API server
go run main.go -api -port=8080

# Test database connection
go run main.go -test-db
```

### Frontend Commands

```bash
# Development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Type checking
npm run type-check
```

### API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Health check |
| `/api/v1/companies` | GET | List all companies |
| `/api/v1/companies/{ticker}` | GET | Get company by ticker |
| `/api/v1/brokerages` | GET | List all brokerages |
| `/api/v1/recommendations` | GET | List recommendations (paginated) |
| `/api/v1/recommendations/company/{ticker}` | GET | Get recommendations for company |
| `/api/v1/recommendations/brokerage/{id}` | GET | Get recommendations from brokerage |

## 🎨 Screenshots

### Dashboard


### Recommendations Table


### Interactive Heatmap



## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

## 📧 Contact

Your Name - your.email@example.com

Project Link: [https://github.com/your-username/stock-investment](https://github.com/your-username/stock-investment)
```

## **2. Backend .env.example**

```env:backend/.env.example
# Database Configuration
DATABASE_USER=your_postgres_username
DATABASE_PASSWORD=your_postgres_password
DATABASE_URL=@localhost:5432/stock_investment

# External API Configuration
API_URL=https://your-external-api.com/api/recommendations
BEARER_TOKEN=your_api_bearer_token

# Server Configuration
PORT=8080
HOST=localhost

# Environment
ENVIRONMENT=development

# CORS Configuration
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000

# Logging
LOG_LEVEL=info

# Rate Limiting (requests per minute)
RATE_LIMIT=1000

# Database Pool Settings
DB_MAX_CONNECTIONS=25
DB_MAX_IDLE_CONNECTIONS=5
DB_CONNECTION_MAX_LIFETIME=5m

# SSL Mode for database connection
# Options: disable, require, verify-ca, verify-full
SSL_MODE=verify-full

# Timeout settings (in seconds)
HTTP_TIMEOUT=30
DB_TIMEOUT=10
```

## **3. Frontend .env.example**

```env:frontend/frontend/.env.example
# API Configuration
VITE_API_BASE_URL=http://localhost:8080
VITE_API_VERSION=v1
VITE_API_TIMEOUT=10000

# Application Configuration
VITE_APP_NAME=AnalystHub
VITE_APP_VERSION=1.0.0
VITE_APP_DESCRIPTION=Stock Investment Analyst Hub

# Environment
VITE_NODE_ENV=development

# Features Flags
VITE_ENABLE_ANALYTICS=false
VITE_ENABLE_DEBUG=true
VITE_ENABLE_MOCK_DATA=false

# UI Configuration
VITE_DEFAULT_PAGE_SIZE=25
VITE_MAX_PAGE_SIZE=100
VITE_ANIMATION_DURATION=300

# External Services (if needed)
VITE_GOOGLE_ANALYTICS_ID=
VITE_SENTRY_DSN=

# Development Settings
VITE_HOT_RELOAD=true
VITE_SOURCE_MAPS=true

# Build Configuration
VITE_BUILD_TARGET=es2020
VITE_MINIFY=true
```

## **4. Additional Documentation Files**

### **DEPLOYMENT.md**

```markdown:DEPLOYMENT.md
# 🚀 Deployment Guide

## Production Deployment

### Backend Deployment

1. **Build the application:**
```bash
go build -o stock-investment-backend main.go
```

2. **Set production environment variables:**
```bash
export ENVIRONMENT=production
export DATABASE_URL=your_production_db_url
export API_URL=your_production_api_url
```

3. **Run the application:**
```bash
./stock-investment-backend -api -port=8080
```

### Frontend Deployment

1. **Build for production:**
```bash
npm run build
```

2. **Deploy to static hosting:**
```bash
# Upload dist/ folder to your hosting provider
# Examples: Netlify, Vercel, AWS S3, etc.
```

### Docker Deployment

See `docker-compose.yml` for containerized deployment.

### Environment-Specific Configurations

- **Development**: Use `.env.local`
- **Staging**: Use `.env.staging`
- **Production**: Use `.env.production`
```

### **API.md**

```markdown:API.md
# 📡 API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
Include Bearer token in Authorization header:
```
Authorization: Bearer your_token_here
```

## Endpoints

### Health Check
```http
GET /health
```

### Companies
```http
GET /companies
GET /companies/{ticker}
```

### Brokerages
```http
GET /brokerages
```

### Recommendations
```http
GET /recommendations?limit=50&offset=0&ticker=AAPL&brokerage_id=uuid
GET /recommendations/company/{ticker}
GET /recommendations/brokerage/{id}
```

## Response Format
```json
{
  "success": true,
  "data": [...],
  "meta": {
    "total": 2961,
    "limit": 50,
    "offset": 0
  }
}
```
```

