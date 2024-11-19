FROM nginx:alpine

# Install SSL certificates
COPY certs/ /etc/nginx/certs/

# Copy nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Create directory for frontend files
RUN mkdir -p /usr/share/nginx/html

# The actual frontend files will be mounted as a volume
