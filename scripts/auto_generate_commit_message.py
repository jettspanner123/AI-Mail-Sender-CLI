import os
from http.client import responses

from ollama import chat
from dataclasses import dataclass
import json
import subprocess
from groq import Groq

@dataclass()
class ResponseData:
    type: str
    description: str


def ask_to_stage_changes():
    # ANSI color codes
    YELLOW = "\033[93m"
    GREEN = "\033[92m"
    RED = "\033[91m"
    CYAN = "\033[96m"
    RESET = "\033[0m"

    print(f"{YELLOW}No staged changes detected.{RESET}")
    print(f"{CYAN}Here are your unstaged changes:{RESET}")
    subprocess.run(["git", "status", "--short"])
    while True:
        print(f"{YELLOW}Would you like to stage ALL changes and continue? (y/n){RESET}")
        user_input = input(f"{CYAN}Your choice: {RESET}").strip().lower()
        if user_input in ["y", "yes"]:
            subprocess.run(["git", "add", "."])
            print(f"{GREEN}All changes staged. Continuing...{RESET}")
            break
        elif user_input in ["n", "no"]:
            print(f"{RED}No changes staged. Exiting.{RESET}")
            exit(0)
        else:
            print(f"{RED}Invalid input. Please enter 'y' or 'n'.{RESET}")
    pass


def gather_git_info():
    diff = subprocess.check_output(["git", "diff", "--cached"]).decode()

    if not diff or diff.strip() == "":
        ask_to_stage_changes()
        diff = subprocess.check_output(["git", "diff", "--cached"]).decode()

    return diff

def generate_commit_message_ollama(git_diff) -> ResponseData:
    prompt = f"""
    You are a commit message generator.

    Given the following git diff, generate a concise commit message.

    Git diff:
    {git_diff}

    Output MUST be valid JSON only.

    Format:
    {{
        "type": "<ENUM: ADDED | FIXED | CHANGED | REMOVED | REFACTORED | DOCUMENTED | TESTED>",
        "description": "<exactly 10 words summary>"
    }}

    STRICT RULES:
    1. Description MUST contain exactly 10 words.
    2. Description MUST NOT contain ANY word (or variation) from the type.
       - ADDED → forbid: add, added, adding
       - FIXED → forbid: fix, fixed, fixing
       - CHANGED → forbid: change, changed, changing
       - REMOVED → forbid: remove, removed, removing
       - REFACTORED → forbid: refactor, refactored, refactoring
       - DOCUMENTED → forbid: document, documented, documenting
       - TESTED → forbid: test, tested, testing
    3. Do NOT repeat or rephrase the type in description.
    4. Use completely different vocabulary for description.
    5. Output ONLY JSON. No explanation.

    If rules are violated, the output is INVALID.
    """

    response = chat(
        model="phi3:mini",
        messages=[
            {
                "role": "user",
                "content": prompt
            }
        ],
        format="json",
    )


    data = json.loads(response.message.content)
    res = ResponseData(**data)
    return res

def generate_commit_message_groq(git_diff) -> ResponseData:
    prompt = f"""
    You are a commit message generator.

    Given the following git diff, generate a concise commit message.

    Git diff:
    {git_diff}

    Output MUST be valid JSON only.

    Format:
    {{
        "type": "<ENUM: ADDED | FIXED | CHANGED | REMOVED | REFACTORED | DOCUMENTED | TESTED>",
        "description": "<exactly 10 words summary>"
    }}

    STRICT RULES:
    1. Description MUST contain exactly 10 words.
    2. Description MUST NOT contain ANY word (or variation) from the type.
       - ADDED → forbid: add, added, adding
       - FIXED → forbid: fix, fixed, fixing
       - CHANGED → forbid: change, changed, changing
       - REMOVED → forbid: remove, removed, removing
       - REFACTORED → forbid: refactor, refactored, refactoring
       - DOCUMENTED → forbid: document, documented, documenting
       - TESTED → forbid: test, tested, testing
    3. Do NOT repeat or rephrase the type in description.
    4. Use completely different vocabulary for description.
    5. Output ONLY JSON. No explanation.

    If rules are violated, the output is INVALID.
    """

    ai_key = os.getenv("AI_KEY")

    if ai_key is None:
        raise ValueError("AI_KEY environment variable is not set.")

    client = Groq(api_key=ai_key)
    completion = client.chat.completions.create(
        model="llama-3.1-8b-instant",
        messages=[
            {
                "role": "user",
                "content": prompt
            }
        ],
    )
    data = json.loads(completion.choices[0].message.content)
    res = ResponseData(**data)
    return res




def get_refined_commit_message(data: ResponseData) -> str:
    return f"{data.type.upper()}: {data.description.strip()}"

if __name__ == "__main__":
    git_diff = gather_git_info()
    commit_message = get_refined_commit_message(generate_commit_message_groq(git_diff))
    subprocess.run(["git", "add", "*"])
    subprocess.run(["git", "commit", "-m", commit_message])
    subprocess.run(["git", "push", "-u", "origin", "main"])
