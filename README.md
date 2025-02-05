# Rule-hit-counter-report-plus
A tool that enhances rule hit counter reports by adding detailed rule information (sources, destinations, labels, etc.) instead of just showing the rule href.

## Usage Steps
1. Export rules using Workloader:
   ```bash
   ./workloader rule-export
   ```

2. Generate a rule hit counter report from your PCE

3. Run the tool:
   ```bash
   ./rhc_plus -report <rule_hit_counter_report.csv> -rules <rules.csv>
   ```
   
   Optional: Add `-json` flag to also generate JSON output
   ```bash
   ./rhc_plus -report <rule_hit_counter_report.csv> -rules <rules.csv> -json
   ```

4. The tool will generate:
   - A new CSV file named `<original_report>_plus.csv`
   - If JSON flag is used, also generates `<original_report>_plus.json`

## Features
- Combines rule hit counter data with detailed rule information
- Supports both CSV and JSON output formats
- Handles both active and draft rules

## Limitations
Override deny rules are not currently supported as Workloader cannot export these rules at this time.
This functionality will be added in a future release.
