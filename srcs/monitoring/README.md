# Monitoring System

This directory contains the configuration for our monitoring system using Prometheus and Grafana.

## Data Retention Policies

### Prometheus
- Metric data is retained for 15 days.
- Data is stored in a Docker volume named `prometheus_data` for persistence.

### Grafana
- Dashboard versions are limited to the last 5 versions.
- Grafana's internal data (dashboards, users, etc.) is stored in a Docker volume named `grafana-data`.

## Security Measures
- User sign-ups are disabled.
- Anonymous access is disabled.
- External snapshots are disabled.
- Cookies are set to secure and same-site strict.

To modify these settings, please update the respective configuration files and rebuild the Docker containers.