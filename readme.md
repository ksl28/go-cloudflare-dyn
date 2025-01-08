# Cloudflare DNS Updater

This Go program interacts with the Cloudflare API to monitor and update DNS records to ensure they point to the correct public IP address. It is particularly useful for managing dynamic DNS (DDNS) when your public IP address changes over time.

## Features
- Fetches DNS records for a specified Cloudflare zone.
- Checks if any specified DNS records are pointing to the correct public IP address.
- Automatically updates the DNS records with the current public IP if they are incorrect.
- Periodically refreshes to recheck and update DNS records based on a configurable time interval.

## Requirements
- Go 1.23+.
- A Cloudflare account and API key.
- A Cloudflare zone ID.

## Usage

### Command Line Flags:
- `-records`: Specify one or more DNS records to monitor. This flag can be repeated to specify multiple records.
- `-apiKey`: Your Cloudflare API key.
- `-zoneID`: Your Cloudflare zone ID.
- `-refreshSeconds`: The number of seconds to wait before checking DNS records again. Default is 300 seconds (5 minutes).

### Example:

To run the program, use the following command:
```bash
go run main.go -apiKey <YOUR_API_KEY> -zoneID <YOUR_ZONE_ID> -records "example.com" -records "sub.example.com" -refreshSeconds 600
