# Start with a stage named 'builder' using official Node.js LTS version image
FROM node:lts AS node-builder

# Set the working directory
WORKDIR /app

# Copy necessary Node.js files, assets, components for Tailwind to look at, and the Tailwind config
COPY package*.json ./
COPY ./assets ./assets
COPY ./html ./html
COPY tailwind.config.js ./

# Install production dependencies and run Tailwind CSS compilation
RUN npm install -g pnpm
RUN pnpm install --production && pnpm exec tailwindcss -i ./assets/global.css -o ./public/global.css --minify

# Start new stage with official Go image
FROM golang:1.21.3 AS builder

# Install the templ tool
RUN go install github.com/a-h/templ/cmd/templ@v0.2.408

# Set the working directory inside the Docker image
WORKDIR /app

# Copy local project files into Docker image
COPY . .

# Copy generated global.css from previous stage (node-builder)
COPY --from=node-builder /app/public ./public

# Run templ generate
RUN templ generate

# Build the Go application
RUN go build -tags netgo -ldflags '-s -w' -o app

# Start a new final stage for a smaller and cleaner final image using debian:stretch-slim
FROM debian:stretch-slim

# Copy Go binary and public directory from builder stage
COPY --from=builder /app/app /app/app
COPY --from=builder /app/public /public

# Expose application port (assuming its 8080)
EXPOSE 8080

# Run the application
CMD ["/app/app"]
