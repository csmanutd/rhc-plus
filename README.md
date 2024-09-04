# Rule-hit-counter-report-plus
This tool will add necessary info (src/dst, labels etc.) of rules for the rule hit counter instead of just showing the href

1. Export rules csv file via Workloader
2. run the tool    
   ./rhc_plus -report <rule_hit_counter_report.csv> -rules <rules.csv>
3. A new csv file <report_new.csv> will be generated

## Caveat
Override deny rules are not supported yet as Workloader can't export those for now.      
It will be added in a future release.
