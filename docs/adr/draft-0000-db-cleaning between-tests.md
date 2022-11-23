# Database cleaning between tests

# Problem

# Approaches considered

# Decision

truncate table + having any initial data (ej. entries for "enum" tables such a
post\_types) in an seed.sql file that we execute after truncation each time.
