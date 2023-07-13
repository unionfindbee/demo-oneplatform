Golang Weather Application with Mayhem CI/CD Integration
This repository offers a hands-on demonstration of how to leverage the power of ForAllSecure's Mayhem for Code in a GitHub Actions workflow. Mayhem for Code is an automated tool that checks your application's binary, packaged as a containerized Docker image, for performance, reliability, and security issues within a CI/CD pipeline.

You can access more detailed instructions about the integration process on the Mayhem for Code GitHub Action page.

Example GitHub Actions Integration
In this updated example, we're working with a Golang-based weather application. This application is built and fuzzed using Mayhem for Code, in order to illustrate how Mayhem can be integrated into a CI/CD workflow.

The Golang weather application is a simple HTTP server that exposes various endpoints to create, retrieve, update, and delete weather records, as well as stream the latest weather updates via WebSockets.

Repository Structure
This repository consists of two branches: main and vulnerable.

The main branch contains the fixed Golang weather application, while the vulnerable branch contains a version of the same application that is deliberately flawed, to illustrate how Mayhem identifies and reports vulnerabilities.

When running the Mayhem for Code GitHub Action, the Golang weather application will be compiled into a Docker image. This image is then pushed to the GitHub Container Registry and ingested by Mayhem for fuzzing and testing. The building of the Docker image leverages a multi-stage Docker image build process.

Getting Started
To get started with this repository:

Fork the repository and register for a Mayhem account.
Visit app.mayhem.security and log in.
Navigate to the bottom left corner and click on your profile icon.
Go to Account Settings -> API Tokens to retrieve your Mayhem account API token.
Save this token in your forked repository's GitHub Secrets as MAYHEM_TOKEN.
In the main branch of your forked repository, go to the Actions tab and execute a CI pipeline. This action will compile the Golang weather application into a Docker image and push it to the GitHub Container Registry. Mayhem will then fuzz this image. As there are no vulnerable versions in the main branch, no issues will be reported in the Security tab.

Note: It may be necessary to adjust your package visibility settings to Public to allow Mayhem to access your Docker image from the GitHub Container Registry. Go to your package in your GitHub repository and select Package Settings. Then, under Package Visibility, set the package to Public.

Switch to the vulnerable branch and create a pull request to merge it with the main branch of your forked repository. The Mayhem for Code GitHub Action will automatically begin, building a Docker image of the vulnerable Golang weather application. Mayhem will fuzz the image and conduct regression and behavior testing on the updated target applications. Detailed results can be found on the PR or the Mayhem server, as well as in the Security tab.

About Us
ForAllSecure was established with a goal to secure the world’s critical software. The company utilizes patented technology from over a decade of CMU research to ensure software safety. ForAllSecure collaborates with Fortune 1000 companies across various sectors, including aerospace, automotive, and high-tech, as well as the US Department of Defense. The company integrates Mayhem into software development cycles for continuous security assurance. ForAllSecure is profitable and funded through revenue, and is expanding rapidly.

Find out more about us here or check out our code security here.