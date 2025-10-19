#!/bin/bash
set -e

cleanup() {
    echo -e "\n🛑 Stopping dev servers..."

    kill "$FRONTEND_PID" 2>/dev/null || true
    kill "$BACKEND_PID" 2>/dev/null || true

    exit 0
}
trap cleanup SIGINT SIGTERM

echo "⚡️ Starting frontend..."
cd frontend
npm install
npm run dev &
FRONTEND_PID=$!
cd ..

echo "🖥️ Starting backend..."
cd backend
go run . &
BACKEND_PID=$!
cd ..

while ! curl -s http://localhost:8000 > /dev/null 2>&1; do
    sleep 1
done

wait $FRONTEND_PID $BACKEND_PID