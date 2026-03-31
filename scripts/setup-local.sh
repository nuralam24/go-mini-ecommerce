#!/bin/bash

echo "🚀 Setting up Local Development Environment"
echo "=========================================="
echo ""

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Check PostgreSQL
echo -e "${BLUE}Step 1: Checking PostgreSQL...${NC}"
if command -v psql &> /dev/null; then
    echo -e "   ✅ ${GREEN}PostgreSQL installed${NC}"
    PSQL_VERSION=$(psql --version | awk '{print $3}')
    echo "   Version: $PSQL_VERSION"
else
    echo -e "   ${RED}❌ PostgreSQL not found${NC}"
    echo ""
    echo "Install PostgreSQL:"
    echo "  macOS:  brew install postgresql@15"
    echo "  Ubuntu: sudo apt-get install postgresql"
    echo ""
    exit 1
fi
echo ""

# Check if PostgreSQL is running
echo -e "${BLUE}Step 2: Checking if PostgreSQL is running...${NC}"
if pg_isready &> /dev/null; then
    echo -e "   ✅ ${GREEN}PostgreSQL is running${NC}"
else
    echo -e "   ${YELLOW}⚠️  PostgreSQL not running${NC}"
    echo "   Starting PostgreSQL..."
    
    # Try to start (macOS)
    if [[ "$OSTYPE" == "darwin"* ]]; then
        brew services start postgresql@15 &> /dev/null || brew services start postgresql &> /dev/null
        sleep 2
        if pg_isready &> /dev/null; then
            echo -e "   ✅ ${GREEN}PostgreSQL started${NC}"
        fi
    else
        echo "   Run: sudo service postgresql start"
    fi
fi
echo ""

# Create database
echo -e "${BLUE}Step 3: Creating database...${NC}"
if psql postgres -lqt | cut -d \| -f 1 | grep -qw go_commerce; then
    echo -e "   ✅ ${GREEN}Database 'go_commerce' already exists${NC}"
else
    echo "   Creating database 'go_commerce'..."
    createdb go_commerce 2>/dev/null
    if [ $? -eq 0 ]; then
        echo -e "   ✅ ${GREEN}Database created${NC}"
    else
        echo -e "   ${YELLOW}⚠️  Database may already exist or permission issue${NC}"
    fi
fi
echo ""

# Update .env
echo -e "${BLUE}Step 4: Updating .env file...${NC}"
cat > .env << 'EOF'
DATABASE_URL='postgresql://postgres@localhost:5432/go_commerce?sslmode=disable'
PORT=8080
ENV=development
JWT_SECRET='dev-secret-key-change-in-production'
EOF
echo -e "   ✅ ${GREEN}.env file updated${NC}"
echo ""

# Run migrations
echo -e "${BLUE}Step 5: Running database migrations...${NC}"
if [ -f "db/migrations/001_schema.sql" ]; then
    psql postgresql://postgres@localhost:5432/go_commerce -f db/migrations/001_schema.sql > /dev/null 2>&1
    if [ $? -eq 0 ]; then
        echo -e "   ✅ ${GREEN}Migrations completed${NC}"
    else
        echo -e "   ${YELLOW}⚠️  Migrations may have already run${NC}"
    fi
else
    echo -e "   ${RED}❌ Migration file not found${NC}"
fi
echo ""

# Install dependencies
echo -e "${BLUE}Step 6: Installing Go dependencies...${NC}"
go mod download > /dev/null 2>&1
echo -e "   ✅ ${GREEN}Dependencies installed${NC}"
echo ""

# Build
echo -e "${BLUE}Step 7: Building application...${NC}"
if go build -o bin/server cmd/server/main.go 2>/dev/null; then
    echo -e "   ✅ ${GREEN}Build successful${NC}"
else
    echo -e "   ${RED}❌ Build failed${NC}"
    exit 1
fi
echo ""

echo "=========================================="
echo -e "${GREEN}✅ Setup Complete!${NC}"
echo "=========================================="
echo ""
echo "To start the server:"
echo -e "  ${BLUE}make run${NC}"
echo ""
echo "Or use live reload:"
echo -e "  ${BLUE}make watch${NC}"
echo ""
echo "To test:"
echo "  curl http://localhost:8080/health"
echo "  open http://localhost:8080/swagger/index.html"
echo ""
echo "=========================================="
