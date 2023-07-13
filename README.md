# Mayhem CI/CD and API Integration

This repository offers a comprehensive example of leveraging the power of ForAllSecure's Mayhem for both Code and API in a GitHub Actions workflow. Mayhem is an automated tool that checks your application's binary, packaged as a containerized Docker image, for performance, reliability, and security issues within a CI/CD pipeline.

The project includes a Golang-based weather application built and fuzzed using Mayhem for Code to illustrate its CI/CD integration. Additionally, it demonstrates how Mayhem's API feature, a dynamic testing tool, can catch reliability, performance, and security bugs in your APIs before they hit production.

For more detailed instructions about the integration process, check out the [Mayhem for Code GitHub Action page](https://github.com/ForAllSecure/mayhem-code-action) and the [Mayhem for API GitHub Action page](https://github.com/ForAllSecure/mayhem-api-action).

## Example GitHub Actions Integration

The Golang weather application is a simple HTTP server that exposes various endpoints to create, retrieve, update, and delete weather records, as well as stream the latest weather updates via WebSockets. This application is tested by both Mayhem for Code and Mayhem for API.

This repository consists of two branches: `main` and `vulnerable`. The `main` branch contains the fixed Golang weather application, while the `vulnerable` branch contains a version of the same application that is deliberately flawed, to illustrate how Mayhem identifies and reports vulnerabilities.

### Launching Mayhem for Code or API in GitHub Actions

1. Fork the repository and register for a Mayhem account.
2. Visit [app.mayhem.security](https://app.mayhem.security) and log in.
3. Navigate to the bottom left corner and click on your profile icon.
4. Go to Account Settings -> API Tokens to retrieve your Mayhem account API token.
5. Save this token in your forked repository's GitHub Secrets as `MAYHEM_TOKEN`.
6. In the main branch of your forked repository, go to the Actions tab and execute a CI pipeline. This action will compile the Golang weather application into a Docker image and push it to the GitHub Container Registry. Mayhem will then fuzz this image. As there are no vulnerable versions in the main branch, no issues will be reported in the Security tab.

Note: You may need to adjust your package visibility settings to Public to allow Mayhem to access your Docker image from the GitHub Container Registry. Go to your package in your GitHub repository and select Package Settings. Then, under Package Visibility, set the package to Public.

7. Switch to the vulnerable branch and create a pull request to merge it with the main branch of your forked repository. The Mayhem for Code and API GitHub Actions will automatically begin, building a Docker image of the vulnerable Golang weather application. Mayhem will fuzz the image and conduct regression and behavior testing on the updated target applications. Detailed results can be found on the PR or the Mayhem server, as well as in the Security tab.

### Reports

Both Mayhem for Code and API generate reports when you pass `sarif-report` or `html-report` to the input. Be sure to include `continue-on-error` in the Mayhem for API step if you want to process the reports in follow-up steps.

To add the report to your build artifacts, include the archive step in your pipeline:

```
# Archive HTML report
- name: Archive Mayhem for API report
  uses: actions/upload-artifact@v3
  with:
    name: mapi-report
    path: mapi.html
```

Uploading SARIF reports to GitHub allows you to view any issue found by Mayhem for API directly in your PR, as well as in the "Security" tab of your repository. This feature currently requires you to have a GitHub Enterprise Plan or a public repository. To upload the SARIF report, include this step in your pipeline:

```
# Upload SARIF file (only available on public repos or GitHub Enterprise)
- name: Upload SARIF file
  uses: github/codeql-action/upload-sarif@v2
  with:
    sarif_file: mapi.sarif
```

If your API server sends back stacktraces in the 500 Internal Server Error responses (for test environments only -- never in production!), Mayhem for API will try to map issues it finds to the exact line of code that triggered the issue.

## About Us

ForAllSecure was established with a goal to secure the world’s critical software. The company utilizes patented technology from over a decade of CMU research to ensure software safety. ForAllSecure collaborates with Fortune 1000 companies across various sectors, including aerospace, automotive, and high-tech, as well as the US Department of Defense. The company integrates Mayhem into software development cycles for continuous security assurance. ForAllSecure is profitable and funded through revenue, and is expanding rapidly.

Find out more about us [here](https://www.mayhem.security/) or check out our [code security here](https://www.mayhem.security/).