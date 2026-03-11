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
export SMTP_SUBJECT="Optimizing Your Hiring Process with Hunar.ai"

# Optional: override defaults used by the Go program
# export CSV_PATH="assets/dataset.csv"
# export ATTACHMENT_PATH="assets/Hunar Conversational AI Agents_Self Serve_V1.pdf"

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
print_var_status "SMTP_SUBJECT"
