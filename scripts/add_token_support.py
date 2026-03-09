#!/usr/bin/env python3
"""
Add GetUserTokenOption() support to all client API calls in feishu-cli.
This script adds user token support to the Go client files.
"""

import os
import re
import sys

# Files to process (relative to project root)
FILES_TO_PROCESS = [
    'internal/client/calendar.go',
    'internal/client/task.go',
    'internal/client/message.go',
    'internal/client/permission.go',
    'internal/client/drive.go',
    'internal/client/export.go',
    'internal/client/comment.go',
    'internal/client/user.go',
    'internal/client/contact.go',
]

def add_token_to_function(content, func_start, func_end):
    """Add token support to a single function."""
    func_content = content[func_start:func_end]

    # Check if already has tokenOpt
    if 'GetUserTokenOption()' in func_content:
        return content

    # Check if function has GetClient() call
    if 'GetClient()' not in func_content:
        return content

    # Check if function has client API call that needs tokenOpt
    if not re.search(r'client\.\w+\.\w+\.\w+\(Context\(\)', func_content):
        return content

    # Pattern to find GetClient() and its error check
    # We need to insert token code after the error check
    patterns = [
        # Pattern 1: return err
        (r'(client, err := GetClient\(\)\n\tif err != nil \{\n\t\treturn ([^\n]+)\n\t\}\n)',
         r'\1\n\ttokenOpt, tokenErr := GetUserTokenOption()\n\tif tokenErr != nil {\n\t\treturn \2\n\t}\n'),
        # Pattern 2: return "", err
        (r'(client, err := GetClient\(\)\n\tif err != nil \{\n\t\treturn "", err\n\t\}\n)',
         r'\1\n\ttokenOpt, tokenErr := GetUserTokenOption()\n\tif tokenErr != nil {\n\t\treturn "", tokenErr\n\t}\n'),
        # Pattern 3: return nil, err
        (r'(client, err := GetClient\(\)\n\tif err != nil \{\n\t\treturn nil, err\n\t\}\n)',
         r'\1\n\ttokenOpt, tokenErr := GetUserTokenOption()\n\tif tokenErr != nil {\n\t\treturn nil, tokenErr\n\t}\n'),
        # Pattern 4: return nil, "", false, err
        (r'(client, err := GetClient\(\)\n\tif err != nil \{\n\t\treturn nil, "", false, err\n\t\}\n)',
         r'\1\n\ttokenOpt, tokenErr := GetUserTokenOption()\n\tif tokenErr != nil {\n\t\treturn nil, "", false, tokenErr\n\t}\n'),
        # Pattern 5: return nil, "", err
        (r'(client, err := GetClient\(\)\n\tif err != nil \{\n\t\treturn nil, "", err\n\t\}\n)',
         r'\1\n\ttokenOpt, tokenErr := GetUserTokenOption()\n\tif tokenErr != nil {\n\t\treturn nil, "", tokenErr\n\t}\n'),
        # Pattern 6: return "", "", err
        (r'(client, err := GetClient\(\)\n\tif err != nil \{\n\t\treturn "", "", err\n\t\}\n)',
         r'\1\n\ttokenOpt, tokenErr := GetUserTokenOption()\n\tif tokenErr != nil {\n\t\treturn "", "", tokenErr\n\t}\n'),
    ]

    modified = func_content
    for pattern, replacement in patterns:
        modified = re.sub(pattern, replacement, modified)

    if modified != func_content:
        return content[:func_start] + modified + content[func_end:]
    return content


def add_token_to_api_calls(content):
    """Add tokenOpt parameter to API calls."""
    # Pattern: client.X.Y.Z(Context(), req) -> client.X.Y.Z(Context(), req, tokenOpt)
    # Pattern: client.X.Y.Z(Context(), reqBuilder.Build()) -> client.X.Y.Z(Context(), reqBuilder.Build(), tokenOpt)

    patterns = [
        (r'client\.(\w+)\.(\w+)\.(\w+)\(Context\(\), ([^)]+)\)',
         r'client.\1.\2.\3(Context(), \4, tokenOpt)'),
    ]

    for pattern, replacement in patterns:
        content = re.sub(pattern, replacement, content)

    return content


def process_file(filepath):
    """Process a single file."""
    print(f"Processing {filepath}...")

    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()

    original = content

    # Add tokenOpt to API calls
    content = add_token_to_api_calls(content)

    # Find all function definitions
    func_pattern = r'// (\w+[^\n]*)\nfunc (\w+)\('

    # Process each function
    # This is a simplified approach - we'll process the whole file

    if content != original:
        # Write backup
        with open(filepath + '.bak', 'w', encoding='utf-8') as f:
            f.write(original)

        # Write modified content
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)

        print(f"  Updated {filepath}")
        return True
    else:
        print(f"  No changes needed for {filepath}")
        return False


def main():
    # Get project root
    script_dir = os.path.dirname(os.path.abspath(__file__))
    project_root = os.path.dirname(script_dir)

    os.chdir(project_root)

    updated = 0
    for filepath in FILES_TO_PROCESS:
        if os.path.exists(filepath):
            if process_file(filepath):
                updated += 1
        else:
            print(f"File not found: {filepath}")

    print(f"\nUpdated {updated} files")


if __name__ == '__main__':
    main()
