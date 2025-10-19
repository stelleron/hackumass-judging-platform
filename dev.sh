#!/bin/bash
set -e

# Start frontend in background
sh frontend.sh &
FRONTEND_PID=$!

# Start backend
sh backend.sh

# Optional: wait for frontend to exit if backend terminates
wait $FRONTEND_PID $BACKEND_PID