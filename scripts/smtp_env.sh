#!/usr/bin/env sh

# SMTP server settings
export SMTP_HOST="smtp.gmail.com"
export SMTP_PORT="587"

# SMTP auth credentials
export SMTP_USERNAME="mohitchitkara@hunar.ai"
export SMTP_PASSWORD="zheg gytc glqb oogo"

# Sender and email metadata
export SMTP_FROM_NAME="Mohit Chitkara"
export SMTP_FROM="mohitchitkara@hunar.ai"

# AI Options
export AI_KEY="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"

print_var_status() {
	var_name="$1"
	eval "var_value=\${$var_name}"
	if [ -n "$var_value" ]; then
		echo "$var_name set ✅"
	else
		echo "$var_name missing ❌"
	fi
}

print_var_status "SMTP_HOST"
print_var_status "SMTP_PORT"
print_var_status "SMTP_USERNAME"
print_var_status "SMTP_PASSWORD"
print_var_status "SMTP_FROM_NAME"
print_var_status "SMTP_FROM"
print_var_status "AI_KEY"

