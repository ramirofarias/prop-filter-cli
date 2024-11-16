# Prop Filter CLI

CLI tool for filtering and transforming property data from JSON or CSV files

## Features

- Filter properties by multiple criteria:
  - Square footage
  - Number of bathrooms
  - Price
  - Distance from coordinates
  - Lighting
  - Description keywords
  - Available ammenities
- Support for both JSON and CSV input/output
- Flexible comparison operators (greater than, less than, equals, etc.)
- Distance-based filtering using geographical coordinates
- Keyword search in property descriptions

## Installation

1. Go to the [Releases](https://github.com/ramirofarias/prop-filter-cli/releases) page
2. Download the appropriate binary for your operating system:
   - Windows: `prop-filter-cli_windows_amd64.exe`
   - macOS: `prop-filter-cli_darwin_amd64`
   - Linux: `prop-filter-cli_linux_amd64`

For Linux/macOS users, make the binary executable:

```bash
chmod +x prop-filter-cli_*
```

## Usage

```bash
prop-filter-cli --input <input-file> [flags]
```

### Required Flags

- `--input`: Path to JSON or CSV input file (required)

### Optional Flags

- `--sqft`: Filter by square footage
  - Examples: "gt 1500", "eq 1500", "lt 1500", "lte 1500", "in 1500,2000"
- `--bathrooms`: Filter by number of bathrooms
  - Examples: "gt 1", "eq 1", "lt 3", "lte 3", "gte 3", "in 1,3"
- `--rooms`: Filter by number of rooms
  - Examples: "gt 1", "eq 1", "lt 3", "lte 3", "gte 3", "in 1,3"
- `--distance`: Filter by distance in km (requires --lat and --long)
  - Examples: "gt 100", "eq 100", "lt 100", "lte 100", "gte 100", "in 150,200"
- `--price`: Filter by price
  - Examples: "gt 1000", "eq 1000", "lt 1000", "lte 1000", "gte 1000"
- `--lat`: Latitude for distance calculations
- `--long`: Longitude for distance calculations
- `--lighting`: Filter by lighting type
  - Possible values: "low", "medium", "high"
- `--keywords`: Search keywords in description (comma-separated)
  - Example: "spacious,big"
- `--ammenities`: Required amenities (comma-separated)
  - Example: "garage,yard"
- `--output`: Output file path (.csv or .json)
  - Example: "output.json" or "output.csv"

## Examples

### Basic Filtering

```bash
# Filter properties larger than 1500 sq ft
prop-filter-cli --input properties.json --sqft "gt 1500"

# Filter properties with 2 or more bathrooms
prop-filter-cli --input properties.csv --bathrooms "gte 2"

# Filter properties under $500,000
prop-filter-cli --input properties.json --price "lt 500000"
```

### Combined Filters

```bash
prop-filter-cli --input properties.json \
  --lighting "high" \
  --ammenities "pool"

prop-filter-cli --input properties.json \
  --sqft "gt 2000" \
  --price "lte 600000" \
  --keywords "spacious,modern"
```

### Distance Filtering

```bash
# Find properties within 10km of coordinates
prop-filter-cli --input properties.json \
  --lat 34.0522 \
  --long -118.2437 \
  --distance "lte 10"
```

### Output to File

```bash
# Export filtered results to CSV
prop-filter-cli --input properties.json \
  --price "lt 400000" \
  --output "affordable_properties.csv"

# Export filtered results to JSON
prop-filter-cli --input properties.csv \
  --lighting "high" \
  --ammenities "garage,pool" \
  --output "luxury_properties.json"
```

## Comparison Operators

- `gt`: Greater than
- `gte`: Greater than or equal to
- `lt`: Less than
- `lte`: Less than or equal to
- `eq`: Equal to
- `in`: Within range (comma-separated values)

## Input File Format

### JSON Format

```json
[
  {
    "squareFootage": 1500,
    "lighting": "medium",
    "price": 300000,
    "rooms": 3,
    "bathrooms": 2,
    "location": [34.0522, -118.2437],
    "description": "Charming 3-bedroom home",
    "ammenities": {
      "yard": true,
      "garage": true,
      "pool": false
    }
  }
]
```

### CSV Format

```csv
squareFootage,lighting,price,rooms,bathrooms,latitude,longitude,description,ammenities
200,medium,250000.00,3,2,34.052200,-118.243700,Charming 3-bedroom home in a quiet neighborhood with easy access to parks and schools.,"{""garage"":true,""pool"":false,""yard"":true}"
```
