# Security Policy

## Supported Versions

We actively support and provide security updates for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| Latest  | Yes                |
| < Latest| No                 |

## Reporting a Vulnerability

We take the security of Vibecheck seriously. If you discover a security vulnerability, please follow these steps:

### How to Report

1. **Do not** open a public GitHub issue for security vulnerabilities
2. Email security details to: [security contact email]
3. Include the following information:
   - Description of the vulnerability
   - Steps to reproduce the issue
   - Potential impact of the vulnerability
   - Suggested fix (if you have one)
   - Your contact information

### What to Expect

- You will receive an acknowledgment of your report within 48 hours
- We will provide an initial assessment within 7 days
- We will keep you informed of our progress
- We will notify you when the vulnerability has been addressed
- We will credit you for the discovery (unless you prefer to remain anonymous)

### Disclosure Policy

- We will work with you to understand and resolve the issue quickly
- Security vulnerabilities will be disclosed publicly after a fix has been released
- We will coordinate with you on the timing of the disclosure
- We will credit you in the security advisory (unless you prefer otherwise)

## Security Best Practices

When using Vibecheck:

- Keep your API keys secure and never commit them to version control
- Use environment variables or the `vibecheck keys` command to manage API keys
- Regularly update to the latest version of Vibecheck
- Review generated commit messages before committing
- Be cautious when using custom prompts that might expose sensitive information

## Known Security Considerations

- API keys are stored locally in `~/.vibecheck_keys.json` - ensure proper file permissions
- Environment variables may be visible in process lists
- Generated commit messages may contain information from your codebase - review before committing
- Network requests are made to external LLM providers - ensure you trust the provider

## Security Updates

Security updates will be released as soon as possible after a vulnerability is confirmed and fixed. We recommend:

- Enabling automatic updates where possible
- Regularly checking for new releases
- Subscribing to security advisories if available

## Scope

The following are considered in-scope for security reporting:

- Remote code execution vulnerabilities
- Authentication and authorization bypasses
- Sensitive data exposure
- API key leakage or improper handling
- Injection vulnerabilities
- Path traversal issues
- Denial of service vulnerabilities

The following are generally considered out-of-scope:

- Issues requiring physical access to the device
- Issues requiring social engineering
- Issues in third-party dependencies (please report to the upstream project)
- Issues that require already compromised user accounts
- Self-XSS vulnerabilities
- Issues requiring very unlikely user interaction

## Contact

For security-related concerns, please contact the maintainers through the appropriate channels as outlined above.

Thank you for helping keep Vibecheck and its users safe.

