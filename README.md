![Static Badge](https://img.shields.io/badge/backend-golang-blue)
![Static Badge](https://img.shields.io/badge/web-react-red)
![Static Badge](https://img.shields.io/badge/mobile-flutter-purple)
![Static Badge](https://img.shields.io/badge/containerization-docker-yellow)

[![Quality Gate Status](http://mathisbrehm.fr:9000/api/project_badges/measure?project=AREA&metric=alert_status&token=sqb_816c237053393841f07d61210c8f0cb66c77c87d)](http://mathisbrehm.fr:9000/dashboard?id=AREA)
[![Security Rating](http://mathisbrehm.fr:9000/api/project_badges/measure?project=AREA&metric=security_rating&token=sqb_816c237053393841f07d61210c8f0cb66c77c87d)](http://mathisbrehm.fr:9000/dashboard?id=AREA)
[![Maintainability Rating](http://mathisbrehm.fr:9000/api/project_badges/measure?project=AREA&metric=sqale_rating&token=sqb_816c237053393841f07d61210c8f0cb66c77c87d)](http://mathisbrehm.fr:9000/dashboard?id=AREA)

# AREA

AREA is an application development project designed to replicate the functionality of IFTTT. Its principle is straightforward: connecting actions and reactions across different services to automate workflows.

## Project Description

AREA's sole purpose is to make your life easier by automating tasks that can be troublesome or easily forgotten. For example, if you need to send yourself a weekly email to remind you to clean your desk, AREA can handle this task automatically.


> [!NOTE]
> Our product is friendly to the visually impared and the color blind.

## Installation

Our project relies on a Docker container to build. You only have to run these two commands to launch AREA:

```bash
sudo docker-compose build
```

```bash
sudo docker-compose up
```

> [!NOTE]
> If not working, make sure that your docker is running and that no other container is running on the port 8080 or 8081

## Documentation

- The user documentation for web version -> [`Web User Documentation`](docs/webUserDocumentation.pdf)
- The user documentation for mobile version -> [`Mobile User Documentation`](docs/mobileUserDocumentation.pdf)
- Implement a new service, action, reaction -> [`Developer Documentation`](backend/README.md)
- Database structure and relationships -> [`Conceptual Data Model`](backend/docs/conceptualDataModel.pdf)
- REST API overview -> [`REST API Diagram`](backend/docs/restAPIDiagram.pdf)

## Authors

* **Mathis Brehm** - [@Mathis Brehm](https://github.com/mathisbdkr)
* **Eva Legrand** - [@Eva Legrand](https://github.com/Chaegnal)
* **Wilson Bordichon** - [@Wilson Bordichon](https://github.com/Wilson-Epitech)
* **Yanis Harkouk** - [@Yanis Harkouk](https://github.com/yanishk)