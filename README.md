# AI Mail Sender CLI

![Go](https://img.shields.io/badge/Go-1.25%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Node.js](https://img.shields.io/badge/Node.js-Required_for_scripts-339933?style=for-the-badge&logo=node.js&logoColor=white)
![SMTP](https://img.shields.io/badge/SMTP-Enabled-FF6B6B?style=for-the-badge&logo=gmail&logoColor=white)
![Concurrency](https://img.shields.io/badge/Concurrency-10_workers-7A5CFA?style=for-the-badge)
![Status](https://img.shields.io/badge/Status-Production_Ready-2EA043?style=for-the-badge)

This project sends personalized emails from a CSV list, attaches a PDF file, and logs each send result to output.log.

## Requirements ✅

- 🐹 Go 1.25+
- 🟢 Node.js and npm (only used to run helper scripts)
- 📬 SMTP credentials for your mail provider

## Project Setup ⚙️

1. Install npm dependencies (if needed for your environment):

   npm install

2. Open and update SMTP values in [smtp_env.sh](smtp_env.sh):
   - SMTP_HOST
   - SMTP_PORT
   - SMTP_USERNAME
   - SMTP_PASSWORD
   - SMTP_FROM_NAME
   - SMTP_FROM
   - SMTP_SUBJECT

3. Make sure your data and attachment files are present:
   - CSV file: [assets/dataset.csv](assets/dataset.csv)
   - PDF attachment: [assets/Hunar Conversational AI Agents_Self Serve_V1.pdf](assets/Hunar%20Conversational%20AI%20Agents_Self%20Serve_V1.pdf)

## NPM Scripts 📜

This repo uses two main npm commands:

1. npm run dev:secrets

   What it does:
   - 🔐 Loads environment variables from [smtp_env.sh](smtp_env.sh)
   - 🧪 Opens a new shell session with those variables available

   Use this first.

2. npm run dev:start

   What it does:
   - 🚀 Runs the Go mail sender app

   Use this after running dev:secrets.

## Typical Run Flow ▶️

In terminal 1:

npm run dev:secrets

In the new shell session:

npm run dev:start

## Optional Arguments 🧩

You can also run the app directly with custom paths:

go run . <csv_file_path> <attachment_file_path>

Example:

go run . assets/dataset_test.csv "assets/Hunar Conversational AI Agents_Self Serve_V1.pdf"

## Output 🧾

- 📺 Console prints per-email success or failure with time taken.
- 🗂️ [output.log](output.log) stores the same result lines.

## Notes 💡

- ⚡ Emails are sent concurrently in batches of 10.
- ⏱️ After every 10 successful emails, the app waits 1 minute before continuing.
- 🔒 Keep [smtp_env.sh](smtp_env.sh) secure, since it contains SMTP credentials.
