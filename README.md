<div id="top"></div>
<!--
*** Thanks for checking out the Best-README-Template. If you have a suggestion
*** that would make this better, please fork the repo and create a pull request
*** or simply open an issue with the tag "enhancement".
*** Don't forget to give the project a star!
*** Thanks again! Now go create something AMAZING! :D
-->

<!-- PROJECT SHIELDS -->
<!--
*** I'm using markdown "reference style" links for readability.
*** Reference links are enclosed in brackets [ ] instead of parentheses ( ).
*** See the bottom of this document for the declaration of the reference variables
*** for contributors-url, forks-url, etc. This is an optional, concise syntax you may use.
*** https://www.markdownguide.org/basic-syntax/#reference-style-links
-->

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]
[![LinkedIn][linkedin-shield]][linkedin-url]

[![Continuous Integration][github-actions-shield-ci]][github-actions-url-ci]
[![Continuous Delivery][github-actions-shield-cd]][github-actions-url-cd]

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <h3 align="center">Forum Group Discussion API</h3>

  <p align="center">
   Forum Group Discussion API: API that provide everything about forum group discussion
    <br />
    <a href="https://moot-rest-api.herokuapp.com/"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://moot-rest-api.herokuapp.com/">View Demo</a>
    ·
    <a href="https://github.com/kelompok-22-capstone-project/forum-group-discussion-backend/issues">Report Bug</a>
    ·
    <a href="https://github.com/kelompok-22-capstone-project/forum-group-discussion-backend/issues">Request Feature</a>
  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->

## About The Project

Forum Group Discussion API is an open source project that provides API (Application Programming Interface) related to forum group discussion.

<p align="right">(<a href="#top">back to top</a>)</p>

### Built With

- [Go](https://golang.org/)
- [PostgreSQL](https://www.postgresql.org/)
- [Echo](https://echo.labstack.com/)

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- GETTING STARTED -->

## Getting Started

### Prerequisites

This is an example of how to list things you need to use the software and how to install them.

- Go
  ```sh
  sudo apt install golang
  ```

### Installation

1. Get a API Documentation at [https://moot-rest-api.herokuapp.com/](https://moot-rest-api.herokuapp.com/)
2. Clone the repo
   ```sh
   git clone git@github.com:kelompok-22-capstone-project/forum-group-discussion-backend.git
   ```
3. Install required dependencies
   ```sh
   go mod tidy
   ```
4. Enter your environment variables in `.env`
   ```bash
   ENV=<ENV>
   PORT=<PORT>
   DB_HOST=<POSTGRESQL_DB_HOST>
   DB_PORT=<POSTGRESQL_PORT>
   DB_USER=<POSTGRESQL_DB_USER>
   DB_PASSWORD=<POSTGRESQL_DB_PASSWORD>
   DB_NAME=<POSTGRESQL_DB_NAME>
   JWT_SECRET=<JWT_SECRET>
   API_KEY=<API_KEY>
   ```
5. Run
   ```sh
   go run main.go
   ```

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- USAGE EXAMPLES -->

## Usage

_For more examples, please refer to the [Documentation](https://moot-rest-api.herokuapp.com/)

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- ROADMAP -->

## Roadmap

- [x] Add Changelog
- [x] Add back to top links
- [ ] Add Additional Templates w/ Examples
- [ ] Add "components" document to easily copy & paste sections of the readme
- [ ] Multi-language Support
  - [ ] Chinese
  - [ ] Spanish

See the [open issues](https://github.com/erikrios/reog-apps-apis/issues) for a full list of proposed features (
and known issues).

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- CONTRIBUTING -->

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any
contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also
simply open an issue with the tag "enhancement". Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- LICENSE -->

## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- CONTACT -->

## Contact

Erik Rio Setiawan - [github](https://github.com/erikrios) - erikriosetiawan15@gmail.com

Project Link: [https://github.com/kelompok-22-capstone-project/forum-group-discussion-backend](https://github.com/kelompok-22-capstone-project/forum-group-discussion-backend)

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- ACKNOWLEDGMENTS -->

## Acknowledgments

Use this space to list resources you find helpful and would like to give credit to. I've included a few of my favorites
to kick things off!

- [Choose an Open Source License](https://choosealicense.com)
- [GitHub Emoji Cheat Sheet](https://www.webpagefx.com/tools/emoji-cheat-sheet)
- [Malven's Flexbox Cheatsheet](https://flexbox.malven.co/)
- [Malven's Grid Cheatsheet](https://grid.malven.co/)
- [Img Shields](https://shields.io)
- [GitHub Pages](https://pages.github.com)
- [Font Awesome](https://fontawesome.com)
- [React Icons](https://react-icons.github.io/react-icons/search)

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[github-actions-shield-ci]: https://github.com/kelompok-22-capstone-project/forum-group-discussion-backend/actions/workflows/ci.yml/badge.svg
[github-actions-url-ci]: https://github.com/kelompok-22-capstone-project/forum-group-discussion-backend/actions/workflows/ci.yml
[github-actions-shield-cd]: https://github.com/kelompok-22-capstone-project/forum-group-discussion-backend/actions/workflows/cd.yml/badge.svg
[github-actions-url-cd]: https://github.com/kelompok-22-capstone-project/forum-group-discussion-backend/actions/workflows/cd.yml
[contributors-shield]: https://img.shields.io/github/contributors/kelompok-22-capstone-project/forum-group-discussion-backend.svg?style=for-the-badge
[contributors-url]: https://github.com/kelompok-22-capstone-project/forum-group-discussion-backend/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/kelompok-22-capstone-project/forum-group-discussion-backend.svg?style=for-the-badge
[forks-url]: https://github.com/kelompok-22-capstone-project/forum-group-discussion-backend/network/members
[stars-shield]: https://img.shields.io/github/stars/kelompok-22-capstone-project/forum-group-discussion-backend.svg?style=for-the-badge
[stars-url]: https://github.com/kelompok-22-capstone-project/forum-group-discussion-backend/stargazers
[issues-shield]: https://img.shields.io/github/issues/kelompok-22-capstone-project/forum-group-discussion-backend.svg?style=for-the-badge
[issues-url]: https://github.com/kelompok-22-capstone-project/forum-group-discussion-backend/issues
[license-shield]: https://img.shields.io/github/license/kelompok-22-capstone-project/forum-group-discussion-backend.svg?style=for-the-badge
[license-url]: https://github.com/kelompok-22-capstone-project/forum-group-discussion-backend/blob/master/LICENSE.txt
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://linkedin.com/in/erikriosetiawan
[product-screenshot]: images/screenshot.png
